// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
	"github.com/sander0542/terraform-provider-nginxproxymanager/internal/provider/models"
	"os"
	"strconv"
)

var _ resource.Resource = &CertificateCustomResource{}
var _ resource.ResourceWithImportState = &CertificateCustomResource{}

func NewCertificateCustomResource() resource.Resource {
	return &CertificateCustomResource{}
}

type CertificateCustomResource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (r *CertificateCustomResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate_custom"
}

func (r *CertificateCustomResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "SSL Certificates --- This resource can be used to manage a custom certificate.",
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
			"owner_user_id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the user that owns the proxy host.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the certificate.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"certificate": schema.StringAttribute{
				MarkdownDescription: "The contents of the certificate.",
				Required:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"certificate_key": schema.StringAttribute{
				MarkdownDescription: "The contents of the certificate key.",
				Required:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"domain_names": schema.ListAttribute{
				MarkdownDescription: "The domain names associated with the certificate.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"expires_on": schema.StringAttribute{
				MarkdownDescription: "The date and time the certificate expires.",
				Computed:            true,
			},
		},
	}
}

func (r *CertificateCustomResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *CertificateCustomResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *models.CertificateCustom

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.validateCertificate(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to validate certificate, got error: %s", err))
		return
	}

	certificateRequest := data.ToCreateRequest(ctx, &resp.Diagnostics)
	certificate, _, err := r.client.CertificatesAPI.CreateCertificate(r.auth).CreateCertificateRequest(*certificateRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create certificate, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), certificate.GetId())...)

	err = r.uploadCertificate(certificate.GetId(), data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to upload certificate, got error: %s", err))
		return
	}

	certificate, _, err = r.client.CertificatesAPI.GetCertificate(r.auth, certificate.GetId()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read certificate, got error: %s", err))
		return
	}

	data.Write(ctx, certificate, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CertificateCustomResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *models.CertificateCustom

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	certificate, _, err := r.client.CertificatesAPI.GetCertificate(r.auth, data.Id.ValueInt64()).Execute()
	if err != nil {
		if err.Error() == "404 Not Found" {
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

func (r *CertificateCustomResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Client Error", "This resource does not support updates.")
}

func (r *CertificateCustomResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *models.CertificateCustom

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

func (r *CertificateCustomResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not convert ID to number, got error: %s", err))
		return
	}

	certificate, _, err := r.client.CertificatesAPI.GetCertificate(r.auth, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read certificate, got error: %s", err))
		return
	}
	if certificate.GetProvider() != "other" {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("The certificate is not a custom certificate, got provider: %s", certificate.GetProvider()))
	}

	diags := resp.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(certificate.GetId()))
	resp.Diagnostics.Append(diags...)
}

func (r *CertificateCustomResource) validateCertificate(data *models.CertificateCustom) error {
	certFile, err := os.CreateTemp("", "certificate")
	if err != nil {
		return err
	}
	defer certFile.Close()
	certKeyFile, err := os.CreateTemp("", "certificate_key")
	if err != nil {
		return err
	}
	defer certKeyFile.Close()

	_, err = certFile.WriteString(data.Certificate.ValueString())
	if err != nil {
		return err
	}
	_, err = certFile.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = certKeyFile.WriteString(data.CertificateKey.ValueString())
	if err != nil {
		return err
	}
	_, err = certKeyFile.Seek(0, 0)
	if err != nil {
		return err
	}

	_, _, err = r.client.CertificatesAPI.ValidateCertificates(r.auth).Certificate(certFile).CertificateKey(certKeyFile).Execute()

	return err
}

func (r *CertificateCustomResource) uploadCertificate(certId int64, data *models.CertificateCustom) error {
	certFile, err := os.CreateTemp("", "certificate")
	if err != nil {
		return err
	}
	defer certFile.Close()
	certKeyFile, err := os.CreateTemp("", "certificate_key")
	if err != nil {
		return err
	}
	defer certKeyFile.Close()

	_, err = certFile.WriteString(data.Certificate.ValueString())
	if err != nil {
		return err
	}
	_, err = certFile.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = certKeyFile.WriteString(data.CertificateKey.ValueString())
	if err != nil {
		return err
	}
	_, err = certKeyFile.Seek(0, 0)
	if err != nil {
		return err
	}

	_, _, err = r.client.CertificatesAPI.UploadCertificate(r.auth, certId).Certificate(certFile).CertificateKey(certKeyFile).Execute()

	return err
}
