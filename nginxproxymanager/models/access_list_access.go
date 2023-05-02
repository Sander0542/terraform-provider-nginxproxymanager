package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/models"
)

type AccessListAccess struct {
	ID         types.Int64  `tfsdk:"id"`
	CreatedOn  types.String `tfsdk:"created_on"`
	ModifiedOn types.String `tfsdk:"modified_on"`
	Meta       types.Map    `tfsdk:"meta"`

	Address   types.String `tfsdk:"address"`
	Directive types.String `tfsdk:"directive"`
}

func (m *AccessListAccess) Load(ctx context.Context, resource *models.AccessListClientResource) diag.Diagnostics {
	meta, diags := types.MapValueFrom(ctx, types.StringType, resource.Meta.Map())

	m.ID = types.Int64Value(resource.ID)
	m.CreatedOn = types.StringValue(resource.CreatedOn)
	m.ModifiedOn = types.StringValue(resource.ModifiedOn)
	m.Meta = meta

	m.Address = types.StringValue(resource.Address)
	m.Directive = types.StringValue(resource.Directive)

	return diags
}
