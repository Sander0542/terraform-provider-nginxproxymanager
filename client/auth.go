package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
	"net/http"
	"strings"
)

type AuthStruct struct {
	Username string `json:"identity"`
	Password string `json:"secret"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (c *Client) Authenticate(ctx context.Context, username, password *string) (*AuthResponse, error) {
	if *username == "" || *password == "" {
		return nil, fmt.Errorf("username and password must be set")
	}

	auth := AuthStruct{
		Username: *username,
		Password: *password,
	}

	rb, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/tokens", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.Header.Set(headers.ContentType, "application/json")

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := AuthResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}
