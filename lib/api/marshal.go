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
