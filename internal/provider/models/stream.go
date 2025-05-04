// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
)

type Stream struct {
	Id          types.Int64  `tfsdk:"id"`
	CreatedOn   types.String `tfsdk:"created_on"`
	ModifiedOn  types.String `tfsdk:"modified_on"`
	OwnerUserId types.Int64  `tfsdk:"owner_user_id"`
	Meta        types.Map    `tfsdk:"meta"`

	IncomingPort   types.Int64  `tfsdk:"incoming_port"`
	ForwardingHost types.String `tfsdk:"forwarding_host"`
	ForwardingPort types.Int64  `tfsdk:"forwarding_port"`
	TcpForwarding  types.Bool   `tfsdk:"tcp_forwarding"`
	UdpForwarding  types.Bool   `tfsdk:"udp_forwarding"`
	CertificateId  types.Int64  `tfsdk:"certificate_id"`
	Enabled        types.Bool   `tfsdk:"enabled"`
}

func (Stream) GetType() attr.Type {
	return types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"id":              types.Int64Type,
		"created_on":      types.StringType,
		"modified_on":     types.StringType,
		"owner_user_id":   types.Int64Type,
		"meta":            types.MapType{ElemType: types.StringType},
		"incoming_port":   types.Int64Type,
		"forwarding_host": types.StringType,
		"forwarding_port": types.Int64Type,
		"tcp_forwarding":  types.BoolType,
		"udp_forwarding":  types.BoolType,
		"certificate_id":  types.Int64Type,
		"enabled":         types.BoolType,
	})
}

func (m *Stream) Write(ctx context.Context, stream *nginxproxymanager.GetStreams200ResponseInner, diags *diag.Diagnostics) {
	var tmpDiags diag.Diagnostics

	m.Id = types.Int64Value(stream.GetId())
	m.CreatedOn = types.StringValue(stream.GetCreatedOn())
	m.ModifiedOn = types.StringValue(stream.GetModifiedOn())
	m.OwnerUserId = types.Int64Value(stream.GetOwnerUserId())

	m.IncomingPort = types.Int64Value(stream.GetIncomingPort())
	m.ForwardingHost = types.StringValue(*stream.GetForwardingHost().String)
	m.ForwardingPort = types.Int64Value(stream.GetForwardingPort())
	m.TcpForwarding = types.BoolValue(stream.GetTcpForwarding())
	m.UdpForwarding = types.BoolValue(stream.GetUdpForwarding())
	if *stream.GetCertificateId().Int64 != 0 {
		m.CertificateId = types.Int64Value(*stream.GetCertificateId().Int64)
	} else {
		m.CertificateId = types.Int64Null()
	}
	m.Enabled = types.BoolValue(stream.GetEnabled())

	m.Meta, tmpDiags = MapMetaFrom(ctx, stream.GetMeta())
	diags.Append(tmpDiags...)
}

func (m *Stream) ToCreateRequest(_ context.Context, _ *diag.Diagnostics) *nginxproxymanager.CreateStreamRequest {
	forwardHost := nginxproxymanager.GetStreams200ResponseInnerForwardingHost{}
	forwardHost.String = m.ForwardingHost.ValueStringPointer()

	request := nginxproxymanager.NewCreateStreamRequest(
		m.IncomingPort.ValueInt64(),
		forwardHost,
		m.ForwardingPort.ValueInt64(),
	)

	request.SetTcpForwarding(m.TcpForwarding.ValueBool())
	request.SetUdpForwarding(m.UdpForwarding.ValueBool())
	certificateId := m.CertificateId.ValueInt64()
	request.SetCertificateId(nginxproxymanager.GetProxyHosts200ResponseInnerCertificateId{
		Int64: &certificateId,
	})
	request.SetMeta(map[string]interface{}{})

	return request
}

func (m *Stream) ToUpdateRequest(_ context.Context, _ *diag.Diagnostics) *nginxproxymanager.UpdateStreamRequest {
	forwardHost := nginxproxymanager.GetStreams200ResponseInnerForwardingHost{}
	forwardHost.String = m.ForwardingHost.ValueStringPointer()

	request := nginxproxymanager.NewUpdateStreamRequest()
	request.SetIncomingPort(m.IncomingPort.ValueInt64())
	request.SetForwardingHost(forwardHost)
	request.SetForwardingPort(m.ForwardingPort.ValueInt64())
	request.SetTcpForwarding(m.TcpForwarding.ValueBool())
	request.SetUdpForwarding(m.UdpForwarding.ValueBool())
	certificateId := m.CertificateId.ValueInt64()
	request.SetCertificateId(nginxproxymanager.GetProxyHosts200ResponseInnerCertificateId{
		Int64: &certificateId,
	})

	return request
}
