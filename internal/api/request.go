package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jameshiester/terraform-provider-bland/common"
	"github.com/jameshiester/terraform-provider-bland/internal/constants"
	utils "github.com/jameshiester/terraform-provider-bland/internal/util"
)

func (client *Client) doRequest(ctx context.Context, token string, request *http.Request, headers http.Header) (*Response, error) {
	if headers != nil {
		request.Header = headers
	}

	if token == "" {
		return nil, errors.New("API key is empty")
	}

	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}

	httpClient := http.DefaultClient

	if request.Header.Get("Authorization") == "" {
		request.Header.Set("Authorization", token)
	}

	ua := client.buildUserAgent(ctx)
	request.Header.Set("User-Agent", ua)
	sessionId, requestId := client.buildCorrelationHeaders(ctx)
	request.Header.Set("X-Correlation-Id", sessionId)
	request.Header.Set("X-Request-Id", requestId)

	apiResponse, err := httpClient.Do(request)
	resp := &Response{
		HttpResponse: apiResponse,
	}

	if err != nil {
		return resp, err
	}

	if apiResponse == nil {
		return resp, errors.New("unexpected nil response without error")
	}

	defer apiResponse.Body.Close()
	body, err := io.ReadAll(apiResponse.Body)
	resp.BodyAsBytes = body

	return resp, err
}

type Response struct {
	HttpResponse *http.Response
	BodyAsBytes  []byte
}

func (apiResponse *Response) MarshallTo(obj any) error {
	// Ensure obj is a pointer to avoid silent failures
	rv := reflect.ValueOf(obj)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("MarshallTo requires a non-nil pointer, got %T", obj)
	}

	return json.NewDecoder(bytes.NewReader(apiResponse.BodyAsBytes)).Decode(obj)
}

func (apiResponse *Response) GetHeader(name string) string {
	return apiResponse.HttpResponse.Header.Get(name)
}

func retryAfter(ctx context.Context, resp *http.Response) time.Duration {
	retryHeader := resp.Header.Get(constants.HEADER_RETRY_AFTER)
	if retryHeader == "" {
		return DefaultRetryAfter()
	}
	tflog.Debug(ctx, "Retry Header: "+retryHeader)

	// Check if the header is a delta-seconds value (integer)
	if deltaSeconds, err := strconv.Atoi(retryHeader); err == nil {
		return time.Duration(deltaSeconds) * time.Second
	}

	// Check if the header is an HTTP-date
	if retryTime, err := http.ParseTime(retryHeader); err == nil {
		// Calculate duration until the retry time
		duration := time.Until(retryTime)
		if duration > 0 {
			return duration
		}
	}

	// Try to parse as a duration string (non-standard but sometimes used)
	if retryAfter, err := time.ParseDuration(retryHeader); err == nil {
		return retryAfter
	}

	// Fallback to a default retry duration
	tflog.Debug(ctx, "Invalid Retry-After header, falling back to default")
	return DefaultRetryAfter()
}

func (client *Client) buildCorrelationHeaders(ctx context.Context) (sessionId string, requestId string) {
	sessionId = ""
	requestId = uuid.New().String() // Generate a new request ID for each request
	requestContext, ok := ctx.Value("").(utils.RequestContextValue)
	if ok {
		// If the request context is available, use the session ID from the request context
		sessionId = requestContext.RequestId
	}
	return sessionId, requestId
}

func (client *Client) buildUserAgent(ctx context.Context) string {
	userAgent := fmt.Sprintf("terraform-provider-bland/%s (%s; %s) terraform/%s go/%s", common.ProviderVersion, runtime.GOOS, runtime.GOARCH, client.Config.TerraformVersion, runtime.Version())

	requestContext, ok := ctx.Value("requestContext").(utils.RequestContextValue)
	if ok {
		userAgent += fmt.Sprintf(" %s %s", requestContext.ObjectName, requestContext.RequestType)
	}

	return userAgent
}

var _ error = UnexpectedHttpStatusCodeError{}

type UnexpectedHttpStatusCodeError struct {
	ExpectedStatusCodes []int
	StatusCode          int
	StatusText          string
	Body                []byte
}

func (e UnexpectedHttpStatusCodeError) Error() string {
	return fmt.Sprintf("Unexpected HTTP status code. Expected: %v, received: [%d] %s | %s", e.ExpectedStatusCodes, e.StatusCode, e.StatusText, e.Body)
}

func NewUnexpectedHttpStatusCodeError(expectedStatusCodes []int, statusCode int, statusText string, body []byte) error {
	return UnexpectedHttpStatusCodeError{
		ExpectedStatusCodes: expectedStatusCodes,
		StatusCode:          statusCode,
		StatusText:          statusText,
		Body:                body,
	}
}
