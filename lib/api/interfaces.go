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

package api

import (
	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	"github.com/SENERGY-Platform/marshaller/lib/converter"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
)

// Controller is the operation surface of this service: one method per endpoint, free of
// http types. The api implements the endpoints on top of it and lib/client implements the
// same interface over http, so a consumer can hold either without knowing which it has.
//
// The error code accompanying every result is the http status the operation would answer
// with, the way the device-repository controller reports it.
type Controller interface {
	Marshal(request messages.MarshallingRequest) (result map[string]string, err error, code int)
	MarshalForService(serviceId string, characteristicId string, request messages.MarshallingRequest) (result map[string]string, err error, code int)

	Unmarshal(request messages.UnmarshallingRequest) (result interface{}, err error, code int)
	UnmarshalForService(serviceId string, characteristicId string, request messages.UnmarshallingRequest) (result interface{}, err error, code int)

	MarshalV2(request messages.MarshallingV2Request) (result map[string]string, err error, code int)
	MarshalV2ForService(serviceId string, request messages.MarshallingV2Request) (result map[string]string, err error, code int)

	UnmarshalV2(request messages.UnmarshallingV2Request) (result interface{}, err error, code int)
	UnmarshalV2ForService(serviceId string, request messages.UnmarshallingV2Request) (result interface{}, err error, code int)

	GetCharacteristicPaths(serviceId string, characteristicId string) (result marshaller.CharacteristicsPathResponse, err error, code int)

	GetPathOptions(query messages.PathOptionsQuery) (result map[string][]marshaller.PathOptionsResultElement, err error, code int)

	FindConfigurables(characteristicId string, services []model.Service) (result configurables.Configurables, err error, code int)
	FindConfigurablesForServiceIds(characteristicId string, serviceIds []string) (result configurables.Configurables, err error, code int)

	TryConverterExtension(call converter.ExtensionCall) (result converter.ExtensionCallResponse, err error, code int)
}
