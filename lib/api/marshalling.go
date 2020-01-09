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
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/marshaller-service/lib/config"
	"github.com/SENERGY-Platform/marshaller-service/lib/configurables"
	"github.com/SENERGY-Platform/marshaller-service/lib/marshaller"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"log"
	"net/http"
)

func init() {
	endpoints = append(endpoints, Marshalling)
}

func Marshalling(router *jwt_http_router.Router, conf config.Config, marshaller *marshaller.Marshaller, configurableService *configurables.ConfigurableService, deviceRepo DeviceRepository) {
	resource := "/marshalling"

	normalizeRequest := func(request *MarshallingRequest, jwt jwt_http_router.Jwt) error {
		if request.Protocol != nil {
			protocol, err := deviceRepo.GetProtocol(config.Impersonate((jwt.Impersonate)), request.Service.ProtocolId)
			if err != nil {
				return err
			}
			request.Protocol = &protocol
		} else if request.Service.ProtocolId != request.Protocol.Id {
			return errors.New("expect service to reference given protocol")
		}
		return nil
	}

	marshal := func(request MarshallingRequest) (map[string]string, error) {
		return marshaller.MarshalInputs(*request.Protocol, request.Service, request.Data, request.CharacteristicId, request.Configurables...)
	}

	router.POST(resource+"/:serviceId/:characteristicId", func(writer http.ResponseWriter, request *http.Request, params jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		msg := MarshallingRequest{}
		serviceId := params.ByName("serviceId")
		if serviceId == "" {
			http.Error(writer, "expect serviceId as parameter in path", http.StatusBadRequest)
			return
		}
		characteristicId := params.ByName("characteristicId")
		if characteristicId == "" {
			http.Error(writer, "expect characteristicId as parameter in path", http.StatusBadRequest)
			return
		}
		err := json.NewDecoder(request.Body).Decode(&msg)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		msg.CharacteristicId = characteristicId
		msg.Service, err = deviceRepo.GetService(config.Impersonate(jwt.Impersonate), serviceId)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		err = normalizeRequest(&msg, jwt)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := marshal(msg)
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
		msg := MarshallingRequest{}
		err := json.NewDecoder(request.Body).Decode(&msg)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		err = normalizeRequest(&msg, jwt)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := marshal(msg)
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
