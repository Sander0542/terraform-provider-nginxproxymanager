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

type DeadHost struct {
	Id          types.Int64  `tfsdk:"id"`
	CreatedOn   types.String `tfsdk:"created_on"`
	ModifiedOn  types.String `tfsdk:"modified_on"`
	OwnerUserId types.Int64  `tfsdk:"owner_user_id"`
	Meta        types.Map    `tfsdk:"meta"`

	DomainNames    types.List   `tfsdk:"domain_names"`
	CertificateId  types.Int64  `tfsdk:"certificate_id"`
	SslForced      types.Bool   `tfsdk:"ssl_forced"`
	HstsEnabled    types.Bool   `tfsdk:"hsts_enabled"`
	HstsSubdomains types.Bool   `tfsdk:"hsts_subdomains"`
	Http2Support   types.Bool   `tfsdk:"http2_support"`
	AdvancedConfig types.String `tfsdk:"advanced_config"`
	Enabled        types.Bool   `tfsdk:"enabled"`
}

func (_ DeadHost) GetType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"id":              types.Int64Type,
		"created_on":      types.StringType,
		"modified_on":     types.StringType,
		"owner_user_id":   types.Int64Type,
		"meta":            types.MapType{ElemType: types.StringType},
		"domain_names":    types.ListType{ElemType: types.StringType},
		"certificate_id":  types.Int64Type,
		"ssl_forced":      types.BoolType,
		"hsts_enabled":    types.BoolType,
		"hsts_subdomains": types.BoolType,
		"http2_support":   types.BoolType,
		"advanced_config": types.StringType,
		"enabled":         types.BoolType,
	})
}

func (m *DeadHost) Write(ctx context.Context, deadHost *nginxproxymanager.GetDeadHosts200ResponseInner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	m.Id = types.Int64Value(deadHost.GetId())
	m.CreatedOn = types.StringValue(deadHost.GetCreatedOn())
	m.ModifiedOn = types.StringValue(deadHost.GetModifiedOn())
	m.OwnerUserId = types.Int64Value(deadHost.GetOwnerUserId())

	if *deadHost.GetCertificateId().Int64 != 0 {
		m.CertificateId = types.Int64Value(*deadHost.GetCertificateId().Int64)
	} else {
		m.CertificateId = types.Int64Null()
	}
	m.SslForced = types.BoolValue(deadHost.GetSslForced())
	m.HstsEnabled = types.BoolValue(deadHost.GetHstsEnabled())
	m.HstsSubdomains = types.BoolValue(deadHost.GetHstsSubdomains())
	m.Http2Support = types.BoolValue(deadHost.GetHttp2Support())
	m.AdvancedConfig = types.StringValue(deadHost.GetAdvancedConfig())
	m.Enabled = types.BoolValue(deadHost.GetEnabled())

	m.Meta, tmpDiags = MapMetaFrom(ctx, deadHost.GetMeta())
	diags.Append(tmpDiags...)

	m.DomainNames, tmpDiags = ListDomainNamesFrom(ctx, deadHost.GetDomainNames())
	diags.Append(tmpDiags...)
}

func (m *DeadHost) ToCreateRequest(ctx context.Context, diags *diag.Diagnostics) *nginxproxymanager.Create404HostRequest {
	domainNames, tmpDiags := DomainNameElementsAs(ctx, m.DomainNames)
	diags.Append(tmpDiags...)

	request := nginxproxymanager.NewCreate404HostRequest(domainNames)
	certificateId := m.CertificateId.ValueInt64()
	request.SetCertificateId(nginxproxymanager.GetProxyHosts200ResponseInnerCertificateId{
		Int64: &certificateId,
	})
	request.SetSslForced(m.SslForced.ValueBool())
	request.SetHstsEnabled(m.HstsEnabled.ValueBool())
	request.SetHstsSubdomains(m.HstsSubdomains.ValueBool())
	request.SetHttp2Support(m.Http2Support.ValueBool())
	request.SetAdvancedConfig(m.AdvancedConfig.ValueString())
	request.SetMeta(map[string]interface{}{})

	return request
}

func (m *DeadHost) ToUpdateRequest(ctx context.Context, diags *diag.Diagnostics) *nginxproxymanager.UpdateDeadHostRequest {
	domainNames, tmpDiags := DomainNameElementsAs(ctx, m.DomainNames)
	diags.Append(tmpDiags...)

	request := nginxproxymanager.NewUpdateDeadHostRequest()
	request.SetDomainNames(domainNames)
	certificateId := m.CertificateId.ValueInt64()
	request.SetCertificateId(nginxproxymanager.GetProxyHosts200ResponseInnerCertificateId{
		Int64: &certificateId,
	})
	request.SetSslForced(m.SslForced.ValueBool())
	request.SetHstsEnabled(m.HstsEnabled.ValueBool())
	request.SetHstsSubdomains(m.HstsSubdomains.ValueBool())
	request.SetHttp2Support(m.Http2Support.ValueBool())
	request.SetAdvancedConfig(m.AdvancedConfig.ValueString())
	request.SetMeta(map[string]interface{}{})

	return request
}
