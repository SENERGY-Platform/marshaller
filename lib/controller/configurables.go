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

package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	"github.com/SENERGY-Platform/marshaller/lib/converter"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
)

func (this *Controller) FindConfigurables(characteristicId string, services []model.Service) (result configurables.Configurables, err error, code int) {
	if characteristicId == "" {
		return result, errors.New("expect characteristic_id as field in body"), http.StatusBadRequest
	}
	result, err = this.configurableService.Find(characteristicId, services)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	return result, nil, http.StatusOK
}

// FindConfigurablesForServiceIds resolves the services first, which is what the
// query-parameter form of the endpoint does.
func (this *Controller) FindConfigurablesForServiceIds(characteristicId string, serviceIds []string) (result configurables.Configurables, err error, code int) {
	if characteristicId == "" {
		return result, errors.New("expect characteristicId as query-parameter"), http.StatusBadRequest
	}
	if len(serviceIds) == 0 {
		return result, errors.New("expect serviceIds as query-parameter"), http.StatusBadRequest
	}
	services := []model.Service{}
	for _, id := range serviceIds {
		service, err := this.deviceRepo.GetService(strings.TrimSpace(id))
		if err != nil {
			return result, err, http.StatusInternalServerError
		}
		services = append(services, service)
	}
	return this.FindConfigurables(characteristicId, services)
}

func (this *Controller) TryConverterExtension(call converter.ExtensionCall) (result converter.ExtensionCallResponse, err error, code int) {
	if this.converter == nil {
		return result, errors.New("api initialized without converter"), http.StatusInternalServerError
	}
	result, err = this.converter.TryExtension(call)
	if err != nil {
		return result, err, http.StatusBadRequest
	}
	return result, nil, http.StatusOK
}
