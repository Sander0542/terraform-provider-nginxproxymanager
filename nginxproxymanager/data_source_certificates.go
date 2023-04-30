package nginxproxymanager

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/common"
)

var (
	_ common.IDataSource                 = &certificatesDataSource{}
	_ datasource.DataSourceWithConfigure = &certificatesDataSource{}
)

func NewCertificatesDataSource() datasource.DataSource {
	b := &common.DataSource{}
	d := &certificatesDataSource{b, nil}
	b.IDataSource = d
	return d
}

type certificatesDataSource struct {
	*common.DataSource
	client *client.Client
}

type certificatesDataSourceModel struct {
	Certificates []certificateItem `tfsdk:"certificates"`
}

type certificateItem struct {
	ID          types.Int64  `tfsdk:"id"`
	CreatedOn   types.String `tfsdk:"created_on"`
	ModifiedOn  types.String `tfsdk:"modified_on"`
	Provider    types.String `tfsdk:"provider_name"`
	NiceName    types.String `tfsdk:"nice_name"`
	DomainNames types.List   `tfsdk:"domain_names"`
	ExpiresOn   types.String `tfsdk:"expires_on"`
	Meta        types.Map    `tfsdk:"meta"`
}

func (d *certificatesDataSource) MetadataImpl(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificates"
}

func (d *certificatesDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Certificates data source",
		Attributes: map[string]schema.Attribute{
			"certificates": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "The ID of the certificate.",
							Computed:    true,
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
				},
			},
		},
	}
}

func (d *certificatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *certificatesDataSource) ReadImpl(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data certificatesDataSourceModel

	certificates, err := d.client.GetCertificates()
	if err != nil {
		resp.Diagnostics.AddError("Error reading certificate", "Could not read certificate, unexpected error: "+err.Error())
		return
	}

	for _, value := range *certificates {
		domainNames, diags := types.ListValueFrom(ctx, types.StringType, value.DomainNames)
		resp.Diagnostics.Append(diags...)
		meta, diags := types.MapValueFrom(ctx, types.StringType, value.Meta.Map())
		resp.Diagnostics.Append(diags...)

		certificate := certificateItem{
			ID:          types.Int64Value(value.ID),
			CreatedOn:   types.StringValue(value.CreatedOn),
			ModifiedOn:  types.StringValue(value.ModifiedOn),
			Provider:    types.StringValue(value.Provider),
			NiceName:    types.StringValue(value.NiceName),
			DomainNames: domainNames,
			ExpiresOn:   types.StringValue(value.ExpiresOn),
			Meta:        meta,
		}
		data.Certificates = append(data.Certificates, certificate)
	}

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
