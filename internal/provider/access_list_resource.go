// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
	"github.com/sander0542/terraform-provider-nginxproxymanager/internal/provider/models"
	"strconv"
)

var _ resource.Resource = &AccessListResource{}
var _ resource.ResourceWithImportState = &AccessListResource{}

func NewAccessListResource() resource.Resource {
	return &AccessListResource{}
}

type AccessListResource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (r *AccessListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_list"
}

func (r *AccessListResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Access Lists --- This resource can be used to manage an access list.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The Id of the access list.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created_on": schema.StringAttribute{
				Description: "The date and time the access list was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified_on": schema.StringAttribute{
				Description: "The date and time the access list was last modified.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the access list.",
				Required:    true,
			},
			"owner_user_id": schema.Int64Attribute{
				Description: "The ID of the user that owns the access list.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"authorizations": schema.SetNestedAttribute{
				Description: "The authorization items of the access list.",
				Computed:    true,
				Optional:    true,
				Default:     setdefault.StaticValue(types.SetValueMust(models.AccessListAuthorizationResource{}.GetType(), []attr.Value{})),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"username": schema.StringAttribute{
							Description: "The username of the authorization item.",
							Required:    true,
						},
						"password": schema.StringAttribute{
							Description: "The password hint of the authorization item.",
							Required:    true,
							Sensitive:   true,
						},
					},
				},
			},
			"access": schema.SetNestedAttribute{
				Description: "The access items of the access list.",
				Computed:    true,
				Optional:    true,
				Default:     setdefault.StaticValue(types.SetValueMust(models.AccessListAccessResource{}.GetType(), []attr.Value{})),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"directive": schema.StringAttribute{
							Description: "The directive of the access item.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("allow", "deny"),
							},
						},
						"address": schema.StringAttribute{
							Description: "The address of the access item.",
							Required:    true,
						},
					},
				},
			},
			"pass_auth": schema.BoolAttribute{
				Description: "Whether or not to pass the authorization header to the upstream server.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"satisfy_any": schema.BoolAttribute{
				Description: "Whether or not to satisfy any of the authorization items.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"meta": schema.MapAttribute{
				Description: "The meta data associated with the access list.",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (r *AccessListResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *AccessListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *models.AccessListResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	request := data.ToCreateRequest(ctx, &resp.Diagnostics)
	accessList, _, err := r.client.AccessListsAPI.CreateAccessList(r.auth).CreateAccessListRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create access list, got error: %s", err))
		return
	}

	data.Write(ctx, accessList, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AccessListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *models.AccessListResource

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	accessList, _, err := r.client.AccessListsAPI.GetAccessList(r.auth, data.Id.ValueInt64()).Expand("clients,items").Execute()
	if err != nil {
		if err.Error() == "404 Not Found" {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read access list, got error: %s", err))
			return
		}
	}

	data.Write(ctx, accessList, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AccessListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *models.AccessListResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	request := data.ToUpdateRequest(ctx, &resp.Diagnostics)
	accessList, _, err := r.client.AccessListsAPI.UpdateAccessList(r.auth, data.Id.ValueInt64()).UpdateAccessListRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update access list, got error: %s", err))
		return
	}

	data.Write(ctx, accessList, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *AccessListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *models.AccessListResource

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	success, _, err := r.client.AccessListsAPI.DeleteAccessList(r.auth, data.Id.ValueInt64()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete access list, got error: %s", err))
		return
	}

	if !success {
		resp.Diagnostics.AddError("Server Error", "Unable to delete access list.")
		return
	}
}

func (r *AccessListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not convert id to number, got error: %s", err))
		return
	}

	accessList, _, err := r.client.AccessListsAPI.GetAccessList(r.auth, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read access list, got error: %s", err))
		return
	}

	diags := resp.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(accessList.GetId()))
	resp.Diagnostics.Append(diags...)
}
