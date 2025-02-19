// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"
	"github.com/sander0542/terraform-provider-nginxproxymanager/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/sander0542/nginxproxymanager-go"
)

var _ datasource.DataSource = &VersionDataSource{}

func NewVersionDataSource() datasource.DataSource {
	return &VersionDataSource{}
}

type VersionDataSource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (d *VersionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_version"
}

func (d *VersionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This data source can be used to get the current version of nginx proxy manager.",
		Attributes: map[string]schema.Attribute{
			"major": schema.Int64Attribute{
				MarkdownDescription: "The major version number.",
				Computed:            true,
			},
			"minor": schema.Int64Attribute{
				MarkdownDescription: "The minor version number.",
				Computed:            true,
			},
			"revision": schema.Int64Attribute{
				MarkdownDescription: "The revision version number.",
				Computed:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "The full version.",
				Computed:            true,
			},
		},
	}
}

func (d *VersionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *VersionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *models.Version

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get health information
	response, _, err := d.client.PublicAPI.Health(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read version, got error: %s", err))
		return
	}

	data.Write(ctx, &response.Version, &resp.Diagnostics)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
