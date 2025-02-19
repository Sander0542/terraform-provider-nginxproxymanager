// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"reflect"
)

func MapMetaFrom(ctx context.Context, meta map[string]interface{}) (basetypes.MapValue, diag.Diagnostics) {
	length := len(meta)
	if length == 0 {
		return types.MapNull(types.StringType), diag.Diagnostics{}
	}

	elements := make(map[string]string, length)
	for key, value := range meta {
		if reflect.ValueOf(value).Kind() == reflect.Ptr {
			elements[key] = fmt.Sprintf("%v", reflect.Indirect(reflect.ValueOf(value)))
		} else {
			elements[key] = fmt.Sprintf("%v", value)
		}
	}

	return types.MapValueFrom(ctx, types.StringType, elements)
}
