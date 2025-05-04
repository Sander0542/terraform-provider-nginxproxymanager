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

type AccessList struct {
	Id          types.Int64  `tfsdk:"id"`
	CreatedOn   types.String `tfsdk:"created_on"`
	ModifiedOn  types.String `tfsdk:"modified_on"`
	OwnerUserId types.Int64  `tfsdk:"owner_user_id"`
	Meta        types.Map    `tfsdk:"meta"`

	Name           types.String `tfsdk:"name"`
	Authorizations types.Set    `tfsdk:"authorizations"`
	Access         types.Set    `tfsdk:"access"`
	PassAuth       types.Bool   `tfsdk:"pass_auth"`
	SatisfyAny     types.Bool   `tfsdk:"satisfy_any"`
}

func (AccessList) GetType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"id":             types.Int64Type,
		"created_on":     types.StringType,
		"modified_on":    types.StringType,
		"owner_user_id":  types.Int64Type,
		"meta":           types.MapType{ElemType: types.StringType},
		"name":           types.StringType,
		"authorizations": types.SetType{ElemType: AccessListAuthorization{}.GetType()},
		"access":         types.SetType{ElemType: AccessListAccess{}.GetType()},
		"pass_auth":      types.BoolType,
		"satisfy_any":    types.BoolType,
	})
}

func (m *AccessList) Write(ctx context.Context, accessList *nginxproxymanager.GetAccessLists200ResponseInner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	m.Id = types.Int64Value(accessList.GetId())
	m.CreatedOn = types.StringValue(accessList.GetCreatedOn())
	m.ModifiedOn = types.StringValue(accessList.GetModifiedOn())
	m.OwnerUserId = types.Int64Value(accessList.GetOwnerUserId())

	m.Name = types.StringValue(accessList.GetName())
	m.PassAuth = types.BoolValue(accessList.GetPassAuth())
	m.SatisfyAny = types.BoolValue(accessList.GetSatisfyAny())

	m.Authorizations, tmpDiags = SetAccessListAuthorizationsFrom(ctx, accessList.GetItems())
	diags.Append(tmpDiags...)

	m.Access, tmpDiags = SetAccessListAccessFrom(ctx, accessList.GetClients())
	diags.Append(tmpDiags...)

	m.Meta, tmpDiags = MapMetaFrom(ctx, accessList.GetMeta())
	diags.Append(tmpDiags...)
}
