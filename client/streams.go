package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/resources"
	"net/http"
)

func (c *Client) GetStreams(ctx context.Context) (*resources.StreamCollection, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/nginx/streams", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := resources.StreamCollection{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

func (c *Client) GetStream(ctx context.Context, id *int64) (*resources.Stream, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/nginx/streams/%d", c.HostURL, *id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := resources.Stream{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	if ar.ID == 0 {
		return nil, nil
	}

	return &ar, nil
}
