package nginxproxymanager

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/common"
)

var (
	_ common.IDataSource                 = &certificateDataSource{}
	_ datasource.DataSourceWithConfigure = &certificateDataSource{}
)

func NewCertificateDataSource() datasource.DataSource {
	b := &common.DataSource{Name: "certificate"}
	d := &certificateDataSource{b, nil}
	b.IDataSource = d
	return d
}

type certificateDataSource struct {
	*common.DataSource
	client *client.Client
}

type certificateDataSourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	CreatedOn   types.String `tfsdk:"created_on"`
	ModifiedOn  types.String `tfsdk:"modified_on"`
	Provider    types.String `tfsdk:"provider_name"`
	NiceName    types.String `tfsdk:"nice_name"`
	DomainNames types.List   `tfsdk:"domain_names"`
	ExpiresOn   types.String `tfsdk:"expires_on"`
	Meta        types.Map    `tfsdk:"meta"`
}

func (d *certificateDataSource) MetadataImpl(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate"
}

func (d *certificateDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Certificate data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The ID of the certificate.",
				Required:    true,
			},
			"created_on": schema.StringAttribute{
				Description: "The date and time the certificate was created.",
				Computed:    true,
			},
			"modified_on": schema.StringAttribute{
				Description: "The date and time the certificate was last modified.",
				Computed:    true,
			},
			"provider_name": schema.StringAttribute{
				Description: "The provider of the certificate.",
				Computed:    true,
			},
			"nice_name": schema.StringAttribute{
				Description: "The nice name of the certificate.",
				Computed:    true,
			},
			"domain_names": schema.ListAttribute{
				Description: "The domain names associated with the certificate.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"expires_on": schema.StringAttribute{
				Description: "The date and time the certificate expires.",
				Computed:    true,
			},
			"meta": schema.MapAttribute{
				Description: "The meta data associated with the certificate.",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (d *certificateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *certificateDataSource) ReadImpl(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data certificateDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	certificate, err := d.client.GetCertificate(data.ID.ValueInt64Pointer())
	if err != nil {
		resp.Diagnostics.AddError("Error reading certificate", "Could not read certificate, unexpected error: "+err.Error())
		return
	}
	if certificate == nil {
		resp.Diagnostics.AddError(
			"Error reading certificate",
			fmt.Sprintf("No certificate found with ID: %d", data.ID))
		return
	}

	domainNames, diags := types.ListValueFrom(ctx, types.StringType, certificate.DomainNames)
	resp.Diagnostics.Append(diags...)
	meta, diags := types.MapValueFrom(ctx, types.StringType, certificate.Meta.Map())
	resp.Diagnostics.Append(diags...)

	data.ID = types.Int64Value(certificate.ID)
	data.CreatedOn = types.StringValue(certificate.CreatedOn)
	data.ModifiedOn = types.StringValue(certificate.ModifiedOn)
	data.Provider = types.StringValue(certificate.Provider)
	data.NiceName = types.StringValue(certificate.NiceName)
	data.DomainNames = domainNames
	data.ExpiresOn = types.StringValue(certificate.ExpiresOn)
	data.Meta = meta

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
