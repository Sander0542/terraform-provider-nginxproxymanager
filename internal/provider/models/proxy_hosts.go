// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
)

type ProxyHosts struct {
	ProxyHosts types.Set `tfsdk:"proxy_hosts"`
}

func (m *ProxyHosts) Write(ctx context.Context, proxyHosts *[]nginxproxymanager.GetProxyHosts200ResponseInner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	elements := make([]ProxyHost, 0, len(*proxyHosts))
	for _, g := range *proxyHosts {
		item := ProxyHost{}
		item.Write(ctx, &g, diags)
		elements = append(elements, item)
	}

	m.ProxyHosts, tmpDiags = types.SetValueFrom(ctx, ProxyHost{}.GetType(), elements)
	diags.Append(tmpDiags...)
}
