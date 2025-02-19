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

type Certificate struct {
	Id          types.Int64  `tfsdk:"id"`
	CreatedOn   types.String `tfsdk:"created_on"`
	ModifiedOn  types.String `tfsdk:"modified_on"`
	OwnerUserId types.Int64  `tfsdk:"owner_user_id"`
	Meta        types.Map    `tfsdk:"meta"`

	Provider    types.String `tfsdk:"provider_name"`
	NiceName    types.String `tfsdk:"nice_name"`
	DomainNames types.List   `tfsdk:"domain_names"`
	ExpiresOn   types.String `tfsdk:"expires_on"`
}

func (_ Certificate) GetType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"id":            types.Int64Type,
		"created_on":    types.StringType,
		"modified_on":   types.StringType,
		"owner_user_id": types.Int64Type,
		"meta":          types.MapType{ElemType: types.StringType},
		"provider_name": types.StringType,
		"nice_name":     types.StringType,
		"domain_names":  types.ListType{ElemType: types.StringType},
		"expires_on":    types.StringType,
	})
}

func (m *Certificate) Write(ctx context.Context, certificate *nginxproxymanager.GetCertificates200ResponseInner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	m.Id = types.Int64Value(certificate.GetId())
	m.CreatedOn = types.StringValue(certificate.GetCreatedOn())
	m.ModifiedOn = types.StringValue(certificate.GetModifiedOn())
	m.OwnerUserId = types.Int64Value(certificate.GetOwnerUserId())

	m.Provider = types.StringValue(certificate.GetProvider())
	m.NiceName = types.StringValue(certificate.GetNiceName())
	m.ExpiresOn = types.StringValue(certificate.GetExpiresOn())

	m.DomainNames, tmpDiags = ListDomainNamesFrom(ctx, certificate.GetDomainNames())
	diags.Append(tmpDiags...)

	meta, err := certificate.GetMeta().ToMap()
	if err == nil {
		m.Meta, tmpDiags = MapMetaFrom(ctx, meta)
		diags.Append(tmpDiags...)
	} else {
		m.Meta = types.MapNull(types.StringType)
	}
}
