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

type User struct {
	Id         types.Int64  `tfsdk:"id"`
	CreatedOn  types.String `tfsdk:"created_on"`
	ModifiedOn types.String `tfsdk:"modified_on"`

	Name        types.String `tfsdk:"name"`
	Nickname    types.String `tfsdk:"nickname"`
	Email       types.String `tfsdk:"email"`
	Avatar      types.String `tfsdk:"avatar"`
	IsDisabled  types.Bool   `tfsdk:"is_disabled"`
	Roles       types.Set    `tfsdk:"roles"`
	Permissions types.Object `tfsdk:"permissions"`
}

func (_ User) GetType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"id":          types.Int64Type,
		"created_on":  types.StringType,
		"modified_on": types.StringType,
		"name":        types.StringType,
		"nickname":    types.StringType,
		"email":       types.StringType,
		"avatar":      types.StringType,
		"is_disabled": types.BoolType,
		"roles":       types.SetType{ElemType: types.StringType},
		"permissions": types.ObjectType{AttrTypes: UserPermissions{}.GetType().AttributeTypes()},
	})
}

func (m *User) Write(ctx context.Context, user *nginxproxymanager.GetAccessLists200ResponseInnerOwner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	m.Id = types.Int64Value(user.GetId())
	m.CreatedOn = types.StringValue(user.GetCreatedOn())
	m.ModifiedOn = types.StringValue(user.GetModifiedOn())

	m.Name = types.StringValue(user.GetName())
	m.Nickname = types.StringValue(user.GetNickname())
	m.Email = types.StringValue(user.GetEmail())
	m.Avatar = types.StringValue(user.GetAvatar())
	m.IsDisabled = types.BoolValue(user.GetIsDisabled())

	m.Roles, tmpDiags = types.SetValueFrom(ctx, types.StringType, user.GetRoles())
	diags.Append(tmpDiags...)

	if user.HasPermissions() {
		m.Permissions, tmpDiags = ObjectUserPermissionsFrom(ctx, user.GetPermissions())
		diags.Append(tmpDiags...)
	} else {
		m.Permissions = types.ObjectNull(UserPermissions{}.GetType().AttributeTypes())
	}
}
