package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var DeadHost = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the dead host.",
		Required:    true,
	},
	"created_on": schema.StringAttribute{
		Description: "The date and time the dead host was created.",
		Computed:    true,
	},
	"modified_on": schema.StringAttribute{
		Description: "The date and time the dead host was last modified.",
		Computed:    true,
	},
	"owner_user_id": schema.Int64Attribute{
		Description: "The ID of the user that owns the dead host.",
		Computed:    true,
	},
	"domain_names": schema.ListAttribute{
		Description: "The domain names associated with the dead host.",
		Computed:    true,
		ElementType: types.StringType,
	},
	"certificate_id": schema.StringAttribute{
		Description: "The ID of the certificate used by the dead host.",
		Computed:    true,
	},
	"ssl_forced": schema.BoolAttribute{
		Description: "Whether SSL is forced for the dead host.",
		Computed:    true,
	},
	"hsts_enabled": schema.BoolAttribute{
		Description: "Whether HSTS is enabled for the dead host.",
		Computed:    true,
	},
	"hsts_subdomains": schema.BoolAttribute{
		Description: "Whether HSTS is enabled for subdomains of the dead host.",
		Computed:    true,
	},
	"http2_support": schema.BoolAttribute{
		Description: "Whether HTTP/2 is supported for the dead host.",
		Computed:    true,
	},
	"advanced_config": schema.StringAttribute{
		Description: "The advanced configuration used by the dead host.",
		Computed:    true,
	},
	"enabled": schema.BoolAttribute{
		Description: "Whether the dead host is enabled.",
		Computed:    true,
	},
	"meta": schema.MapAttribute{
		Description: "The meta data associated with the dead host.",
		ElementType: types.StringType,
		Computed:    true,
	},
}
