package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/utils"
)

var nestedDeadHost = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the proxy host.",
		Computed:    true,
	},
}

var DeadHosts = map[string]schema.Attribute{
	"dead_hosts": schema.ListNestedAttribute{
		Description: "The dead hosts.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: utils.MergeMaps(ProxyHost, nestedDeadHost),
		},
	},
}
