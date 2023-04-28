package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ProxyHostResponse struct {
	ID                    int64           `json:"id"`
	CreatedAt             string          `json:"created_at"`
	ModifiedOn            string          `json:"modified_on"`
	DomainNames           []string        `json:"domain_names"`
	ForwardScheme         string          `json:"forward_scheme"`
	ForwardHost           string          `json:"forward_host"`
	ForwardPort           uint16          `json:"forward_port"`
	CertificateID         json.RawMessage `json:"certificate_id"`
	SSLForced             boolAsInt       `json:"ssl_forced"`
	HSTSEnabled           boolAsInt       `json:"hsts_enabled"`
	HSTSSubdomains        boolAsInt       `json:"hsts_subdomains"`
	HTTP2Support          boolAsInt       `json:"http2_support"`
	BlockExploits         boolAsInt       `json:"block_exploits"`
	CachingEnabled        boolAsInt       `json:"caching_enabled"`
	AllowWebsocketUpgrade boolAsInt       `json:"allow_websocket_upgrade"`
	AccessListID          int64           `json:"access_list_id"`
	AdvancedConfig        string          `json:"advanced_config"`
	Enabled               boolAsInt       `json:"enabled"`
	Meta                  Meta            `json:"meta"`
}

type ProxyHostsResponse []ProxyHostResponse

func (c *Client) GetProxyHosts() (*ProxyHostsResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/nginx/proxy-hosts", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := ProxyHostsResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

func (c *Client) GetProxyHost(id *int64) (*ProxyHostResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/nginx/proxy-hosts/%d", c.HostURL, *id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := ProxyHostResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	if ar.ID == 0 {
		return nil, fmt.Errorf("proxy host with id %d not found", *id)
	}

	return &ar, nil
}
