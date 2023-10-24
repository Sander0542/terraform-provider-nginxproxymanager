package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var CertificateCustom = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the certificate.",
		Computed:    true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
	},
	"name": schema.StringAttribute{
		Description: "The name of the certificate.",
		Required:    true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	},
	"certificate": schema.StringAttribute{
		Description: "The contents of the certificate.",
		Required:    true,
		Sensitive:   true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	},
	"certificate_key": schema.StringAttribute{
		Description: "The contents of the certificate key.",
		Required:    true,
		Sensitive:   true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
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
	"created_on": schema.StringAttribute{
		Description: "The date and time the certificate was created.",
		Computed:    true,
	},
	"modified_on": schema.StringAttribute{
		Description: "The date and time the certificate was last modified.",
		Computed:    true,
	},
}
