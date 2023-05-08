package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var UserPermissions = map[string]schema.Attribute{
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
}
