// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
	"slices"
)

type UserResource struct {
	Id         types.Int64  `tfsdk:"id"`
	CreatedOn  types.String `tfsdk:"created_on"`
	ModifiedOn types.String `tfsdk:"modified_on"`

	Name        types.String `tfsdk:"name"`
	Nickname    types.String `tfsdk:"nickname"`
	Email       types.String `tfsdk:"email"`
	Avatar      types.String `tfsdk:"avatar"`
	IsDisabled  types.Bool   `tfsdk:"is_disabled"`
	IsAdmin     types.Bool   `tfsdk:"is_admin"`
	Permissions types.Object `tfsdk:"permissions"`
}

func (m *UserResource) Write(ctx context.Context, user *nginxproxymanager.GetAccessLists200ResponseInnerOwner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	m.Id = types.Int64Value(user.GetId())
	m.CreatedOn = types.StringValue(user.GetCreatedOn())
	m.ModifiedOn = types.StringValue(user.GetModifiedOn())

	m.Name = types.StringValue(user.GetName())
	m.Nickname = types.StringValue(user.GetNickname())
	m.Email = types.StringValue(user.GetEmail())
	m.Avatar = types.StringValue(user.GetAvatar())
	m.IsDisabled = types.BoolValue(user.GetIsDisabled())
	m.IsAdmin = types.BoolValue(slices.Contains(user.GetRoles(), "admin"))

	if user.HasPermissions() {
		m.Permissions, tmpDiags = ObjectUserPermissionsFrom(ctx, user.GetPermissions())
		diags.Append(tmpDiags...)
	} else {
		m.Permissions = types.ObjectNull(UserPermissions{}.GetType().AttributeTypes())
	}
}

func (m *UserResource) ToCreateRequest(ctx context.Context, diags *diag.Diagnostics) *nginxproxymanager.CreateUserRequest {
	request := nginxproxymanager.NewCreateUserRequest(
		m.Name.ValueString(),
		m.Nickname.ValueString(),
		m.Email.ValueString(),
	)

	request.SetIsDisabled(m.IsDisabled.ValueBool())
	if m.IsAdmin.ValueBool() {
		request.SetRoles([]string{"admin"})
	} else {
		request.SetRoles([]string{})
	}

	return request
}

func (m *UserResource) ToUpdateRequest(ctx context.Context, diags *diag.Diagnostics) *nginxproxymanager.UpdateUserRequest {
	request := nginxproxymanager.NewUpdateUserRequest()
	request.SetName(m.Name.ValueString())
	request.SetNickname(m.Nickname.ValueString())
	request.SetEmail(m.Email.ValueString())
	request.SetIsDisabled(m.IsDisabled.ValueBool())
	if m.IsAdmin.ValueBool() {
		request.SetRoles([]string{"admin"})
	} else {
		request.SetRoles([]string{})
	}

	return request
}
