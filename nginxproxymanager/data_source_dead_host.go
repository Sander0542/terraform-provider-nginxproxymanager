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
	_ common.IDataSource                 = &deadHostDataSource{}
	_ datasource.DataSourceWithConfigure = &deadHostDataSource{}
)

func NewDeadHostDataSource() datasource.DataSource {
	b := &common.DataSource{Name: "dead_host"}
	d := &deadHostDataSource{b, nil}
	b.IDataSource = d
	return d
}

type deadHostDataSource struct {
	*common.DataSource
	client *client.Client
}

func (d *deadHostDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Hosts --- Fetches a dead host by ID.",
		Attributes:  attributes.DeadHost,
	}
}

func (d *deadHostDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *deadHostDataSource) ReadImpl(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DeadHost

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	deadHost, err := d.client.GetDeadHost(ctx, data.ID.ValueInt64Pointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading dead host",
			"Could not read dead host, unexpected error: "+err.Error())
		return
	}
	if deadHost == nil {
		resp.Diagnostics.AddError(
			"Error reading dead host",
			fmt.Sprintf("No dead host found with ID: %d", data.ID))
		return
	}

	resp.Diagnostics.Append(data.Load(ctx, deadHost)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
