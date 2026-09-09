/*
 * Copyright 2026 InfAI (CC SES)
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

	"github.com/SENERGY-Platform/marshaller/lib/config"
)

// writeJson sends a successful result. An encoding failure is logged rather than returned:
// the status line is already on the wire at that point.
func writeJson(config config.Config, writer http.ResponseWriter, result interface{}) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(writer).Encode(result)
	if err != nil {
		config.GetLogger().Error("unable to encode response", "error", err)
	}
}

// decodeBody reads a json request body. It answers 400 itself and reports whether the
// caller may continue, because every endpoint treats an undecodable body the same way.
func decodeBody[T any](writer http.ResponseWriter, request *http.Request) (result T, ok bool) {
	err := json.NewDecoder(request.Body).Decode(&result)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return result, false
	}
	return result, true
}
