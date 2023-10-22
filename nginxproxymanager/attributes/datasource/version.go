package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var Version = map[string]schema.Attribute{
	"major": schema.Int64Attribute{
		Description: "The major version number.",
		Computed:    true,
	},
	"minor": schema.Int64Attribute{
		Description: "The minor version number.",
		Computed:    true,
	},
	"revision": schema.Int64Attribute{
		Description: "The revision version number.",
		Computed:    true,
	},
	"version": schema.StringAttribute{
		Description: "The full version.",
		Computed:    true,
	},
}
