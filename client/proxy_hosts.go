package client

import (
	"encoding/json"
	"fmt"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/models"
	"net/http"
	"strings"
)

func (c *Client) CreateProxyHost(proxyHost *models.ProxyHostCreate) (*models.ProxyHostResponse, error) {
	rb, err := json.Marshal(proxyHost)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/nginx/proxy-hosts", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := models.ProxyHostResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

func (c *Client) GetProxyHosts() (*models.ProxyHostsResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/nginx/proxy-hosts", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := models.ProxyHostsResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

func (c *Client) GetProxyHost(id *int64) (*models.ProxyHostResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/nginx/proxy-hosts/%d", c.HostURL, *id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := models.ProxyHostResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	if ar.ID == 0 {
		return nil, nil
	}

	return &ar, nil
}

func (c *Client) UpdateProxyHost(id *int64, proxyHost *models.ProxyHostUpdate) (*models.ProxyHostResponse, error) {
	rb, err := json.Marshal(proxyHost)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/nginx/proxy-hosts/%d", c.HostURL, *id), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := models.ProxyHostResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

func (c *Client) DeleteProxyHost(id *int64) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/nginx/proxy-hosts/%d", c.HostURL, *id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req, nil)
	if err != nil {
		return err
	}

	return nil
}
