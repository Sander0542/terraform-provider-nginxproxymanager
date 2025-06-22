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

var _ datasource.DataSource = &CertificateDataSource{}

func NewCertificateDataSource() datasource.DataSource {
	return &CertificateDataSource{}
}

type CertificateDataSource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (d *CertificateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate"
}

func (d *CertificateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "SSL Certificates --- This data source can be used to get information about a specific certificate.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the certificate.",
				Required:            true,
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
	}
}

func (d *CertificateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if data := dataSourceConfigure(ctx, req, resp); data != nil {
		d.client = data.Client
		d.auth = data.Auth
	}
}

func (d *CertificateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *models.Certificate

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, _, err := d.client.CertificatesAPI.GetCertificate(d.auth, data.Id.ValueInt64()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read certificate, got error: %s", err))
		return
	}

	data.Write(ctx, response, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
