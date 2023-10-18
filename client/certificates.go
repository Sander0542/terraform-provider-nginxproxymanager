package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/inputs"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/resources"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

func (c *Client) CreateCertificateCustom(ctx context.Context, certificateCustom *inputs.CertificateCustom) (*resources.Certificate, error) {
	_, err := c.ValidateCertificate(ctx, certificateCustom)
	if err != nil {
		return nil, err
	}

	certificate := &inputs.Certificate{
		NiceName: certificateCustom.Name,
		Provider: "other",
	}

	rb, err := json.Marshal(certificate)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/nginx/certificates", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := resources.Certificate{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	certificateUpload := &inputs.CertificateUpload{
		CertificateId:  ar.ID,
		Certificate:    certificateCustom.Certificate,
		CertificateKey: certificateCustom.CertificateKey,
	}

	uploaded, err := c.UploadCertificate(ctx, certificateUpload)
	if err != nil {
		return nil, err
	}

	ar.Meta = *uploaded

	return &ar, nil
}

func (c *Client) GetCertificates(ctx context.Context) (*resources.CertificateCollection, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/nginx/certificates", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := resources.CertificateCollection{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

func (c *Client) GetCertificate(ctx context.Context, id *int64) (*resources.Certificate, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/nginx/certificates/%d", c.HostURL, *id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := resources.Certificate{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	if ar.ID == 0 {
		return nil, nil
	}

	return &ar, nil
}

func (c *Client) ValidateCertificate(ctx context.Context, certificate *inputs.CertificateCustom) (*resources.CertificateValidated, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	certPart, err := writer.CreateFormFile("certificate", "certificate.crt")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(certPart, strings.NewReader(certificate.Certificate))
	if err != nil {
		return nil, err
	}

	keyPart, err := writer.CreateFormFile("certificate_key", "certificate.key")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(keyPart, strings.NewReader(certificate.CertificateKey))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/nginx/certificates/validate", c.HostURL), payload)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := resources.CertificateValidated{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

func (c *Client) UploadCertificate(ctx context.Context, certificate *inputs.CertificateUpload) (*resources.Meta, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	certPart, err := writer.CreateFormFile("certificate", "certificate.crt")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(certPart, strings.NewReader(certificate.Certificate))
	if err != nil {
		return nil, err
	}

	keyPart, err := writer.CreateFormFile("certificate_key", "certificate.key")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(keyPart, strings.NewReader(certificate.CertificateKey))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/nginx/certificates/%d/upload", c.HostURL, certificate.CertificateId), payload)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := resources.Meta{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

func (c *Client) DeleteCertificate(ctx context.Context, id *int64) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/api/nginx/certificates/%d", c.HostURL, *id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req, nil)
	if err != nil {
		return err
	}

	return nil
}
