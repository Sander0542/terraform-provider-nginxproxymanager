package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/utils"
)

var nestedStream = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the stream.",
		Computed:    true,
	},
}

var Streams = map[string]schema.Attribute{
	"streams": schema.ListNestedAttribute{
		Description: "The streams.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: utils.MergeMaps(Stream, nestedStream),
		},
	},
}
