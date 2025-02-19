// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
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

func (_ UserPermissions) getType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"visibility":        types.StringType,
		"access_lists":      types.StringType,
		"certificates":      types.StringType,
		"dead_hosts":        types.StringType,
		"proxy_hosts":       types.StringType,
		"redirection_hosts": types.StringType,
		"streams":           types.StringType,
	})
}

func (m *UserPermissions) Write(ctx context.Context, permissions *nginxproxymanager.GetAccessLists200ResponseInnerOwnerPermissions, diags *diag.Diagnostics) {
	m.AccessLists = types.StringValue(permissions.GetAccessLists())
	m.Certificates = types.StringValue(permissions.GetCertificates())
	m.DeadHosts = types.StringValue(permissions.GetDeadHosts())
	m.ProxyHosts = types.StringValue(permissions.GetProxyHosts())
	m.RedirectionHosts = types.StringValue(permissions.GetRedirectionHosts())
	m.Streams = types.StringValue(permissions.GetStreams())
	m.Visibility = types.StringValue(permissions.GetVisibility())
}

func ObjectUserPermissionsFrom(ctx context.Context, permissions nginxproxymanager.GetAccessLists200ResponseInnerOwnerPermissions) (types.Object, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	attributes := UserPermissions{}
	attributes.Write(ctx, &permissions, &diags)

	object, objectDiags := types.ObjectValueFrom(ctx, UserPermissions{}.getType().(attr.TypeWithAttributeTypes).AttributeTypes(), attributes)
	diags.Append(objectDiags...)

	return object, diags
}
