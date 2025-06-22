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

var _ datasource.DataSource = &AccessListDataSource{}

func NewAccessListDataSource() datasource.DataSource {
	return &AccessListDataSource{}
}

type AccessListDataSource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (d *AccessListDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_list"
}

func (d *AccessListDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Access Lists --- This data source can be used to get information about a specific access list.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the access list.",
				Required:            true,
			},
			"created_on": schema.StringAttribute{
				MarkdownDescription: "The date and time the access list was created.",
				Computed:            true,
			},
			"modified_on": schema.StringAttribute{
				MarkdownDescription: "The date and time the access list was last modified.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the access list.",
				Computed:            true,
			},
			"owner_user_id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the user that owns the access list.",
				Computed:            true,
			},
			"authorizations": schema.SetNestedAttribute{
				MarkdownDescription: "The authorization items of the access list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "The Id of the authorization item.",
							Computed:    true,
						},
						"created_on": schema.StringAttribute{
							Description: "The date and time the authorization item was created.",
							Computed:    true,
						},
						"modified_on": schema.StringAttribute{
							Description: "The date and time the authorization item was last modified.",
							Computed:    true,
						},
						"username": schema.StringAttribute{
							Description: "The username of the authorization item.",
							Computed:    true,
						},
						"password_hint": schema.StringAttribute{
							Description: "The password hint of the authorization item.",
							Computed:    true,
							Sensitive:   true,
						},
						"meta": schema.MapAttribute{
							Description: "The meta data associated with the authorization item.",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
			"access": schema.SetNestedAttribute{
				MarkdownDescription: "The access items of the access list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "The Id of the access item.",
							Computed:            true,
						},
						"created_on": schema.StringAttribute{
							MarkdownDescription: "The date and time the access item was created.",
							Computed:            true,
						},
						"modified_on": schema.StringAttribute{
							MarkdownDescription: "The date and time the access item was last modified.",
							Computed:            true,
						},
						"address": schema.StringAttribute{
							MarkdownDescription: "The address of the access item.",
							Computed:            true,
						},
						"directive": schema.StringAttribute{
							MarkdownDescription: "The directive of the access item.",
							Computed:            true,
						},
						"meta": schema.MapAttribute{
							MarkdownDescription: "The meta data associated with the authorization item.",
							ElementType:         types.StringType,
							Computed:            true,
						},
					},
				},
			},
			"pass_auth": schema.BoolAttribute{
				MarkdownDescription: "Whether or not to pass the authorization header to the upstream server.",
				Computed:            true,
			},
			"satisfy_any": schema.BoolAttribute{
				MarkdownDescription: "Whether or not to satisfy any of the authorization items.",
				Computed:            true,
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "The meta data associated with the access list.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (d *AccessListDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if data := dataSourceConfigure(ctx, req, resp); data != nil {
		d.client = data.Client
		d.auth = data.Auth
	}
}

func (d *AccessListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *models.AccessList

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, _, err := d.client.AccessListsAPI.GetAccessList(d.auth, data.Id.ValueInt64()).Expand("clients,items").Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read access list, got error: %s", err))
		return
	}

	data.Write(ctx, (*nginxproxymanager.GetAccessLists200ResponseInner)(response), &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
