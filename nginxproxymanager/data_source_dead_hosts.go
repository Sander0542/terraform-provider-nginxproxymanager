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
	_ common.IDataSource                 = &deadHostsDataSource{}
	_ datasource.DataSourceWithConfigure = &deadHostsDataSource{}
)

func NewDeadHostsDataSource() datasource.DataSource {
	b := &common.DataSource{Name: "dead_hosts"}
	d := &deadHostsDataSource{b, nil}
	b.IDataSource = d
	return d
}

type deadHostsDataSource struct {
	*common.DataSource
	client *client.Client
}

func (d *deadHostsDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Proxy Hosts data source",
		Attributes:  attributes.DeadHosts,
	}
}

func (d *deadHostsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *deadHostsDataSource) ReadImpl(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	deadHosts, err := d.client.GetDeadHosts(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading dead hosts", "Could not read dead hosts, unexpected error: "+err.Error())
		return
	}

	var data models.DeadHosts
	resp.Diagnostics.Append(data.Load(ctx, deadHosts)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
