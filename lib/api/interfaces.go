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
	"context"
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
//
// ctx carries the trace of the caller: the api passes the context of the request, the client
// sends it along with the request.
type Controller interface {
	Marshal(ctx context.Context, request messages.MarshallingRequest) (result map[string]string, err error, code int)
	MarshalForService(ctx context.Context, serviceId string, characteristicId string, request messages.MarshallingRequest) (result map[string]string, err error, code int)

	Unmarshal(ctx context.Context, request messages.UnmarshallingRequest) (result interface{}, err error, code int)
	UnmarshalForService(ctx context.Context, serviceId string, characteristicId string, request messages.UnmarshallingRequest) (result interface{}, err error, code int)

	MarshalV2(ctx context.Context, request messages.MarshallingV2Request) (result map[string]string, err error, code int)
	MarshalV2ForService(ctx context.Context, serviceId string, request messages.MarshallingV2Request) (result map[string]string, err error, code int)

	UnmarshalV2(ctx context.Context, request messages.UnmarshallingV2Request) (result interface{}, err error, code int)
	UnmarshalV2ForService(ctx context.Context, serviceId string, request messages.UnmarshallingV2Request) (result interface{}, err error, code int)

	GetCharacteristicPaths(ctx context.Context, serviceId string, characteristicId string) (result marshaller.CharacteristicsPathResponse, err error, code int)

	GetPathOptions(ctx context.Context, query messages.PathOptionsQuery) (result map[string][]marshaller.PathOptionsResultElement, err error, code int)

	FindConfigurables(ctx context.Context, characteristicId string, services []model.Service) (result configurables.Configurables, err error, code int)
	FindConfigurablesForServiceIds(ctx context.Context, characteristicId string, serviceIds []string) (result configurables.Configurables, err error, code int)

	TryConverterExtension(ctx context.Context, call converter.ExtensionCall) (result converter.ExtensionCallResponse, err error, code int)
}
