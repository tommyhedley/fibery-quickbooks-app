package app

import "github.com/tommyhedley/fibery-quickbooks-app/pkgs/fibery"

type ActionRegistry map[string]fibery.Action

func (ar ActionRegistry) Register(a fibery.Action) {
	ar[a.ActionId] = a
}

func (ar ActionRegistry) Get(id string) (fibery.Action, bool) {
	if action, exists := ar[id]; exists {
		return action, true
	}
	return fibery.Action{}, false
}

func (ar ActionRegistry) GetAll() []fibery.Action {
	actions := make([]fibery.Action, 0, len(ar))
	for _, action := range ar {
		actions = append(actions, action)
	}
	return actions
}

var Actions = make(ActionRegistry)

var testAction = fibery.Action{
	ActionId:    "test",
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
}

func init() {
	Actions.Register(testAction)
}
