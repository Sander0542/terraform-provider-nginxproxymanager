package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/resources"
)

type User struct {
	ID         types.Int64  `tfsdk:"id"`
	CreatedOn  types.String `tfsdk:"created_on"`
	ModifiedOn types.String `tfsdk:"modified_on"`
	Meta       types.Map    `tfsdk:"meta"`

	Name        types.String    `tfsdk:"name"`
	Nickname    types.String    `tfsdk:"nickname"`
	Email       types.String    `tfsdk:"email"`
	Avatar      types.String    `tfsdk:"avatar"`
	IsDisabled  types.Bool      `tfsdk:"is_disabled"`
	Roles       []types.String  `tfsdk:"roles"`
	Permissions UserPermissions `tfsdk:"permissions"`
}

func (m *User) Load(ctx context.Context, resource *resources.User) diag.Diagnostics {
	meta, diags := types.MapValueFrom(ctx, types.StringType, resource.Meta.Map())

	m.ID = types.Int64Value(resource.ID)
	m.CreatedOn = types.StringValue(resource.CreatedOn)
	m.ModifiedOn = types.StringValue(resource.ModifiedOn)
	m.Meta = meta

	m.Name = types.StringValue(resource.Name)
	m.Nickname = types.StringValue(resource.Nickname)
	m.Email = types.StringValue(resource.Email)
	m.Avatar = types.StringValue(resource.Avatar)
	m.IsDisabled = types.BoolValue(resource.IsDisabled.Bool())

	m.Roles = make([]types.String, len(resource.Roles))
	for i, v := range resource.Roles {
		m.Roles[i] = types.StringValue(v)
	}

	diags.Append(m.Permissions.Load(ctx, &resource.Permissions)...)

	return diags
}
