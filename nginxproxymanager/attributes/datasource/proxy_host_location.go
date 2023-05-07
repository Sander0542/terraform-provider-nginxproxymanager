package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var ProxyHostLocation = map[string]schema.Attribute{
	"path": schema.StringAttribute{
		Description: "The path associated with the location.",
		Computed:    true,
	},
	"forward_scheme": schema.StringAttribute{
		Description: "The scheme used to forward requests to the location. Can be either `http` or `https`.",
		Computed:    true,
	},
	"forward_host": schema.StringAttribute{
		Description: "The host used to forward requests to the location.",
		Computed:    true,
	},
	"forward_port": schema.Int64Attribute{
		Description: "The port used to forward requests to the location.",
		Computed:    true,
	},
	"advanced_config": schema.StringAttribute{
		Description: "The advanced configuration used by the location.",
		Computed:    true,
	},
}
