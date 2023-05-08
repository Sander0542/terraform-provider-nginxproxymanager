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
	_ common.IDataSource                 = &redirectionHostDataSource{}
	_ datasource.DataSourceWithConfigure = &redirectionHostDataSource{}
)

func NewRedirectionHostDataSource() datasource.DataSource {
	b := &common.DataSource{Name: "redirection_host"}
	d := &redirectionHostDataSource{b, nil}
	b.IDataSource = d
	return d
}

type redirectionHostDataSource struct {
	*common.DataSource
	client *client.Client
}

func (d *redirectionHostDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a redirection host by ID.",
		Attributes:  attributes.RedirectionHost,
	}
}

func (d *redirectionHostDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *redirectionHostDataSource) ReadImpl(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.RedirectionHost

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	redirectionHost, err := d.client.GetRedirectionHost(ctx, data.ID.ValueInt64Pointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading redirection host",
			"Could not read redirection host, unexpected error: "+err.Error())
		return
	}
	if redirectionHost == nil {
		resp.Diagnostics.AddError(
			"Error reading redirection host",
			fmt.Sprintf("No redirection host found with ID: %d", data.ID))
		return
	}

	resp.Diagnostics.Append(data.Load(ctx, redirectionHost)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
