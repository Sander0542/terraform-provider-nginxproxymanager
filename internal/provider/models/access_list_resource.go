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

type AccessListResource struct {
	Id          types.Int64  `tfsdk:"id"`
	CreatedOn   types.String `tfsdk:"created_on"`
	ModifiedOn  types.String `tfsdk:"modified_on"`
	OwnerUserId types.Int64  `tfsdk:"owner_user_id"`
	Meta        types.Map    `tfsdk:"meta"`

	Name           types.String `tfsdk:"name"`
	Authorizations types.Set    `tfsdk:"authorizations"`
	Access         types.Set    `tfsdk:"access"`
	PassAuth       types.Bool   `tfsdk:"pass_auth"`
	SatisfyAny     types.Bool   `tfsdk:"satisfy_any"`
}

func (_ AccessListResource) GetType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"id":             types.Int64Type,
		"created_on":     types.StringType,
		"modified_on":    types.StringType,
		"owner_user_id":  types.Int64Type,
		"meta":           types.MapType{ElemType: types.StringType},
		"name":           types.StringType,
		"authorizations": types.SetType{ElemType: AccessListAuthorizationResource{}.GetType()},
		"access":         types.SetType{ElemType: AccessListAccessResource{}.GetType()},
		"pass_auth":      types.BoolType,
		"satisfy_any":    types.BoolType,
	})
}

func (m *AccessListResource) Write(ctx context.Context, accessList *nginxproxymanager.CreateAccessList201Response, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	m.Id = types.Int64Value(accessList.GetId())
	m.CreatedOn = types.StringValue(accessList.GetCreatedOn())
	m.ModifiedOn = types.StringValue(accessList.GetModifiedOn())
	m.OwnerUserId = types.Int64Value(accessList.GetOwnerUserId())

	m.Name = types.StringValue(accessList.GetName())
	m.PassAuth = types.BoolValue(accessList.GetPassAuth())
	m.SatisfyAny = types.BoolValue(accessList.GetSatisfyAny())

	m.Authorizations, tmpDiags = SetAccessListAuthorizationResourcesFrom(ctx, m.Authorizations, accessList.GetItems())
	diags.Append(tmpDiags...)

	m.Access, tmpDiags = SetAccessListAccessResourcesFrom(ctx, accessList.GetClients())
	diags.Append(tmpDiags...)

	m.Meta, tmpDiags = MapMetaFrom(ctx, accessList.GetMeta())
	diags.Append(tmpDiags...)
}

func (m *AccessListResource) ToCreateRequest(ctx context.Context, diags *diag.Diagnostics) *nginxproxymanager.CreateAccessListRequest {
	authorizations, tmpDiags := AccessListAuthorizationResourceElementsAs(ctx, m.Authorizations)
	diags.Append(tmpDiags...)

	access, tmpDiags := AccessListAccessResourceElementsAs(ctx, m.Access)
	diags.Append(tmpDiags...)

	request := nginxproxymanager.NewCreateAccessListRequest(m.Name.ValueString())

	request.SetSatisfyAny(m.SatisfyAny.ValueBool())
	request.SetPassAuth(m.PassAuth.ValueBool())
	request.SetMeta(map[string]interface{}{})

	requestAuthorizations := make([]nginxproxymanager.CreateAccessListRequestItemsInner, 0, len(authorizations))
	for _, authorization := range authorizations {
		item := authorization.Read(ctx, diags)
		requestAuthorizations = append(requestAuthorizations, *item)
	}
	request.SetItems(requestAuthorizations)

	requestAccess := make([]nginxproxymanager.CreateAccessListRequestClientsInner, 0, len(access))
	for _, accessItem := range access {
		item := accessItem.Read(ctx, diags)
		requestAccess = append(requestAccess, *item)
	}
	request.SetClients(requestAccess)

	return request
}

func (m *AccessListResource) ToUpdateRequest(ctx context.Context, diags *diag.Diagnostics) *nginxproxymanager.UpdateAccessListRequest {
	authorizations, tmpDiags := AccessListAuthorizationResourceElementsAs(ctx, m.Authorizations)
	diags.Append(tmpDiags...)

	access, tmpDiags := AccessListAccessResourceElementsAs(ctx, m.Access)
	diags.Append(tmpDiags...)

	request := nginxproxymanager.NewUpdateAccessListRequest()

	request.SetName(m.Name.ValueString())
	request.SetSatisfyAny(m.SatisfyAny.ValueBool())
	request.SetPassAuth(m.PassAuth.ValueBool())

	requestAuthorizations := make([]nginxproxymanager.CreateAccessListRequestItemsInner, 0, len(authorizations))
	for _, authorization := range authorizations {
		item := authorization.Read(ctx, diags)
		requestAuthorizations = append(requestAuthorizations, *item)
	}
	request.SetItems(requestAuthorizations)

	requestAccess := make([]nginxproxymanager.CreateAccessListRequestClientsInner, 0, len(access))
	for _, accessItem := range access {
		item := accessItem.Read(ctx, diags)
		requestAccess = append(requestAccess, *item)
	}
	request.SetClients(requestAccess)

	return request
}
