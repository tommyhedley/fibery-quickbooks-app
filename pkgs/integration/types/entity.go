package types

import (
	"fmt"

	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"
	"github.com/tommyhedley/fibery-quickbooks-app/pkgs/integration"
)

var entity = integration.NewUnionType(
	[]integration.DualType{customer, vendor, employee},
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
				case "customer":
					if id, ok = m["id"].(string); !ok {
						return nil, fmt.Errorf("unable to extract 'id' from Customer item")
					}
					return "c:" + id, nil
				case "vendor":
					if id, ok = m["id"].(string); !ok {
						return nil, fmt.Errorf("unable to extract 'id' from Customer item")
					}
					return "v:" + id, nil
				case "employee":
					if id, ok = m["id"].(string); !ok {
						return nil, fmt.Errorf("unable to extract 'id' from Customer item")
					}
					return "e:" + id, nil
				default:
					return nil, nil
				}
			},
		},
		"QBOId": {
			Params: fibery.Field{
				Name: "QBO ID",
				Type: fibery.Text,
			},
			Convert: func(s string, m map[string]any) (any, error) {
				id, ok := m["id"].(string)
				if !ok {
					return nil, fmt.Errorf("unable to extract 'id' from Customer item")
				}
				return id, nil
			},
		},
		"Name": {
			Params: fibery.Field{
				Name:    "Name",
				Type:    fibery.Text,
				SubType: fibery.Title,
			},
			Convert: func(s string, m map[string]any) (any, error) {
				var name string
				var ok bool

				switch s {
				case "customer":
					if name, ok = m["name"].(string); !ok {
						return nil, fmt.Errorf("unable to extract 'name' from Customer item")
					}
					return name + "(Customer)", nil
				case "vendor":
					if name, ok = m["name"].(string); !ok {
						return nil, fmt.Errorf("unable to extract 'name' from Customer item")
					}
					return name + "(Vendor)", nil
				case "employee":
					if name, ok = m["name"].(string); !ok {
						return nil, fmt.Errorf("unable to extract 'name' from Customer item")
					}
					return name + "(Employee)", nil
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
				syncAction, ok := m["__syncAction"].(string)
				if !ok {
					return nil, fmt.Errorf("unable to extract 'id' from Customer item")
				}
				return syncAction, nil
			},
		},
		"CustomerId": {
			Params: fibery.Field{
				Name: "Customer ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.OTO,
					Name:          "Customer",
					TargetName:    "Entity",
					TargetType:    "vustomer",
					TargetFieldID: "id",
				},
			},
			Convert: func(s string, m map[string]any) (any, error) {
				if s == "customer" {
					id, ok := m["id"].(string)
					if !ok {
						return nil, fmt.Errorf("unable to extract 'id' from Customer item")
					}
					return id, nil
				} else {
					return nil, nil
				}
			},
		},
		"EmployeeId": {
			Params: fibery.Field{
				Name: "Employee ID",
				Type: fibery.Text,
				Relation: &fibery.Relation{
					Cardinality:   fibery.OTO,
					Name:          "Employee",
					TargetName:    "Entity",
					TargetType:    "vmployee",
					TargetFieldID: "id",
				},
			},
			Convert: func(s string, m map[string]any) (any, error) {
				if s == "employee" {
					id, ok := m["id"].(string)
					if !ok {
						return nil, fmt.Errorf("unable to extract 'id' from Customer item")
					}
					return id, nil
				} else {
					return nil, nil
				}
			},
		},
		"VendorId": {
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
				if s == "vendor" {
					id, ok := m["id"].(string)
					if !ok {
						return nil, fmt.Errorf("unable to extract 'id' from Customer item")
					}
					return id, nil
				} else {
					return nil, nil
				}
			},
		},
	},
)
