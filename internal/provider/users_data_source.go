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

var _ datasource.DataSource = &UsersDataSource{}

func NewUsersDataSource() datasource.DataSource {
	return &UsersDataSource{}
}

type UsersDataSource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (d *UsersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *UsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Users --- This data source can be used to get information on all users.",
		Attributes: map[string]schema.Attribute{
			"users": schema.SetNestedAttribute{
				MarkdownDescription: "The users.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "The Id of the user.",
							Computed:            true,
						},
						"created_on": schema.StringAttribute{
							MarkdownDescription: "The date and time the user was created.",
							Computed:            true,
						},
						"modified_on": schema.StringAttribute{
							MarkdownDescription: "The date and time the user was last modified.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the user.",
							Computed:            true,
						},
						"nickname": schema.StringAttribute{
							MarkdownDescription: "The nickname of the user.",
							Computed:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "The email of the user.",
							Computed:            true,
						},
						"avatar": schema.StringAttribute{
							MarkdownDescription: "The avatar of the user.",
							Computed:            true,
						},
						"is_disabled": schema.BoolAttribute{
							MarkdownDescription: "Whether the user is disabled.",
							Computed:            true,
						},
						"roles": schema.SetAttribute{
							MarkdownDescription: "The roles of the user.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"permissions": schema.SingleNestedAttribute{
							MarkdownDescription: "The permissions of the user.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"access_lists": schema.StringAttribute{
									MarkdownDescription: "The permission value for the access lists.",
									Computed:            true,
								},
								"certificates": schema.StringAttribute{
									MarkdownDescription: "The permission value for the certificates.",
									Computed:            true,
								},
								"dead_hosts": schema.StringAttribute{
									MarkdownDescription: "The permission value for the dead hosts.",
									Computed:            true,
								},
								"proxy_hosts": schema.StringAttribute{
									MarkdownDescription: "The permission value for the proxy hosts.",
									Computed:            true,
								},
								"redirection_hosts": schema.StringAttribute{
									MarkdownDescription: "The permission value for the redirection hosts.",
									Computed:            true,
								},
								"streams": schema.StringAttribute{
									MarkdownDescription: "The permission value for the streams.",
									Computed:            true,
								},
								"visibility": schema.StringAttribute{
									MarkdownDescription: "The level of visibility for the user.",
									Computed:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *UsersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *UsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *models.Users

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, _, err := d.client.UsersAPI.GetUsers(d.auth).Expand("permissions").Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read users, got error: %s", err))
		return
	}

	data.Write(ctx, &response, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
