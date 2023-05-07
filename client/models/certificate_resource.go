package models

type CertificateResource struct {
	resource
	Provider    string   `json:"provider"`
	NiceName    string   `json:"nice_name"`
	DomainNames []string `json:"domain_names"`
	ExpiresOn   string   `json:"expires_on"`
}

type CertificateResourceCollection []CertificateResource
