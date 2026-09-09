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
	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
)

func (c *Client) Marshal(request messages.MarshallingRequest) (result map[string]string, err error, code int) {
	return post[map[string]string](c.baseUrl, "/marshal", request, c.optionalAuthTokenForApiGatewayRequest)
}

func (c *Client) MarshalForService(serviceId string, characteristicId string, request messages.MarshallingRequest) (result map[string]string, err error, code int) {
	return post[map[string]string](c.baseUrl, "/marshal/"+pathSegment(serviceId)+"/"+pathSegment(characteristicId), request, c.optionalAuthTokenForApiGatewayRequest)
}

func (c *Client) MarshalV2(request messages.MarshallingV2Request) (result map[string]string, err error, code int) {
	return post[map[string]string](c.baseUrl, "/v2/marshal", request, c.optionalAuthTokenForApiGatewayRequest)
}

func (c *Client) MarshalV2ForService(serviceId string, request messages.MarshallingV2Request) (result map[string]string, err error, code int) {
	return post[map[string]string](c.baseUrl, "/v2/marshal/"+pathSegment(serviceId), request, c.optionalAuthTokenForApiGatewayRequest)
}

func (c *Client) Unmarshal(request messages.UnmarshallingRequest) (result interface{}, err error, code int) {
	return post[interface{}](c.baseUrl, "/unmarshal", request, c.optionalAuthTokenForApiGatewayRequest)
}

func (c *Client) UnmarshalForService(serviceId string, characteristicId string, request messages.UnmarshallingRequest) (result interface{}, err error, code int) {
	return post[interface{}](c.baseUrl, "/unmarshal/"+pathSegment(serviceId)+"/"+pathSegment(characteristicId), request, c.optionalAuthTokenForApiGatewayRequest)
}

func (c *Client) UnmarshalV2(request messages.UnmarshallingV2Request) (result interface{}, err error, code int) {
	return post[interface{}](c.baseUrl, "/v2/unmarshal", request, c.optionalAuthTokenForApiGatewayRequest)
}

func (c *Client) UnmarshalV2ForService(serviceId string, request messages.UnmarshallingV2Request) (result interface{}, err error, code int) {
	return post[interface{}](c.baseUrl, "/v2/unmarshal/"+pathSegment(serviceId), request, c.optionalAuthTokenForApiGatewayRequest)
}
