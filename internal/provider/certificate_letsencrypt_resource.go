// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
	"github.com/sander0542/terraform-provider-nginxproxymanager/internal/provider/models"
	"strconv"
)

var _ resource.Resource = &CertificateLetsencryptResource{}
var _ resource.ResourceWithImportState = &CertificateLetsencryptResource{}

func NewCertificateLetsencryptResource() resource.Resource {
	return &CertificateLetsencryptResource{}
}

type CertificateLetsencryptResource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (r *CertificateLetsencryptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate_letsencrypt"
}

func (r *CertificateLetsencryptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "SSL Certificates --- This resource can be used to manage a Let's Encrypt certificate.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the certificate.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created_on": schema.StringAttribute{
				MarkdownDescription: "The date and time the certificate was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified_on": schema.StringAttribute{
				MarkdownDescription: "The date and time the certificate was last modified.",
				Computed:            true,
			},
			"expires_on": schema.StringAttribute{
				MarkdownDescription: "The date and time the certificate expires.",
				Computed:            true,
			},
			"owner_user_id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the user that owns the certificate.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"domain_names": schema.SetAttribute{
				MarkdownDescription: "The domain names associated with the certificate.",
				Required:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"letsencrypt_email": schema.StringAttribute{
				MarkdownDescription: "The email address to use for the Let's Encrypt certificate.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"letsencrypt_agree": schema.BoolAttribute{
				MarkdownDescription: "Whether you agree to the [Let's Encrypt Terms of Service](https://letsencrypt.org/repository/).",
				Required:            true,
				Validators: []validator.Bool{
					boolvalidator.Equals(true),
				},
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"dns_challenge": schema.BoolAttribute{
				MarkdownDescription: "Whether to use DNS validation to request the Let's Encrypt certificate.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"dns_provider": schema.StringAttribute{
				MarkdownDescription: "The DNS provider to use for the DNS challenge.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dns_provider_credentials": schema.StringAttribute{
				MarkdownDescription: "The credentials to use for the provider in the DNS challenge.",
				Computed:            true,
				Optional:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"propagation_seconds": schema.Int64Attribute{
				MarkdownDescription: "The number of seconds to wait for DNS to propagate before asking the ACME server to verify the DNS record.",
				Computed:            true,
				Optional:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *CertificateLetsencryptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if data := resourceConfigure(ctx, req, resp); data != nil {
		r.client = data.Client
		r.auth = data.Auth
	}
}

func (r *CertificateLetsencryptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *models.CertificateLetsencrypt

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	certificateRequest := data.ToCreateRequest(ctx, &resp.Diagnostics)
	certificate, _, err := r.client.CertificatesAPI.CreateCertificate(r.auth).CreateCertificateRequest(*certificateRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create certificate, got error: %s", err))
		return
	}

	data.Write(ctx, certificate, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CertificateLetsencryptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *models.CertificateLetsencrypt

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	certificate, response, err := r.client.CertificatesAPI.GetCertificate(r.auth, data.Id.ValueInt64()).Execute()
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read certificate, got error: %s", err))
			return
		}
	}

	data.Write(ctx, certificate, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CertificateLetsencryptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Client Error", "This resource does not support updates.")
}

func (r *CertificateLetsencryptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *models.CertificateLetsencrypt

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, _, err := r.client.CertificatesAPI.DeleteCertificate(r.auth, data.Id.ValueInt64()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete certificate, got error: %s", err))
		return
	}

	if !response {
		resp.Diagnostics.AddError("Server Error", "Unable to delete certificate.")
		return
	}
}

func (r *CertificateLetsencryptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not convert id to number, got error: %s", err))
		return
	}

	certificate, _, err := r.client.CertificatesAPI.GetCertificate(r.auth, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read certificate, got error: %s", err))
		return
	}
	if certificate.GetProvider() != "letsencrypt" {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("The certificate is not a letsencrypt certificate, got provider: %s", certificate.GetProvider()))
	}

	diags := resp.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(certificate.GetId()))
	resp.Diagnostics.Append(diags...)
}
