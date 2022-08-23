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

package v2

import (
	convertermodel "github.com/SENERGY-Platform/converter/lib/model"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"sort"
	"strings"
)

func New(config config.Config, converter Converter, characteristics CharacteristicsRepo) *Marshaller {
	return &Marshaller{
		config:          config,
		converter:       converter,
		characteristics: characteristics,
	}
}

type Marshaller struct {
	config          config.Config
	converter       Converter
	characteristics CharacteristicsRepo
}

type CharacteristicsRepo interface {
	GetCharacteristic(id string) (characteristic model.Characteristic, err error)
	GetConcept(id string) (concept model.Concept, err error)
	GetConceptIdOfFunction(id string) string
}

type CharacteristicId = string

type Converter interface {
	Cast(in interface{}, from CharacteristicId, to CharacteristicId) (out interface{}, err error)
	CastWithExtension(in interface{}, from CharacteristicId, to CharacteristicId, extensions []convertermodel.ConverterExtension) (out interface{}, err error)
}

func (this *Marshaller) GetInputPaths(service model.Service, functionId string, aspectNode *model.AspectNode) (result []string) {
	return this.getPathsFromContentsByCriteria(service.Inputs, functionId, aspectNode)
}

func (this *Marshaller) GetOutputPaths(service model.Service, functionId string, aspectNode *model.AspectNode) (result []string) {
	return this.getPathsFromContentsByCriteria(service.Outputs, functionId, aspectNode)
}

func (this *Marshaller) getPathsFromContentsByCriteria(contents []model.Content, functionId string, aspectNode *model.AspectNode) (result []string) {
	withDistance := []pathWithDistance{}
	for _, c := range contents {
		subResults := this.getPathsFromVariableByCriteriaWithDistance(c.ContentVariable, functionId, aspectNode, []string{})
		if len(subResults) > 0 {
			withDistance = append(withDistance, subResults...)
		}
	}
	sort.Slice(withDistance, func(i, j int) bool {
		return withDistance[i].distance < withDistance[j].distance
	})
	for _, element := range withDistance {
		result = append(result, element.path)
	}
	return result
}

type pathWithDistance struct {
	path     string
	distance int
}

func (this *Marshaller) getPathsFromVariableByCriteriaWithDistance(variable model.ContentVariable, functionId string, aspectNode *model.AspectNode, currentPath []string) (result []pathWithDistance) {
	currentPath = append(currentPath, variable.Name)
	result = []pathWithDistance{}
	aspectDistanceLevel := -1
	if aspectNode == nil {
		aspectDistanceLevel = 0
	} else if variable.AspectId == aspectNode.Id {
		aspectDistanceLevel = 0
	} else if contains(aspectNode.ChildIds, variable.AspectId) {
		aspectDistanceLevel = 1
	} else if contains(aspectNode.DescendentIds, variable.AspectId) {
		aspectDistanceLevel = 2
	}
	functionMatches := false
	if functionId == "" || variable.FunctionId == functionId {
		functionMatches = true
	}
	if aspectDistanceLevel > -1 && functionMatches {
		return []pathWithDistance{
			{
				path:     strings.Join(currentPath, "."),
				distance: aspectDistanceLevel,
			},
		}
	}
	for _, sub := range variable.SubContentVariables {
		subResults := this.getPathsFromVariableByCriteriaWithDistance(sub, functionId, aspectNode, currentPath)
		if len(subResults) > 0 {
			result = append(result, subResults...)
		}
	}
	return result
}

func contains(ids []string, id string) bool {
	for _, element := range ids {
		if element == id {
			return true
		}
	}
	return false
}
