// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/sander0542/nginxproxymanager-go"
)

var _ datasource.DataSource = &StreamsDataSource{}

func NewStreamsDataSource() datasource.DataSource {
	return &StreamsDataSource{}
}

type StreamsDataSource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (d *StreamsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_streams"
}

func (d *StreamsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Hosts --- This data source can be used to get information on all streams.",
		Attributes: map[string]schema.Attribute{
			"streams": schema.SetNestedAttribute{
				MarkdownDescription: "The streams.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "The Id of the stream.",
							Computed:            true,
						},
						"created_on": schema.StringAttribute{
							MarkdownDescription: "The date and time the stream was created.",
							Computed:            true,
						},
						"modified_on": schema.StringAttribute{
							MarkdownDescription: "The date and time the stream was last modified.",
							Computed:            true,
						},
						"owner_user_id": schema.Int64Attribute{
							MarkdownDescription: "The Id of the user that owns the stream.",
							Computed:            true,
						},
						"incoming_port": schema.Int64Attribute{
							MarkdownDescription: "The incoming port of the stream.",
							Computed:            true,
						},
						"forwarding_host": schema.StringAttribute{
							MarkdownDescription: "The forwarding host of the stream.",
							Computed:            true,
						},
						"forwarding_port": schema.Int64Attribute{
							MarkdownDescription: "The forwarding port of the stream.",
							Computed:            true,
						},
						"tcp_forwarding": schema.BoolAttribute{
							MarkdownDescription: "Whether TCP forwarding is enabled.",
							Computed:            true,
						},
						"udp_forwarding": schema.BoolAttribute{
							MarkdownDescription: "Whether UDP forwarding is enabled.",
							Computed:            true,
						},
						"certificate_id": schema.Int64Attribute{
							MarkdownDescription: "The Id of the certificate used by the stream.",
							Computed:            true,
						},
						"enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether the stream is enabled.",
							Computed:            true,
						},
						"meta": schema.MapAttribute{
							MarkdownDescription: "The meta data associated with the stream.",
							ElementType:         types.StringType,
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *StreamsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *StreamsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *models.Streams

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, _, err := d.client.StreamsAPI.GetStreams(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read streams, got error: %s", err))
		return
	}

	data.Write(ctx, &response, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
