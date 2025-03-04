// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func ListDomainNamesFrom(ctx context.Context, domainNames []string) (basetypes.ListValue, diag.Diagnostics) {
	return types.ListValueFrom(ctx, types.StringType, domainNames)
}

func DomainNameElementsAs(ctx context.Context, list types.List) ([]string, diag.Diagnostics) {
	domainNames := make([]string, len(list.Elements()))
	diags := list.ElementsAs(ctx, &domainNames, false)

	return domainNames, diags
}
