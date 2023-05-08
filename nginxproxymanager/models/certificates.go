package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/resources"
)

type Certificates struct {
	Certificates []Certificate `tfsdk:"certificates"`
}

func (m *Certificates) Load(ctx context.Context, resource *resources.CertificateCollection) diag.Diagnostics {
	diags := diag.Diagnostics{}
	m.Certificates = make([]Certificate, len(*resource))
	for i, certificate := range *resource {
		diags.Append(m.Certificates[i].Load(ctx, &certificate)...)
	}

	return diags
}
