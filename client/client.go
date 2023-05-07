package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
	UserAgent  string
}

func NewClient(host *string, username *string, password *string, version string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: http.DefaultTransport,
		},
		UserAgent: fmt.Sprintf("terraform-provider-nginxproxymanager/%s", version),
	}

	if host != nil {
		c.HostURL = *host
	}

	if username == nil || password == nil {
		return &c, nil
	}

	ar, err := c.Authenticate(context.Background(), username, password)
	if err != nil {
		return nil, err
	}

	c.Token = ar.Token

	return &c, nil
}

func (c *Client) doRequest(req *http.Request, authToken *string) ([]byte, error) {
	token := c.Token

	if authToken != nil {
		token = *authToken
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
