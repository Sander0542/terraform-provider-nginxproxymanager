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

type CertificateLetsencrypt struct {
	Id          types.Int64  `tfsdk:"id"`
	CreatedOn   types.String `tfsdk:"created_on"`
	ModifiedOn  types.String `tfsdk:"modified_on"`
	ExpiresOn   types.String `tfsdk:"expires_on"`
	OwnerUserId types.Int64  `tfsdk:"owner_user_id"`

	DomainNames            types.Set    `tfsdk:"domain_names"`
	LetsencryptEmail       types.String `tfsdk:"letsencrypt_email"`
	LetsencryptAgree       types.Bool   `tfsdk:"letsencrypt_agree"`
	DnsChallenge           types.Bool   `tfsdk:"dns_challenge"`
	DnsProvider            types.String `tfsdk:"dns_provider"`
	DnsProviderCredentials types.String `tfsdk:"dns_provider_credentials"`
	PropagationSeconds     types.Int64  `tfsdk:"propagation_seconds"`
}

func (_ CertificateLetsencrypt) GetType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"id":                       types.Int64Type,
		"created_on":               types.StringType,
		"modified_on":              types.StringType,
		"expires_on":               types.StringType,
		"owner_user_id":            types.Int64Type,
		"domain_names":             types.SetType{ElemType: types.StringType},
		"letsencrypt_email":        types.StringType,
		"letsencrypt_agree":        types.BoolType,
		"dns_challenge":            types.BoolType,
		"dns_provider":             types.StringType,
		"dns_provider_credentials": types.StringType,
		"propagation_seconds":      types.Int64Type,
	})
}

func (m *CertificateLetsencrypt) Write(ctx context.Context, certificate *nginxproxymanager.GetCertificates200ResponseInner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	m.Id = types.Int64Value(certificate.GetId())
	m.CreatedOn = types.StringValue(certificate.GetCreatedOn())
	m.ModifiedOn = types.StringValue(certificate.GetModifiedOn())
	m.ExpiresOn = types.StringValue(certificate.GetExpiresOn())
	m.OwnerUserId = types.Int64Value(certificate.GetOwnerUserId())

	meta := certificate.GetMeta()
	if meta.HasLetsencryptAgree() {
		m.LetsencryptAgree = types.BoolValue(meta.GetLetsencryptAgree())
	} else {
		m.LetsencryptAgree = types.BoolNull()
	}
	if meta.HasLetsencryptEmail() {
		m.LetsencryptEmail = types.StringValue(meta.GetLetsencryptEmail())
	} else {
		m.LetsencryptEmail = types.StringNull()
	}
	if meta.HasDnsChallenge() {
		m.DnsChallenge = types.BoolValue(meta.GetDnsChallenge())
	} else {
		m.DnsChallenge = types.BoolNull()
	}
	if meta.HasDnsProvider() {
		m.DnsProvider = types.StringValue(meta.GetDnsProvider())
	} else {
		m.DnsProvider = types.StringNull()
	}
	if meta.HasDnsProviderCredentials() {
		m.DnsProviderCredentials = types.StringValue(meta.GetDnsProviderCredentials())
	} else {
		m.DnsProviderCredentials = types.StringNull()
	}
	if meta.HasPropagationSeconds() {
		m.PropagationSeconds = types.Int64Value(meta.GetPropagationSeconds())
	} else {
		m.PropagationSeconds = types.Int64Null()
	}

	m.DomainNames, tmpDiags = SetDomainNamesFrom(ctx, certificate.GetDomainNames())
	diags.Append(tmpDiags...)
}

func (m *CertificateLetsencrypt) ToCreateRequest(ctx context.Context, diags *diag.Diagnostics) *nginxproxymanager.CreateCertificateRequest {
	domainNames, tmpDiags := DomainNameElementsAs(ctx, m.DomainNames)
	diags.Append(tmpDiags...)

	meta := *nginxproxymanager.NewGetCertificates200ResponseInnerMeta()
	meta.SetLetsencryptEmail(m.LetsencryptEmail.ValueString())
	meta.SetLetsencryptAgree(m.LetsencryptAgree.ValueBool())
	meta.SetDnsChallenge(m.DnsChallenge.ValueBool())
	if m.DnsChallenge.ValueBool() {
		meta.SetDnsProvider(m.DnsProvider.ValueString())
		meta.SetDnsProviderCredentials(m.DnsProviderCredentials.ValueString())
		if !m.PropagationSeconds.IsNull() {
			meta.SetPropagationSeconds(m.PropagationSeconds.ValueInt64())
		}
	}

	request := nginxproxymanager.NewCreateCertificateRequest("letsencrypt")
	request.SetDomainNames(domainNames)
	request.SetMeta(meta)

	return request
}
