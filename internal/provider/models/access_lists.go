// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
)

type AccessLists struct {
	AccessLists types.Set `tfsdk:"access_lists"`
}

func (m *AccessLists) Write(ctx context.Context, accessLists *[]nginxproxymanager.GetAccessLists200ResponseInner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	elements := make([]AccessList, 0, len(*accessLists))
	for _, g := range *accessLists {
		item := AccessList{}
		item.Write(ctx, &g, diags)
		elements = append(elements, item)
	}

	m.AccessLists, tmpDiags = types.SetValueFrom(ctx, AccessList{}.GetType(), elements)
	diags.Append(tmpDiags...)
}
