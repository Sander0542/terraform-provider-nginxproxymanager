package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/resources"
)

type UserPermissions struct {
	ID         types.Int64  `tfsdk:"id"`
	CreatedOn  types.String `tfsdk:"created_on"`
	ModifiedOn types.String `tfsdk:"modified_on"`
	Meta       types.Map    `tfsdk:"meta"`

	AccessLists      types.String `tfsdk:"access_lists"`
	Certificates     types.String `tfsdk:"certificates"`
	DeadHosts        types.String `tfsdk:"dead_hosts"`
	ProxyHosts       types.String `tfsdk:"proxy_hosts"`
	RedirectionHosts types.String `tfsdk:"redirection_hosts"`
	Streams          types.String `tfsdk:"streams"`
	Visibility       types.String `tfsdk:"visibility"`
}

func (m *UserPermissions) Load(ctx context.Context, resource *resources.UserPermissions) diag.Diagnostics {
	meta, diags := types.MapValueFrom(ctx, types.StringType, resource.Meta.Map())

	m.ID = types.Int64Value(resource.ID)
	m.CreatedOn = types.StringValue(resource.CreatedOn)
	m.ModifiedOn = types.StringValue(resource.ModifiedOn)
	m.Meta = meta

	m.AccessLists = types.StringValue(resource.AccessLists)
	m.Certificates = types.StringValue(resource.Certificates)
	m.DeadHosts = types.StringValue(resource.DeadHosts)
	m.ProxyHosts = types.StringValue(resource.ProxyHosts)
	m.RedirectionHosts = types.StringValue(resource.RedirectionHosts)
	m.Streams = types.StringValue(resource.Streams)
	m.Visibility = types.StringValue(resource.Visibility)

	return diags
}
