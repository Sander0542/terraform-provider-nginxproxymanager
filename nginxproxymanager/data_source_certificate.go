package nginxproxymanager

import (
	"context"
	"fmt"
	attributes "github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/attributes/datasource"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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

func (d *certificateDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "SSL Certificates --- Certificate data source",
		Attributes:  attributes.Certificate,
	}
}

func (d *certificateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *certificateDataSource) ReadImpl(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.Certificate

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	certificate, err := d.client.GetCertificate(ctx, data.ID.ValueInt64Pointer())
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

	resp.Diagnostics.Append(data.Load(ctx, certificate)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
