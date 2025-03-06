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

type RedirectionHost struct {
	Id          types.Int64  `tfsdk:"id"`
	CreatedOn   types.String `tfsdk:"created_on"`
	ModifiedOn  types.String `tfsdk:"modified_on"`
	OwnerUserId types.Int64  `tfsdk:"owner_user_id"`
	Meta        types.Map    `tfsdk:"meta"`

	DomainNames       types.Set    `tfsdk:"domain_names"`
	ForwardScheme     types.String `tfsdk:"forward_scheme"`
	ForwardDomainName types.String `tfsdk:"forward_domain_name"`
	ForwardHttpCode   types.Int64  `tfsdk:"forward_http_code"`
	CertificateId     types.Int64  `tfsdk:"certificate_id"`
	SslForced         types.Bool   `tfsdk:"ssl_forced"`
	HstsEnabled       types.Bool   `tfsdk:"hsts_enabled"`
	HstsSubdomains    types.Bool   `tfsdk:"hsts_subdomains"`
	Http2Support      types.Bool   `tfsdk:"http2_support"`
	PreservePath      types.Bool   `tfsdk:"preserve_path"`
	BlockExploits     types.Bool   `tfsdk:"block_exploits"`
	AdvancedConfig    types.String `tfsdk:"advanced_config"`
	Enabled           types.Bool   `tfsdk:"enabled"`
}

func (_ RedirectionHost) GetType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"id":                  types.Int64Type,
		"created_on":          types.StringType,
		"modified_on":         types.StringType,
		"owner_user_id":       types.Int64Type,
		"meta":                types.MapType{ElemType: types.StringType},
		"domain_names":        types.SetType{ElemType: types.StringType},
		"forward_scheme":      types.StringType,
		"forward_domain_name": types.StringType,
		"forward_http_code":   types.Int64Type,
		"certificate_id":      types.Int64Type,
		"ssl_forced":          types.BoolType,
		"hsts_enabled":        types.BoolType,
		"hsts_subdomains":     types.BoolType,
		"http2_support":       types.BoolType,
		"preserve_path":       types.BoolType,
		"block_exploits":      types.BoolType,
		"advanced_config":     types.StringType,
		"enabled":             types.BoolType,
	})
}

func (m *RedirectionHost) Write(ctx context.Context, redirectionHost *nginxproxymanager.GetRedirectionHosts200ResponseInner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	m.Id = types.Int64Value(redirectionHost.GetId())
	m.CreatedOn = types.StringValue(redirectionHost.GetCreatedOn())
	m.ModifiedOn = types.StringValue(redirectionHost.GetModifiedOn())
	m.OwnerUserId = types.Int64Value(redirectionHost.GetOwnerUserId())

	m.ForwardScheme = types.StringValue(redirectionHost.GetForwardScheme())
	m.ForwardDomainName = types.StringValue(redirectionHost.GetForwardDomainName())
	m.ForwardHttpCode = types.Int64Value(redirectionHost.GetForwardHttpCode())
	if *redirectionHost.GetCertificateId().Int64 != 0 {
		m.CertificateId = types.Int64Value(*redirectionHost.GetCertificateId().Int64)
	} else {
		m.CertificateId = types.Int64Null()
	}
	m.SslForced = types.BoolValue(redirectionHost.GetSslForced())
	m.HstsEnabled = types.BoolValue(redirectionHost.GetHstsEnabled())
	m.HstsSubdomains = types.BoolValue(redirectionHost.GetHstsSubdomains())
	m.Http2Support = types.BoolValue(redirectionHost.GetHttp2Support())
	m.PreservePath = types.BoolValue(redirectionHost.GetPreservePath())
	m.BlockExploits = types.BoolValue(redirectionHost.GetBlockExploits())
	m.AdvancedConfig = types.StringValue(redirectionHost.GetAdvancedConfig())
	m.Enabled = types.BoolValue(redirectionHost.GetEnabled())

	m.Meta, tmpDiags = MapMetaFrom(ctx, redirectionHost.GetMeta())
	diags.Append(tmpDiags...)

	m.DomainNames, tmpDiags = SetDomainNamesFrom(ctx, redirectionHost.GetDomainNames())
	diags.Append(tmpDiags...)
}

func (m *RedirectionHost) ToCreateRequest(ctx context.Context, diags *diag.Diagnostics) *nginxproxymanager.CreateRedirectionHostRequest {
	domainNames, tmpDiags := DomainNameElementsAs(ctx, m.DomainNames)
	diags.Append(tmpDiags...)

	request := nginxproxymanager.NewCreateRedirectionHostRequest(
		domainNames,
		m.ForwardHttpCode.ValueInt64(),
		m.ForwardScheme.ValueString(),
		m.ForwardDomainName.ValueString(),
	)

	certificateId := m.CertificateId.ValueInt64()
	request.SetCertificateId(nginxproxymanager.GetProxyHosts200ResponseInnerCertificateId{
		Int64: &certificateId,
	})
	request.SetSslForced(m.SslForced.ValueBool())
	request.SetHstsEnabled(m.HstsEnabled.ValueBool())
	request.SetHstsSubdomains(m.HstsSubdomains.ValueBool())
	request.SetHttp2Support(m.Http2Support.ValueBool())
	request.SetPreservePath(m.PreservePath.ValueBool())
	request.SetBlockExploits(m.BlockExploits.ValueBool())
	request.SetAdvancedConfig(m.AdvancedConfig.ValueString())
	request.SetMeta(map[string]interface{}{})

	return request
}

func (m *RedirectionHost) ToUpdateRequest(ctx context.Context, diags *diag.Diagnostics) *nginxproxymanager.UpdateRedirectionHostRequest {
	domainNames, tmpDiags := DomainNameElementsAs(ctx, m.DomainNames)
	diags.Append(tmpDiags...)

	request := nginxproxymanager.NewUpdateRedirectionHostRequest()

	request.SetDomainNames(domainNames)
	request.SetForwardHttpCode(m.ForwardHttpCode.ValueInt64())
	request.SetForwardScheme(m.ForwardScheme.ValueString())
	request.SetForwardDomainName(m.ForwardDomainName.ValueString())
	certificateId := m.CertificateId.ValueInt64()
	request.SetCertificateId(nginxproxymanager.GetProxyHosts200ResponseInnerCertificateId{
		Int64: &certificateId,
	})
	request.SetSslForced(m.SslForced.ValueBool())
	request.SetHstsEnabled(m.HstsEnabled.ValueBool())
	request.SetHstsSubdomains(m.HstsSubdomains.ValueBool())
	request.SetHttp2Support(m.Http2Support.ValueBool())
	request.SetPreservePath(m.PreservePath.ValueBool())
	request.SetBlockExploits(m.BlockExploits.ValueBool())
	request.SetAdvancedConfig(m.AdvancedConfig.ValueString())

	return request
}
