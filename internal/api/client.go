// Copyright (c) James Hiester.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jameshiester/terraform-provider-bland/internal/config"
	utils "github.com/jameshiester/terraform-provider-bland/internal/util"
	arrays "github.com/jameshiester/terraform-provider-bland/internal/util/array"
)

// Client is a base client for specific API clients implemented in services.
type Client struct {
	Config   *config.ProviderConfig
	BaseAuth *Auth
}

// GetConfig returns the provider configuration.
func (client *Client) GetConfig() *config.ProviderConfig {
	return client.Config
}

// ProviderClient is a wrapper around the API client that provides additional helper methods.
type ProviderClient struct {
	Config *config.ProviderConfig
	Api    *Client
}

// ApiHttpResponse is a wrapper around http.Response that provides additional helper methods.
func NewApiClientBase(providerConfig *config.ProviderConfig, baseAuth *Auth) *Client {
	return &Client{
		Config:   providerConfig,
		BaseAuth: baseAuth,
	}
}

var retryableStatusCodes = []int{
	http.StatusUnauthorized,        // 401 is retryable because the token may have expired.
	http.StatusRequestTimeout,      // 408 is retryable because the request may have timed out.
	http.StatusTooEarly,            // 425 is retryable because the request may have been rate limited.
	http.StatusTooManyRequests,     // 429 is retryable because the request may have been rate limited.
	http.StatusInternalServerError, // 500 is retryable because the server may be overloaded.
	http.StatusBadGateway,          // 502 is retryable because the server may be overloaded.
	http.StatusServiceUnavailable,  // 503 is retryable because the server may be overloaded.
	http.StatusGatewayTimeout,      // 504 is retryable because the server may be overloaded.
	499,                            // 499 is retryable because the client may have closed the connection.
}

// Execute executes an HTTP request with the given method, url, headers, and body.
//
// Parameters:
//   - ctx: context.Context - Provides context for the request, allowing for timeout and cancellation control.
//   - scopes: []string - A list of scopes that the request should be associated with. If no scopes are provided, the method attempts to infer the scope from the URL.
//   - method: string - Specifies the HTTP method to be used for the request (e.g., "GET", "POST", "PATCH").
//   - url: string - The URL to which the request is sent. This includes the scheme, host, path, and query parameters. The URL must be absolute and properly formatted.
//   - headers: http.Header - A collection of HTTP headers to include in the request. Headers provide additional information about the request, such as content type, authorization tokens, and custom metadata.
//   - body: any - The body of the request, which can be of any type. This is typically used for methods like POST and PATCH, where data needs to be sent to the server.
//   - acceptableStatusCodes: []int - A list of HTTP status codes that are considered acceptable for the response. If the response status code is not in this list, the method treats it as an error.
//   - responseObj: any - An optional parameter where the response body can be unmarshaled into. This is useful for directly obtaining a structured representation of the response data.
//
// Returns:
//   - *Response: The response from the HTTP request.
//   - error: An error if the request fails. Possible error types include:
//   - UrlFormatError: Returned if the URL is invalid or not absolute.
//   - UnexpectedHttpStatusCodeError: Returned if the response status code is not acceptable.
//
// If no scopes are provided, the method attempts to infer the scope from the URL. The URL is validated to ensure it is absolute and properly formatted.
// The HTTP request is then prepared and executed. The response status code is checked against the list of acceptable status codes. If the status code
// is not acceptable, an error is returned. If a responseObj is provided, the response body is unmarshaled into this object.
func (client *Client) Execute(ctx context.Context, scopes []string, method, url string, headers http.Header, body any, acceptableStatusCodes []int, responseObj any) (*Response, error) {

	for {

		bodyBuffer, err := prepareRequestBody(body)
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequestWithContext(ctx, method, url, bodyBuffer)
		if err != nil {
			return nil, err
		}

		resp, err := client.doRequest(ctx, client.Config.APIKey, request, headers)
		if err != nil {
			return resp, fmt.Errorf("Error making %s request to %s. %w", request.Method, request.RequestURI, err)
		}

		isAcceptable := len(acceptableStatusCodes) > 0 && arrays.Contains(acceptableStatusCodes, resp.HttpResponse.StatusCode)
		if isAcceptable {
			if responseObj != nil && len(resp.BodyAsBytes) > 0 {
				err = resp.MarshallTo(responseObj)
				if err != nil {
					return resp, fmt.Errorf("Error marshalling response to json. %w", err)
				}
			}

			return resp, nil
		}

		isRetryable := arrays.Contains(retryableStatusCodes, resp.HttpResponse.StatusCode)
		if !isRetryable {
			return resp, NewUnexpectedHttpStatusCodeError(acceptableStatusCodes, resp.HttpResponse.StatusCode, resp.HttpResponse.Status, resp.BodyAsBytes)
		}

		waitFor := retryAfter(ctx, resp.HttpResponse)

		tflog.Debug(ctx, fmt.Sprintf("Received status code %d for request %s, retrying after %s", resp.HttpResponse.StatusCode, url, waitFor))

		err = client.SleepWithContext(ctx, waitFor)
		if err != nil {
			return resp, err
		}
	}
}

func (client *Client) HandleNotFoundResponse(resp *Response) error {
	if resp.HttpResponse.StatusCode == http.StatusNotFound {
		return fmt.Errorf("resource not found at '%s'", resp.HttpResponse.Request.URL)
	}
	return nil
}

func (client *Client) HandleForbiddenResponse(resp *Response) error {
	if resp.HttpResponse.StatusCode == http.StatusForbidden {
		return fmt.Errorf("access denied to resource at '%s'. Please validate your permissions", resp.HttpResponse.Request.URL)
	}
	return nil
}

// RetryAfterDefault returns a random duration between 10 and 20 seconds.
func DefaultRetryAfter() time.Duration {
	return time.Duration((rand.Intn(10) + 10)) * time.Second
}

// SleepWithContext sleeps for the given duration or until the context is canceled.
func (client *Client) SleepWithContext(ctx context.Context, duration time.Duration) error {
	if utils.IsTestContext(ctx) {
		return nil
	}
	if client.Config.TestMode {
		// Don't sleep during testing.
		return nil
	}
	select {
	case <-time.After(duration):
		// Time has elapsed.
		return nil
	case <-ctx.Done():
		// Context was canceled.
		return ctx.Err()
	}
}

func prepareRequestBody(body any) (io.Reader, error) {
	var bodyBuffer io.Reader
	if body != nil && (reflect.ValueOf(body).Kind() != reflect.Ptr || !reflect.ValueOf(body).IsNil()) {
		if strp, ok := body.(*string); ok {
			bodyBuffer = strings.NewReader(*strp)
		} else {
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			bodyBuffer = bytes.NewBuffer(bodyBytes)
		}
	}

	return bodyBuffer, nil
}
