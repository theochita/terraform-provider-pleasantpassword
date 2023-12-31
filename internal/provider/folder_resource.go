// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	PPSClient "github.com/theochita/go-pleasant-password"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FolderResource{}
var _ resource.ResourceWithImportState = &FolderResource{}

func NewFolderResource() resource.Resource {
	return &FolderResource{}
}

// ExampleResource defines the resource implementation.
type FolderResource struct {
	client *PPSClient.APIClient
	ctx    *context.Context
}

// ExampleResourceModel describes the resource data model.
type FolderResourceModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	ParentID types.String `tfsdk:"parent_id"`
	Notes    types.String `tfsdk:"notes"`
}

func (r *FolderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folder"
}

func (r *FolderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The `folder` resource allows you to create and manage folders in Pleasant Password Server.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the folder.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the folder.",
				Required:            true,
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "The identifier of the parent folder.",
				Optional:            true,
				Computed:            true,
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "Additional notes for the folder.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *FolderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerclient, ok := req.ProviderData.(ProviderClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *PPSClient.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = &providerclient.Client
	r.ctx = &providerclient.Ctx

}

func (r *FolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FolderResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	param := PPSClient.NewV6CredentialGroupInputWithDefaults()
	param.Name = data.Name.ValueStringPointer()
	param.Notes = data.Notes.ValueStringPointer()
	param.ParentId = data.ParentID.ValueStringPointer()

	res, httpres, err := r.client.DefaultAPI.PostV6Folders(*r.ctx).V6CredentialGroupInput(*param).Execute()

	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API: ", err.Error())
		return
	}
	if httpres.StatusCode != 200 {
		resp.Diagnostics.AddError("Got an unexpected response code", fmt.Sprintf("Got an unexpected response code %v", httpres.StatusCode))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	sanityresult, err := strconv.Unquote(res)
	if err != nil {
		sanityresult = res
	}

	data.Id = types.StringValue(sanityresult)
	data.Name = types.StringValue(param.GetName())
	data.Notes = types.StringValue(param.GetNotes())
	data.ParentID = types.StringValue(param.GetParentId())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FolderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FolderResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, httpres, err := r.client.DefaultAPI.GetV6FoldersByID(*r.ctx, data.Id.ValueString()).Execute()

	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	if httpres.StatusCode != 200 {
		resp.Diagnostics.AddError("Got an unexpected response code", fmt.Sprintf("Got an unexpected response code %v", httpres.StatusCode))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}
	data.Id = types.StringValue(res.GetId())
	data.Name = types.StringValue(res.GetName())
	data.Notes = types.StringValue(res.GetNotes())
	data.ParentID = types.StringValue(res.GetParentId())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FolderResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	param := PPSClient.NewV6CredentialGroupInputWithDefaults()
	param.Name = data.Name.ValueStringPointer()
	param.Notes = data.Notes.ValueStringPointer()
	param.ParentId = data.ParentID.ValueStringPointer()

	httpres, err := r.client.DefaultAPI.PatchV6FoldersByID(*r.ctx, data.Id.ValueString()).V6CredentialGroupInput(*param).Execute()

	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API: ", err.Error())
		return
	}
	if httpres.StatusCode != 204 {
		resp.Diagnostics.AddError("Got an unexpected response code", fmt.Sprintf("Got an unexpected response code %v", httpres.StatusCode))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FolderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FolderResourceModel

	//Delete removes current folder and all subfolders and credentials

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpres, err := r.client.DefaultAPI.DeleteV6FoldersByID(*r.ctx, data.Id.ValueString()).Execute()

	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API: ", err.Error())
		return
	}
	if httpres.StatusCode != 204 {
		resp.Diagnostics.AddError("Got an unexpected response code", fmt.Sprintf("Got an unexpected response code %v", httpres.StatusCode))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *FolderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// not implemented
}
