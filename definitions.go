package main

import (
	"github.com/tommyhedley/fibery/fibery-qbo-integration/pkgs/qbo"
)

type integrationAccountInfo struct {
	Name    string `json:"name,omitempty"`
	RealmID string `json:"realmId,omitempty"`
	qbo.BearerToken
}
