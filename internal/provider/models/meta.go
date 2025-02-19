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

	elements := make(map[string]basetypes.StringValue, length)
	for key, value := range meta {
		if value == nil {
			elements[key] = types.StringNull()
		} else if reflect.ValueOf(value).Kind() == reflect.Ptr {
			elements[key] = types.StringValue(fmt.Sprintf("%v", reflect.Indirect(reflect.ValueOf(value))))
		} else {
			elements[key] = types.StringValue(fmt.Sprintf("%v", value))
		}
	}

	return types.MapValueFrom(ctx, types.StringType, elements)
}

func MetaElementsAs(ctx context.Context, metaMap types.Map) (map[string]interface{}, diag.Diagnostics) {
	meta := make(map[string]interface{}, len(metaMap.Elements()))
	for key, value := range metaMap.Elements() {
		meta[key] = value.String()
	}

	return meta, diag.Diagnostics{}
}
