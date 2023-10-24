package inputs

type CertificateUpload struct {
	CertificateId  int64  `json:"certificate_id"`
	Certificate    string `json:"certificate"`
	CertificateKey string `json:"certificate_key"`
}
