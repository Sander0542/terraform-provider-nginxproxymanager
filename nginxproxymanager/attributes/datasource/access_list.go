package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var AccessList = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the access list.",
		Required:    true,
	},
	"created_on": schema.StringAttribute{
		Description: "The date and time the access list was created.",
		Computed:    true,
	},
	"modified_on": schema.StringAttribute{
		Description: "The date and time the access list was last modified.",
		Computed:    true,
	},
	"name": schema.StringAttribute{
		Description: "The name of the access list.",
		Computed:    true,
	},
	"owner_user_id": schema.Int64Attribute{
		Description: "The ID of the user that owns the access list.",
		Computed:    true,
	},
	"authorization": schema.ListNestedAttribute{
		Description: "The authorization items of the access list.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: AccessListAuthorization,
		},
	},
	"access": schema.ListNestedAttribute{
		Description: "The access items of the access list.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: AccessListAccess,
		},
	},
	"pass_auth": schema.BoolAttribute{
		Description: "Whether or not to pass the authorization header to the upstream server.",
		Computed:    true,
	},
	"satisfy_any": schema.BoolAttribute{
		Description: "Whether or not to satisfy any of the authorization items.",
		Computed:    true,
	},
	"meta": schema.MapAttribute{
		Description: "The meta data associated with the access list.",
		ElementType: types.StringType,
		Computed:    true,
	},
}
