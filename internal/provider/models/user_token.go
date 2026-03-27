// Copyright Sander Jochems 2025, 2026, 0
// SPDX-License-Identifier: MIT

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type UserToken struct {
	Token types.String `tfsdk:"token"`
}
