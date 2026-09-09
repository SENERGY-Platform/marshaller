/*
 * Copyright 2019 InfAI (CC SES)
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
	"net/http"
	"strings"

	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	"github.com/SENERGY-Platform/marshaller/lib/api/metrics"
	"github.com/SENERGY-Platform/marshaller/lib/config"
)

func init() {
	endpoints = append(endpoints, &Configurables{})
}

type Configurables struct{}

// GetConfigurables godoc
// @Summary      find configurables for a characteristic
// @Description  returns the values that can be configured alongside the given characteristic in the listed services
// @Tags         configurables
// @Produce      json
// @Security     Bearer
// @Param        characteristicId query string true "id of the characteristic that is already set"
// @Param        serviceIds query string true "comma separated service ids"
// @Success      200 {object} configurables.Configurables
// @Failure      400 {string} string "characteristicId or serviceIds is missing"
// @Failure      500 {string} string
// @Router       /configurables [GET]
func (this *Configurables) GetConfigurables(config config.Config, router *http.ServeMux, ctrl Controller, m *metrics.Metrics) {
	router.HandleFunc("GET /configurables", func(writer http.ResponseWriter, request *http.Request) {
		serviceIds := []string{}
		if serviceIdsStr := request.URL.Query().Get("serviceIds"); serviceIdsStr != "" {
			serviceIds = strings.Split(serviceIdsStr, ",")
		}
		result, err, code := ctrl.FindConfigurablesForServiceIds(request.URL.Query().Get("characteristicId"), serviceIds)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writeJson(config, writer, result)
	})
}

// FindConfigurables godoc
// @Summary      find configurables for a characteristic
// @Description  the request-body form of GET /configurables, for callers that hold the services already
// @Tags         configurables
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        query body messages.FindConfigurablesRequest true "the characteristic and the services to search"
// @Success      200 {object} configurables.Configurables
// @Failure      400 {string} string "the body is not valid json, or characteristic_id is missing"
// @Failure      500 {string} string
// @Router       /configurables [POST]
func (this *Configurables) FindConfigurables(config config.Config, router *http.ServeMux, ctrl Controller, m *metrics.Metrics) {
	router.HandleFunc("POST /configurables", func(writer http.ResponseWriter, request *http.Request) {
		msg, ok := decodeBody[messages.FindConfigurablesRequest](writer, request)
		if !ok {
			return
		}
		result, err, code := ctrl.FindConfigurables(msg.CharacteristicId, msg.Services)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writeJson(config, writer, result)
	})
}
