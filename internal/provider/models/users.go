// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
)

type Users struct {
	Users types.Set `tfsdk:"users"`
}

func (m *Users) Write(ctx context.Context, users *[]nginxproxymanager.GetAccessLists200ResponseInnerOwner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	elements := make([]User, 0, len(*users))
	for _, g := range *users {
		item := User{}
		item.Write(ctx, &g, diags)
		elements = append(elements, item)
	}

	m.Users, tmpDiags = types.SetValueFrom(ctx, User{}.GetType(), elements)
	diags.Append(tmpDiags...)
}
