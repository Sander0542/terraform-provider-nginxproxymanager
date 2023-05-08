package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/resources"
)

type Users struct {
	Users []User `tfsdk:"users"`
}

func (m *Users) Load(ctx context.Context, resource *resources.UserCollection) diag.Diagnostics {
	diags := diag.Diagnostics{}
	m.Users = make([]User, len(*resource))
	for i, user := range *resource {
		diags.Append(m.Users[i].Load(ctx, &user)...)
	}

	return diags
}
