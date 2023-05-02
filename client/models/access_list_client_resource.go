package models

type AccessListClientResource struct {
	resource
	AccessListId int64  `json:"access_list_id"`
	Address      string `json:"address"`
	Directive    string `json:"directive"`
}

type AccessListClientResourceCollection []AccessListClientResource
