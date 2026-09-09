/*
 * Copyright 2022 InfAI (CC SES)
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
	endpoints = append(endpoints, &MarshallingV2{})
}

type MarshallingV2 struct{}

// MarshalV2 godoc
// @Summary      marshal values into a protocol message
// @Description  marshals every entry of data into the message format of the service; an entry without paths is placed by its function and aspects
// @Tags         marshal, v2
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        message body messages.MarshallingV2Request true "service, protocol and the values to marshal"
// @Success      200 {object} map[string]string "the protocol segments of the message"
// @Failure      400 {string} string "the body is not valid json, or the service does not reference the given protocol"
// @Failure      500 {string} string
// @Router       /v2/marshal [POST]
func (this *MarshallingV2) MarshalV2(config config.Config, router *http.ServeMux, ctrl Controller, m *metrics.Metrics) {
	resource := "/v2/marshal"
	router.HandleFunc("POST "+resource, func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		msg, ok := decodeBody[messages.MarshallingV2Request](writer, request)
		if !ok {
			return
		}
		result, err, code := ctrl.MarshalV2(msg)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writeJson(config, writer, result)
		m.LogMarshallingRequest(request, resource, msg, time.Since(start))
	})
}

// MarshalV2ForService godoc
// @Summary      marshal values into a protocol message, service from the path
// @Description  like POST /v2/marshal, but the service is read from the device-repository by the id in the path
// @Tags         marshal, v2
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        serviceId path string true "id of the service the message is for"
// @Param        message body messages.MarshallingV2Request true "the values to marshal"
// @Success      200 {object} map[string]string "the protocol segments of the message"
// @Failure      400 {string} string "the body is not valid json, or the service does not reference the given protocol"
// @Failure      500 {string} string
// @Router       /v2/marshal/{serviceId} [POST]
func (this *MarshallingV2) MarshalV2ForService(config config.Config, router *http.ServeMux, ctrl Controller, m *metrics.Metrics) {
	resource := "/v2/marshal"
	router.HandleFunc("POST "+resource+"/{serviceId}", func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		msg, ok := decodeBody[messages.MarshallingV2Request](writer, request)
		if !ok {
			return
		}
		result, err, code := ctrl.MarshalV2ForService(request.PathValue("serviceId"), msg)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writeJson(config, writer, result)
		//the metric label keeps the old httprouter spelling on purpose: it is a label
		//value existing dashboards query by, and renaming it would start a new time
		//series and silently empty every panel built on the old one
		m.LogMarshallingRequest(request, resource+"/:serviceId", msg, time.Since(start))
	})
}
