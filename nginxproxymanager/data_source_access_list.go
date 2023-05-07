package nginxproxymanager

import (
	"context"
	"fmt"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/models"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client"
	attributes "github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/attributes/datasource"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/common"
)

var (
	_ common.IDataSource                 = &accessListDataSource{}
	_ datasource.DataSourceWithConfigure = &accessListDataSource{}
)

func NewAccessListDataSource() datasource.DataSource {
	b := &common.DataSource{Name: "access_list"}
	d := &accessListDataSource{b, nil}
	b.IDataSource = d
	return d
}

type accessListDataSource struct {
	*common.DataSource
	client *client.Client
}

func (d *accessListDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Access List data source",
		Attributes:  attributes.AccessList,
	}
}

func (d *accessListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *accessListDataSource) ReadImpl(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.AccessList

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	accessList, err := d.client.GetAccessList(ctx, data.ID.ValueInt64Pointer())
	if err != nil {
		sentry.CaptureException(err)
		resp.Diagnostics.AddError("Error reading access list", "Could not read access list, unexpected error: "+err.Error())
		return
	}
	if accessList == nil {
		resp.Diagnostics.AddError(
			"Error reading access list",
			fmt.Sprintf("No access list found with ID: %d", data.ID))
		return
	}

	resp.Diagnostics.Append(data.Load(ctx, accessList)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
