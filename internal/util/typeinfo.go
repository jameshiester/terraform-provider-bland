// Copyright (c) James Hiester.
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jameshiester/terraform-provider-bland/common"
	"github.com/jameshiester/terraform-provider-bland/internal/constants"
)

type RequestContextValue struct {
	ObjectName  string
	RequestType string
	RequestId   string
}

// TypeInfo represents a managed object type in the provider such as a resource or data source.
// Resource and data source types can inherit from TypeInfo to provide a consistent way to reference the type.
type TypeInfo struct {
	ProviderTypeName string
	TypeName         string
}

// FullTypeName returns the full type name in the format provider_type.
func (t *TypeInfo) FullTypeName() string {
	if t.ProviderTypeName == "" {
		return fmt.Sprintf("bland_%s", t.TypeName)
	}

	return fmt.Sprintf("%s_%s", t.ProviderTypeName, t.TypeName)
}

// ContextKey is a custom type for context keys.
// A custom type is needed to avoid collisions with other packages that use the same key.
type ContextKey string

type ExecutionContextValue struct {
	ProviderVersion string
	OperatingSystem string
	Architecture    string
	GoVersion       string
}

// TestContextValue is a struct that holds the test name for a given test.
type TestContextValue struct {
	IsTestMode bool
	TestName   string
}

// Context keys for the execution and request context.
const (
	EXECUTION_CONTEXT_KEY ContextKey = "executionContext"
	REQUEST_CONTEXT_KEY   ContextKey = "requestContext"
	TEST_CONTEXT_KEY      ContextKey = "testContext"
)

func UnitTestContext(ctx context.Context, testName string) context.Context {
	return context.WithValue(ctx, TEST_CONTEXT_KEY, TestContextValue{IsTestMode: true, TestName: testName})
}

// EnterRequestScope is a helper function that logs the start of a request scope and returns a closure that can be used to defer the exit of the request scope
// This function should be called at the start of a resource or data source request function
// The returned closure should be deferred at the start of the function
// The closure will log the end of the request scope
// The context is updated with the request context so that it can be accessed in lower level functions.
func EnterRequestContext[T AllowedRequestTypes](ctx context.Context, typ TypeInfo, req T) (context.Context, func()) {
	reqId := uuid.New().String()
	reqType := reflect.TypeOf(req).String()
	name := typ.FullTypeName()

	tflog.Debug(ctx, fmt.Sprintf("%s START: %s", reqType, name), map[string]any{
		"requestId":       reqId,
		"providerVersion": common.ProviderVersion,
	})

	// Add the request context to the context so that we can access it in lower level functions.
	ctx = context.WithValue(ctx, REQUEST_CONTEXT_KEY, RequestContextValue{RequestType: reqType, ObjectName: name, RequestId: reqId})
	ctx = tflog.SetField(ctx, "request_id", reqId)
	ctx = tflog.SetField(ctx, "request_type", reqType)

	ctx, cancel := enterTimeoutContext(ctx, req)

	// This returns a closure that can be used to defer the exit of the request scope.
	return ctx, func() {
		tflog.Debug(ctx, fmt.Sprintf("%s END: %s", reqType, name))
		if cancel != nil {
			(*cancel)()
		}
	}
}

// EnterTimeoutContext is a helper function that enters a timeout context based on the request type and the timeouts set in the plan or state.
func enterTimeoutContext[T AllowedRequestTypes](ctx context.Context, req T) (context.Context, *context.CancelFunc) {
	var tos timeouts.Value
	switch req := any(req).(type) {
	case resource.CreateRequest:
		diag := req.Plan.GetAttribute(ctx, path.Root("timeouts"), &tos)
		if diag.HasError() {
			return ctx, nil
		}

		dur, err := tos.Create(ctx, constants.DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES)
		if err != nil {
			// function returns default timeout even if error occurs
			tflog.Debug(ctx, "Could not retrieve create timeout, using default")
		}

		ctx, cancel := context.WithTimeout(ctx, dur)
		return ctx, &cancel
	case resource.ReadRequest:
		diag := req.State.GetAttribute(ctx, path.Root("timeouts"), &tos)
		if diag.HasError() {
			return ctx, nil
		}

		dur, err := tos.Read(ctx, constants.DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES)
		if err != nil {
			// function returns default timeout even if error occurs
			tflog.Debug(ctx, "Could not retrieve read timeout, using default")
		}

		ctx, cancel := context.WithTimeout(ctx, dur)
		return ctx, &cancel
	case resource.UpdateRequest:
		diag := req.Plan.GetAttribute(ctx, path.Root("timeouts"), &tos)
		if diag.HasError() {
			return ctx, nil
		}

		dur, err := tos.Update(ctx, constants.DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES)
		if err != nil {
			// function returns default timeout even if error occurs
			tflog.Debug(ctx, "Could not retrieve update timeout, using default")
		}

		ctx, cancel := context.WithTimeout(ctx, dur)
		return ctx, &cancel
	case resource.DeleteRequest:
		diag := req.State.GetAttribute(ctx, path.Root("timeouts"), &tos)
		if diag.HasError() {
			return ctx, nil
		}

		dur, err := tos.Delete(ctx, constants.DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES)
		if err != nil {
			// function returns default timeout even if error occurs
			tflog.Debug(ctx, "Could not retrieve delete timeout, using default")
		}

		ctx, cancel := context.WithTimeout(ctx, dur)
		return ctx, &cancel
	default:
		return ctx, nil
	}
}

// EnterProviderScope is a helper function that logs the start of a provider scope and returns a closure that can be used to defer the loging of the exit of the provider scope.
func EnterProviderContext[T AllowedProviderRequestTypes](ctx context.Context, req T) (context.Context, func()) {
	reqType := reflect.TypeOf(req).String()

	tflog.Debug(ctx, fmt.Sprintf("%s START", reqType), map[string]any{
		"providerVersion": common.ProviderVersion,
	})

	// This returns a closure that can be used to defer the exit of the provider scope.
	return ctx, func() {
		tflog.Debug(ctx, fmt.Sprintf("%s END", reqType), map[string]any{
			"providerVersion": common.ProviderVersion,
		})
	}
}

// AllowedRequestTypes is an interface that defines the allowed request types for the getRequestTypeName function.
type AllowedRequestTypes interface {
	resource.CreateRequest |
		resource.MetadataRequest |
		resource.ReadRequest |
		resource.UpdateRequest |
		resource.DeleteRequest |
		resource.SchemaRequest |
		resource.ConfigureRequest |
		resource.ModifyPlanRequest |
		resource.ImportStateRequest |
		resource.UpgradeStateRequest |
		resource.ValidateConfigRequest |
		datasource.ReadRequest |
		datasource.SchemaRequest |
		datasource.ConfigureRequest |
		datasource.MetadataRequest |
		datasource.ValidateConfigRequest
}

// AllowedProviderRequestTypes is an interface that defines the allowed request types for the EnterProviderContext function.
type AllowedProviderRequestTypes interface {
	provider.ConfigureRequest |
		provider.MetaSchemaRequest |
		provider.MetadataRequest |
		provider.SchemaRequest |
		provider.ValidateConfigRequest
}

// TestContext creates a new context with the test context value.
func TestContext(ctx context.Context, testName string) context.Context {
	return context.WithValue(ctx, TEST_CONTEXT_KEY, TestContextValue{IsTestMode: true, TestName: testName})
}

// IsTestContext returns true if the context is a test context.
func IsTestContext(ctx context.Context) bool {
	testContext, ok := ctx.Value(TEST_CONTEXT_KEY).(TestContextValue)
	return ok && testContext.IsTestMode
}

// CheckContextTimeout checks if the context has been cancelled or deadline exceeded.
// Returns an error if the context is done, nil otherwise.
// This is commonly used in polling loops to respect timeouts.
func CheckContextTimeout(ctx context.Context, operation string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("timed out during %s: %w", operation, ctx.Err())
	default:
		return nil
	}
}
