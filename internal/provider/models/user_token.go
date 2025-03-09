// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type UserToken struct {
	Token types.String `tfsdk:"token"`
}
