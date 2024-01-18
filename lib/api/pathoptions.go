/*
 * Copyright 2021 InfAI (CC SES)
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
	"encoding/json"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	"github.com/SENERGY-Platform/marshaller/lib/converter"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	v2 "github.com/SENERGY-Platform/marshaller/lib/marshaller/v2"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"strings"
)

func init() {
	endpoints = append(endpoints, PathOptions)
}

func PathOptions(router *httprouter.Router, config config.Config, marshaller *marshaller.Marshaller, marshallerV2 *v2.Marshaller, service *configurables.ConfigurableService, repo DeviceRepository, converter *converter.Converter, metrics *Metrics) {

	router.GET("/path-options", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		deviceTypeIdsStr := request.URL.Query().Get("device-type-ids")
		if deviceTypeIdsStr == "" {
			writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(writer).Encode(map[string]interface{}{})
			return
		}
		deviceTypeIds := strings.Split(strings.ReplaceAll(deviceTypeIdsStr, " ", ""), ",")

		characteristicIdFilterStr := request.URL.Query().Get("characteristic-filter")
		if characteristicIdFilterStr == "" {
			writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(writer).Encode(map[string]interface{}{})
			return
		}
		characteristicIdFilter := strings.Split(strings.ReplaceAll(characteristicIdFilterStr, " ", ""), ",")

		functionId := strings.TrimSpace(request.URL.Query().Get("function-id"))
		aspectId := strings.TrimSpace(request.URL.Query().Get("aspect-id"))

		withoutEnvelope := false
		withoutEnvelopeStr := request.URL.Query().Get("function-id")
		var err error
		if withoutEnvelopeStr != "" {
			withoutEnvelope, err = strconv.ParseBool(strings.TrimSpace(withoutEnvelopeStr))
		}

		result, err, code := marshaller.GetPathOption(deviceTypeIds, functionId, aspectId, characteristicIdFilter, !withoutEnvelope)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		} else {
			writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(writer).Encode(result)
			return
		}
	})

	router.POST("/query/path-options", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		query := PathOptionsQuery{}
		err := json.NewDecoder(request.Body).Decode(&query)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err, code := marshaller.GetPathOption(query.DeviceTypeIds, query.FunctionId, query.AspectId, query.CharacteristicIdFilter, !query.WithoutEnvelope)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		} else {
			writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(writer).Encode(result)
			return
		}
	})
}
