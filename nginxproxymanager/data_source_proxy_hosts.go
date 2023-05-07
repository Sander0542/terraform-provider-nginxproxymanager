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
	_ common.IDataSource                 = &proxyHostsDataSource{}
	_ datasource.DataSourceWithConfigure = &proxyHostsDataSource{}
)

func NewProxyHostsDataSource() datasource.DataSource {
	b := &common.DataSource{Name: "proxy_hosts"}
	d := &proxyHostsDataSource{b, nil}
	b.IDataSource = d
	return d
}

type proxyHostsDataSource struct {
	*common.DataSource
	client *client.Client
}

func (d *proxyHostsDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Proxy Hosts data source",
		Attributes:  attributes.ProxyHosts,
	}
}

func (d *proxyHostsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *proxyHostsDataSource) ReadImpl(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	proxyHosts, err := d.client.GetProxyHosts(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading proxy hosts", "Could not read proxy hosts, unexpected error: "+err.Error())
		return
	}

	var data models.ProxyHosts
	resp.Diagnostics.Append(data.Load(ctx, proxyHosts)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
