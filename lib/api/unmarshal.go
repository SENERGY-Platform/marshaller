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

	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	"github.com/SENERGY-Platform/marshaller/lib/api/metrics"
	"github.com/SENERGY-Platform/marshaller/lib/config"
)

func init() {
	endpoints = append(endpoints, &Unmarshalling{})
}

type Unmarshalling struct{}

// Unmarshal godoc
// @Summary      unmarshal a protocol message into a value
// @Description  reads the message of the given service and returns its value in the requested characteristic
// @Tags         unmarshal
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        message body messages.UnmarshallingRequest true "service, protocol and the message to unmarshal"
// @Success      200 {object} interface{} "the value in the requested characteristic"
// @Failure      400 {string} string "the body is not valid json, or the service does not reference the given protocol"
// @Failure      500 {string} string
// @Router       /unmarshal [POST]
func (this *Unmarshalling) Unmarshal(config config.Config, router *http.ServeMux, ctrl Controller, m *metrics.Metrics) {
	router.HandleFunc("POST /unmarshal", func(writer http.ResponseWriter, request *http.Request) {
		msg, ok := decodeBody[messages.UnmarshallingRequest](writer, request)
		if !ok {
			return
		}
		result, err, code := ctrl.Unmarshal(msg)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writeJson(config, writer, result)
	})
}

// UnmarshalForService godoc
// @Summary      unmarshal a protocol message, service and characteristic from the path
// @Description  like POST /unmarshal, but the service is read from the device-repository by the id in the path and the characteristic comes from the path as well
// @Tags         unmarshal
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        serviceId path string true "id of the service the message came from"
// @Param        characteristicId path string true "id of the characteristic the value is wanted in"
// @Param        message body messages.UnmarshallingRequest true "the message to unmarshal"
// @Success      200 {object} interface{} "the value in the requested characteristic"
// @Failure      400 {string} string "the body is not valid json, or the service does not reference the given protocol"
// @Failure      500 {string} string
// @Router       /unmarshal/{serviceId}/{characteristicId} [POST]
func (this *Unmarshalling) UnmarshalForService(config config.Config, router *http.ServeMux, ctrl Controller, m *metrics.Metrics) {
	router.HandleFunc("POST /unmarshal/{serviceId}/{characteristicId}", func(writer http.ResponseWriter, request *http.Request) {
		msg, ok := decodeBody[messages.UnmarshallingRequest](writer, request)
		if !ok {
			return
		}
		result, err, code := ctrl.UnmarshalForService(request.PathValue("serviceId"), request.PathValue("characteristicId"), msg)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writeJson(config, writer, result)
	})
}
