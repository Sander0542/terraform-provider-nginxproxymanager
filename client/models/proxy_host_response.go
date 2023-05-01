package models

import (
	"encoding/json"
)

type ProxyHostResponse struct {
	ID                    int64               `json:"id"`
	CreatedOn             string              `json:"created_on"`
	ModifiedOn            string              `json:"modified_on"`
	OwnerUserID           int64               `json:"owner_user_id"`
	DomainNames           []string            `json:"domain_names"`
	ForwardScheme         string              `json:"forward_scheme"`
	ForwardHost           string              `json:"forward_host"`
	ForwardPort           uint16              `json:"forward_port"`
	CertificateID         json.RawMessage     `json:"certificate_id"`
	SSLForced             boolAsInt           `json:"ssl_forced"`
	HSTSEnabled           boolAsInt           `json:"hsts_enabled"`
	HSTSSubdomains        boolAsInt           `json:"hsts_subdomains"`
	HTTP2Support          boolAsInt           `json:"http2_support"`
	BlockExploits         boolAsInt           `json:"block_exploits"`
	CachingEnabled        boolAsInt           `json:"caching_enabled"`
	AllowWebsocketUpgrade boolAsInt           `json:"allow_websocket_upgrade"`
	AccessListID          int64               `json:"access_list_id"`
	AdvancedConfig        string              `json:"advanced_config"`
	Enabled               boolAsInt           `json:"enabled"`
	Meta                  Meta                `json:"meta"`
	Locations             []ProxyHostLocation `json:"locations"`
}
