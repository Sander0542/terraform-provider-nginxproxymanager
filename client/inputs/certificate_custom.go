package inputs

type CertificateCustom struct {
	Name           string `json:"name"`
	Certificate    string `json:"certificate"`
	CertificateKey string `json:"certificate_key"`
}
