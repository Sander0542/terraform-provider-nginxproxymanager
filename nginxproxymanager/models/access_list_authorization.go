package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/models"
)

type AccessListAuthorization struct {
	ID         types.Int64  `tfsdk:"id"`
	CreatedOn  types.String `tfsdk:"created_on"`
	ModifiedOn types.String `tfsdk:"modified_on"`
	Meta       types.Map    `tfsdk:"meta"`

	Username     types.String `tfsdk:"username"`
	PasswordHint types.String `tfsdk:"password_hint"`
}

func (m *AccessListAuthorization) Load(ctx context.Context, resource *models.AccessListAuthResource) diag.Diagnostics {
	meta, diags := types.MapValueFrom(ctx, types.StringType, resource.Meta.Map())

	m.ID = types.Int64Value(resource.ID)
	m.CreatedOn = types.StringValue(resource.CreatedOn)
	m.ModifiedOn = types.StringValue(resource.ModifiedOn)
	m.Meta = meta

	m.Username = types.StringValue(resource.Username)
	m.PasswordHint = types.StringValue(resource.Hint)

	return diags
}
