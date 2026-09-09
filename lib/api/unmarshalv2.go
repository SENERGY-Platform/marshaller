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
