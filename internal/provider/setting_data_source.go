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

var _ datasource.DataSource = &SettingDataSource{}

func NewSettingDataSource() datasource.DataSource {
	return &SettingDataSource{}
}

type SettingDataSource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (d *SettingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_setting"
}

func (d *SettingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Settings --- This data source can be used to get information about a specific setting.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Id of the setting.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the setting.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the setting.",
				Computed:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The value of the setting.",
				Computed:            true,
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "The meta data associated with the value.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (d *SettingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *SettingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *models.Setting

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, _, err := d.client.SettingsAPI.GetSetting(d.auth, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read setting, got error: %s", err))
		return
	}

	data.Write(ctx, response, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
