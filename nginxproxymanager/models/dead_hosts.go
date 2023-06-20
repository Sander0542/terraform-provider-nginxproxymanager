package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/resources"
)

type DeadHosts struct {
	DeadHosts []DeadHost `tfsdk:"dead_hosts"`
}

func (m *DeadHosts) Load(ctx context.Context, resource *resources.DeadHostCollection) diag.Diagnostics {
	diags := diag.Diagnostics{}
	m.DeadHosts = make([]DeadHost, len(*resource))
	for i, deadHost := range *resource {
		diags.Append(m.DeadHosts[i].Load(ctx, &deadHost)...)
	}

	return diags
}
