// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package secret

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jameshiester/terraform-provider-bland/internal/api"
	utils "github.com/jameshiester/terraform-provider-bland/internal/util"
)

var _ resource.ResourceWithValidateConfig = &SecretResource{}

type SecretResource struct {
	utils.TypeInfo
	SecretClient *SecretClient
}

func NewSecretResource() resource.Resource {
	return &SecretResource{
		TypeInfo: utils.TypeInfo{
			TypeName: "secret",
		},
	}
}

func (r SecretResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data SecretModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If attribute_one is not configured, return without warning.
	if data.Static.ValueBool() && data.Value.ValueStringPointer() != nil {
		return
	}

	// If attribute_two is not null, return without warning.
	if !data.Static.ValueBool() && data.Config != nil {
		return
	}
	if data.Static.ValueBool() {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("attribute_two"),
			"Missing Attribute Configuration",
			"Expected 'value' to be set when static is set to true. "+
				"The resource may return unexpected results.",
		)
	} else {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("attribute_two"),
			"Missing Attribute Configuration",
			"Expected 'config' to be set when static is set to true. "+
				"The resource may return unexpected results.",
		)
	}

}

func (r *SecretResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *SecretResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.SecretClient = newSecretClient(client.Api)
}

func (r *SecretResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a secret.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the secret.",
			},
			"static": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Defines if secret is static or refreshes.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The value for a static secret.",
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("config")),
				},
			},
			"config": schema.SingleNestedAttribute{
				MarkdownDescription: "Configuration for refreshable secret.",
				Optional:            true,
				Validators: []validator.Object{
					objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("value")),
				},
				Attributes: map[string]schema.Attribute{
					"body": schema.StringAttribute{
						MarkdownDescription: "JSON body for the refresh request.",
						Optional:            true,
					},
					"method": schema.StringAttribute{
						MarkdownDescription: "HTTP method for the refresh request.",
						Required:            true,
					},
					"refresh_interval": schema.Int32Attribute{
						MarkdownDescription: "Refresh interval for the refresh request.",
						Optional:            true,
						Computed:            true,
						Default:             int32default.StaticInt32(60),
					},
					"response": schema.StringAttribute{
						MarkdownDescription: "Value to extract from the refresh request response.",
						Required:            true,
					},
					"url": schema.StringAttribute{
						MarkdownDescription: "URL for the refresh request.",
						Required:            true,
					},
					"headers": schema.MapAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "Headers for the refresh request.",
						Optional:            true,
					},
				},
			},
		},
	}
}

func (r *SecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SecretModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	dto, err := ConvertToSecretDto(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing secret", err.Error())
		return
	}
	model := createSecretDto{
		Name:  dto.Name,
		Value: dto.Value,
	}
	created, err := r.SecretClient.CreateSecret(ctx, model)
	if err != nil {
		resp.Diagnostics.AddError("Error creating secret", err.Error())
		return
	}
	createdModel := ConvertFromSecretDto(*created)
	resp.State.Set(ctx, createdModel)
}

func (r *SecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SecretModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	read, err := r.SecretClient.ReadSecret(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading secret", err.Error())
		return
	}
	model := ConvertFromSecretDto(*read)
	model.Value = state.Value
	resp.State.Set(ctx, model)
}

func (r *SecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SecretModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	dto, err := ConvertToSecretDto(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing secret", err.Error())
		return
	}
	if plan.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Error parsing secret ID", "failed to parse Secret ID")
		return
	}
	updateDto := updateSecretDto{
		Name:   dto.Name,
		Config: dto.Config,
		Value:  dto.Value,
	}
	updated, err := r.SecretClient.UpdateSecret(ctx, plan.ID.ValueString(), updateDto)
	if err != nil {
		resp.Diagnostics.AddError("Error updating secret", err.Error())
		return
	}
	model := ConvertFromSecretDto(*updated)
	resp.State.Set(ctx, model)
}

func (r *SecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SecretModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.SecretClient.DeleteSecret(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting secret", err.Error())
		return
	}
	resp.State.RemoveResource(ctx)
}
