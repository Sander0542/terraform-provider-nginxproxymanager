package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sander0542/terraform-provider-nginxproxymanager/client/models"
)

func (c *Client) GetAccessLists(ctx context.Context) (*models.AccessListResourceCollection, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/nginx/access-lists?expand=items,clients", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := models.AccessListResourceCollection{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

func (c *Client) GetAccessList(ctx context.Context, id *int64) (*models.AccessListResource, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/nginx/access-lists/%d?expand=items,clients", c.HostURL, *id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := models.AccessListResource{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	if ar.ID == 0 {
		return nil, nil
	}

	return &ar, nil
}
