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

type AccessListAccessResource struct {
	Address   types.String `tfsdk:"address"`
	Directive types.String `tfsdk:"directive"`
}

func (_ AccessListAccessResource) GetType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"address":   types.StringType,
		"directive": types.StringType,
	})
}

func (m *AccessListAccessResource) Write(ctx context.Context, access *nginxproxymanager.GetAccessLists200ResponseInnerClientsInner, diags *diag.Diagnostics) {
	m.Address = types.StringValue(*access.GetAddress().String)
	m.Directive = types.StringValue(access.GetDirective())
}

func (m *AccessListAccessResource) Read(ctx context.Context, diags *diag.Diagnostics) *nginxproxymanager.CreateAccessListRequestClientsInner {
	address := nginxproxymanager.GetAccessLists200ResponseInnerClientsInnerAddress{}
	address.String = m.Address.ValueStringPointer()

	location := nginxproxymanager.NewCreateAccessListRequestClientsInner()
	location.SetDirective(m.Directive.ValueString())
	location.SetAddress(address)

	return location
}

func SetAccessListAccessResourcesFrom(ctx context.Context, accessList []nginxproxymanager.GetAccessLists200ResponseInnerClientsInner) (types.Set, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	elements := make([]AccessListAccessResource, 0, len(accessList))
	for _, g := range accessList {
		item := AccessListAccessResource{}
		item.Write(ctx, &g, &diags)
		elements = append(elements, item)
	}

	set, setDiags := types.SetValueFrom(ctx, AccessListAccessResource{}.GetType(), elements)

	diags.Append(setDiags...)

	return set, diags
}

func AccessListAccessResourceElementsAs(ctx context.Context, set types.Set) ([]AccessListAccessResource, diag.Diagnostics) {
	access := make([]AccessListAccessResource, len(set.Elements()))
	diags := set.ElementsAs(ctx, &access, false)

	return access, diags
}
