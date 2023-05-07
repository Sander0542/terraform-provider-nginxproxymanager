package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var AccessListAccess = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the access item.",
		Computed:    true,
	},
	"created_on": schema.StringAttribute{
		Description: "The date and time the access item was created.",
		Computed:    true,
	},
	"modified_on": schema.StringAttribute{
		Description: "The date and time the access item was last modified.",
		Computed:    true,
	},
	"address": schema.StringAttribute{
		Description: "The address of the access item.",
		Computed:    true,
	},
	"directive": schema.StringAttribute{
		Description: "The directive of the access item.",
		Computed:    true,
	},
	"meta": schema.MapAttribute{
		Description: "The meta data associated with the authorization item.",
		ElementType: types.StringType,
		Computed:    true,
	},
}
