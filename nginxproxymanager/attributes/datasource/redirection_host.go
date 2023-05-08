package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var RedirectionHost = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the redirection host.",
		Required:    true,
	},
	"created_on": schema.StringAttribute{
		Description: "The date and time the redirection host was created.",
		Computed:    true,
	},
	"modified_on": schema.StringAttribute{
		Description: "The date and time the redirection host was last modified.",
		Computed:    true,
	},
	"owner_user_id": schema.Int64Attribute{
		Description: "The ID of the user that owns the redirection host.",
		Computed:    true,
	},
	"domain_names": schema.ListAttribute{
		Description: "The domain names associated with the redirection host.",
		Computed:    true,
		ElementType: types.StringType,
	},
	"forward_scheme": schema.StringAttribute{
		Description: "The scheme used to forward requests to the redirection host. Can be either `auto`, `http` or `https`.",
		Computed:    true,
	},
	"forward_domain_name": schema.StringAttribute{
		Description: "The domain name used to forward requests to the redirection host.",
		Computed:    true,
	},
	"forward_http_code": schema.Int64Attribute{
		Description: "The HTTP code used to forward requests to the redirection host.",
		Computed:    true,
	},
	"certificate_id": schema.StringAttribute{
		Description: "The ID of the certificate used by the redirection host.",
		Computed:    true,
	},
	"ssl_forced": schema.BoolAttribute{
		Description: "Whether SSL is forced for the redirection host.",
		Computed:    true,
	},
	"hsts_enabled": schema.BoolAttribute{
		Description: "Whether HSTS is enabled for the redirection host.",
		Computed:    true,
	},
	"hsts_subdomains": schema.BoolAttribute{
		Description: "Whether HSTS is enabled for subdomains of the redirection host.",
		Computed:    true,
	},
	"http2_support": schema.BoolAttribute{
		Description: "Whether HTTP/2 is supported for the redirection host.",
		Computed:    true,
	},
	"preserve_path": schema.BoolAttribute{
		Description: "Whether the path is preserved for the redirection host.",
		Computed:    true,
	},
	"block_exploits": schema.BoolAttribute{
		Description: "Whether exploits are blocked for the redirection host.",
		Computed:    true,
	},
	"advanced_config": schema.StringAttribute{
		Description: "The advanced configuration used by the redirection host.",
		Computed:    true,
	},
	"enabled": schema.BoolAttribute{
		Description: "Whether the redirection host is enabled.",
		Computed:    true,
	},
	"meta": schema.MapAttribute{
		Description: "The meta data associated with the redirection host.",
		ElementType: types.StringType,
		Computed:    true,
	},
}
