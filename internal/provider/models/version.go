// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
)

type Version struct {
	Major    types.Int64  `tfsdk:"major"`
	Minor    types.Int64  `tfsdk:"minor"`
	Revision types.Int64  `tfsdk:"revision"`
	Version  types.String `tfsdk:"version"`
}

func (m *Version) Write(_ context.Context, version *nginxproxymanager.Health200ResponseVersion, _ *diag.Diagnostics) {
	m.Major = types.Int64Value(version.GetMajor())
	m.Minor = types.Int64Value(version.GetMinor())
	m.Revision = types.Int64Value(version.GetRevision())
	m.Version = types.StringValue(fmt.Sprintf("%d.%d.%d", version.GetMajor(), version.GetMinor(), version.GetRevision()))
}
