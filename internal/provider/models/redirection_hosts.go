// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
)

type RedirectionHosts struct {
	RedirectionHosts types.Set `tfsdk:"redirection_hosts"`
}

func (m *RedirectionHosts) Write(ctx context.Context, redirectionHosts *[]nginxproxymanager.GetRedirectionHosts200ResponseInner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	elements := make([]RedirectionHost, 0, len(*redirectionHosts))
	for _, g := range *redirectionHosts {
		item := RedirectionHost{}
		item.Write(ctx, &g, diags)
		elements = append(elements, item)
	}

	m.RedirectionHosts, tmpDiags = types.SetValueFrom(ctx, RedirectionHost{}.GetType(), elements)
	diags.Append(tmpDiags...)
}
