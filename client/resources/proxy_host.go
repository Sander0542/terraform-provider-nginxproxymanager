package resources

import (
	"encoding/json"
)

type ProxyHost struct {
	resourceWithOwner
	DomainNames           []string                    `json:"domain_names"`
	ForwardScheme         string                      `json:"forward_scheme"`
	ForwardHost           string                      `json:"forward_host"`
	ForwardPort           uint16                      `json:"forward_port"`
	CertificateID         json.RawMessage             `json:"certificate_id"`
	SSLForced             bool                        `json:"ssl_forced"`
	HSTSEnabled           bool                        `json:"hsts_enabled"`
	HSTSSubdomains        bool                        `json:"hsts_subdomains"`
	HTTP2Support          bool                        `json:"http2_support"`
	BlockExploits         bool                        `json:"block_exploits"`
	CachingEnabled        bool                        `json:"caching_enabled"`
	AllowWebsocketUpgrade bool                        `json:"allow_websocket_upgrade"`
	AccessListID          int64                       `json:"access_list_id"`
	AdvancedConfig        string                      `json:"advanced_config"`
	Enabled               bool                        `json:"enabled"`
	Locations             ProxyHostLocationCollection `json:"locations"`
}

type ProxyHostCollection []ProxyHost
