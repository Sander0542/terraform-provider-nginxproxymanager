// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"
	"github.com/sander0542/terraform-provider-nginxproxymanager/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
)

var _ datasource.DataSource = &ProxyHostDataSource{}

func NewProxyHostDataSource() datasource.DataSource {
	return &ProxyHostDataSource{}
}

type ProxyHostDataSource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (d *ProxyHostDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_proxy_host"
}

func (d *ProxyHostDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This data source can be used to get information about a specific proxy host.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the proxy host.",
				Required:            true,
			},
			"created_on": schema.StringAttribute{
				MarkdownDescription: "The date and time the proxy host was created.",
				Computed:            true,
			},
			"modified_on": schema.StringAttribute{
				MarkdownDescription: "The date and time the proxy host was last modified.",
				Computed:            true,
			},
			"owner_user_id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the user that owns the proxy host.",
				Computed:            true,
			},
			"domain_names": schema.ListAttribute{
				MarkdownDescription: "The domain names associated with the proxy host.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"forward_scheme": schema.StringAttribute{
				MarkdownDescription: "The scheme used to forward requests to the proxy host. Can be either `http` or `https`.",
				Computed:            true,
			},
			"forward_host": schema.StringAttribute{
				MarkdownDescription: "The host used to forward requests to the proxy host.",
				Computed:            true,
			},
			"forward_port": schema.Int64Attribute{
				MarkdownDescription: "The port used to forward requests to the proxy host.",
				Computed:            true,
			},
			"certificate_id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the certificate used by the proxy host.",
				Computed:            true,
			},
			"ssl_forced": schema.BoolAttribute{
				MarkdownDescription: "Whether SSL is forced for the proxy host.",
				Computed:            true,
			},
			"hsts_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether HSTS is enabled for the proxy host.",
				Computed:            true,
			},
			"hsts_subdomains": schema.BoolAttribute{
				MarkdownDescription: "Whether HSTS is enabled for subdomains of the proxy host.",
				Computed:            true,
			},
			"http2_support": schema.BoolAttribute{
				MarkdownDescription: "Whether HTTP/2 is supported for the proxy host.",
				Computed:            true,
			},
			"block_exploits": schema.BoolAttribute{
				MarkdownDescription: "Whether exploits are blocked for the proxy host.",
				Computed:            true,
			},
			"caching_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether caching is enabled for the proxy host.",
				Computed:            true,
			},
			"allow_websocket_upgrade": schema.BoolAttribute{
				MarkdownDescription: "Whether websocket upgrades are allowed for the proxy host.",
				Computed:            true,
			},
			"access_list_id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the access list used by the proxy host.",
				Computed:            true,
			},
			"advanced_config": schema.StringAttribute{
				MarkdownDescription: "The advanced configuration used by the proxy host.",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the proxy host is enabled.",
				Computed:            true,
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "The meta data associated with the proxy host.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"locations": schema.SetNestedAttribute{
				MarkdownDescription: "The locations associated with the proxy host.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"path": schema.StringAttribute{
							MarkdownDescription: "The path associated with the location.",
							Computed:            true,
						},
						"forward_scheme": schema.StringAttribute{
							MarkdownDescription: "The scheme used to forward requests to the location. Can be either `http` or `https`.",
							Computed:            true,
						},
						"forward_host": schema.StringAttribute{
							MarkdownDescription: "The host used to forward requests to the location.",
							Computed:            true,
						},
						"forward_port": schema.Int64Attribute{
							MarkdownDescription: "The port used to forward requests to the location.",
							Computed:            true,
						},
						"advanced_config": schema.StringAttribute{
							MarkdownDescription: "The advanced configuration used by the location.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *ProxyHostDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *ProxyHostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *models.ProxyHost

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, _, err := d.client.ProxyHostsAPI.GetProxyHost(d.auth, data.Id.ValueInt64()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read proxy host, got error: %s", err))
		return
	}

	data.Write(ctx, response, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
