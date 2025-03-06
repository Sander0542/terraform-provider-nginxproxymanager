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

type CertificateCustom struct {
	Id          types.Int64  `tfsdk:"id"`
	CreatedOn   types.String `tfsdk:"created_on"`
	ModifiedOn  types.String `tfsdk:"modified_on"`
	OwnerUserId types.Int64  `tfsdk:"owner_user_id"`

	Name           types.String `tfsdk:"name"`
	Certificate    types.String `tfsdk:"certificate"`
	CertificateKey types.String `tfsdk:"certificate_key"`
	DomainNames    types.Set    `tfsdk:"domain_names"`
	ExpiresOn      types.String `tfsdk:"expires_on"`
}

func (_ CertificateCustom) GetType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"id":              types.Int64Type,
		"created_on":      types.StringType,
		"modified_on":     types.StringType,
		"owner_user_id":   types.Int64Type,
		"name":            types.StringType,
		"certificate":     types.StringType,
		"certificate_key": types.StringType,
		"domain_names":    types.SetType{ElemType: types.StringType},
		"expires_on":      types.StringType,
	})
}

func (m *CertificateCustom) Write(ctx context.Context, certificate *nginxproxymanager.GetCertificates200ResponseInner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	m.Id = types.Int64Value(certificate.GetId())
	m.CreatedOn = types.StringValue(certificate.GetCreatedOn())
	m.ModifiedOn = types.StringValue(certificate.GetModifiedOn())
	m.OwnerUserId = types.Int64Value(certificate.GetOwnerUserId())

	m.Name = types.StringValue(certificate.GetNiceName())
	meta := certificate.GetMeta()
	if meta.HasCertificate() {
		m.Certificate = types.StringValue(meta.GetCertificate())
	} else {
		m.Certificate = types.StringNull()
	}
	if meta.HasCertificateKey() {
		m.CertificateKey = types.StringValue(meta.GetCertificateKey())
	} else {
		m.CertificateKey = types.StringNull()
	}
	m.ExpiresOn = types.StringValue(certificate.GetExpiresOn())

	m.DomainNames, tmpDiags = SetDomainNamesFrom(ctx, certificate.GetDomainNames())
	diags.Append(tmpDiags...)
}

func (m *CertificateCustom) ToCreateRequest(ctx context.Context, diags *diag.Diagnostics) *nginxproxymanager.CreateCertificateRequest {
	request := nginxproxymanager.NewCreateCertificateRequest("other")

	request.SetNiceName(m.Name.ValueString())

	return request
}
