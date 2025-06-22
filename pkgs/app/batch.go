package app

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tommyhedley/quickbooks-go"
)

func EncodeQueryBID(entityType string, page int, attachable bool) string {
	if attachable {
		return fmt.Sprintf("attachable:%s:%d", entityType, page)
	} else {
		return fmt.Sprintf("%s:%d", entityType, page)
	}
}

func DecodeQueryBID(BID string) (string, int, bool, error) {
	parts := strings.Split(BID, ":")
	if len(parts) > 3 || len(parts) < 2 {
		return "", 0, false, fmt.Errorf("invalid BID passed, expecting 2 or 3 parts, received: %d", len(parts))
	}

	if len(parts) == 3 && parts[0] == "attachable" {
		startPosition, err := strconv.Atoi(parts[2])
		if err != nil {
			return "", 0, false, fmt.Errorf("unable to convert part %d in  BID to int: %w", 3, err)
		}

		return parts[1], startPosition, true, nil
	}

	startPosition, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, false, fmt.Errorf("unable to convert part %d in  BID to int: %w", 3, err)
	}

	return parts[0], startPosition, false, nil
}

func batchQueryRequest(entityType string, ids []string, page, pageSize int, attachable bool) quickbooks.BatchItemRequest {
	if attachable {
		if len(ids) > 0 {
			idString := "'" + strings.Join(ids, "','") + "'"
			return quickbooks.BatchItemRequest{
				BID:   EncodeQueryBID(entityType, page, true),
				Query: fmt.Sprintf("Select Id, AttachableRef From Attachable Where AttachableRef.EntityRef.Type = '%s' And AttachableRef.EntityRef.Value in (%s) STARTPOSITION %d MAXRESULTS %d", entityType, idString, startPosition(page, pageSize), pageSize),
			}
		} else {
			return quickbooks.BatchItemRequest{
				BID:   EncodeQueryBID(entityType, page, true),
				Query: fmt.Sprintf("Select Id, AttachableRef From Attachable Where AttachableRef.EntityRef.Type = '%s' STARTPOSITION %d MAXRESULTS %d", entityType, startPosition(page, pageSize), pageSize),
			}
		}
	} else {
		if len(ids) > 0 {
			idString := "'" + strings.Join(ids, "','") + "'"
			return quickbooks.BatchItemRequest{
				BID:   EncodeQueryBID(entityType, page, false),
				Query: fmt.Sprintf("Select * From %s Where Id in (%s) ORDERBY Id STARTPOSITION %d MAXRESULTS %d", entityType, idString, startPosition(page, pageSize), pageSize),
			}
		} else {
			return quickbooks.BatchItemRequest{
				BID:   EncodeQueryBID(entityType, page, false),
				Query: fmt.Sprintf("Select * From %s ORDERBY Id STARTPOSITION %d MAXRESULTS %d", entityType, startPosition(page, pageSize), pageSize),
			}
		}
	}
}
