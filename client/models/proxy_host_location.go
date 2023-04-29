package models

type ProxyHostLocation struct {
	Path           string `json:"path"`
	ForwardScheme  string `json:"forward_scheme"`
	ForwardHost    string `json:"forward_host"`
	ForwardPort    uint16 `json:"forward_port"`
	AdvancedConfig string `json:"advanced_config"`
}
