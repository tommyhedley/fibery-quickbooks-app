package app

import (
	"fmt"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/quickbooks-go"
)

func attachableField(schema map[string]fibery.Field, attachablesFieldId string) bool {
	for fieldId, field := range schema {
		if field.SubType == fibery.File && fieldId == attachablesFieldId {
			return true
		}
	}
	return false
}

func AttachableURL(a quickbooks.Attachable) string {
	return fmt.Sprintf("app://resource?type=%s&id=%s", "attachable", a.Id)
}

func indexAttachables(
	entityType string,
	attachables []quickbooks.Attachable,
	idx map[string][]quickbooks.Attachable,
	pageSize int,
) (map[string][]quickbooks.Attachable, bool) {
	if idx == nil {
		idx = make(map[string][]quickbooks.Attachable, len(attachables))
	}

	for _, att := range attachables {
		for _, ref := range att.AttachableRef {
			if ref.EntityRef.Type != entityType {
				continue
			}
			key := ref.EntityRef.Value
			idx[key] = append(idx[key], att)
		}
	}

	return idx, len(attachables) == pageSize
}
