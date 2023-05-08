package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/utils"
)

var nestedRedirectionHost = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the redirection host.",
		Computed:    true,
	},
}

var RedirectionHosts = map[string]schema.Attribute{
	"redirection_hosts": schema.ListNestedAttribute{
		Description: "The redirection hosts.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: utils.MergeMaps(RedirectionHost, nestedRedirectionHost),
		},
	},
}
