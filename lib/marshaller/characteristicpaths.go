/*
 * Copyright 2020 InfAI (CC SES)
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

package marshaller

import (
	"errors"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"net/http"
)

type CharacteristicsPathResponse struct {
	Path                    string `json:"path"`
	ServiceCharacteristicId string `json:"service_characteristic_id"`
}

var ErrCharacteristicNotFoundInService = errors.New("characteristic not in service")

func (this *Marshaller) GetServiceCharacteristicPath(service model.Service, characteristicId string) (result CharacteristicsPathResponse, err error, code int) {
	matchingServiceCharacteristicId, _, err := this.getMatchingOutputRootCharacteristic(service.Outputs, characteristicId)
	if err == ErrorNoMatchFound {
		return result, ErrCharacteristicNotFoundInService, http.StatusNotFound
	}
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	result.ServiceCharacteristicId = matchingServiceCharacteristicId
	result.Path, err = getPathOfCharacteristic(service.Outputs, result.ServiceCharacteristicId)
	if err == ErrorNoMatchFound {
		return result, ErrCharacteristicNotFoundInService, http.StatusNotFound
	}
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	return result, nil, http.StatusOK
}

func getPathOfCharacteristic(outputs []model.Content, characteristicId string) (path string, err error) {
	for _, output := range outputs {
		path, found := getPathOfCharacteristicInContentVariable(output.ContentVariable, characteristicId)
		if found {
			return path, nil
		}
	}
	return path, ErrCharacteristicNotFoundInService
}

func getPathOfCharacteristicInContentVariable(variable model.ContentVariable, characteristicId string) (path string, found bool) {
	if variable.CharacteristicId == characteristicId {
		return variable.Name, true
	}
	for _, sub := range variable.SubContentVariables {
		subPath, found := getPathOfCharacteristicInContentVariable(sub, characteristicId)
		if found {
			path = variable.Name + "." + subPath
			return path, true
		}
	}
	return "", false
}
