package common

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/sentry"
)

var (
	_ resource.Resource = &Resource{}
)

type IResource interface {
	SchemaImpl(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse)
	ReadImpl(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse)
	CreateImpl(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse)
	UpdateImpl(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse)
	DeleteImpl(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse)
}

type Resource struct {
	IResource

	Name string
}

func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, r.Name)
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.SchemaImpl(ctx, req, resp)
	sentry.CaptureDiagnostics(resp.Diagnostics)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	span := sentry.StartResource(ctx, "read", r.Name)
	defer span.Finish()
	r.ReadImpl(span.Context(), req, resp)
	sentry.CaptureDiagnostics(resp.Diagnostics)
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	span := sentry.StartResource(ctx, "create", r.Name)
	defer span.Finish()
	r.CreateImpl(span.Context(), req, resp)
	sentry.CaptureDiagnostics(resp.Diagnostics)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	span := sentry.StartResource(ctx, "update", r.Name)
	defer span.Finish()
	r.UpdateImpl(span.Context(), req, resp)
	sentry.CaptureDiagnostics(resp.Diagnostics)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	span := sentry.StartResource(ctx, "delete", r.Name)
	defer span.Finish()
	r.DeleteImpl(span.Context(), req, resp)
	sentry.CaptureDiagnostics(resp.Diagnostics)
}
