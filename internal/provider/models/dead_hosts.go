// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
)

type DeadHosts struct {
	DeadHosts types.Set `tfsdk:"dead_hosts"`
}

func (m *DeadHosts) Write(ctx context.Context, deadHosts *[]nginxproxymanager.GetDeadHosts200ResponseInner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	elements := make([]DeadHost, 0, len(*deadHosts))
	for _, g := range *deadHosts {
		item := DeadHost{}
		item.Write(ctx, &g, diags)
		elements = append(elements, item)
	}

	m.DeadHosts, tmpDiags = types.SetValueFrom(ctx, DeadHost{}.GetType(), elements)
	diags.Append(tmpDiags...)
}
