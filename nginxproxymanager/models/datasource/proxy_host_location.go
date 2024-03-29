package datasource

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/resources"
)

type ProxyHostLocation struct {
	Path           types.String `tfsdk:"path"`
	ForwardScheme  types.String `tfsdk:"forward_scheme"`
	ForwardHost    types.String `tfsdk:"forward_host"`
	ForwardPort    types.Int64  `tfsdk:"forward_port"`
	AdvancedConfig types.String `tfsdk:"advanced_config"`
}

func (m *ProxyHostLocation) Load(_ context.Context, resource *resources.ProxyHostLocation) diag.Diagnostics {
	m.Path = types.StringValue(resource.Path)
	m.ForwardScheme = types.StringValue(resource.ForwardScheme)
	m.ForwardHost = types.StringValue(resource.ForwardHost)
	m.ForwardPort = types.Int64Value(int64(resource.ForwardPort))
	m.AdvancedConfig = types.StringValue(resource.AdvancedConfig)

	return diag.Diagnostics{}
}
