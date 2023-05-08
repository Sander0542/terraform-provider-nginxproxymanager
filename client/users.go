package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sander0542/terraform-provider-nginxproxymanager/client/resources"
)

func (c *Client) GetUsers(ctx context.Context) (*resources.UserCollection, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/users?expand=items,clients", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := resources.UserCollection{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

func (c *Client) GetUser(ctx context.Context, id *int64) (*resources.User, error) {
	return c.getUser(ctx, strconv.FormatInt(*id, 10))
}

func (c *Client) GetMe(ctx context.Context) (*resources.User, error) {
	return c.getUser(ctx, "me")
}

func (c *Client) getUser(ctx context.Context, resource string) (*resources.User, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/users/%d?expand=items,clients", c.HostURL, resource), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := resources.User{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	if ar.ID == 0 {
		return nil, nil
	}

	return &ar, nil
}
