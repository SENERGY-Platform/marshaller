/*
 * Copyright 2023 InfAI (CC SES)
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
	"net/http"

	"github.com/SENERGY-Platform/marshaller/lib/api/metrics"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/converter"
)

func init() {
	endpoints = append(endpoints, &ConverterExtension{})
}

type ConverterExtension struct{}

// TryConverterExtension godoc
// @Summary      try a converter extension
// @Description  runs a single converter extension against an input and returns its output, for developing an extension
// @Tags         converter
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        call body converter.ExtensionCall true "the extension and the input to run it on"
// @Success      200 {object} converter.ExtensionCallResponse
// @Failure      400 {string} string "the body is not valid json, or the extension failed"
// @Failure      500 {string} string "the service runs without a converter"
// @Router       /converter/extension-call [POST]
func (this *ConverterExtension) TryConverterExtension(config config.Config, router *http.ServeMux, ctrl Controller, m *metrics.Metrics) {
	router.HandleFunc("POST /converter/extension-call", func(writer http.ResponseWriter, request *http.Request) {
		call := converter.ExtensionCall{}
		//this endpoint answers a decode failure with its own wording rather than the
		//decoder's, the way it did before
		if err := json.NewDecoder(request.Body).Decode(&call); err != nil {
			http.Error(writer, "expect valid json in request body", http.StatusBadRequest)
			return
		}
		result, err, code := ctrl.TryConverterExtension(call)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writeJson(config, writer, result)
	})
}
