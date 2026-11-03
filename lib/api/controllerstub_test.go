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
	"net/http"

	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	"github.com/SENERGY-Platform/marshaller/lib/converter"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
)

// controllerStub answers everything with an empty result. The routing tests only ask
// whether a route exists, so what it returns does not matter — that it implements
// Controller does.
type controllerStub struct{}

func (this *controllerStub) Marshal(ctx context.Context, request messages.MarshallingRequest) (map[string]string, error, int) {
	return map[string]string{}, nil, http.StatusOK
}

func (this *controllerStub) MarshalForService(ctx context.Context, serviceId string, characteristicId string, request messages.MarshallingRequest) (map[string]string, error, int) {
	return map[string]string{}, nil, http.StatusOK
}

func (this *controllerStub) Unmarshal(ctx context.Context, request messages.UnmarshallingRequest) (interface{}, error, int) {
	return nil, nil, http.StatusOK
}

func (this *controllerStub) UnmarshalForService(ctx context.Context, serviceId string, characteristicId string, request messages.UnmarshallingRequest) (interface{}, error, int) {
	return nil, nil, http.StatusOK
}

func (this *controllerStub) MarshalV2(ctx context.Context, request messages.MarshallingV2Request) (map[string]string, error, int) {
	return map[string]string{}, nil, http.StatusOK
}

func (this *controllerStub) MarshalV2ForService(ctx context.Context, serviceId string, request messages.MarshallingV2Request) (map[string]string, error, int) {
	return map[string]string{}, nil, http.StatusOK
}

func (this *controllerStub) UnmarshalV2(ctx context.Context, request messages.UnmarshallingV2Request) (interface{}, error, int) {
	return nil, nil, http.StatusOK
}

func (this *controllerStub) UnmarshalV2ForService(ctx context.Context, serviceId string, request messages.UnmarshallingV2Request) (interface{}, error, int) {
	return nil, nil, http.StatusOK
}

func (this *controllerStub) GetCharacteristicPaths(ctx context.Context, serviceId string, characteristicId string) (marshaller.CharacteristicsPathResponse, error, int) {
	return marshaller.CharacteristicsPathResponse{}, nil, http.StatusOK
}

func (this *controllerStub) GetPathOptions(ctx context.Context, query messages.PathOptionsQuery) (map[string][]marshaller.PathOptionsResultElement, error, int) {
	return map[string][]marshaller.PathOptionsResultElement{}, nil, http.StatusOK
}

func (this *controllerStub) FindConfigurables(ctx context.Context, characteristicId string, services []model.Service) (configurables.Configurables, error, int) {
	return configurables.Configurables{}, nil, http.StatusOK
}

func (this *controllerStub) FindConfigurablesForServiceIds(ctx context.Context, characteristicId string, serviceIds []string) (configurables.Configurables, error, int) {
	return configurables.Configurables{}, nil, http.StatusOK
}

func (this *controllerStub) TryConverterExtension(ctx context.Context, call converter.ExtensionCall) (converter.ExtensionCallResponse, error, int) {
	return converter.ExtensionCallResponse{}, nil, http.StatusOK
}
