package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ProxyHost = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the proxy host.",
		Computed:    true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
	},
	"created_on": schema.StringAttribute{
		Description: "The date and time the proxy host was created.",
		Computed:    true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	},
	"modified_on": schema.StringAttribute{
		Description: "The date and time the proxy host was last modified.",
		Computed:    true,
	},
	"owner_user_id": schema.Int64Attribute{
		Description: "The ID of the user that owns the proxy host.",
		Computed:    true,
	},
	"domain_names": schema.ListAttribute{
		Description: "The domain names associated with the proxy host.",
		Required:    true,
		ElementType: types.StringType,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
	},
	"forward_scheme": schema.StringAttribute{
		Description: "The scheme used to forward requests to the proxy host. Can be either `http` or `https`.",
		Required:    true,
		Validators: []validator.String{
			stringvalidator.OneOf("http", "https"),
		},
	},
	"forward_host": schema.StringAttribute{
		Description: "The host used to forward requests to the proxy host.",
		Required:    true,
	},
	"forward_port": schema.Int64Attribute{
		Description: "The port used to forward requests to the proxy host. Must be between 1 and 65535.",
		Required:    true,
		Validators: []validator.Int64{
			int64validator.Between(1, 65535),
		},
	},
	"certificate_id": schema.StringAttribute{
		Description: "The ID of the certificate used by the proxy host.",
		Computed:    true,
		Optional:    true,
		Default:     stringdefault.StaticString("0"),
	},
	"ssl_forced": schema.BoolAttribute{
		Description: "Whether SSL is forced for the proxy host.",
		Computed:    true,
		Optional:    true,
		Default:     booldefault.StaticBool(false),
	},
	"hsts_enabled": schema.BoolAttribute{
		Description: "Whether HSTS is enabled for the proxy host.",
		Computed:    true,
		Optional:    true,
		Default:     booldefault.StaticBool(false),
	},
	"hsts_subdomains": schema.BoolAttribute{
		Description: "Whether HSTS is enabled for subdomains of the proxy host.",
		Computed:    true,
		Optional:    true,
		Default:     booldefault.StaticBool(false),
	},
	"http2_support": schema.BoolAttribute{
		Description: "Whether HTTP/2 is supported for the proxy host.",
		Computed:    true,
		Optional:    true,
		Default:     booldefault.StaticBool(false),
	},
	"block_exploits": schema.BoolAttribute{
		Description: "Whether exploits are blocked for the proxy host.",
		Computed:    true,
		Optional:    true,
		Default:     booldefault.StaticBool(false),
	},
	"caching_enabled": schema.BoolAttribute{
		Description: "Whether caching is enabled for the proxy host.",
		Computed:    true,
		Optional:    true,
		Default:     booldefault.StaticBool(false),
	},
	"allow_websocket_upgrade": schema.BoolAttribute{
		Description: "Whether websocket upgrades are allowed for the proxy host.",
		Computed:    true,
		Optional:    true,
		Default:     booldefault.StaticBool(true),
	},
	"access_list_id": schema.Int64Attribute{
		Description: "The ID of the access list used by the proxy host.",
		Computed:    true,
		Optional:    true,
		Default:     int64default.StaticInt64(0),
	},
	"advanced_config": schema.StringAttribute{
		Description: "The advanced configuration used by the proxy host.",
		Computed:    true,
		Optional:    true,
		Default:     stringdefault.StaticString(""),
	},
	"enabled": schema.BoolAttribute{
		Description: "Whether the proxy host is enabled.",
		Computed:    true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
	},
	"meta": schema.MapAttribute{
		Description: "The meta data associated with the proxy host.",
		ElementType: types.StringType,
		Computed:    true,
	},
}
