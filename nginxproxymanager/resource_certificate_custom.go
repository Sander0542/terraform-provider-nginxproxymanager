package nginxproxymanager

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client"
	"github.com/sander0542/terraform-provider-nginxproxymanager/client/inputs"
	attributes "github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/attributes/resources"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/common"
	"github.com/sander0542/terraform-provider-nginxproxymanager/nginxproxymanager/models"
	"strconv"
)

var (
	_ common.IResource                    = &certificateCustomResource{}
	_ resource.ResourceWithConfigure      = &certificateCustomResource{}
	_ resource.ResourceWithValidateConfig = &certificateCustomResource{}
	_ resource.ResourceWithImportState    = &certificateCustomResource{}
)

func NewCertificateCustomResource() resource.Resource {
	b := &common.Resource{Name: "certificate_custom"}
	r := &certificateCustomResource{b, nil}
	b.IResource = r
	return r
}

type certificateCustomResource struct {
	*common.Resource
	client *client.Client
}

func (r *certificateCustomResource) SchemaImpl(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a custom certificate.",
		Attributes:  attributes.CertificateCustom,
	}
}

func (r *certificateCustomResource) Configure(ctx context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Client)
}

func (r *certificateCustomResource) CreateImpl(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.Certificate
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var item inputs.CertificateCustom
	diags.Append(plan.Save(ctx, &item)...)

	certificate, err := r.client.CreateCertificateCustom(ctx, &item)
	if err != nil {
		resp.Diagnostics.AddError("Error creating certificate custom", "Could not create certificate custom, unexpected error: "+err.Error())
		return
	}

	resp.Diagnostics.Append(plan.Load(ctx, certificate)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *certificateCustomResource) ReadImpl(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *models.Certificate
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	certificate, err := r.client.GetCertificate(ctx, state.ID.ValueInt64Pointer())
	if err != nil {
		resp.Diagnostics.AddError("Error reading certificate custom", "Could not read certificate custom, unexpected error: "+err.Error())
		return
	}
	if certificate == nil {
		state = nil
	} else {
		resp.Diagnostics.Append(state.Load(ctx, certificate)...)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *certificateCustomResource) UpdateImpl(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *certificateCustomResource) DeleteImpl(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *models.Certificate
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCertificate(ctx, state.ID.ValueInt64Pointer())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting certificate custom", "Could not delete certificate custom, unexpected error: "+err.Error())
		return
	}
}

func (r *certificateCustomResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data models.Certificate

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Provider.ValueString() != "custom" {
		resp.Diagnostics.AddAttributeError(
			path.Root("provider"),
			"Invalid provider",
			"Only 'custom' is allowed as provider for a certificate custom.")
	}
}

func (r *certificateCustomResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Error importing certificate custom", "Could not import certificate custom, unexpected error: "+err.Error())
		return
	}

	diags := resp.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(id))
	resp.Diagnostics.Append(diags...)
}
