// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/sander0542/nginxproxymanager-go"
)

type SettingDefaultSite struct {
	Page     types.String `tfsdk:"page"`
	Redirect types.String `tfsdk:"redirect"`
	Html     types.String `tfsdk:"html"`
}

func (_ SettingDefaultSite) GetType() attr.TypeWithAttributeTypes {
	return types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"page":     types.StringType,
		"redirect": types.StringType,
		"html":     types.StringType,
	})
}

func (m *SettingDefaultSite) Write(ctx context.Context, setting *nginxproxymanager.GetSettings200ResponseInner, diags *diag.Diagnostics) {
	m.Page = types.StringPointerValue(setting.GetValue().String)

	if redirect, ok := setting.GetMeta()["redirect"]; ok {
		m.Redirect = types.StringValue(fmt.Sprintf("%v", redirect))
	} else {
		m.Redirect = types.StringNull()
	}
	if html, ok := setting.GetMeta()["html"]; ok {
		m.Html = types.StringValue(fmt.Sprintf("%v", html))
	} else {
		m.Html = types.StringNull()
	}
}

func (m *SettingDefaultSite) ToRequest(ctx context.Context, diags *diag.Diagnostics) *nginxproxymanager.UpdateSettingRequest {
	meta := nginxproxymanager.UpdateSettingRequestMeta{}
	meta.SetRedirect(m.Redirect.ValueString())
	meta.SetHtml(m.Html.ValueString())

	request := nginxproxymanager.NewUpdateSettingRequest()
	request.SetValue(m.Page.ValueString())
	request.SetMeta(meta)

	return request
}

func ObjectSettingDefaultSiteFrom(ctx context.Context, setting nginxproxymanager.GetSettings200ResponseInner) (types.Object, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	attributes := SettingDefaultSite{}
	attributes.Write(ctx, &setting, &diags)

	object, objectDiags := types.ObjectValueFrom(ctx, SettingDefaultSite{}.GetType().AttributeTypes(), attributes)
	diags.Append(objectDiags...)

	return object, diags
}

func SettingDefaultSiteAs(ctx context.Context, object types.Object) (SettingDefaultSite, diag.Diagnostics) {
	setting := SettingDefaultSite{}
	diags := object.As(ctx, &setting, basetypes.ObjectAsOptions{})

	return setting, diags
}
