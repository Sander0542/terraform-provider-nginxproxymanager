package resources

import "encoding/json"

type RedirectionHost struct {
	resourceWithOwner
	DomainNames       []string        `json:"domain_names"`
	ForwardScheme     string          `json:"forward_scheme"`
	ForwardDomainName string          `json:"forward_domain_name"`
	ForwardHTTPCode   uint16          `json:"forward_http_code"`
	CertificateID     json.RawMessage `json:"certificate_id"`
	SSLForced         boolAsInt       `json:"ssl_forced"`
	HSTSEnabled       boolAsInt       `json:"hsts_enabled"`
	HSTSSubdomains    boolAsInt       `json:"hsts_subdomains"`
	HTTP2Support      boolAsInt       `json:"http2_support"`
	PreservePath      boolAsInt       `json:"preserve_path"`
	BlockExploits     boolAsInt       `json:"block_exploits"`
	AdvancedConfig    string          `json:"advanced_config"`
	Enabled           boolAsInt       `json:"enabled"`
}

type RedirectionHostCollection []RedirectionHost
