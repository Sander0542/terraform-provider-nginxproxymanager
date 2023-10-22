package nginxproxymanager

import (
	"context"
	"fmt"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client"
	attributes "github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/attributes/datasource"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/common"
)

var (
	_ common.IDataSource                 = &versionDataSource{}
	_ datasource.DataSourceWithConfigure = &versionDataSource{}
)

func NewVersionDataSource() datasource.DataSource {
	b := &common.DataSource{Name: "version"}
	d := &versionDataSource{b, nil}
	b.IDataSource = d
	return d
}

type versionDataSource struct {
	*common.DataSource
	client *client.Client
}

func (d *versionDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Meta --- Version data source",
		Attributes:  attributes.Version,
	}
}

func (d *versionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *versionDataSource) ReadImpl(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.Version

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	api, err := d.client.GetApi(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading version", "Could not read version, unexpected error: "+err.Error())
		return
	}
	if api == nil {
		resp.Diagnostics.AddError("Error reading version", fmt.Sprintf("No version found"))
		return
	}

	resp.Diagnostics.Append(data.Load(ctx, api)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
