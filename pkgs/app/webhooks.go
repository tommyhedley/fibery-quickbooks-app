package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

type WebhookRequest struct {
	Params struct {
		Connection                      string    `json:"connection"`
		XForwardedPort                  string    `json:"x-forwarded-port"`
		XForwardedPath                  string    `json:"x-forwarded-path"`
		XForwardedPrefix                string    `json:"x-forwarded-prefix"`
		XRealIP                         string    `json:"x-real-ip"`
		UserAgent                       string    `json:"user-agent"`
		ContentType                     string    `json:"content-type"`
		Accept                          string    `json:"accept"`
		IntuitSignature                 string    `json:"intuit-signature"`
		IntuitCreatedTime               time.Time `json:"intuit-created-time"`
		IntuitTId                       string    `json:"intuit-t-id"`
		IntuitNotificationSchemaVersion string    `json:"intuit-notification-schema-version"`
		AcceptEncoding                  string    `json:"accept-encoding"`
		Authorization                   string    `json:"authorization"`
	} `json:"params"`
	Types   []string              `json:"types"`
	Filter  map[string]any        `json:"filter"`
	Account QuickBooksAccountInfo `json:"account"`
	Payload struct {
		EventNotifications []struct {
			RealmId         string `json:"realmId"`
			DataChangeEvent struct {
				Entities []struct {
					Id          string    `json:"id"`
					Operation   string    `json:"operation"`
					Name        string    `json:"name"`
					LastUpdated time.Time `json:"lastUpdated"`
				} `json:"entities"`
			} `json:"dataChangeEvent"`
		} `json:"eventNotifications"`
	} `json:"payload"`
}

type WebhookUpdatedSource struct {
	ids           []string
	batchData     *quickbooks.BatchItemResponse
	getAttachable bool
	attachables   map[string][]quickbooks.Attachable
}

type RelatedType struct {
	typ           CDCType
	getAttachable bool
	attachables   map[string][]quickbooks.Attachable
}

type WebhookGroup struct {
	sync.Mutex
	webhookTypes      map[string]fibery.Type
	relatedTypes      map[string]*RelatedType
	updatedSources    map[string]*WebhookUpdatedSource
	deletedSources    map[string][]string
	changeDataCapture *quickbooks.ChangeDataCapture
	oldestChange      time.Time
	cdcLookback       time.Duration
	idCache           *IdCache
	account           QuickBooksAccountInfo
	integration       *Integration
}

func buildWebhookGroup(req WebhookRequest, i *Integration, cdcLookback time.Duration) (*WebhookGroup, error) {
	idCache, exists := i.idStore.GetOrCreateIdCache(req.Account.RealmId)
	if !exists {
		return nil, fmt.Errorf("no idCache was found for realmId: %s, perform a full sync before enabling webhooks", req.Account.RealmId)
	}

	group := &WebhookGroup{
		webhookTypes:   make(map[string]fibery.Type),
		relatedTypes:   make(map[string]*RelatedType),
		updatedSources: make(map[string]*WebhookUpdatedSource),
		deletedSources: make(map[string][]string),
		oldestChange:   time.Now(),
		cdcLookback:    cdcLookback,
		idCache:        idCache,
		account:        req.Account,
		integration:    i,
	}

	allSources := make(map[string]bool)
	relatedTypesBySource := make(map[string]map[string]*RelatedType)

	attachableFieldId := i.config.AttachableFieldId

	reqTypes := make(map[string]struct{}, len(req.Types))

	for _, typeId := range req.Types {
		reqTypes[typeId] = struct{}{}
	}

	for typeId := range reqTypes {
		regType, ok := i.types.Get(typeId)
		if !ok {
			return nil, fmt.Errorf("requestedType: %s not found", typeId)
		} else {
			switch t := regType.(type) {
			case UnionType:
				if t.Webhook() {
					group.webhookTypes[t.Id()] = regType

					for _, sourceType := range t.Types() {
						_, ok = allSources[sourceType.Type()]
						if !ok {
							allSources[sourceType.Type()] = false
						}
					}
				}
			case WebhookDependentType:
				group.webhookTypes[t.Id()] = regType

				_, ok = allSources[t.SourceType()]
				if !ok {
					allSources[t.SourceType()] = false
				}
			case WebhookType:
				group.webhookTypes[t.Id()] = regType

				_, ok = allSources[t.Type()]
				if !ok {
					allSources[t.Type()] = t.Attachables(attachableFieldId)
				}

				if t.Attachables(attachableFieldId) {
					allSources[t.Type()] = true
				}

				relatedTypes, ok := relatedTypesBySource[t.Type()]
				if !ok {
					relatedTypes = make(map[string]*RelatedType)
				}

				for typeId, relatedType := range t.RelatedTypes() {
					if _, ok := reqTypes[typeId]; ok {
						relatedTypes[typeId] = &RelatedType{
							typ:           relatedType,
							getAttachable: relatedType.Attachables(attachableFieldId),
						}
					}
				}

				relatedTypesBySource[t.Type()] = relatedTypes
			}
		}
	}

	for _, event := range req.Payload.EventNotifications {
		if event.RealmId != req.Account.RealmId {
			continue
		}
		for _, entity := range event.DataChangeEvent.Entities {
			if getAttach, exists := allSources[entity.Name]; exists {
				switch entity.Operation {
				case "Create", "Update", "Emailed", "Void":
					updateSource, ok := group.updatedSources[entity.Name]
					if !ok {
						updateSource = &WebhookUpdatedSource{}
					}

					updateSource.ids = append(updateSource.ids, entity.Id)

					if getAttach {
						updateSource.getAttachable = getAttach
					}

					group.updatedSources[entity.Name] = updateSource

				case "Delete", "Merge":
					deleteSource, ok := group.deletedSources[entity.Name]
					if !ok {
						deleteSource = make([]string, 0, 1)
					}

					deleteSource = append(deleteSource, entity.Id)
					group.deletedSources[entity.Name] = deleteSource
				}

				if relatedTypes, exists := relatedTypesBySource[entity.Name]; exists {
					for typeId, relatedType := range relatedTypes {
						group.relatedTypes[typeId] = relatedType
					}

					if group.oldestChange.IsZero() || entity.LastUpdated.Before(group.oldestChange) {
						group.oldestChange = entity.LastUpdated
					}
				}
			}
		}
	}

	return group, nil
}

func (wg *WebhookGroup) indexBatch(batch []quickbooks.BatchItemResponse) error {
	wg.Lock()
	defer wg.Unlock()
	for _, resp := range batch {
		faults := resp.Fault.Faults
		if len(faults) > 0 {
			return fmt.Errorf("fault for %s: %w", resp.BID, quickbooks.BatchError{Faults: faults})
		}

		entityType, _, isAttachable, err := DecodeQueryBID(resp.BID)
		if err != nil {
			return fmt.Errorf("error decoding query BID: %w", err)
		}

		relatedType, exists := wg.relatedTypes[entityType]

		if exists && isAttachable {
			attachables := resp.QueryResponse.Attachable
			if len(attachables) > 0 {
				existing := relatedType.attachables
				updated, more := indexAttachables(
					entityType, attachables, existing, wg.integration.config.QuickBooks.PageSize,
				)
				relatedType.attachables = updated
				if more {
					return fmt.Errorf("more than %d attachables returned for: %s", wg.integration.config.QuickBooks.PageSize, entityType)
				}

			}
			continue
		}

		updatedSource, exists := wg.updatedSources[entityType]
		if !exists {
			return fmt.Errorf("no updatedSources found for: %s", entityType)
		}

		if isAttachable {
			attachables := resp.QueryResponse.Attachable
			if len(attachables) > 0 {
				existing := updatedSource.attachables
				updated, more := indexAttachables(
					entityType, attachables, existing, wg.integration.config.QuickBooks.PageSize,
				)
				updatedSource.attachables = updated
				if more {
					return fmt.Errorf("more than %d attachables returned for: %s", wg.integration.config.QuickBooks.PageSize, entityType)
				}
			}
		} else {
			if updatedSource.batchData != nil {
				return fmt.Errorf("a batch response entry already exists for %s", entityType)
			}

			updatedSource.batchData = &resp
		}
	}

	return nil
}

func (wg *WebhookGroup) doBatch(req []quickbooks.BatchItemRequest, params quickbooks.RequestParameters) error {
	client := wg.integration.client
	batch, err := client.BatchRequest(params, req)
	if err != nil {
		return fmt.Errorf("error fetching inital batch: %w", err)
	}

	err = wg.indexBatch(batch)
	if err != nil {
		return fmt.Errorf("error indexing webhook batch: %w", err)
	}

	return nil
}

func (wg *WebhookGroup) doCDC(req []string, params quickbooks.RequestParameters) error {
	client := wg.integration.client

	reqTime := wg.oldestChange.Add(-wg.cdcLookback)

	cdc, err := client.ChangeDataCapture(params, req, reqTime)
	if err != nil {
		return fmt.Errorf("error fetching cdc: %w", err)
	}

	wg.Lock()
	wg.changeDataCapture = &cdc
	wg.Unlock()

	return nil
}

func (wg *WebhookGroup) fetchAll(ctx context.Context) error {
	var (
		fetch sync.WaitGroup
		once  sync.Once
		first error
	)

	page := 1
	pageSize := wg.integration.config.QuickBooks.PageSize

	batchReq := make([]quickbooks.BatchItemRequest, 0, len(wg.updatedSources))

	for sourceType, source := range wg.updatedSources {
		if source.getAttachable {
			batchReq = append(batchReq, batchQueryRequest(sourceType, source.ids, page, pageSize, true))
		}
		batchReq = append(batchReq, batchQueryRequest(sourceType, source.ids, page, pageSize, false))
	}

	cdcReq := make([]string, 0, len(wg.relatedTypes))

	for _, rlType := range wg.relatedTypes {
		if rlType.getAttachable {
			batchReq = append(batchReq, batchQueryRequest(rlType.typ.Type(), nil, page, pageSize, true))
		}
		cdcReq = append(cdcReq, rlType.typ.Type())
	}

	record := func(err error) {
		once.Do(func() { first = err })
	}

	params := quickbooks.RequestParameters{
		Ctx:             ctx,
		RealmId:         wg.account.RealmId,
		Token:           &wg.account.BearerToken,
		WaitOnRateLimit: true,
	}

	if len(cdcReq) > 0 {
		fmt.Println(cdcReq)
		fetch.Add(1)
		go func() {
			defer fetch.Done()
			if err := wg.doCDC(cdcReq, params); err != nil {
				record(err)
			}
		}()
	}

	if len(batchReq) > 0 {
		fmt.Println(batchReq)
		fetch.Add(1)
		go func() {
			defer fetch.Done()
			if err := wg.doBatch(batchReq, params); err != nil {
				record(err)
			}
		}()
	}

	fetch.Wait()

	return first
}

func (wg *WebhookGroup) process() (map[string][]map[string]any, error) {
	output := map[string][]map[string]any{}
	pageSize := wg.integration.config.QuickBooks.PageSize
	for typeId, regType := range wg.webhookTypes {
		switch t := regType.(type) {
		case UnionType:
			batchResponses := make(map[string]*quickbooks.BatchItemResponse)
			sourceDeletions := make(map[string][]string)
			for _, sourceType := range t.Types() {
				source, ok := wg.updatedSources[sourceType.Type()]
				if !ok {
					continue
				}

				batchResponses[sourceType.Type()] = source.batchData

				deleteIds, ok := wg.deletedSources[sourceType.Type()]
				if !ok {
					continue
				}

				sourceDeletions[sourceType.Type()] = deleteIds
			}

			updateItems, moreSource, err := t.ProcessBatchQuery(batchResponses, pageSize)
			if err != nil {
				return nil, fmt.Errorf("error processing batch query: %w", err)
			}

			if len(moreSource) > 0 {
				return nil, fmt.Errorf("%s batch response returned more than 1 page", typeId)
			}

			deleteItems, err := t.ProcessWebhookDeletions(sourceDeletions)
			if err != nil {
				return nil, fmt.Errorf("error processing deletedItems: %w", err)
			}

			responseLength := len(updateItems) + len(deleteItems)

			if responseLength > 0 {

				typeOutout, ok := output[typeId]
				if !ok {
					typeOutout = make([]map[string]any, 0, responseLength)
				}

				typeOutout = append(typeOutout, updateItems...)
				typeOutout = append(typeOutout, deleteItems...)

				output[typeId] = typeOutout
			}
		case WebhookDependentType:
			source, ok := wg.updatedSources[t.SourceType()]
			if !ok {
				continue
			}

			updateItems, more, err := t.ProcessBatchQuery(source.batchData, wg.idCache, pageSize)
			if err != nil {
				return nil, fmt.Errorf("error processing batch query: %w", err)
			}

			if more {
				return nil, fmt.Errorf("%s batch response returned more than 1 page", typeId)
			}

			deleteIds, ok := wg.deletedSources[t.SourceType()]
			if !ok {
				deleteIds = make([]string, 0)
			}

			deleteItems, err := t.ProcessWebhookDeletions(deleteIds, wg.idCache)
			if err != nil {
				return nil, fmt.Errorf("error processing deletedItems: %w", err)
			}

			typeOutout, ok := output[typeId]
			if !ok {
				typeOutout = make([]map[string]any, 0, len(updateItems)+len(deleteItems))
			}

			typeOutout = append(typeOutout, updateItems...)
			typeOutout = append(typeOutout, deleteItems...)

			output[typeId] = typeOutout
		case WebhookType:
			source, ok := wg.updatedSources[t.Type()]
			if !ok {
				continue
			}

			updateItems, more, err := t.ProcessBatchQuery(source.batchData, source.attachables, pageSize)
			if err != nil {
				return nil, fmt.Errorf("error processing batch query: %w", err)
			}

			if more {
				return nil, fmt.Errorf("%s batch response returned more than 1 page", typeId)
			}

			deleteIds, ok := wg.deletedSources[t.Type()]
			if !ok {
				deleteIds = make([]string, 0)
			}

			deleteItems, err := t.ProcessWebhookDeletions(deleteIds)
			if err != nil {
				return nil, fmt.Errorf("error processing deletedItems: %w", err)
			}

			typeOutout, ok := output[typeId]
			if !ok {
				typeOutout = make([]map[string]any, 0, len(updateItems)+len(deleteItems))
			}

			typeOutout = append(typeOutout, updateItems...)
			typeOutout = append(typeOutout, deleteItems...)

			output[typeId] = typeOutout
		}
	}

	for typeId, rlType := range wg.relatedTypes {
		items, err := rlType.typ.ProcessCDCQuery(wg.changeDataCapture, rlType.attachables, pageSize)
		if err != nil {
			return nil, fmt.Errorf("error processing changeDataCapture for %s", typeId)
		}

		typeOutout, ok := output[typeId]
		if !ok {
			typeOutout = make([]map[string]any, 0, len(items))
		}

		typeOutout = append(typeOutout, items...)
		output[typeId] = typeOutout
	}

	return output, nil
}
