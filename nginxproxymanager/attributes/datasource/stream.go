package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var Stream = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the stream.",
		Required:    true,
	},
	"created_on": schema.StringAttribute{
		Description: "The date and time the stream was created.",
		Computed:    true,
	},
	"modified_on": schema.StringAttribute{
		Description: "The date and time the stream was last modified.",
		Computed:    true,
	},
	"owner_user_id": schema.Int64Attribute{
		Description: "The ID of the user that owns the stream.",
		Computed:    true,
	},
	"incoming_port": schema.Int64Attribute{
		Description: "The incoming port of the stream.",
		Computed:    true,
	},
	"forwarding_host": schema.StringAttribute{
		Description: "The forwarding host of the stream.",
		Computed:    true,
	},
	"forwarding_port": schema.Int64Attribute{
		Description: "The forwarding port of the stream.",
		Computed:    true,
	},
	"tcp_forwarding": schema.BoolAttribute{
		Description: "Whether TCP forwarding is enabled.",
		Computed:    true,
	},
	"udp_forwarding": schema.BoolAttribute{
		Description: "Whether UDP forwarding is enabled.",
		Computed:    true,
	},
	"enabled": schema.BoolAttribute{
		Description: "Whether the stream is enabled.",
		Computed:    true,
	},
	"meta": schema.MapAttribute{
		Description: "The meta data associated with the stream.",
		ElementType: types.StringType,
		Computed:    true,
	},
}
