package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var ProxyHostLocation = map[string]schema.Attribute{
	"path": schema.StringAttribute{
		Description: "The path associated with the location.",
		Required:    true,
	},
	"forward_scheme": schema.StringAttribute{
		Description: "The scheme used to forward requests to the location. Can be either `http` or `https`.",
		Required:    true,
		Validators: []validator.String{
			stringvalidator.OneOf("http", "https"),
		},
	},
	"forward_host": schema.StringAttribute{
		Description: "The host used to forward requests to the location.",
		Required:    true,
	},
	"forward_port": schema.Int64Attribute{
		Description: "The port used to forward requests to the location. Must be between 1 and 65535.",
		Required:    true,
		Validators: []validator.Int64{
			int64validator.Between(1, 65535),
		},
	},
	"advanced_config": schema.StringAttribute{
		Description: "The advanced configuration used by the location.",
		Computed:    true,
		Optional:    true,
		Default:     stringdefault.StaticString(""),
	},
}
