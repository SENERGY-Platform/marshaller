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
