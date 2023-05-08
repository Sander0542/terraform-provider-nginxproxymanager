package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/resources"
)

type Streams struct {
	Streams []Stream `tfsdk:"streams"`
}

func (m *Streams) Load(ctx context.Context, resource *resources.StreamCollection) diag.Diagnostics {
	diags := diag.Diagnostics{}
	m.Streams = make([]Stream, len(*resource))
	for i, stream := range *resource {
		diags.Append(m.Streams[i].Load(ctx, &stream)...)
	}

	return diags
}
