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

var _ datasource.DataSource = &CertificatesDataSource{}

func NewCertificatesDataSource() datasource.DataSource {
	return &CertificatesDataSource{}
}

type CertificatesDataSource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (d *CertificatesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificates"
}

func (d *CertificatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "SSL Certificates --- This data source can be used to get information on all certificates.",
		Attributes: map[string]schema.Attribute{
			"certificates": schema.SetNestedAttribute{
				MarkdownDescription: "The certificates.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "The Id of the certificate.",
							Computed:            true,
						},
						"created_on": schema.StringAttribute{
							MarkdownDescription: "The date and time the certificate was created.",
							Computed:            true,
						},
						"modified_on": schema.StringAttribute{
							MarkdownDescription: "The date and time the certificate was last modified.",
							Computed:            true,
						},
						"owner_user_id": schema.Int64Attribute{
							MarkdownDescription: "The Id of the user that owns the certificate.",
							Computed:            true,
						},
						"provider_name": schema.StringAttribute{
							MarkdownDescription: "The provider of the certificate.",
							Computed:            true,
						},
						"nice_name": schema.StringAttribute{
							MarkdownDescription: "The nice name of the certificate.",
							Computed:            true,
						},
						"domain_names": schema.SetAttribute{
							MarkdownDescription: "The domain names associated with the certificate.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"expires_on": schema.StringAttribute{
							MarkdownDescription: "The date and time the certificate expires.",
							Computed:            true,
						},
						"meta": schema.MapAttribute{
							MarkdownDescription: "The meta data associated with the certificate.",
							ElementType:         types.StringType,
							Computed:            true,
							Sensitive:           true,
						},
					},
				},
			},
		},
	}
}

func (d *CertificatesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if data := dataSourceConfigure(ctx, req, resp); data != nil {
		d.client = data.Client
		d.auth = data.Auth
	}
}

func (d *CertificatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *models.Certificates

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, _, err := d.client.CertificatesAPI.GetCertificates(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read certificates, got error: %s", err))
		return
	}

	data.Write(ctx, &response, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
