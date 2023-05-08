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
	_ common.IDataSource                 = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

func NewUserDataSource() datasource.DataSource {
	b := &common.DataSource{Name: "user"}
	d := &userDataSource{b, nil}
	b.IDataSource = d
	return d
}

type userDataSource struct {
	*common.DataSource
	client *client.Client
}

func (d *userDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "User data source",
		Attributes:  attributes.User,
	}
}

func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *userDataSource) ReadImpl(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.User

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	user, err := d.client.GetUser(ctx, data.ID.ValueInt64Pointer())
	if err != nil {
		resp.Diagnostics.AddError("Error reading user", "Could not read user, unexpected error: "+err.Error())
		return
	}
	if user == nil {
		resp.Diagnostics.AddError(
			"Error reading user",
			fmt.Sprintf("No user found with ID: %d", data.ID))
		return
	}

	resp.Diagnostics.Append(data.Load(ctx, user)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
