package common

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource = &Resource{}
)

type IResource interface {
	MetadataImpl(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse)
	SchemaImpl(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse)
	ReadImpl(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse)
	CreateImpl(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse)
	UpdateImpl(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse)
	DeleteImpl(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse)
}

type Resource struct {
	IResource

	resourceName string
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.MetadataImpl(ctx, req, resp)

	r.resourceName = fmt.Sprintf("Resource %s", resp.TypeName)
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.SchemaImpl(ctx, req, resp)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	span := sentry.StartSpan(ctx, "terraform.resource.read", sentry.TransactionName(r.resourceName))
	defer span.Finish()
	r.ReadImpl(ctx, req, resp)
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	span := sentry.StartSpan(ctx, "terraform.resource.create", sentry.TransactionName(r.resourceName))
	defer span.Finish()
	r.CreateImpl(ctx, req, resp)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	span := sentry.StartSpan(ctx, "terraform.resource.update", sentry.TransactionName(r.resourceName))
	defer span.Finish()
	r.UpdateImpl(ctx, req, resp)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	span := sentry.StartSpan(ctx, "terraform.resource.delete", sentry.TransactionName(r.resourceName))
	defer span.Finish()
	r.DeleteImpl(ctx, req, resp)
}
