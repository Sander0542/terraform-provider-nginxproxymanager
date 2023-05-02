package models

type AccessListResource struct {
	resourceWithOwner
	Name           string                             `json:"name"`
	Items          AccessListAuthResourceCollection   `json:"items"`
	Clients        AccessListClientResourceCollection `json:"clients"`
	PassAuth       boolAsInt                          `json:"pass_auth"`
	SatisfyAny     boolAsInt                          `json:"satisfy_any"`
	ProxyHostCount int64                              `json:"proxy_host_count"`
}

type AccessListResourceCollection []AccessListResource
