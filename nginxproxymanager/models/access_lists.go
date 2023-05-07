package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/models"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/common"
)

type AccessLists struct {
	common.IModel[models.AccessListResourceCollection]
	AccessLists []AccessList `tfsdk:"access_lists"`
}

func (m *AccessLists) Load(ctx context.Context, resource *models.AccessListResourceCollection) diag.Diagnostics {
	diags := diag.Diagnostics{}
	m.AccessLists = make([]AccessList, len(*resource))
	for i, accessList := range *resource {
		diags.Append(m.AccessLists[i].Load(ctx, &accessList)...)
	}

	return diags
}
