/*
 * Copyright 2026 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package client speaks to this service over http. It implements api.Controller, so a
// consumer can hold either this client or the controller itself behind one interface —
// which is what makes an in-process test against the real implementation possible.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/SENERGY-Platform/marshaller/lib/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Interface interface {
	api.Controller
}

type Client struct {
	baseUrl string

	//sends every request; Options.Timeout is set on it
	httpClient *http.Client

	// optionalAuthTokenForApiGatewayRequest:
	//   - may be nil if used internally, without kong routing
	//   - is used for requests that internally don`t need an auth token but are forced to send one if the request is routed over the SENERGY-Platform api-gateway
	optionalAuthTokenForApiGatewayRequest func() (token string, err error)
}

// Options configure a client from NewClient. The zero value sends with http.DefaultClient, without a
// timeout of its own and without a token for the api gateway.
type Options struct {
	// Timeout bounds every request including reading its answer; 0 leaves the bound to the context
	// of the call and to HttpClient
	Timeout time.Duration

	// HttpClient sends the requests; nil means http.DefaultClient. Timeout is set on a copy, the
	// given client stays unchanged.
	HttpClient *http.Client

	// OptionalAuthTokenForApiGatewayRequest:
	//   - may be nil if used internally, without kong routing
	//   - is used for requests that internally don`t need an auth token but are forced to send one if the request is routed over the SENERGY-Platform api-gateway
	OptionalAuthTokenForApiGatewayRequest func() (token string, err error)
}

// NewClient creates a new client
func NewClient(baseUrl string, options Options) Interface {
	httpClient := options.HttpClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if options.Timeout > 0 {
		bounded := *httpClient
		bounded.Timeout = options.Timeout
		httpClient = &bounded
	}
	return &Client{baseUrl: baseUrl, httpClient: httpClient, optionalAuthTokenForApiGatewayRequest: options.OptionalAuthTokenForApiGatewayRequest}
}

// do adds the trace context of req to its headers, so the called service continues the trace of the caller.
// Without an OpenTelemetry setup in the calling service the global propagator is a no-op.
func do[T any](c *Client, req *http.Request) (result T, err error, code int) {
	if c.optionalAuthTokenForApiGatewayRequest != nil && req.Header.Get("Authorization") == "" {
		token, err := c.optionalAuthTokenForApiGatewayRequest()
		if err != nil {
			return result, err, http.StatusInternalServerError
		}
		req.Header.Set("Authorization", token)
	}
	otel.GetTextMapPropagator().Inject(req.Context(), propagation.HeaderCarrier(req.Header))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		temp, _ := io.ReadAll(resp.Body) //read error response and ensure that resp.Body is read to EOF
		return result, fmt.Errorf("unexpected statuscode %v: %v", resp.StatusCode, strings.TrimSpace(string(temp))), resp.StatusCode
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		_, _ = io.ReadAll(resp.Body) //ensure resp.Body is read to EOF
		return result, err, http.StatusInternalServerError
	}
	return result, nil, resp.StatusCode
}

// post sends body as json to path and decodes the answer into T.
func post[T any](ctx context.Context, c *Client, path string, body interface{}) (result T, err error, code int) {
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(body)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseUrl+path, buf)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	req.Header.Set("Content-Type", "application/json")
	return do[T](c, req)
}

func get[T any](ctx context.Context, c *Client, path string) (result T, err error, code int) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseUrl+path, nil)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	return do[T](c, req)
}
