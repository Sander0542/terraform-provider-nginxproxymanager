package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var Certificate = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the certificate.",
		Required:    true,
	},
	"created_on": schema.StringAttribute{
		Description: "The date and time the certificate was created.",
		Computed:    true,
	},
	"modified_on": schema.StringAttribute{
		Description: "The date and time the certificate was last modified.",
		Computed:    true,
	},
	"provider_name": schema.StringAttribute{
		Description: "The provider of the certificate.",
		Computed:    true,
	},
	"nice_name": schema.StringAttribute{
		Description: "The nice name of the certificate.",
		Computed:    true,
	},
	"domain_names": schema.ListAttribute{
		Description: "The domain names associated with the certificate.",
		Computed:    true,
		ElementType: types.StringType,
	},
	"expires_on": schema.StringAttribute{
		Description: "The date and time the certificate expires.",
		Computed:    true,
	},
	"meta": schema.MapAttribute{
		Description: "The meta data associated with the certificate.",
		ElementType: types.StringType,
		Computed:    true,
	},
}
