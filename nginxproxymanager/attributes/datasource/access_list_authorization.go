package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var AccessListAuthorization = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the authorization item.",
		Computed:    true,
	},
	"created_on": schema.StringAttribute{
		Description: "The date and time the authorization item was created.",
		Computed:    true,
	},
	"modified_on": schema.StringAttribute{
		Description: "The date and time the authorization item was last modified.",
		Computed:    true,
	},
	"username": schema.StringAttribute{
		Description: "The username of the authorization item.",
		Computed:    true,
	},
	"password_hint": schema.StringAttribute{
		Description: "The password hint of the authorization item.",
		Computed:    true,
		Sensitive:   true,
	},
	"meta": schema.MapAttribute{
		Description: "The meta data associated with the authorization item.",
		ElementType: types.StringType,
		Computed:    true,
	},
}
