// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pathways

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jameshiester/terraform-provider-bland/internal/api"
	utils "github.com/jameshiester/terraform-provider-bland/internal/util"
)

var _ resource.Resource = &ConversationalPathwayResource{}
var _ resource.ResourceWithImportState = &ConversationalPathwayResource{}

type ConversationalPathwayResource struct {
	utils.TypeInfo
	PathwayClient client
}

func NewConversationalPathwayResource() resource.Resource {
	return &ConversationalPathwayResource{
		TypeInfo: utils.TypeInfo{
			TypeName: "connection",
		},
	}
}

func (r *ConversationalPathwayResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *ConversationalPathwayResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	_, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a [Conversational Pathway](https://docs.bland.ai/tutorials/pathways).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique pathway id",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the pathway",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the pathway",
				Required:            true,
			},
		},
	}
}

func (d *ConversationalPathwayResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("connection_parameters"),
			path.MatchRoot("connection_parameters_set"),
		),
	}
}

func (r *ConversationalPathwayResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	_, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
		return
	}

	client, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.PathwayClient = newPathwayClient(client.Api)
}

func (r *ConversationalPathwayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var plan ConversationalPathwayDataSourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	modelToCreate := createPathwayDto{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	connection, err := r.PathwayClient.CreatePathway(ctx, modelToCreate)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create connection", err.Error())
		return
	}

	responseModel := ConvertFromPathwayDto(*connection)
	plan.ID = responseModel.ID
	plan.Description = responseModel.Description
	plan.Name = responseModel.Name

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ConversationalPathwayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *ConversationalPathwayDataSourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathway, err := r.PathwayClient.GetPathway(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, api.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
		return
	}

	model := ConvertFromPathwayDto(*pathway)
	state.ID = model.ID
	state.Name = model.Name
	state.Description = model.Description
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ConversationalPathwayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan ConversationalPathwayDataSourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state ConversationalPathwayDataSourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateParams := updatePathwayDto{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	updateReponse, err := r.PathwayClient.UpdatePathway(ctx, plan.ID.ValueString(), updateParams)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
		return
	}

	modelState := ConvertFromPathwayDto(*updateReponse)
	plan.ID = modelState.ID
	plan.Name = modelState.Name
	plan.Description = modelState.Description

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ConversationalPathwayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *ConversationalPathwayDataSourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.PathwayClient.DeletePathway(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
		return
	}
}

func (r *ConversationalPathwayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
