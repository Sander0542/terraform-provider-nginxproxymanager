package resources

type AccessListClient struct {
	resource
	AccessListId int64  `json:"access_list_id"`
	Address      string `json:"address"`
	Directive    string `json:"directive"`
}

type AccessListClientCollection []AccessListClient
