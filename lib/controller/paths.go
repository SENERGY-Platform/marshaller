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

package controller

import (
	"errors"
	"net/http"

	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
)

func (this *Controller) GetCharacteristicPaths(serviceId string, characteristicId string) (result marshaller.CharacteristicsPathResponse, err error, code int) {
	if serviceId == "" {
		return result, errors.New("expect serviceId as parameter in path"), http.StatusBadRequest
	}
	if characteristicId == "" {
		return result, errors.New("expect characteristicId as parameter in path"), http.StatusBadRequest
	}
	service, err, code := this.deviceRepo.GetServiceWithErrCode(serviceId)
	if err != nil {
		return result, err, code
	}
	return this.marshaller.GetServiceCharacteristicPath(service, characteristicId)
}

// GetPathOptions answers both the query-parameter and the request-body form of the
// path-options endpoint.
//
// An empty CharacteristicIdFilter means "do not filter" and yields every characteristic of
// the function — it is not the same as an absent characteristic-filter query parameter,
// which the GET endpoint answers with an empty result before it gets here.
func (this *Controller) GetPathOptions(query messages.PathOptionsQuery) (result map[string][]marshaller.PathOptionsResultElement, err error, code int) {
	return this.marshaller.GetPathOption(query.DeviceTypeIds, query.FunctionId, query.GetAspectIds(), query.CharacteristicIdFilter, !query.WithoutEnvelope)
}
