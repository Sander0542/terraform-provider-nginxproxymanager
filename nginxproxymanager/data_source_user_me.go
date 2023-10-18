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
	_ common.IDataSource                 = &userMeDataSource{}
	_ datasource.DataSourceWithConfigure = &userMeDataSource{}
)

func NewUserMeDataSource() datasource.DataSource {
	b := &common.DataSource{Name: "user_me"}
	d := &userMeDataSource{b, nil}
	b.IDataSource = d
	return d
}

type userMeDataSource struct {
	*common.DataSource
	client *client.Client
}

func (d *userMeDataSource) SchemaImpl(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Users --- User me data source",
		Attributes:  attributes.UserMe,
	}
}

func (d *userMeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}

func (d *userMeDataSource) ReadImpl(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	user, err := d.client.GetMe(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading user", "Could not read user, unexpected error: "+err.Error())
		return
	}

	var data models.User
	resp.Diagnostics.Append(data.Load(ctx, user)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
