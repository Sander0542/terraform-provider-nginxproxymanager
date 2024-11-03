package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/resources"
)

type Stream struct {
	ID          types.Int64  `tfsdk:"id"`
	CreatedOn   types.String `tfsdk:"created_on"`
	ModifiedOn  types.String `tfsdk:"modified_on"`
	OwnerUserId types.Int64  `tfsdk:"owner_user_id"`
	Meta        types.Map    `tfsdk:"meta"`

	IncomingPort   types.Int64  `tfsdk:"incoming_port"`
	ForwardingHost types.String `tfsdk:"forwarding_host"`
	ForwardingPort types.Int64  `tfsdk:"forwarding_port"`
	TCPForwarding  types.Bool   `tfsdk:"tcp_forwarding"`
	UDPForwarding  types.Bool   `tfsdk:"udp_forwarding"`
	Enabled        types.Bool   `tfsdk:"enabled"`
}

func (m *Stream) Load(ctx context.Context, resource *resources.Stream) diag.Diagnostics {
	meta, diags := types.MapValueFrom(ctx, types.StringType, resource.Meta.Map())

	m.ID = types.Int64Value(resource.ID)
	m.CreatedOn = types.StringValue(resource.CreatedOn)
	m.ModifiedOn = types.StringValue(resource.ModifiedOn)
	m.OwnerUserId = types.Int64Value(resource.OwnerUserID)
	m.Meta = meta

	m.IncomingPort = types.Int64Value(int64(resource.IncomingPort))
	m.ForwardingHost = types.StringValue(resource.ForwardingHost)
	m.ForwardingPort = types.Int64Value(int64(resource.ForwardingPort))
	m.TCPForwarding = types.BoolValue(resource.TCPForwarding)
	m.UDPForwarding = types.BoolValue(resource.UDPForwarding)
	m.Enabled = types.BoolValue(resource.Enabled)

	return diags
}
