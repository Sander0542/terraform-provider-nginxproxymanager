// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
)

type Streams struct {
	Streams types.Set `tfsdk:"streams"`
}

func (m *Streams) Write(ctx context.Context, streams *[]nginxproxymanager.GetStreams200ResponseInner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	elements := make([]Stream, 0, len(*streams))
	for _, g := range *streams {
		item := Stream{}
		item.Write(ctx, &g, diags)
		elements = append(elements, item)
	}

	m.Streams, tmpDiags = types.SetValueFrom(ctx, Stream{}.GetType(), elements)
	diags.Append(tmpDiags...)
}
