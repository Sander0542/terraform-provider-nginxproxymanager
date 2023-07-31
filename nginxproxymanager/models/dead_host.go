package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/resources"
)

type DeadHost struct {
	ID          types.Int64  `tfsdk:"id"`
	CreatedOn   types.String `tfsdk:"created_on"`
	ModifiedOn  types.String `tfsdk:"modified_on"`
	OwnerUserId types.Int64  `tfsdk:"owner_user_id"`
	Meta        types.Map    `tfsdk:"meta"`

	DomainNames    []types.String `tfsdk:"domain_names"`
	CertificateID  types.String   `tfsdk:"certificate_id"`
	SSLForced      types.Bool     `tfsdk:"ssl_forced"`
	HSTSEnabled    types.Bool     `tfsdk:"hsts_enabled"`
	HSTSSubdomains types.Bool     `tfsdk:"hsts_subdomains"`
	HTTP2Support   types.Bool     `tfsdk:"http2_support"`
	BlockExploits  types.Bool     `tfsdk:"block_exploits"`
	AdvancedConfig types.String   `tfsdk:"advanced_config"`
	Enabled        types.Bool     `tfsdk:"enabled"`
}

func (m *DeadHost) Load(ctx context.Context, resource *resources.DeadHost) diag.Diagnostics {
	meta, diags := types.MapValueFrom(ctx, types.StringType, resource.Meta.Map())

	m.ID = types.Int64Value(resource.ID)
	m.CreatedOn = types.StringValue(resource.CreatedOn)
	m.ModifiedOn = types.StringValue(resource.ModifiedOn)
	m.OwnerUserId = types.Int64Value(resource.OwnerUserID)
	m.Meta = meta

	m.CertificateID = types.StringValue(string(resource.CertificateID))
	m.SSLForced = types.BoolValue(resource.SSLForced.Bool())
	m.HSTSEnabled = types.BoolValue(resource.HSTSEnabled.Bool())
	m.HSTSSubdomains = types.BoolValue(resource.HSTSSubdomains.Bool())
	m.HTTP2Support = types.BoolValue(resource.HTTP2Support.Bool())
	m.AdvancedConfig = types.StringValue(resource.AdvancedConfig)
	m.Enabled = types.BoolValue(resource.Enabled.Bool())

	m.DomainNames = make([]types.String, len(resource.DomainNames))
	for i, v := range resource.DomainNames {
		m.DomainNames[i] = types.StringValue(v)
	}

	return diags
}
