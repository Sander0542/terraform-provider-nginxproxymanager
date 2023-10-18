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
	_ common.IDataSource                 = &streamsDataSource{}
	_ datasource.DataSourceWithConfigure = &streamsDataSource{}
)

func NewStreamsDataSource() datasource.DataSource {
	b := &common.DataSource{Name: "streams"}
	d := &streamsDataSource{b, nil}
	b.IDataSource = d
	return d
}

type streamsDataSource struct {
	*common.DataSource
	client *client.Client
}

func (d *streamsDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Hosts --- Stream data source.",
		Attributes:  attributes.Streams,
	}
}

func (d *streamsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *streamsDataSource) ReadImpl(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	streams, err := d.client.GetStreams(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading streams", "Could not read streams, unexpected error: "+err.Error())
		return
	}

	var data models.Streams
	resp.Diagnostics.Append(data.Load(ctx, streams)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
