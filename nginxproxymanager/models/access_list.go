package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/models"
)

type AccessList struct {
	ID          types.Int64  `tfsdk:"id"`
	CreatedOn   types.String `tfsdk:"created_on"`
	ModifiedOn  types.String `tfsdk:"modified_on"`
	OwnerUserID types.Int64  `tfsdk:"owner_user_id"`
	Meta        types.Map    `tfsdk:"meta"`

	Name          types.String              `tfsdk:"name"`
	Authorization []AccessListAuthorization `tfsdk:"authorization"`
	Access        []AccessListAccess        `tfsdk:"access"`
	PassAuth      types.Bool                `tfsdk:"pass_auth"`
	SatisfyAny    types.Bool                `tfsdk:"satisfy_any"`
}

func (m *AccessList) Load(ctx context.Context, resource *models.AccessListResource) diag.Diagnostics {
	meta, diags := types.MapValueFrom(ctx, types.StringType, resource.Meta.Map())

	m.ID = types.Int64Value(resource.ID)
	m.CreatedOn = types.StringValue(resource.CreatedOn)
	m.ModifiedOn = types.StringValue(resource.ModifiedOn)
	m.OwnerUserID = types.Int64Value(resource.OwnerUserID)
	m.Meta = meta

	m.Name = types.StringValue(resource.Name)
	m.PassAuth = types.BoolValue(resource.PassAuth.Bool())
	m.SatisfyAny = types.BoolValue(resource.SatisfyAny.Bool())

	m.Authorization = make([]AccessListAuthorization, len(resource.Items))
	for i, v := range resource.Items {
		diags.Append(m.Authorization[i].Load(ctx, &v)...)
	}
	m.Access = make([]AccessListAccess, len(resource.Clients))
	for i, v := range resource.Clients {
		diags.Append(m.Access[i].Load(ctx, &v)...)
	}

	return diags
}
