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

type AccessListAccess struct {
	Id         types.Int64  `tfsdk:"id"`
	CreatedOn  types.String `tfsdk:"created_on"`
	ModifiedOn types.String `tfsdk:"modified_on"`
	Meta       types.Map    `tfsdk:"meta"`

	Address   types.String `tfsdk:"address"`
	Directive types.String `tfsdk:"directive"`
}

func (AccessListAccess) GetType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"id":          types.Int64Type,
		"created_on":  types.StringType,
		"modified_on": types.StringType,
		"meta":        types.MapType{ElemType: types.StringType},
		"address":     types.StringType,
		"directive":   types.StringType,
	})
}

func (m *AccessListAccess) Write(ctx context.Context, access *nginxproxymanager.GetAccessLists200ResponseInnerClientsInner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	m.Id = types.Int64Value(access.GetId())
	m.CreatedOn = types.StringValue(access.GetCreatedOn())
	m.ModifiedOn = types.StringValue(access.GetModifiedOn())

	m.Address = types.StringValue(*access.GetAddress().String)
	m.Directive = types.StringValue(access.GetDirective())

	m.Meta, tmpDiags = MapMetaFrom(ctx, access.GetMeta())
	diags.Append(tmpDiags...)
}

func SetAccessListAccessFrom(ctx context.Context, accessList []nginxproxymanager.GetAccessLists200ResponseInnerClientsInner) (types.Set, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	elements := make([]AccessListAccess, 0, len(accessList))
	for _, g := range accessList {
		item := AccessListAccess{}
		item.Write(ctx, &g, &diags)
		elements = append(elements, item)
	}

	set, setDiags := types.SetValueFrom(ctx, AccessListAccess{}.GetType(), elements)

	diags.Append(setDiags...)

	return set, diags
}
