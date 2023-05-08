package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var UserPermissions = map[string]schema.Attribute{
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
	"access_lists": schema.StringAttribute{
		Description: "The permission value for the access lists.",
		Computed:    true,
	},
	"certificates": schema.StringAttribute{
		Description: "The permission value for the certificates.",
		Computed:    true,
	},
	"dead_hosts": schema.StringAttribute{
		Description: "The permission value for the dead hosts.",
		Computed:    true,
	},
	"proxy_hosts": schema.StringAttribute{
		Description: "The permission value for the proxy hosts.",
		Computed:    true,
	},
	"redirection_hosts": schema.StringAttribute{
		Description: "The permission value for the redirection hosts.",
		Computed:    true,
	},
	"streams": schema.StringAttribute{
		Description: "The permission value for the streams.",
		Computed:    true,
	},
	"visibility": schema.StringAttribute{
		Description: "The level of visibility for the user.",
		Computed:    true,
	},
	"meta": schema.MapAttribute{
		Description: "The meta data associated with the user.",
		ElementType: types.StringType,
		Computed:    true,
	},
}
