package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/inputs"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/resources"
)

type ProxyHost struct {
	ID          types.Int64  `tfsdk:"id"`
	CreatedOn   types.String `tfsdk:"created_on"`
	ModifiedOn  types.String `tfsdk:"modified_on"`
	OwnerUserId types.Int64  `tfsdk:"owner_user_id"`
	Meta        types.Map    `tfsdk:"meta"`

	DomainNames           []types.String       `tfsdk:"domain_names"`
	ForwardScheme         types.String         `tfsdk:"forward_scheme"`
	ForwardHost           types.String         `tfsdk:"forward_host"`
	ForwardPort           types.Int64          `tfsdk:"forward_port"`
	CertificateID         types.String         `tfsdk:"certificate_id"`
	SSLForced             types.Bool           `tfsdk:"ssl_forced"`
	HSTSEnabled           types.Bool           `tfsdk:"hsts_enabled"`
	HSTSSubdomains        types.Bool           `tfsdk:"hsts_subdomains"`
	HTTP2Support          types.Bool           `tfsdk:"http2_support"`
	BlockExploits         types.Bool           `tfsdk:"block_exploits"`
	CachingEnabled        types.Bool           `tfsdk:"caching_enabled"`
	AllowWebsocketUpgrade types.Bool           `tfsdk:"allow_websocket_upgrade"`
	AccessListID          types.Int64          `tfsdk:"access_list_id"`
	AdvancedConfig        types.String         `tfsdk:"advanced_config"`
	Enabled               types.Bool           `tfsdk:"enabled"`
	Locations             []*ProxyHostLocation `tfsdk:"location"`
}

func (m *ProxyHost) Load(ctx context.Context, resource *resources.ProxyHost) diag.Diagnostics {
	meta, diags := types.MapValueFrom(ctx, types.StringType, resource.Meta.Map())

	m.ID = types.Int64Value(resource.ID)
	m.CreatedOn = types.StringValue(resource.CreatedOn)
	m.ModifiedOn = types.StringValue(resource.ModifiedOn)
	m.OwnerUserId = types.Int64Value(resource.OwnerUserID)
	m.Meta = meta

	m.ForwardScheme = types.StringValue(resource.ForwardScheme)
	m.ForwardHost = types.StringValue(resource.ForwardHost)
	m.ForwardPort = types.Int64Value(int64(resource.ForwardPort))
	m.CertificateID = types.StringValue(string(resource.CertificateID))
	m.SSLForced = types.BoolValue(resource.SSLForced.Bool())
	m.HSTSEnabled = types.BoolValue(resource.HSTSEnabled.Bool())
	m.HSTSSubdomains = types.BoolValue(resource.HSTSSubdomains.Bool())
	m.HTTP2Support = types.BoolValue(resource.HTTP2Support.Bool())
	m.BlockExploits = types.BoolValue(resource.BlockExploits.Bool())
	m.CachingEnabled = types.BoolValue(resource.CachingEnabled.Bool())
	m.AllowWebsocketUpgrade = types.BoolValue(resource.AllowWebsocketUpgrade.Bool())
	m.AccessListID = types.Int64Value(resource.AccessListID)
	m.AdvancedConfig = types.StringValue(resource.AdvancedConfig)
	m.Enabled = types.BoolValue(resource.Enabled.Bool())

	m.DomainNames = make([]types.String, len(resource.DomainNames))
	for i, v := range resource.DomainNames {
		m.DomainNames[i] = types.StringValue(v)
	}
	var locations []*ProxyHostLocation
	for _, locationResponse := range resource.Locations {
		location := &ProxyHostLocation{}
		location.Load(ctx, &locationResponse)

		locations = append(locations, location)
	}
	m.Locations = locations

	return diags
}

func (m *ProxyHost) Save(ctx context.Context, input *inputs.ProxyHost) diag.Diagnostics {
	diags := diag.Diagnostics{}

	input.ForwardScheme = m.ForwardScheme.ValueString()
	input.ForwardHost = m.ForwardHost.ValueString()
	input.ForwardPort = uint16(m.ForwardPort.ValueInt64())
	input.CertificateID = m.CertificateID.ValueString()
	input.SSLForced = m.SSLForced.ValueBool()
	input.HSTSEnabled = m.HSTSEnabled.ValueBool()
	input.HSTSSubdomains = m.HSTSSubdomains.ValueBool()
	input.HTTP2Support = m.HTTP2Support.ValueBool()
	input.BlockExploits = m.BlockExploits.ValueBool()
	input.CachingEnabled = m.CachingEnabled.ValueBool()
	input.AllowWebsocketUpgrade = m.AllowWebsocketUpgrade.ValueBool()
	input.AccessListID = m.AccessListID.ValueInt64()
	input.AdvancedConfig = m.AdvancedConfig.ValueString()
	input.Meta = map[string]string{}

	input.DomainNames = make([]string, len(m.DomainNames))
	for i, v := range m.DomainNames {
		input.DomainNames[i] = v.ValueString()
	}
	input.Locations = make([]inputs.ProxyHostLocation, len(m.Locations))
	for i, v := range m.Locations {
		diags.Append(v.Save(ctx, &input.Locations[i])...)
	}

	return diags
}
