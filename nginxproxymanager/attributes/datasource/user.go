package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var User = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the user.",
		Required:    true,
	},
	"created_on": schema.StringAttribute{
		Description: "The date and time the user was created.",
		Computed:    true,
	},
	"modified_on": schema.StringAttribute{
		Description: "The date and time the user was last modified.",
		Computed:    true,
	},
	"name": schema.StringAttribute{
		Description: "The name of the user.",
		Computed:    true,
	},
	"nickname": schema.StringAttribute{
		Description: "The nickname of the user.",
		Computed:    true,
	},
	"email": schema.StringAttribute{
		Description: "The email of the user.",
		Computed:    true,
	},
	"avatar": schema.StringAttribute{
		Description: "The avatar of the user.",
		Computed:    true,
	},
	"is_disabled": schema.BoolAttribute{
		Description: "Whether the user is disabled.",
		Computed:    true,
	},
	"roles": schema.ListAttribute{
		Description: "The roles of the user.",
		Computed:    true,
		ElementType: types.StringType,
	},
	"permissions": schema.SingleNestedAttribute{
		Description: "The permissions of the user.",
		Computed:    true,
		Attributes:  UserPermissions,
	},
	"meta": schema.MapAttribute{
		Description: "The meta data associated with the user.",
		ElementType: types.StringType,
		Computed:    true,
	},
}
