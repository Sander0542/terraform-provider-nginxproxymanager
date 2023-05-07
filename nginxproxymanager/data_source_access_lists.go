package nginxproxymanager

import (
	"context"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/models"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client"
	attributes "github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/attributes/datasource"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/common"
)

var (
	_ common.IDataSource                 = &accessListsDataSource{}
	_ datasource.DataSourceWithConfigure = &accessListsDataSource{}
)

func NewAccessListsDataSource() datasource.DataSource {
	b := &common.DataSource{Name: "access_lists"}
	d := &accessListsDataSource{b, nil}
	b.IDataSource = d
	return d
}

type accessListsDataSource struct {
	*common.DataSource
	client *client.Client
}

func (d *accessListsDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Access Lists data source",
		Attributes:  attributes.AccessLists,
	}
}

func (d *accessListsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *accessListsDataSource) ReadImpl(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	accessLists, err := d.client.GetAccessLists(ctx)
	if err != nil {
		sentry.CaptureException(err)
		resp.Diagnostics.AddError("Error reading access lists", "Could not read access lists, unexpected error: "+err.Error())
		return
	}

	var data models.AccessLists
	data.AccessLists = make([]models.AccessList, len(*accessLists))
	for i, v := range *accessLists {
		resp.Diagnostics.Append(data.AccessLists[i].Load(ctx, &v)...)
	}

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
