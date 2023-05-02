package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/utils"
)

var nestedAccessList = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the access list.",
		Computed:    true,
	},
}

var AccessLists = map[string]schema.Attribute{
	"access_lists": schema.ListNestedAttribute{
		Description: "The access lists.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: utils.MergeMaps(AccessList, nestedAccessList),
		},
	},
}
