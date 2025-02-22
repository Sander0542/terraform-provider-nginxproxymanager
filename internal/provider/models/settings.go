// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
)

type Settings struct {
	Settings types.Set `tfsdk:"settings"`
}

func (m *Settings) Write(ctx context.Context, settings *[]nginxproxymanager.GetSettings200ResponseInner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	elements := make([]Setting, 0, len(*settings))
	for _, g := range *settings {
		item := Setting{}
		item.WriteInner(ctx, &g, diags)
		elements = append(elements, item)
	}

	m.Settings, tmpDiags = types.SetValueFrom(ctx, Setting{}.GetType(), elements)
	diags.Append(tmpDiags...)
}
