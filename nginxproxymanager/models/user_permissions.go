package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/resources"
)

type UserPermissions struct {
	AccessLists      types.String `tfsdk:"access_lists"`
	Certificates     types.String `tfsdk:"certificates"`
	DeadHosts        types.String `tfsdk:"dead_hosts"`
	ProxyHosts       types.String `tfsdk:"proxy_hosts"`
	RedirectionHosts types.String `tfsdk:"redirection_hosts"`
	Streams          types.String `tfsdk:"streams"`
	Visibility       types.String `tfsdk:"visibility"`
}

func (m *UserPermissions) Load(_ context.Context, resource *resources.UserPermissions) diag.Diagnostics {
	m.AccessLists = types.StringValue(resource.AccessLists)
	m.Certificates = types.StringValue(resource.Certificates)
	m.DeadHosts = types.StringValue(resource.DeadHosts)
	m.ProxyHosts = types.StringValue(resource.ProxyHosts)
	m.RedirectionHosts = types.StringValue(resource.RedirectionHosts)
	m.Streams = types.StringValue(resource.Streams)
	m.Visibility = types.StringValue(resource.Visibility)

	return diag.Diagnostics{}
}
