// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func SetDomainNamesFrom(ctx context.Context, domainNames []string) (basetypes.SetValue, diag.Diagnostics) {
	return types.SetValueFrom(ctx, types.StringType, domainNames)
}

func DomainNameElementsAs(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
	domainNames := make([]string, len(set.Elements()))
	diags := set.ElementsAs(ctx, &domainNames, false)

	return domainNames, diags
}
