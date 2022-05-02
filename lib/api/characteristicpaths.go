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
	v2 "github.com/SENERGY-Platform/marshaller/lib/marshaller/v2"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func init() {
	endpoints = append(endpoints, CharacteristicPathEndpoint)
}

func CharacteristicPathEndpoint(router *httprouter.Router, config config.Config, marshaller *marshaller.Marshaller, marshallerV2 *v2.Marshaller, configurableService *configurables.ConfigurableService, deviceRepo DeviceRepository) {
	resource := "/characteristic-paths"

	router.GET(resource+"/:serviceId/:characteristicId", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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
		service, err, code := deviceRepo.GetServiceWithErrCode(serviceId)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		result, err, code := marshaller.GetServiceCharacteristicPath(service, characteristicId)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		json.NewEncoder(writer).Encode(result)
	})

}
