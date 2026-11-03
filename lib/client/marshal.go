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

package client

import (
	"context"
	"net/url"

	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
)

func (c *Client) Marshal(ctx context.Context, request messages.MarshallingRequest) (result map[string]string, err error, code int) {
	return post[map[string]string](ctx, c, "/marshal", request)
}

func (c *Client) MarshalForService(ctx context.Context, serviceId string, characteristicId string, request messages.MarshallingRequest) (result map[string]string, err error, code int) {
	return post[map[string]string](ctx, c, "/marshal/"+url.PathEscape(serviceId)+"/"+url.PathEscape(characteristicId), request)
}

func (c *Client) MarshalV2(ctx context.Context, request messages.MarshallingV2Request) (result map[string]string, err error, code int) {
	return post[map[string]string](ctx, c, "/v2/marshal", request)
}

func (c *Client) MarshalV2ForService(ctx context.Context, serviceId string, request messages.MarshallingV2Request) (result map[string]string, err error, code int) {
	return post[map[string]string](ctx, c, "/v2/marshal/"+url.PathEscape(serviceId), request)
}

func (c *Client) Unmarshal(ctx context.Context, request messages.UnmarshallingRequest) (result interface{}, err error, code int) {
	return post[interface{}](ctx, c, "/unmarshal", request)
}

func (c *Client) UnmarshalForService(ctx context.Context, serviceId string, characteristicId string, request messages.UnmarshallingRequest) (result interface{}, err error, code int) {
	return post[interface{}](ctx, c, "/unmarshal/"+url.PathEscape(serviceId)+"/"+url.PathEscape(characteristicId), request)
}

func (c *Client) UnmarshalV2(ctx context.Context, request messages.UnmarshallingV2Request) (result interface{}, err error, code int) {
	return post[interface{}](ctx, c, "/v2/unmarshal", request)
}

func (c *Client) UnmarshalV2ForService(ctx context.Context, serviceId string, request messages.UnmarshallingV2Request) (result interface{}, err error, code int) {
	return post[interface{}](ctx, c, "/v2/unmarshal/"+url.PathEscape(serviceId), request)
}
