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
	"time"

	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	"github.com/SENERGY-Platform/marshaller/lib/api/metrics"
	"github.com/SENERGY-Platform/marshaller/lib/config"
)

func init() {
	endpoints = append(endpoints, &UnmarshallingV2{})
}

type UnmarshallingV2 struct{}

// UnmarshalV2 godoc
// @Summary      unmarshal a protocol message into a value
// @Description  returns the value at the requested path; without a path it is determined from the function and the requested aspects, closest aspect first
// @Tags         unmarshal, v2
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        message body messages.UnmarshallingV2Request true "service, protocol and the message to unmarshal"
// @Success      200 {object} interface{} "the value in the requested characteristic"
// @Failure      400 {string} string "the body is not valid json, the service does not reference the given protocol, or no output path matches the criteria"
// @Failure      500 {string} string
// @Router       /v2/unmarshal [POST]
func (this *UnmarshallingV2) UnmarshalV2(config config.Config, router *http.ServeMux, ctrl Controller, m *metrics.Metrics) {
	resource := "/v2/unmarshal"
	router.HandleFunc("POST "+resource, func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		msg, ok := decodeBody[messages.UnmarshallingV2Request](writer, request)
		if !ok {
			return
		}
		result, err, code := ctrl.UnmarshalV2(msg)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writeJson(config, writer, result)
		m.LogUnmarshallingRequest(request, resource, msg, time.Since(start))
	})
}

// UnmarshalV2ForService godoc
// @Summary      unmarshal a protocol message into a value, service from the path
// @Description  like POST /v2/unmarshal, but the service is read from the device-repository by the id in the path
// @Tags         unmarshal, v2
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        serviceId path string true "id of the service the message came from"
// @Param        message body messages.UnmarshallingV2Request true "the message to unmarshal"
// @Success      200 {object} interface{} "the value in the requested characteristic"
// @Failure      400 {string} string "the body is not valid json, the service does not reference the given protocol, or no output path matches the criteria"
// @Failure      500 {string} string
// @Router       /v2/unmarshal/{serviceId} [POST]
func (this *UnmarshallingV2) UnmarshalV2ForService(config config.Config, router *http.ServeMux, ctrl Controller, m *metrics.Metrics) {
	resource := "/v2/unmarshal"
	router.HandleFunc("POST "+resource+"/{serviceId}", func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		msg, ok := decodeBody[messages.UnmarshallingV2Request](writer, request)
		if !ok {
			return
		}
		result, err, code := ctrl.UnmarshalV2ForService(request.PathValue("serviceId"), msg)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writeJson(config, writer, result)
		//the metric label keeps the old httprouter spelling, see marshalv2.go
		m.LogUnmarshallingRequest(request, resource+"/:serviceId", msg, time.Since(start))
	})
}
