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
	"net/url"
	"strings"

	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	"github.com/SENERGY-Platform/marshaller/lib/converter"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
)

func (c *Client) GetCharacteristicPaths(serviceId string, characteristicId string) (result marshaller.CharacteristicsPathResponse, err error, code int) {
	return get[marshaller.CharacteristicsPathResponse](c.baseUrl, "/characteristic-paths/"+pathSegment(serviceId)+"/"+pathSegment(characteristicId), c.optionalAuthTokenForApiGatewayRequest)
}

// GetPathOptions uses the request-body form of the endpoint, because it is the one whose
// semantics do not depend on which parameters are present: an empty characteristic filter
// means "do not filter" there, while an absent characteristic-filter query parameter makes
// the GET form answer with an empty result.
func (c *Client) GetPathOptions(query messages.PathOptionsQuery) (result map[string][]marshaller.PathOptionsResultElement, err error, code int) {
	return post[map[string][]marshaller.PathOptionsResultElement](c.baseUrl, "/query/path-options", query, c.optionalAuthTokenForApiGatewayRequest)
}

func (c *Client) FindConfigurables(characteristicId string, services []model.Service) (result configurables.Configurables, err error, code int) {
	return post[configurables.Configurables](c.baseUrl, "/configurables", messages.FindConfigurablesRequest{
		CharacteristicId: characteristicId,
		Services:         services,
	}, c.optionalAuthTokenForApiGatewayRequest)
}

func (c *Client) FindConfigurablesForServiceIds(characteristicId string, serviceIds []string) (result configurables.Configurables, err error, code int) {
	query := url.Values{}
	query.Set("characteristicId", characteristicId)
	query.Set("serviceIds", strings.Join(serviceIds, ","))
	return get[configurables.Configurables](c.baseUrl, "/configurables?"+query.Encode(), c.optionalAuthTokenForApiGatewayRequest)
}

func (c *Client) TryConverterExtension(call converter.ExtensionCall) (result converter.ExtensionCallResponse, err error, code int) {
	return post[converter.ExtensionCallResponse](c.baseUrl, "/converter/extension-call", call, c.optionalAuthTokenForApiGatewayRequest)
}
