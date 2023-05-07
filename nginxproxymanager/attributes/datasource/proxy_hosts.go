package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/utils"
)

var nestedProxyHost = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the proxy host.",
		Computed:    true,
	},
}

var ProxyHosts = map[string]schema.Attribute{
	"proxy_hosts": schema.ListNestedAttribute{
		Description: "The proxy hosts.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: utils.MergeMaps(ProxyHost, nestedProxyHost),
		},
	},
}
