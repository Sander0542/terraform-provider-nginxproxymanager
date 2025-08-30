// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

func (UserPermissions) GetType() attr.TypeWithAttributeTypes {
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

func (m *UserPermissions) ToRequest(ctx context.Context, diags *diag.Diagnostics) *nginxproxymanager.GetAccessLists200ResponseInnerOwnerPermissions {
	request := nginxproxymanager.NewGetAccessLists200ResponseInnerOwnerPermissions(
		m.Visibility.ValueString(),
		m.AccessLists.ValueString(),
		m.DeadHosts.ValueString(),
		m.ProxyHosts.ValueString(),
		m.RedirectionHosts.ValueString(),
		m.Streams.ValueString(),
		m.Certificates.ValueString(),
	)

	return request
}

func ObjectUserPermissionsFrom(ctx context.Context, permissions nginxproxymanager.GetAccessLists200ResponseInnerOwnerPermissions) (types.Object, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	attributes := UserPermissions{}
	attributes.Write(ctx, &permissions, &diags)

	object, objectDiags := types.ObjectValueFrom(ctx, UserPermissions{}.GetType().AttributeTypes(), attributes)
	diags.Append(objectDiags...)

	return object, diags
}

func UserPermissionsAs(ctx context.Context, object types.Object) (UserPermissions, diag.Diagnostics) {
	permissions := UserPermissions{}
	diags := object.As(ctx, &permissions, basetypes.ObjectAsOptions{})

	return permissions, diags
}
