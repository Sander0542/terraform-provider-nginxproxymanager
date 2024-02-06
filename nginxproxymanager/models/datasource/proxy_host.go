package datasource

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/resources"
)

type ProxyHost struct {
	ID          types.Int64  `tfsdk:"id"`
	CreatedOn   types.String `tfsdk:"created_on"`
	ModifiedOn  types.String `tfsdk:"modified_on"`
	OwnerUserId types.Int64  `tfsdk:"owner_user_id"`
	Meta        types.Map    `tfsdk:"meta"`

	DomainNames           types.List           `tfsdk:"domain_names"`
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
	Locations             []*ProxyHostLocation `tfsdk:"locations"`
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

	domainNames, domainNamesDiags := types.ListValueFrom(ctx, types.StringType, resource.DomainNames)
	diags.Append(domainNamesDiags...)
	m.DomainNames = domainNames

	var locations []*ProxyHostLocation
	for _, locationResponse := range resource.Locations {
		location := &ProxyHostLocation{}
		location.Load(ctx, &locationResponse)

		locations = append(locations, location)
	}
	m.Locations = locations

	return diags
}
