package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/utils"
)

var nestedUser = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the user.",
		Computed:    true,
	},
}

var Users = map[string]schema.Attribute{
	"users": schema.ListNestedAttribute{
		Description: "The users.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: utils.MergeMaps(User, nestedCertificate),
		},
	},
}
