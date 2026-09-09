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

	"github.com/SENERGY-Platform/marshaller/lib/api/metrics"
	"github.com/SENERGY-Platform/marshaller/lib/config"
)

func init() {
	endpoints = append(endpoints, &CharacteristicPaths{})
}

type CharacteristicPaths struct{}

// GetCharacteristicPaths godoc
// @Summary      list the paths a characteristic can be read from
// @Description  returns the content-variable paths of the service that carry the given characteristic
// @Tags         characteristic-paths
// @Produce      json
// @Security     Bearer
// @Param        serviceId path string true "id of the service"
// @Param        characteristicId path string true "id of the characteristic"
// @Success      200 {object} marshaller.CharacteristicsPathResponse
// @Failure      400 {string} string
// @Failure      404 {string} string "the service is unknown"
// @Failure      500 {string} string
// @Router       /characteristic-paths/{serviceId}/{characteristicId} [GET]
func (this *CharacteristicPaths) GetCharacteristicPaths(config config.Config, router *http.ServeMux, ctrl Controller, m *metrics.Metrics) {
	router.HandleFunc("GET /characteristic-paths/{serviceId}/{characteristicId}", func(writer http.ResponseWriter, request *http.Request) {
		result, err, code := ctrl.GetCharacteristicPaths(request.PathValue("serviceId"), request.PathValue("characteristicId"))
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writeJson(config, writer, result)
	})
}
