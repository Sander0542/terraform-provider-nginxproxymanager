package client

import (
	"encoding/json"
	"fmt"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/models"
	"net/http"
)

func (c *Client) GetCertificates() (*models.CertificatesResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/nginx/certificates", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := models.CertificatesResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

func (c *Client) GetCertificate(id *int64) (*models.CertificateResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/nginx/certificates/%d", c.HostURL, *id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := models.CertificateResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	if ar.ID == 0 {
		return nil, nil
	}

	return &ar, nil
}
