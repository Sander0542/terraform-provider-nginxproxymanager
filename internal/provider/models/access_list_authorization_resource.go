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

type AccessListAuthorizationResource struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

type authorizationCredential struct {
	Username string
	Password string
}

func (m *authorizationCredential) matches(username string, passwordHint string) bool {
	if m.Username != username {
		return false
	}

	if len(m.Password) != len(passwordHint) {
		return false
	}

	if []rune(m.Password)[0] != []rune(passwordHint)[0] {
		return false
	}

	return true
}

func (_ AccessListAuthorizationResource) GetType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"username": types.StringType,
		"password": types.StringType,
	})
}

func (m *AccessListAuthorizationResource) Write(ctx context.Context, authorization *nginxproxymanager.GetAccessLists200ResponseInnerItemsInner, diags *diag.Diagnostics) {
	m.Username = types.StringValue(authorization.GetUsername())
	m.Password = types.StringUnknown()
}

func (m *AccessListAuthorizationResource) Read(ctx context.Context, diags *diag.Diagnostics) *nginxproxymanager.CreateAccessListRequestItemsInner {
	location := nginxproxymanager.NewCreateAccessListRequestItemsInner()
	location.SetUsername(m.Username.ValueString())
	location.SetPassword(m.Password.ValueString())

	return location
}

func SetAccessListAuthorizationResourcesFrom(ctx context.Context, currentSet types.Set, authorizations []nginxproxymanager.GetAccessLists200ResponseInnerItemsInner) (types.Set, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	currentAuthorizations, _ := AccessListAuthorizationResourceElementsAs(ctx, currentSet)
	credentials := make([]authorizationCredential, 0, len(currentAuthorizations))
	for _, currentAuthorization := range currentAuthorizations {
		credentials = append(credentials, authorizationCredential{
			Username: currentAuthorization.Username.ValueString(),
			Password: currentAuthorization.Password.ValueString(),
		})
	}

	elements := make([]AccessListAuthorizationResource, 0, len(authorizations))
	for _, g := range authorizations {
		item := AccessListAuthorizationResource{}
		item.Write(ctx, &g, &diags)

		for i, credential := range credentials {
			if credential.matches(g.GetUsername(), g.GetHint()) {
				item.Password = types.StringValue(credential.Password)
				credentials = append(credentials[:i], credentials[i+1:]...)
				break
			}
		}

		elements = append(elements, item)
	}

	set, setDiags := types.SetValueFrom(ctx, AccessListAuthorizationResource{}.GetType(), elements)

	diags.Append(setDiags...)

	return set, diags
}

func AccessListAuthorizationResourceElementsAs(ctx context.Context, set types.Set) ([]AccessListAuthorizationResource, diag.Diagnostics) {
	authorizations := make([]AccessListAuthorizationResource, len(set.Elements()))
	diags := set.ElementsAs(ctx, &authorizations, false)

	return authorizations, diags
}
