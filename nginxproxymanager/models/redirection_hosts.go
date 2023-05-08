package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/resources"
)

type RedirectionHosts struct {
	RedirectionHosts []RedirectionHost `tfsdk:"redirection_hosts"`
}

func (m *RedirectionHosts) Load(ctx context.Context, resource *resources.RedirectionHostCollection) diag.Diagnostics {
	diags := diag.Diagnostics{}
	m.RedirectionHosts = make([]RedirectionHost, len(*resource))
	for i, redirectionHost := range *resource {
		diags.Append(m.RedirectionHosts[i].Load(ctx, &redirectionHost)...)
	}

	return diags
}
