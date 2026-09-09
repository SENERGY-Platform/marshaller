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
	endpoints = append(endpoints, &Marshalling{})
}

type Marshalling struct{}

// Marshal godoc
// @Summary      marshal a value into a protocol message
// @Description  marshals the given data into the message format of the service, using the characteristic to interpret it
// @Tags         marshal
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        message body messages.MarshallingRequest true "service, protocol and the data to marshal"
// @Success      200 {object} map[string]string "the protocol segments of the message"
// @Failure      400 {string} string "the body is not valid json, or the service does not reference the given protocol"
// @Failure      500 {string} string
// @Router       /marshal [POST]
func (this *Marshalling) Marshal(config config.Config, router *http.ServeMux, ctrl Controller, m *metrics.Metrics) {
	router.HandleFunc("POST /marshal", func(writer http.ResponseWriter, request *http.Request) {
		msg, ok := decodeBody[messages.MarshallingRequest](writer, request)
		if !ok {
			return
		}
		result, err, code := ctrl.Marshal(msg)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writeJson(config, writer, result)
	})
}

// MarshalForService godoc
// @Summary      marshal a value into a protocol message, service and characteristic from the path
// @Description  like POST /marshal, but the service is read from the device-repository by the id in the path and the characteristic comes from the path as well
// @Tags         marshal
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        serviceId path string true "id of the service the message is for"
// @Param        characteristicId path string true "id of the characteristic the data is given in"
// @Param        message body messages.MarshallingRequest true "the data to marshal"
// @Success      200 {object} map[string]string "the protocol segments of the message"
// @Failure      400 {string} string "the body is not valid json, or the service does not reference the given protocol"
// @Failure      500 {string} string
// @Router       /marshal/{serviceId}/{characteristicId} [POST]
func (this *Marshalling) MarshalForService(config config.Config, router *http.ServeMux, ctrl Controller, m *metrics.Metrics) {
	router.HandleFunc("POST /marshal/{serviceId}/{characteristicId}", func(writer http.ResponseWriter, request *http.Request) {
		msg, ok := decodeBody[messages.MarshallingRequest](writer, request)
		if !ok {
			return
		}
		result, err, code := ctrl.MarshalForService(request.PathValue("serviceId"), request.PathValue("characteristicId"), msg)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writeJson(config, writer, result)
	})
}
