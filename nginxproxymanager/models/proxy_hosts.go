package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/models"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/common"
)

type ProxyHosts struct {
	common.IModel[models.ProxyHostResourceCollection]
	ProxyHosts []ProxyHost `tfsdk:"proxy_hosts"`
}

func (m *ProxyHosts) Load(ctx context.Context, resource *models.ProxyHostResourceCollection) diag.Diagnostics {
	diags := diag.Diagnostics{}
	m.ProxyHosts = make([]ProxyHost, len(*resource))
	for i, proxyHost := range *resource {
		diags.Append(m.ProxyHosts[i].Load(ctx, &proxyHost)...)
	}

	return diags
}
