/*
 * Copyright 2020 InfAI (CC SES)
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
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	v2 "github.com/SENERGY-Platform/marshaller/lib/marshaller/v2"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strings"
)

func init() {
	endpoints = append(endpoints, Configurables)
}

func Configurables(router *httprouter.Router, config config.Config, marshaller *marshaller.Marshaller, marshallerV2 *v2.Marshaller, configurableService *configurables.ConfigurableService, deviceRepo DeviceRepository) {
	resource := "/configurables"

	router.GET(resource, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		characteristicId := request.URL.Query().Get("characteristicId")
		if characteristicId == "" {
			http.Error(writer, "expect characteristicId as query-parameter", http.StatusBadRequest)
			return
		}
		serviceIdsStr := request.URL.Query().Get("serviceIds")
		if serviceIdsStr == "" {
			http.Error(writer, "expect serviceIds as query-parameter", http.StatusBadRequest)
			return
		}
		serviceIds := strings.Split(serviceIdsStr, ",")
		services := []model.Service{}
		for _, id := range serviceIds {
			service, err := deviceRepo.GetService(strings.TrimSpace(id))
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			services = append(services, service)
		}
		result, err := configurableService.Find(characteristicId, services)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(writer).Encode(result)
		if err != nil {
			log.Println("ERROR: unable to encode response", err)
		}
	})

	router.POST(resource, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		msg := FindConfigurablesRequest{}
		err := json.NewDecoder(request.Body).Decode(&msg)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if msg.CharacteristicId == "" {
			http.Error(writer, "expect characteristic_id as field in body", http.StatusBadRequest)
			return
		}
		result, err := configurableService.Find(msg.CharacteristicId, msg.Services)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(writer).Encode(result)
		if err != nil {
			log.Println("ERROR: unable to encode response", err)
		}
	})

}
