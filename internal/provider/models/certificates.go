// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
)

type Certificates struct {
	Certificates types.Set `tfsdk:"certificates"`
}

func (m *Certificates) Write(ctx context.Context, certificates *[]nginxproxymanager.GetCertificates200ResponseInner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	elements := make([]Certificate, 0, len(*certificates))
	for _, g := range *certificates {
		item := Certificate{}
		item.Write(ctx, &g, diags)
		elements = append(elements, item)
	}

	m.Certificates, tmpDiags = types.SetValueFrom(ctx, Certificate{}.getType(), elements)
	diags.Append(tmpDiags...)
}
