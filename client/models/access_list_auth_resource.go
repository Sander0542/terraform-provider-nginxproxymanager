package models

type AccessListAuthResource struct {
	resource
	AccessListId int64  `json:"access_list_id"`
	Hint         string `json:"hint"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

type AccessListAuthResourceCollection []AccessListAuthResource
