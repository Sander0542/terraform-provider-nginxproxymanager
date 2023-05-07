package common

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/sentry"
)

var (
	_ datasource.DataSource = &DataSource{}
)

type IDataSource interface {
	SchemaImpl(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	ReadImpl(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

type DataSource struct {
	IDataSource

	Name string
}

func (d *DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, d.Name)
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	d.SchemaImpl(ctx, req, resp)
	sentry.CaptureDiagnostics(resp.Diagnostics)
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	span := sentry.StartDataSource(ctx, "read", d.Name)
	defer span.Finish()
	d.ReadImpl(span.Context(), req, resp)
	sentry.CaptureDiagnostics(resp.Diagnostics)
}
