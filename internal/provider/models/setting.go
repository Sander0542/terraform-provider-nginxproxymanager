// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
)

type Setting struct {
	Id          types.String `tfsdk:"id"`
	Meta        types.Map    `tfsdk:"meta"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Value       types.String `tfsdk:"value"`
}

func (_ Setting) GetType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"id":          types.StringType,
		"meta":        types.MapType{ElemType: types.StringType},
		"name":        types.StringType,
		"description": types.StringType,
		"value":       types.StringType,
	})
}

func (m *Setting) Write(ctx context.Context, setting *nginxproxymanager.GetSetting200Response, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	m.Id = types.StringValue(setting.GetId())
	m.Name = types.StringValue(setting.GetName())
	m.Description = types.StringValue(setting.GetDescription())

	value := setting.GetValue()
	if value.String != nil {
		m.Value = types.StringValue(*value.String)
	} else if value.Int64 != nil {
		m.Value = types.StringValue(fmt.Sprintf("%d", *value.Int64))
	} else if value.Float32 != nil {
		m.Value = types.StringValue(fmt.Sprintf("%f", *value.Float32))
	} else {
		diags.AddWarning("Unsupported settings value", fmt.Sprintf("Setting with id %s contains an unsupported value", setting.GetId()))
		m.Value = types.StringNull()
	}

	m.Meta, tmpDiags = MapMetaFrom(ctx, setting.GetMeta())
	diags.Append(tmpDiags...)
}

func (m *Setting) WriteInner(ctx context.Context, setting *nginxproxymanager.GetSettings200ResponseInner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	m.Id = types.StringValue(setting.GetId())
	m.Name = types.StringValue(setting.GetName())
	m.Description = types.StringValue(setting.GetDescription())

	value := setting.GetValue()
	if value.String != nil {
		m.Value = types.StringValue(*value.String)
	} else if value.Int64 != nil {
		m.Value = types.StringValue(fmt.Sprintf("%d", *value.Int64))
	} else if value.Float32 != nil {
		m.Value = types.StringValue(fmt.Sprintf("%f", *value.Float32))
	} else {
		diags.AddWarning("Unsupported settings value", fmt.Sprintf("Setting with id %s contains an unsupported value", setting.GetId()))
		m.Value = types.StringNull()
	}

	m.Meta, tmpDiags = MapMetaFrom(ctx, setting.GetMeta())
	diags.Append(tmpDiags...)
}
