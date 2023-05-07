package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ProxyHost = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the proxy host.",
		Required:    true,
	},
	"created_on": schema.StringAttribute{
		Description: "The date and time the proxy host was created.",
		Computed:    true,
	},
	"modified_on": schema.StringAttribute{
		Description: "The date and time the proxy host was last modified.",
		Computed:    true,
	},
	"owner_user_id": schema.Int64Attribute{
		Description: "The ID of the user that owns the proxy host.",
		Computed:    true,
	},
	"domain_names": schema.ListAttribute{
		Description: "The domain names associated with the proxy host.",
		Computed:    true,
		ElementType: types.StringType,
	},
	"forward_scheme": schema.StringAttribute{
		Description: "The scheme used to forward requests to the proxy host. Can be either `http` or `https`.",
		Computed:    true,
	},
	"forward_host": schema.StringAttribute{
		Description: "The host used to forward requests to the proxy host.",
		Computed:    true,
	},
	"forward_port": schema.Int64Attribute{
		Description: "The port used to forward requests to the proxy host.",
		Computed:    true,
	},
	"certificate_id": schema.StringAttribute{
		Description: "The ID of the certificate used by the proxy host.",
		Computed:    true,
	},
	"ssl_forced": schema.BoolAttribute{
		Description: "Whether SSL is forced for the proxy host.",
		Computed:    true,
	},
	"hsts_enabled": schema.BoolAttribute{
		Description: "Whether HSTS is enabled for the proxy host.",
		Computed:    true,
	},
	"hsts_subdomains": schema.BoolAttribute{
		Description: "Whether HSTS is enabled for subdomains of the proxy host.",
		Computed:    true,
	},
	"http2_support": schema.BoolAttribute{
		Description: "Whether HTTP/2 is supported for the proxy host.",
		Computed:    true,
	},
	"block_exploits": schema.BoolAttribute{
		Description: "Whether exploits are blocked for the proxy host.",
		Computed:    true,
	},
	"caching_enabled": schema.BoolAttribute{
		Description: "Whether caching is enabled for the proxy host.",
		Computed:    true,
	},
	"allow_websocket_upgrade": schema.BoolAttribute{
		Description: "Whether websocket upgrades are allowed for the proxy host.",
		Computed:    true,
	},
	"access_list_id": schema.Int64Attribute{
		Description: "The ID of the access list used by the proxy host.",
		Computed:    true,
	},
	"advanced_config": schema.StringAttribute{
		Description: "The advanced configuration used by the proxy host.",
		Computed:    true,
	},
	"enabled": schema.BoolAttribute{
		Description: "Whether the proxy host is enabled.",
		Computed:    true,
	},
	"meta": schema.MapAttribute{
		Description: "The meta data associated with the proxy host.",
		ElementType: types.StringType,
		Computed:    true,
	},
	"locations": schema.ListNestedAttribute{
		Description: "The locations associated with the proxy host.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: ProxyHostLocation,
		},
	},
}
