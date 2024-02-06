package datasource

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/resources"
)

type ProxyHosts struct {
	ProxyHosts []ProxyHost `tfsdk:"proxy_hosts"`
}

func (m *ProxyHosts) Load(ctx context.Context, resource *resources.ProxyHostCollection) diag.Diagnostics {
	diags := diag.Diagnostics{}
	m.ProxyHosts = make([]ProxyHost, len(*resource))
	for i, proxyHost := range *resource {
		diags.Append(m.ProxyHosts[i].Load(ctx, &proxyHost)...)
	}

	return diags
}
