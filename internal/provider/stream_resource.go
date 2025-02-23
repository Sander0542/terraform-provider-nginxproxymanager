// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
	"github.com/sander0542/terraform-provider-nginxproxymanager/internal/provider/models"
	"strconv"
)

var _ resource.Resource = &StreamResource{}
var _ resource.ResourceWithImportState = &StreamResource{}

func NewStreamResource() resource.Resource {
	return &StreamResource{}
}

type StreamResource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (r *StreamResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stream"
}

func (r *StreamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Hosts --- This resource can be used to manage a stream.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the stream.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created_on": schema.StringAttribute{
				MarkdownDescription: "The date and time the stream was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified_on": schema.StringAttribute{
				MarkdownDescription: "The date and time the stream was last modified.",
				Computed:            true,
			},
			"owner_user_id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the user that owns the stream.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"incoming_port": schema.Int64Attribute{
				MarkdownDescription: "The incoming port of the stream.",
				Required:            true,
			},
			"forwarding_host": schema.StringAttribute{
				MarkdownDescription: "The forwarding host of the stream.",
				Required:            true,
			},
			"forwarding_port": schema.Int64Attribute{
				MarkdownDescription: "The forwarding port of the stream.",
				Required:            true,
			},
			"tcp_forwarding": schema.BoolAttribute{
				MarkdownDescription: "Whether TCP forwarding is enabled.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
			},
			"udp_forwarding": schema.BoolAttribute{
				MarkdownDescription: "Whether UDP forwarding is enabled.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"certificate_id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the certificate used by the stream.",
				Optional:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the stream is enabled.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
			},
			"meta": schema.MapAttribute{
				MarkdownDescription: "The meta data associated with the stream.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (r *StreamResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *StreamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *models.Stream

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	streamEnabled := data.Enabled.ValueBool()

	request := data.ToCreateRequest(ctx, &resp.Diagnostics)
	stream, _, err := r.client.StreamsAPI.CreateStream(r.auth).CreateStreamRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create stream, got error: %s", err))
		return
	}

	data.Write(ctx, stream, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), data.Id)...)

	err = r.toggleStream(stream.GetId(), stream.GetEnabled(), streamEnabled)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("enabled"), "Client Error", fmt.Sprintf("Unable to update stream, got err: %s", err))
		return
	}

	data.Enabled = types.BoolValue(streamEnabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StreamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *models.Stream

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	stream, _, err := r.client.StreamsAPI.GetStream(r.auth, data.Id.ValueInt64()).Execute()
	if err != nil {
		if err.Error() == "404 Not Found" {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read stream, got error: %s", err))
			return
		}
	}

	data.Write(ctx, stream, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StreamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *models.Stream

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	streamEnabled := data.Enabled.ValueBool()

	request := data.ToUpdateRequest(ctx, &resp.Diagnostics)
	stream, _, err := r.client.StreamsAPI.UpdateStream(r.auth, data.Id.ValueInt64()).UpdateStreamRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update stream, got error: %s", err))
		return
	}

	data.Write(ctx, stream, &resp.Diagnostics)

	err = r.toggleStream(stream.GetId(), stream.GetEnabled(), streamEnabled)
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("enabled"), "Client Error", fmt.Sprintf("Unable to update stream, got err: %s", err))
		return
	}

	data.Enabled = types.BoolValue(streamEnabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *StreamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *models.Stream

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	success, _, err := r.client.StreamsAPI.DeleteStream(r.auth, data.Id.ValueInt64()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete stream, got error: %s", err))
		return
	}

	if !success {
		resp.Diagnostics.AddError("Server Error", "Unable to delete stream.")
		return
	}
}

func (r *StreamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not convert id to number, got error: %s", err))
		return
	}

	stream, _, err := r.client.StreamsAPI.GetStream(r.auth, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read stream, got error: %s", err))
		return
	}

	diags := resp.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(stream.GetId()))
	resp.Diagnostics.Append(diags...)
}

func (r *StreamResource) toggleStream(streamId int64, current bool, desired bool) error {
	if desired && !current {
		enableResponse, _, err := r.client.StreamsAPI.EnableStream(r.auth, streamId).Execute()
		if err != nil {
			return err
		} else if !enableResponse {
			return errors.New("unable to enable stream")
		}
	} else if !desired && current {
		disableResponse, _, err := r.client.StreamsAPI.DisableStream(r.auth, streamId).Execute()
		if err != nil {
			return err
		} else if !disableResponse {
			return errors.New("unable to disable stream")
		}
	}

	return nil
}
