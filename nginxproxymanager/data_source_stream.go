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
	_ common.IDataSource                 = &streamDataSource{}
	_ datasource.DataSourceWithConfigure = &streamDataSource{}
)

func NewStreamDataSource() datasource.DataSource {
	b := &common.DataSource{Name: "stream"}
	d := &streamDataSource{b, nil}
	b.IDataSource = d
	return d
}

type streamDataSource struct {
	*common.DataSource
	client *client.Client
}

func (d *streamDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a stream by ID.",
		Attributes:  attributes.Stream,
	}
}

func (d *streamDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *streamDataSource) ReadImpl(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.Stream

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	stream, err := d.client.GetStream(ctx, data.ID.ValueInt64Pointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading stream",
			"Could not read stream, unexpected error: "+err.Error())
		return
	}
	if stream == nil {
		resp.Diagnostics.AddError(
			"Error reading stream",
			fmt.Sprintf("No stream found with ID: %d", data.ID))
		return
	}

	resp.Diagnostics.Append(data.Load(ctx, stream)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
