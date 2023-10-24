package resources

type CertificateValidated struct {
	Certificate    certificate `json:"certificate"`
	CertificateKey bool        `json:"certificate_key"`
}

type certificate struct {
	CN     string           `json:"cn"`
	Issuer string           `json:"issuer"`
	Dates  certificateDates `json:"dates"`
}

type certificateDates struct {
	From int `json:"from"`
	To   int `json:"to"`
}
