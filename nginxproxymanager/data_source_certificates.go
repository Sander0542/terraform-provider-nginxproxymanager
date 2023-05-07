package nginxproxymanager

import (
	"context"
	attributes "github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/attributes/datasource"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/common"
)

var (
	_ common.IDataSource                 = &certificatesDataSource{}
	_ datasource.DataSourceWithConfigure = &certificatesDataSource{}
)

func NewCertificatesDataSource() datasource.DataSource {
	b := &common.DataSource{Name: "certificates"}
	d := &certificatesDataSource{b, nil}
	b.IDataSource = d
	return d
}

type certificatesDataSource struct {
	*common.DataSource
	client *client.Client
}

func (d *certificatesDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Certificates data source",
		Attributes:  attributes.Certificates,
	}
}

func (d *certificatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *certificatesDataSource) ReadImpl(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	certificates, err := d.client.GetCertificates(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading certificate", "Could not read certificate, unexpected error: "+err.Error())
		return
	}

	var data models.Certificates
	resp.Diagnostics.Append(data.Load(ctx, certificates)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
