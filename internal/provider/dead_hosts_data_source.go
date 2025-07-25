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

var _ datasource.DataSource = &DeadHostsDataSource{}

func NewDeadHostsDataSource() datasource.DataSource {
	return &DeadHostsDataSource{}
}

type DeadHostsDataSource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (d *DeadHostsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dead_hosts"
}

func (d *DeadHostsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Hosts --- This data source can be used to get information on all 404 hosts.",
		Attributes: map[string]schema.Attribute{
			"dead_hosts": schema.SetNestedAttribute{
				MarkdownDescription: "The 404 hosts.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "The Id of the dead host.",
							Computed:            true,
						},
						"created_on": schema.StringAttribute{
							MarkdownDescription: "The date and time the dead host was created.",
							Computed:            true,
						},
						"modified_on": schema.StringAttribute{
							MarkdownDescription: "The date and time the dead host was last modified.",
							Computed:            true,
						},
						"owner_user_id": schema.Int64Attribute{
							MarkdownDescription: "The Id of the user that owns the dead host.",
							Computed:            true,
						},
						"domain_names": schema.SetAttribute{
							MarkdownDescription: "The domain names associated with the dead host.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"certificate_id": schema.Int64Attribute{
							MarkdownDescription: "The Id of the certificate used by the dead host.",
							Computed:            true,
						},
						"ssl_forced": schema.BoolAttribute{
							MarkdownDescription: "Whether SSL is forced for the dead host.",
							Computed:            true,
						},
						"hsts_enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether HSTS is enabled for the dead host.",
							Computed:            true,
						},
						"hsts_subdomains": schema.BoolAttribute{
							MarkdownDescription: "Whether HSTS is enabled for subdomains of the dead host.",
							Computed:            true,
						},
						"http2_support": schema.BoolAttribute{
							MarkdownDescription: "Whether HTTP/2 is supported for the dead host.",
							Computed:            true,
						},
						"advanced_config": schema.StringAttribute{
							MarkdownDescription: "The advanced configuration used by the dead host.",
							Computed:            true,
						},
						"enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether the dead host is enabled.",
							Computed:            true,
						},
						"meta": schema.MapAttribute{
							MarkdownDescription: "The meta data associated with the dead host.",
							ElementType:         types.StringType,
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *DeadHostsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if data := dataSourceConfigure(ctx, req, resp); data != nil {
		d.client = data.Client
		d.auth = data.Auth
	}
}

func (d *DeadHostsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *models.DeadHosts

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, _, err := d.client.Class404HostsAPI.GetDeadHosts(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read dead hosts, got error: %s", err))
		return
	}

	data.Write(ctx, &response, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
