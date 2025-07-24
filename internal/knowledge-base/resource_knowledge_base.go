// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package knowledgebase

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jameshiester/terraform-provider-bland/internal/api"
	utils "github.com/jameshiester/terraform-provider-bland/internal/util"
)

var _ resource.Resource = &KnowledgeBaseResource{}

type KnowledgeBaseResource struct {
	utils.TypeInfo
	KnowledgeBaseClient *KnowledgeBaseClient
}

func NewKnowledgeBaseResource() resource.Resource {
	return &KnowledgeBaseResource{
		TypeInfo: utils.TypeInfo{
			TypeName: "knowledge_base",
		},
	}
}

func (r *KnowledgeBaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *KnowledgeBaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	_, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	if req.ProviderData == nil {
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
	r.KnowledgeBaseClient = NewKnowledgeBaseClient(client.Api)
}

func (r *KnowledgeBaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a knowledge base.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique knowledge base id",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the knowledge base",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the knowledge base",
				Required:            true,
			},
			"file": schema.StringAttribute{
				MarkdownDescription: "Base64 encoded file content for the knowledge base",
				Optional:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("text")),
				},
			},
			"text": schema.StringAttribute{
				MarkdownDescription: "Input text from the knowledge base",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("file")),
				},
			},
			"extracted_text": schema.StringAttribute{
				MarkdownDescription: "Extracted text from the knowledge base",
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *KnowledgeBaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	_, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan KnowledgeBaseModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dto := ConvertToCreateKnowledgeBaseDto(plan)
	vectorID, err := r.KnowledgeBaseClient.CreateKnowledgeBase(ctx, dto)
	if err != nil {
		resp.Diagnostics.AddError("Error creating knowledge base", err.Error())
		return
	}

	read, err := r.KnowledgeBaseClient.ReadKnowledgeBase(ctx, *vectorID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading knowledge base", err.Error())
		return
	}

	model := ConvertFromKnowledgeBaseDto(*read)
	model.File = plan.File
	model.Text = plan.Text
	resp.State.Set(ctx, model)
}

func (r *KnowledgeBaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	_, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state KnowledgeBaseModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	read, err := r.KnowledgeBaseClient.ReadKnowledgeBase(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading knowledge base", err.Error())
		return
	}

	model := ConvertFromKnowledgeBaseDto(*read)
	model.File = state.File
	model.Text = state.Text
	resp.State.Set(ctx, model)
}

func (r *KnowledgeBaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	_, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan KnowledgeBaseModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dto := ConvertToUpdateKnowledgeBaseDto(plan)
	_, err := r.KnowledgeBaseClient.UpdateKnowledgeBase(ctx, plan.ID.ValueString(), dto)
	if err != nil {
		resp.Diagnostics.AddError("Error updating knowledge base", err.Error())
		return
	}
	read, err := r.KnowledgeBaseClient.ReadKnowledgeBase(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading knowledge base", err.Error())
		return
	}

	model := ConvertFromKnowledgeBaseDto(*read)
	model.File = plan.File
	model.Text = plan.Text
	resp.State.Set(ctx, model)
}

func (r *KnowledgeBaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	_, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state KnowledgeBaseModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.KnowledgeBaseClient.DeleteKnowledgeBase(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting knowledge base", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
