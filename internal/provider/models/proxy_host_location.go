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

type ProxyHostLocation struct {
	Path           types.String `tfsdk:"path"`
	ForwardScheme  types.String `tfsdk:"forward_scheme"`
	ForwardHost    types.String `tfsdk:"forward_host"`
	ForwardPort    types.Int64  `tfsdk:"forward_port"`
	AdvancedConfig types.String `tfsdk:"advanced_config"`
}

func (_ ProxyHostLocation) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"path":            types.StringType,
		"forward_scheme":  types.StringType,
		"forward_host":    types.StringType,
		"forward_port":    types.Int64Type,
		"advanced_config": types.StringType,
	})
}

func (m *ProxyHostLocation) Write(_ context.Context, location *nginxproxymanager.GetProxyHosts200ResponseInnerLocationsInner, _ *diag.Diagnostics) {
	m.Path = types.StringValue(location.GetPath())
	m.ForwardScheme = types.StringValue(location.GetForwardScheme())
	m.ForwardHost = types.StringValue(location.GetForwardHost())
	m.ForwardPort = types.Int64Value(location.GetForwardPort())
	m.AdvancedConfig = types.StringValue(location.GetAdvancedConfig())
}

func SetProxyHostLocationsFrom(ctx context.Context, locations []nginxproxymanager.GetProxyHosts200ResponseInnerLocationsInner) (types.Set, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	elements := make([]ProxyHostLocation, 0, len(locations))
	for _, g := range locations {
		item := ProxyHostLocation{}
		item.Write(ctx, &g, &diags)
		elements = append(elements, item)
	}

	set, setDiags := types.SetValueFrom(ctx, ProxyHostLocation{}.getType(), elements)

	diags.Append(setDiags...)

	return set, diags
}
