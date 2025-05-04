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

type AccessListAuthorization struct {
	Id         types.Int64  `tfsdk:"id"`
	CreatedOn  types.String `tfsdk:"created_on"`
	ModifiedOn types.String `tfsdk:"modified_on"`
	Meta       types.Map    `tfsdk:"meta"`

	Username     types.String `tfsdk:"username"`
	PasswordHint types.String `tfsdk:"password_hint"`
}

func (AccessListAuthorization) GetType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"id":            types.Int64Type,
		"created_on":    types.StringType,
		"modified_on":   types.StringType,
		"meta":          types.MapType{ElemType: types.StringType},
		"username":      types.StringType,
		"password_hint": types.StringType,
	})
}

func (m *AccessListAuthorization) Write(ctx context.Context, authorization *nginxproxymanager.GetAccessLists200ResponseInnerItemsInner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	m.Id = types.Int64Value(authorization.GetId())
	m.CreatedOn = types.StringValue(authorization.GetCreatedOn())
	m.ModifiedOn = types.StringValue(authorization.GetModifiedOn())

	m.Username = types.StringValue(authorization.GetUsername())
	m.PasswordHint = types.StringValue(authorization.GetHint())

	m.Meta, tmpDiags = MapMetaFrom(ctx, authorization.GetMeta())
	diags.Append(tmpDiags...)
}

func SetAccessListAuthorizationsFrom(ctx context.Context, authorizations []nginxproxymanager.GetAccessLists200ResponseInnerItemsInner) (types.Set, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	elements := make([]AccessListAuthorization, 0, len(authorizations))
	for _, g := range authorizations {
		item := AccessListAuthorization{}
		item.Write(ctx, &g, &diags)
		elements = append(elements, item)
	}

	set, setDiags := types.SetValueFrom(ctx, AccessListAuthorization{}.GetType(), elements)

	diags.Append(setDiags...)

	return set, diags
}
