package resources

import (
	"encoding/json"
)

type DeadHost struct {
	resourceWithOwner
	DomainNames    []string        `json:"domain_names"`
	CertificateID  json.RawMessage `json:"certificate_id"`
	SSLForced      boolAsInt       `json:"ssl_forced"`
	HSTSEnabled    boolAsInt       `json:"hsts_enabled"`
	HSTSSubdomains boolAsInt       `json:"hsts_subdomains"`
	HTTP2Support   boolAsInt       `json:"http2_support"`
	AdvancedConfig string          `json:"advanced_config"`
	Enabled        bool            `json:"enabled"`
}

type DeadHostCollection []DeadHost
