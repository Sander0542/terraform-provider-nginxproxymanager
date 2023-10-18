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
	_ common.IDataSource                 = &redirectionHostsDataSource{}
	_ datasource.DataSourceWithConfigure = &redirectionHostsDataSource{}
)

func NewRedirectionHostsDataSource() datasource.DataSource {
	b := &common.DataSource{Name: "redirection_hosts"}
	d := &redirectionHostsDataSource{b, nil}
	b.IDataSource = d
	return d
}

type redirectionHostsDataSource struct {
	*common.DataSource
	client *client.Client
}

func (d *redirectionHostsDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Hosts --- Redirection Hosts data source",
		Attributes:  attributes.RedirectionHosts,
	}
}

func (d *redirectionHostsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *redirectionHostsDataSource) ReadImpl(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	redirectionHosts, err := d.client.GetRedirectionHosts(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading redirection hosts", "Could not read redirection hosts, unexpected error: "+err.Error())
		return
	}

	var data models.RedirectionHosts
	resp.Diagnostics.Append(data.Load(ctx, redirectionHosts)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
