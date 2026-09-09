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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/SENERGY-Platform/marshaller/lib/api"
)

type Interface interface {
	api.Controller
}

type Client struct {
	baseUrl string

	// optionalAuthTokenForApiGatewayRequest:
	//   - may be nil if used internally, without kong routing
	//   - is used for requests that internally don`t need an auth token but are forced to send one if the request is routed over the SENERGY-Platform api-gateway
	optionalAuthTokenForApiGatewayRequest func() (token string, err error)
}

// NewClient creates a new client
// optionalAuthTokenForApiGatewayRequest:
//   - may be nil if used internally, without kong routing
//   - is used for requests that internally don`t need an auth token but are forced to send one if the request is routed over the SENERGY-Platform api-gateway
func NewClient(baseUrl string, optionalAuthTokenForApiGatewayRequest func() (token string, err error)) Interface {
	return &Client{baseUrl: baseUrl, optionalAuthTokenForApiGatewayRequest: optionalAuthTokenForApiGatewayRequest}
}

func do[T any](req *http.Request, optionalAuthTokenForApiGatewayRequest func() (token string, err error)) (result T, err error, code int) {
	if optionalAuthTokenForApiGatewayRequest != nil && req.Header.Get("Authorization") == "" {
		token, err := optionalAuthTokenForApiGatewayRequest()
		if err != nil {
			return result, err, http.StatusInternalServerError
		}
		req.Header.Set("Authorization", token)
	}
	resp, err := http.DefaultClient.Do(req)
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
func post[T any](baseUrl string, path string, body interface{}, optionalAuthTokenForApiGatewayRequest func() (token string, err error)) (result T, err error, code int) {
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(body)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	req, err := http.NewRequest(http.MethodPost, baseUrl+path, buf)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	req.Header.Set("Content-Type", "application/json")
	return do[T](req, optionalAuthTokenForApiGatewayRequest)
}

func get[T any](baseUrl string, path string, optionalAuthTokenForApiGatewayRequest func() (token string, err error)) (result T, err error, code int) {
	req, err := http.NewRequest(http.MethodGet, baseUrl+path, nil)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	return do[T](req, optionalAuthTokenForApiGatewayRequest)
}

// pathSegment escapes an id for use in a request path. Service and characteristic ids are
// urns, so they carry characters that have to survive the round trip.
func pathSegment(id string) string {
	return url.PathEscape(id)
}
