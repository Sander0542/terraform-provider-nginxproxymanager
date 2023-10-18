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
	_ common.IDataSource                 = &proxyHostDataSource{}
	_ datasource.DataSourceWithConfigure = &proxyHostDataSource{}
)

func NewProxyHostDataSource() datasource.DataSource {
	b := &common.DataSource{Name: "proxy_host"}
	d := &proxyHostDataSource{b, nil}
	b.IDataSource = d
	return d
}

type proxyHostDataSource struct {
	*common.DataSource
	client *client.Client
}

func (d *proxyHostDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Hosts --- Fetches a proxy host by ID.",
		Attributes:  attributes.ProxyHost,
	}
}

func (d *proxyHostDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *proxyHostDataSource) ReadImpl(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.ProxyHost

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	proxyHost, err := d.client.GetProxyHost(ctx, data.ID.ValueInt64Pointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading proxy host",
			"Could not read proxy host, unexpected error: "+err.Error())
		return
	}
	if proxyHost == nil {
		resp.Diagnostics.AddError(
			"Error reading proxy host",
			fmt.Sprintf("No proxy host found with ID: %d", data.ID))
		return
	}

	resp.Diagnostics.Append(data.Load(ctx, proxyHost)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
