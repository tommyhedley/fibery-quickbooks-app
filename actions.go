package main

import "github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"

func BuildActions() []fibery.Action {
	return []fibery.Action{
		{
			ActionID:    "test",
			Name:        "Test Action",
			Description: "This is a text action",
			Args: []fibery.ActionArg{
				{
					Id:           "a1",
					Name:         "Arg 1",
					Description:  "The first argument",
					ArgType:      fibery.TextArg,
					TextTemplate: false,
				},
				{
					Id:           "a2",
					Name:         "Arg 2",
					Description:  "The second argument",
					ArgType:      fibery.TextArg,
					TextTemplate: true,
				},
				{
					Id:           "a3",
					Name:         "Arg 3",
					Description:  "The third argument",
					ArgType:      fibery.TextAreaArg,
					TextTemplate: false,
				},
				{
					Id:           "a4",
					Name:         "Arg 4",
					Description:  "The fourth argument",
					ArgType:      fibery.TextAreaArg,
					TextTemplate: true,
				},
			},
		},
	}
}
