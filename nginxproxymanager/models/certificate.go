package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/resources"
)

type Certificate struct {
	ID         types.Int64  `tfsdk:"id"`
	CreatedOn  types.String `tfsdk:"created_on"`
	ModifiedOn types.String `tfsdk:"modified_on"`
	Meta       types.Map    `tfsdk:"meta"`

	Provider    types.String `tfsdk:"provider_name"`
	NiceName    types.String `tfsdk:"nice_name"`
	DomainNames types.List   `tfsdk:"domain_names"`
	ExpiresOn   types.String `tfsdk:"expires_on"`
}

func (m *Certificate) Load(ctx context.Context, resource *resources.Certificate) diag.Diagnostics {
	meta, diags := types.MapValueFrom(ctx, types.StringType, resource.Meta.Map())
	domainNames, diags2 := types.ListValueFrom(ctx, types.StringType, resource.DomainNames)
	diags.Append(diags2...)

	m.ID = types.Int64Value(resource.ID)
	m.CreatedOn = types.StringValue(resource.CreatedOn)
	m.ModifiedOn = types.StringValue(resource.ModifiedOn)
	m.Meta = meta

	m.Provider = types.StringValue(resource.Provider)
	m.NiceName = types.StringValue(resource.NiceName)
	m.DomainNames = domainNames
	m.ExpiresOn = types.StringValue(resource.ExpiresOn)

	return diags
}
