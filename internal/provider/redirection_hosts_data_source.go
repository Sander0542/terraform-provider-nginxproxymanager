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

var _ datasource.DataSource = &RedirectionHostsDataSource{}

func NewRedirectionHostsDataSource() datasource.DataSource {
	return &RedirectionHostsDataSource{}
}

type RedirectionHostsDataSource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (d *RedirectionHostsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_redirection_hosts"
}

func (d *RedirectionHostsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This data source can be used to get information on all redirection hosts.",
		Attributes: map[string]schema.Attribute{
			"redirection_hosts": schema.SetNestedAttribute{
				MarkdownDescription: "The redirection hosts.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "The Id of the redirection host.",
							Computed:            true,
						},
						"created_on": schema.StringAttribute{
							MarkdownDescription: "The date and time the redirection host was created.",
							Computed:            true,
						},
						"modified_on": schema.StringAttribute{
							MarkdownDescription: "The date and time the redirection host was last modified.",
							Computed:            true,
						},
						"owner_user_id": schema.Int64Attribute{
							MarkdownDescription: "The Id of the user that owns the redirection host.",
							Computed:            true,
						},
						"domain_names": schema.ListAttribute{
							MarkdownDescription: "The domain names associated with the redirection host.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"forward_scheme": schema.StringAttribute{
							MarkdownDescription: "The scheme used to forward requests to the redirection host. Can be either `auto`, `http` or `https`.",
							Computed:            true,
						},
						"forward_domain_name": schema.StringAttribute{
							MarkdownDescription: "The domain name used to forward requests to the redirection host.",
							Computed:            true,
						},
						"forward_http_code": schema.Int64Attribute{
							MarkdownDescription: "The HTTP code used to forward requests to the redirection host.",
							Computed:            true,
						},
						"certificate_id": schema.Int64Attribute{
							MarkdownDescription: "The Id of the certificate used by the redirection host.",
							Computed:            true,
						},
						"ssl_forced": schema.BoolAttribute{
							MarkdownDescription: "Whether SSL is forced for the redirection host.",
							Computed:            true,
						},
						"hsts_enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether HSTS is enabled for the redirection host.",
							Computed:            true,
						},
						"hsts_subdomains": schema.BoolAttribute{
							MarkdownDescription: "Whether HSTS is enabled for subdomains of the redirection host.",
							Computed:            true,
						},
						"http2_support": schema.BoolAttribute{
							MarkdownDescription: "Whether HTTP/2 is supported for the redirection host.",
							Computed:            true,
						},
						"preserve_path": schema.BoolAttribute{
							MarkdownDescription: "Whether the path is preserved for the redirection host.",
							Computed:            true,
						},
						"block_exploits": schema.BoolAttribute{
							MarkdownDescription: "Whether exploits are blocked for the redirection host.",
							Computed:            true,
						},
						"advanced_config": schema.StringAttribute{
							MarkdownDescription: "The advanced configuration used by the redirection host.",
							Computed:            true,
						},
						"enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether the redirection host is enabled.",
							Computed:            true,
						},
						"meta": schema.MapAttribute{
							MarkdownDescription: "The meta data associated with the redirection host.",
							ElementType:         types.StringType,
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *RedirectionHostsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *RedirectionHostsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *models.RedirectionHosts

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, _, err := d.client.RedirectionHostsAPI.GetRedirectionHosts(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read redirection hosts, got error: %s", err))
		return
	}

	data.Write(ctx, &response, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
