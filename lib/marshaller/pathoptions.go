/*
 * Copyright 2021 InfAI (CC SES)
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
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"net/http"
	"strings"
)

type PathOptionsResultElement struct {
	ServiceId string   `json:"service_id"`
	JsonPath  []string `json:"json_path"`
}

func (this *Marshaller) GetPathOption(deviceTypeIds []string, functionId string, aspectId string, characteristicIdFilter []string, withEnvelope bool) (result map[string][]PathOptionsResultElement, err error, code int) {
	result = map[string][]PathOptionsResultElement{}
	for _, deviceTypeId := range deviceTypeIds {
		result[deviceTypeId], err, code = this.getPathOptionForDeviceType(deviceTypeId, functionId, aspectId, characteristicIdFilter, withEnvelope)
		if err != nil {
			return
		}
	}
	return result, nil, http.StatusOK
}

func (this *Marshaller) getPathOptionForDeviceType(deviceTypeId string, functionId string, aspectId string, characteristicIdFilter []string, withEnvelope bool) (result []PathOptionsResultElement, err error, code int) {
	result = []PathOptionsResultElement{}
	dt, err, code := this.devicerepo.GetDeviceType(deviceTypeId)
	if err != nil {
		return nil, err, code
	}
	services := this.filterMatchingServices(dt.Services, functionId, aspectId)
	for _, service := range services {
		paths, err, code := this.getPathOptionsForService(service, functionId, characteristicIdFilter)
		if err != nil {
			return result, err, code
		}
		if withEnvelope {
			for i, e := range paths {
				paths[i] = "value." + e
			}
		}
		if len(paths) > 0 {
			result = append(result, PathOptionsResultElement{
				ServiceId: service.Id,
				JsonPath:  paths,
			})
		}
	}
	return result, nil, http.StatusOK
}

func (this *Marshaller) getPathOptionsForService(service model.Service, functionId string, characteristicIdFilter []string) (paths []string, err error, code int) {
	characteristics, err := this.ConceptRepo.GetCharacteristicsOfFunction(functionId)
	if err != nil {
		return paths, err, http.StatusInternalServerError
	}
	if len(characteristicIdFilter) > 0 {
		characteristics = intersection(characteristics, characteristicIdFilter)
	}
	characteristicsSet := map[string]bool{}
	for _, c := range characteristics {
		characteristicsSet[c] = true
	}
	paths = this.findPathsOfCharacteristicsInContents(service.Outputs, characteristicsSet)
	return paths, nil, http.StatusOK
}

func (this *Marshaller) findPathsOfCharacteristicsInContents(contents []model.Content, characteristics map[string]bool) (result []string) {
	for _, c := range contents {
		result = append(result, this.findPathsInContentVariable([]string{}, c.ContentVariable, characteristics, "{{NAME}}")...)
	}
	return
}

func (this *Marshaller) findPathsInContentVariable(currentPath []string, variable model.ContentVariable, characteristics map[string]bool, pathSegmentPattern string) (result []string) {
	currentPath = append(currentPath, strings.ReplaceAll(pathSegmentPattern, "{{NAME}}", variable.Name))

	if characteristics[variable.CharacteristicId] {
		return []string{strings.Join(currentPath, "")}
	}

	nextPattern := ".{{NAME}}"
	if variable.Type == model.List {
		nextPattern = "[{{NAME}}]"
	}
	for _, v := range variable.SubContentVariables {
		result = append(result, this.findPathsInContentVariable(currentPath, v, characteristics, nextPattern)...)
	}
	return result
}

func intersection(alist []string, blist []string) (result []string) {
	if len(alist) == 0 || len(blist) == 0 {
		return []string{}
	}
	aindex := map[string]bool{}
	for _, a := range alist {
		aindex[a] = true
	}
	for _, b := range blist {
		if aindex[b] {
			result = append(result, b)
		}
	}
	return
}

func (this *Marshaller) filterMatchingServices(services []model.Service, functionId string, aspectId string) (result []model.Service) {
	for _, service := range services {
		if contains(service.FunctionIds, functionId) && contains(service.AspectIds, aspectId) {
			result = append(result, service)
		}
	}
	return
}

func contains(ids []string, id string) bool {
	for _, element := range ids {
		if element == id {
			return true
		}
	}
	return false
}
