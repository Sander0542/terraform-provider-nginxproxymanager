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
	DefaultSite types.Object `tfsdk:"default_site"`
}

func (m *Settings) Write(ctx context.Context, settings []nginxproxymanager.GetSettings200ResponseInner, diags *diag.Diagnostics) {
	for _, setting := range settings {
		var tmpDiags diag.Diagnostics

		switch setting.GetId() {
		case "default-site":
			m.DefaultSite, tmpDiags = ObjectSettingDefaultSiteFrom(ctx, setting)
		}

		diags.Append(tmpDiags...)
	}
}
