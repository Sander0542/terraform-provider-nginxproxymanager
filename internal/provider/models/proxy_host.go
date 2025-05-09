// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
)

type ProxyHost struct {
	Id          types.Int64  `tfsdk:"id"`
	CreatedOn   types.String `tfsdk:"created_on"`
	ModifiedOn  types.String `tfsdk:"modified_on"`
	OwnerUserId types.Int64  `tfsdk:"owner_user_id"`
	Meta        types.Map    `tfsdk:"meta"`

	DomainNames           types.Set    `tfsdk:"domain_names"`
	ForwardScheme         types.String `tfsdk:"forward_scheme"`
	ForwardHost           types.String `tfsdk:"forward_host"`
	ForwardPort           types.Int64  `tfsdk:"forward_port"`
	CertificateId         types.Int64  `tfsdk:"certificate_id"`
	SslForced             types.Bool   `tfsdk:"ssl_forced"`
	HstsEnabled           types.Bool   `tfsdk:"hsts_enabled"`
	HstsSubdomains        types.Bool   `tfsdk:"hsts_subdomains"`
	Http2Support          types.Bool   `tfsdk:"http2_support"`
	BlockExploits         types.Bool   `tfsdk:"block_exploits"`
	CachingEnabled        types.Bool   `tfsdk:"caching_enabled"`
	AllowWebsocketUpgrade types.Bool   `tfsdk:"allow_websocket_upgrade"`
	AccessListId          types.Int64  `tfsdk:"access_list_id"`
	AdvancedConfig        types.String `tfsdk:"advanced_config"`
	Enabled               types.Bool   `tfsdk:"enabled"`
	Locations             types.Set    `tfsdk:"locations"`
}

func (ProxyHost) GetType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"id":                      types.Int64Type,
		"created_on":              types.StringType,
		"modified_on":             types.StringType,
		"owner_user_id":           types.Int64Type,
		"meta":                    types.MapType{ElemType: types.StringType},
		"domain_names":            types.SetType{ElemType: types.StringType},
		"forward_scheme":          types.StringType,
		"forward_host":            types.StringType,
		"forward_port":            types.Int64Type,
		"certificate_id":          types.Int64Type,
		"ssl_forced":              types.BoolType,
		"hsts_enabled":            types.BoolType,
		"hsts_subdomains":         types.BoolType,
		"http2_support":           types.BoolType,
		"block_exploits":          types.BoolType,
		"caching_enabled":         types.BoolType,
		"allow_websocket_upgrade": types.BoolType,
		"access_list_id":          types.Int64Type,
		"advanced_config":         types.StringType,
		"enabled":                 types.BoolType,
		"locations":               types.SetType{ElemType: ProxyHostLocation{}.GetType()},
	})
}

func (m *ProxyHost) Write(ctx context.Context, proxyHost *nginxproxymanager.GetProxyHosts200ResponseInner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	m.Id = types.Int64Value(proxyHost.GetId())
	m.CreatedOn = types.StringValue(proxyHost.GetCreatedOn())
	m.ModifiedOn = types.StringValue(proxyHost.GetModifiedOn())
	m.OwnerUserId = types.Int64Value(proxyHost.GetOwnerUserId())
	m.ForwardScheme = types.StringValue(proxyHost.GetForwardScheme())
	m.ForwardHost = types.StringValue(proxyHost.GetForwardHost())
	m.ForwardPort = types.Int64Value(proxyHost.GetForwardPort())
	if *proxyHost.GetCertificateId().Int64 != 0 {
		m.CertificateId = types.Int64Value(*proxyHost.GetCertificateId().Int64)
	} else {
		m.CertificateId = types.Int64Null()
	}
	m.SslForced = types.BoolValue(proxyHost.GetSslForced())
	m.HstsEnabled = types.BoolValue(proxyHost.GetHstsEnabled())
	m.HstsSubdomains = types.BoolValue(proxyHost.GetHstsSubdomains())
	m.Http2Support = types.BoolValue(proxyHost.GetHttp2Support())
	m.BlockExploits = types.BoolValue(proxyHost.GetBlockExploits())
	m.CachingEnabled = types.BoolValue(proxyHost.GetCachingEnabled())
	m.AllowWebsocketUpgrade = types.BoolValue(proxyHost.GetAllowWebsocketUpgrade())
	if proxyHost.GetAccessListId() != 0 {
		m.AccessListId = types.Int64Value(proxyHost.GetAccessListId())
	} else {
		m.AccessListId = types.Int64Null()
	}
	m.AdvancedConfig = types.StringValue(proxyHost.GetAdvancedConfig())
	m.Enabled = types.BoolValue(proxyHost.GetEnabled())

	m.Meta, tmpDiags = MapMetaFrom(ctx, proxyHost.GetMeta())
	diags.Append(tmpDiags...)

	m.DomainNames, tmpDiags = SetDomainNamesFrom(ctx, proxyHost.GetDomainNames())
	diags.Append(tmpDiags...)

	m.Locations, tmpDiags = SetProxyHostLocationsFrom(ctx, proxyHost.GetLocations())
	diags.Append(tmpDiags...)
}

func (m *ProxyHost) ToCreateRequest(ctx context.Context, diags *diag.Diagnostics) *nginxproxymanager.CreateProxyHostRequest {
	domainNames, tmpDiags := DomainNameElementsAs(ctx, m.DomainNames)
	diags.Append(tmpDiags...)

	locations, tmpDiags := ProxyHostLocationElementsAs(ctx, m.Locations)
	diags.Append(tmpDiags...)

	request := nginxproxymanager.NewCreateProxyHostRequest(
		domainNames,
		m.ForwardScheme.ValueString(),
		m.ForwardHost.ValueString(),
		m.ForwardPort.ValueInt64(),
	)

	certificateId := m.CertificateId.ValueInt64()
	request.SetCertificateId(nginxproxymanager.GetProxyHosts200ResponseInnerCertificateId{
		Int64: &certificateId,
	})
	request.SetAccessListId(m.AccessListId.ValueInt64())
	request.SetSslForced(m.SslForced.ValueBool())
	request.SetHstsEnabled(m.HstsEnabled.ValueBool())
	request.SetHstsSubdomains(m.HstsSubdomains.ValueBool())
	request.SetHttp2Support(m.Http2Support.ValueBool())
	request.SetBlockExploits(m.BlockExploits.ValueBool())
	request.SetCachingEnabled(m.CachingEnabled.ValueBool())
	request.SetAllowWebsocketUpgrade(m.AllowWebsocketUpgrade.ValueBool())
	request.SetAdvancedConfig(m.AdvancedConfig.ValueString())
	request.SetMeta(map[string]interface{}{})

	requestLocations := make([]nginxproxymanager.GetProxyHosts200ResponseInnerLocationsInner, 0, len(locations))
	for _, location := range locations {
		item := location.Read(ctx, diags)
		requestLocations = append(requestLocations, *item)
	}
	request.SetLocations(requestLocations)

	return request
}

func (m *ProxyHost) ToUpdateRequest(ctx context.Context, diags *diag.Diagnostics) *nginxproxymanager.UpdateProxyHostRequest {
	domainNames, tmpDiags := DomainNameElementsAs(ctx, m.DomainNames)
	diags.Append(tmpDiags...)

	locations, tmpDiags := ProxyHostLocationElementsAs(ctx, m.Locations)
	diags.Append(tmpDiags...)

	request := nginxproxymanager.NewUpdateProxyHostRequest()

	request.SetDomainNames(domainNames)
	request.SetForwardScheme(m.ForwardScheme.ValueString())
	request.SetForwardHost(m.ForwardHost.ValueString())
	request.SetForwardPort(m.ForwardPort.ValueInt64())
	certificateId := m.CertificateId.ValueInt64()
	request.SetCertificateId(nginxproxymanager.GetProxyHosts200ResponseInnerCertificateId{
		Int64: &certificateId,
	})
	request.SetAccessListId(m.AccessListId.ValueInt64())
	request.SetSslForced(m.SslForced.ValueBool())
	request.SetHstsEnabled(m.HstsEnabled.ValueBool())
	request.SetHstsSubdomains(m.HstsSubdomains.ValueBool())
	request.SetHttp2Support(m.Http2Support.ValueBool())
	request.SetBlockExploits(m.BlockExploits.ValueBool())
	request.SetCachingEnabled(m.CachingEnabled.ValueBool())
	request.SetAllowWebsocketUpgrade(m.AllowWebsocketUpgrade.ValueBool())
	request.SetAdvancedConfig(m.AdvancedConfig.ValueString())

	requestLocations := make([]nginxproxymanager.GetProxyHosts200ResponseInnerLocationsInner, 0, len(locations))
	for _, location := range locations {
		item := location.Read(ctx, diags)
		requestLocations = append(requestLocations, *item)
	}
	request.SetLocations(requestLocations)

	return request
}
