package types

import (
	"fmt"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/integration"
)

var entity = integration.NewUnionType(
	[]integration.StandardType{customer, vendor, employee},
	"entity",
	"Entity",
	map[string]integration.UnionFieldDef{
		"id": {
			Params: fibery.Field{
				Name: "ID",
				Type: fibery.Id,
			},
			Convert: func(s string, m map[string]any) (any, error) {
				var id string
				var ok bool

				switch s {
				case "Customer":
					if id, ok = m["id"].(string); !ok {
						return nil, fmt.Errorf("unable to extract 'id' from customer item")
					}
					return "c:" + id, nil
				case "Vendor":
					if id, ok = m["id"].(string); !ok {
						return nil, fmt.Errorf("unable to extract 'id' from vendor item")
					}
					return "v:" + id, nil
				case "Employee":
					if id, ok = m["id"].(string); !ok {
						return nil, fmt.Errorf("unable to extract 'id' from employee item")
					}
					return "e:" + id, nil
				default:
					return nil, nil
				}
			},
		},
		"qboId": {
			Params: fibery.Field{
				Name:     "QBO ID",
				Type:     fibery.Text,
				ReadOnly: true,
			},
			Convert: func(s string, m map[string]any) (any, error) {
				id, ok := m["id"].(string)
				if !ok {
					return nil, fmt.Errorf("unable to extract 'id' for qboId from %s item", s)
				}
				return id, nil
			},
		},
		"name": {
			Params: fibery.Field{
				Name:    "Name",
				Type:    fibery.Text,
				SubType: fibery.Title,
			},
			Convert: func(s string, m map[string]any) (any, error) {
				var name string
				var ok bool

				switch s {
				case "Customer":
					if name, ok = m["displayName"].(string); !ok {
						return nil, fmt.Errorf("unable to extract 'displayName' from customer item")
					}
					return name + " (Customer)", nil
				case "Vendor":
					if name, ok = m["displayName"].(string); !ok {
						return nil, fmt.Errorf("unable to extract 'displayName' from vendor item")
					}
					return name + " (Vendor)", nil
				case "Employee":
					if name, ok = m["displayName"].(string); !ok {
						return nil, fmt.Errorf("unable to extract 'displayName' from employee item")
					}
					return name + " (Employee)", nil
				default:
					return nil, nil
				}
			},
		},
		"__syncAction": {
			Params: fibery.Field{
				Type: fibery.Text,
				Name: "Sync Action",
			},
			Convert: func(s string, m map[string]any) (any, error) {
				syncAction, ok := m["__syncAction"].(fibery.SyncAction)
				if !ok {
					return nil, fmt.Errorf("unable to extract '__syncAction' from %s item", s)
				}
				return syncAction, nil
			},
		},
		"customerId": {
			Params: fibery.Field{
				Name: "Customer ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.OTO,
					Name:          "Customer",
					TargetName:    "Entity",
					TargetType:    "customer",
					TargetFieldID: "id",
				},
			},
			Convert: func(s string, m map[string]any) (any, error) {
				if s == "Customer" {
					id, ok := m["id"].(string)
					if !ok {
						return nil, fmt.Errorf("unable to extract 'id' for customerId from customer item")
					}
					return id, nil
				} else {
					return nil, nil
				}
			},
		},
		"employeeId": {
			Params: fibery.Field{
				Name: "Employee ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.OTO,
					Name:          "Employee",
					TargetName:    "Entity",
					TargetType:    "employee",
					TargetFieldID: "id",
				},
			},
			Convert: func(s string, m map[string]any) (any, error) {
				if s == "Employee" {
					id, ok := m["id"].(string)
					if !ok {
						return nil, fmt.Errorf("unable to extract 'id' for employeeId from employee item")
					}
					return id, nil
				} else {
					return nil, nil
				}
			},
		},
		"vendorId": {
			Params: fibery.Field{
				Name: "Vendor ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.OTO,
					Name:          "Vendor",
					TargetName:    "Entity",
					TargetType:    "vendor",
					TargetFieldID: "id",
				},
			},
			Convert: func(s string, m map[string]any) (any, error) {
				if s == "Vendor" {
					id, ok := m["id"].(string)
					if !ok {
						return nil, fmt.Errorf("unable to extract 'id' for vendorId from vendor item")
					}
					return id, nil
				} else {
					return nil, nil
				}
			},
		},
	},
)

func init() {
	integration.Types.Register(entity)
}
