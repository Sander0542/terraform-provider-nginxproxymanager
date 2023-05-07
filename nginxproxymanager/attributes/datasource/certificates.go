package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/utils"
)

var nestedCertificate = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the certificate.",
		Computed:    true,
	},
}

var Certificates = map[string]schema.Attribute{
	"certificates": schema.ListNestedAttribute{
		Description: "The certificates.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: utils.MergeMaps(Certificate, nestedCertificate),
		},
	},
}
