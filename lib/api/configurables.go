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
	"github.com/SENERGY-Platform/marshaller-service/lib/config"
	"github.com/SENERGY-Platform/marshaller-service/lib/configurables"
	"github.com/SENERGY-Platform/marshaller-service/lib/marshaller"
	"github.com/SENERGY-Platform/marshaller-service/lib/marshaller/model"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"log"
	"net/http"
	"strings"
)

func init() {
	endpoints = append(endpoints, Configurables)
}

func Configurables(router *jwt_http_router.Router, conf config.Config, marshaller *marshaller.Marshaller, configurableService *configurables.ConfigurableService, deviceRepo DeviceRepository) {
	resource := "/configurables"

	router.GET(resource, func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
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
			service, err := deviceRepo.GetService(config.Impersonate(jwt.Impersonate), strings.TrimSpace(id))
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

	router.POST(resource, func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
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
