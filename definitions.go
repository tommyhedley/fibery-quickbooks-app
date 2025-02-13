package main

import (
	"github.com/tommyhedley/quickbooks-go"
)

type integrationAccountInfo struct {
	Name    string `json:"name,omitempty"`
	RealmID string `json:"realmId,omitempty"`
	quickbooks.BearerToken
}
