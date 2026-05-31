// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sander0542/nginxproxymanager-go"
	"github.com/sander0542/terraform-provider-nginxproxymanager/internal/provider/models"
	"strconv"
)

var _ resource.Resource = &UserResource{}
var _ resource.ResourceWithImportState = &UserResource{}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

type UserResource struct {
	client *nginxproxymanager.APIClient
	auth   context.Context
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Users --- This resource can be used to manage a user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "The Id of the user.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created_on": schema.StringAttribute{
				MarkdownDescription: "The date and time the user was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified_on": schema.StringAttribute{
				MarkdownDescription: "The date and time the user was last modified.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the user.",
				Required:            true,
			},
			"nickname": schema.StringAttribute{
				MarkdownDescription: "The nickname of the user.",
				Required:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email of the user.",
				Required:            true,
			},
			"avatar": schema.StringAttribute{
				MarkdownDescription: "The avatar of the user.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_disabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the user is disabled.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_admin": schema.BoolAttribute{
				MarkdownDescription: "Whether the user is an administrator.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"permissions": schema.SingleNestedAttribute{
				MarkdownDescription: "The permissions of the user.",
				Computed:            true,
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"access_lists": schema.StringAttribute{
						MarkdownDescription: "The permission value for the access lists. Must be one of `manage`, `view` or `hidden`.",
						Computed:            true,
						Optional:            true,
						Default:             stringdefault.StaticString("manage"),
						Validators: []validator.String{
							stringvalidator.OneOf("manage", "view", "hidden"),
						},
					},
					"certificates": schema.StringAttribute{
						MarkdownDescription: "The permission value for the certificates. Must be one of `manage`, `view` or `hidden`.",
						Computed:            true,
						Optional:            true,
						Default:             stringdefault.StaticString("manage"),
						Validators: []validator.String{
							stringvalidator.OneOf("manage", "view", "hidden"),
						},
					},
					"dead_hosts": schema.StringAttribute{
						MarkdownDescription: "The permission value for the dead hosts. Must be one of `manage`, `view` or `hidden`.",
						Computed:            true,
						Optional:            true,
						Default:             stringdefault.StaticString("manage"),
						Validators: []validator.String{
							stringvalidator.OneOf("manage", "view", "hidden"),
						},
					},
					"proxy_hosts": schema.StringAttribute{
						MarkdownDescription: "The permission value for the proxy hosts. Must be one of `manage`, `view` or `hidden`.",
						Computed:            true,
						Optional:            true,
						Default:             stringdefault.StaticString("manage"),
						Validators: []validator.String{
							stringvalidator.OneOf("manage", "view", "hidden"),
						},
					},
					"redirection_hosts": schema.StringAttribute{
						MarkdownDescription: "The permission value for the redirection hosts. Must be one of `manage`, `view` or `hidden`.",
						Computed:            true,
						Optional:            true,
						Default:             stringdefault.StaticString("manage"),
						Validators: []validator.String{
							stringvalidator.OneOf("manage", "view", "hidden"),
						},
					},
					"streams": schema.StringAttribute{
						MarkdownDescription: "The permission value for the streams. Must be one of `manage`, `view` or `hidden`.",
						Computed:            true,
						Optional:            true,
						Default:             stringdefault.StaticString("manage"),
						Validators: []validator.String{
							stringvalidator.OneOf("manage", "view", "hidden"),
						},
					},
					"visibility": schema.StringAttribute{
						MarkdownDescription: "The level of visibility for the user. Must be one of `user` or `all`.",
						Computed:            true,
						Optional:            true,
						Default:             stringdefault.StaticString("user"),
						Validators: []validator.String{
							stringvalidator.OneOf("user", "all"),
						},
					},
				},
			},
		},
	}
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *models.UserResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	request := data.ToCreateRequest(ctx, &resp.Diagnostics)
	user, _, err := r.client.UsersAPI.CreateUser(r.auth).CreateUserRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create user, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), user.GetId())...)

	if !data.Permissions.IsUnknown() {
		permissions, tmpDiags := models.UserPermissionsAs(ctx, data.Permissions)
		resp.Diagnostics.Append(tmpDiags...)

		err = r.updatePermissions(user.GetId(), *permissions.ToRequest(ctx, &resp.Diagnostics))
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("permissions"), "Client Error", fmt.Sprintf("Unable to update user permissions, got err: %s", err))
		}
	}

	userIdVal := user.GetId()
	userId := nginxproxymanager.Int64AsGetUserUserIDParameter(&userIdVal)
	user, _, err = r.client.UsersAPI.GetUser(r.auth, userId).Expand("permissions").Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user, got error: %s", err))
		return
	}

	data.Write(ctx, user, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *models.UserResource

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	userId := nginxproxymanager.Int64AsGetUserUserIDParameter(data.Id.ValueInt64Pointer())
	user, _, err := r.client.UsersAPI.GetUser(r.auth, userId).Expand("permissions").Execute()
	if err != nil {
		if err.Error() == "404 Not Found" {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user, got error: %s", err))
			return
		}
	}

	data.Write(ctx, user, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *models.UserResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	request := data.ToUpdateRequest(ctx, &resp.Diagnostics)
	userId := nginxproxymanager.Int64AsGetUserUserIDParameter(data.Id.ValueInt64Pointer())
	user, _, err := r.client.UsersAPI.UpdateUser(r.auth, userId).UpdateUserRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update user, got error: %s", err))
		return
	}

	if !data.Permissions.IsUnknown() {
		permissions, tmpDiags := models.UserPermissionsAs(ctx, data.Permissions)
		resp.Diagnostics.Append(tmpDiags...)

		err = r.updatePermissions(user.GetId(), *permissions.ToRequest(ctx, &resp.Diagnostics))
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("permissions"), "Client Error", fmt.Sprintf("Unable to update user permissions, got err: %s", err))
		}
	}

	user, _, err = r.client.UsersAPI.GetUser(r.auth, userId).Expand("permissions").Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user, got error: %s", err))
		return
	}

	data.Write(ctx, user, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *models.UserResource

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	success, _, err := r.client.UsersAPI.DeleteUser(r.auth, data.Id.ValueInt64()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete user, got error: %s", err))
		return
	}

	if !success {
		resp.Diagnostics.AddError("Server Error", "Unable to delete user.")
		return
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Could not convert id to number, got error: %s", err))
		return
	}

	userId := nginxproxymanager.Int64AsGetUserUserIDParameter(&id)
	user, _, err := r.client.UsersAPI.GetUser(r.auth, userId).Expand("permissions").Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user, got error: %s", err))
		return
	}

	diags := resp.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(user.GetId()))
	resp.Diagnostics.Append(diags...)
}

func (r *UserResource) updatePermissions(userId int64, permissions nginxproxymanager.GetAccessLists200ResponseInnerOwnerPermissions) error {
	success, _, err := r.client.UsersAPI.UpdateUserPermissions(r.auth, userId).GetAccessLists200ResponseInnerOwnerPermissions(permissions).Execute()
	if err != nil {
		return err
	}

	if !success {
		return errors.New("unable to update user permissions")
	}

	return nil
}
