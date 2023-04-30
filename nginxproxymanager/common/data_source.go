package common

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var (
	_ datasource.DataSource = &DataSource{}
)

type IDataSource interface {
	MetadataImpl(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse)
	SchemaImpl(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse)
	ReadImpl(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
}

type DataSource struct {
	IDataSource

	dataSourceName string
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.MetadataImpl(ctx, req, resp)

	d.dataSourceName = fmt.Sprintf("Data Source %s", resp.TypeName)
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	d.SchemaImpl(ctx, req, resp)
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	span := sentry.StartSpan(ctx, "terraform.data_source.read", sentry.TransactionName(d.dataSourceName))
	defer span.Finish()
	d.ReadImpl(ctx, req, resp)
}
