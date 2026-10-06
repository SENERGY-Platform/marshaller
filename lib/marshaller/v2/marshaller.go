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
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/models/go/models"
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
	CastWithExtension(in interface{}, from CharacteristicId, to CharacteristicId, extensions []models.ConverterExtension) (out interface{}, err error)
}

// InputPathSelectionMode decides which of the input paths matching the function and the
// aspects of a marshal request receive its value, when the request names no path itself.
type InputPathSelectionMode int

const (
	// ClosestInputPaths writes only the matching paths whose aspects lie closest to the
	// queried aspects in the aspect tree, like an unmarshal reads only the closest output.
	// Paths at the same distance are all written, so a request without aspects still writes
	// every variable of its function.
	ClosestInputPaths InputPathSelectionMode = iota
	// AllInputPaths writes every matching path, including the ones whose aspects only
	// descend from a queried aspect: a request for air sets inside_air and outside_air alike.
	AllInputPaths
)

// InputPathSelection is the mode Marshal uses. It is a variable and not a constant because
// controlling functions only recently got combined with aspects, and which paths a command
// for a parent aspect should set may still change.
var InputPathSelection = ClosestInputPaths

func (this *Marshaller) GetInputPaths(service model.Service, functionId string, aspectNodes []model.AspectNode) (result []string) {
	return this.getPathsFromContentsByCriteria(service.Inputs, functionId, aspectNodes)
}

// selectInputPaths returns the input paths Marshal writes for a request without paths,
// according to InputPathSelection.
func (this *Marshaller) selectInputPaths(service model.Service, functionId string, aspectNodes []model.AspectNode) (result []string) {
	withDistance := this.getPathsWithDistanceFromContentsByCriteria(service.Inputs, functionId, aspectNodes)
	for _, element := range withDistance {
		if InputPathSelection == ClosestInputPaths && element.distance > withDistance[0].distance {
			break
		}
		result = append(result, element.path)
	}
	return result
}

func (this *Marshaller) GetOutputPaths(service model.Service, functionId string, aspectNodes []model.AspectNode) (result []string) {
	return this.getPathsFromContentsByCriteria(service.Outputs, functionId, aspectNodes)
}

func (this *Marshaller) getPathsFromContentsByCriteria(contents []model.Content, functionId string, aspectNodes []model.AspectNode) (result []string) {
	for _, element := range this.getPathsWithDistanceFromContentsByCriteria(contents, functionId, aspectNodes) {
		result = append(result, element.path)
	}
	return result
}

// getPathsWithDistanceFromContentsByCriteria returns the matching paths sorted by their
// aspect distance, closest first.
func (this *Marshaller) getPathsWithDistanceFromContentsByCriteria(contents []model.Content, functionId string, aspectNodes []model.AspectNode) (withDistance []pathWithDistance) {
	withDistance = []pathWithDistance{}
	for _, c := range contents {
		subResults := this.getPathsFromVariableByCriteriaWithDistance(c.ContentVariable, functionId, aspectNodes, []string{})
		if len(subResults) > 0 {
			withDistance = append(withDistance, subResults...)
		}
	}
	sort.SliceStable(withDistance, func(i, j int) bool {
		return withDistance[i].distance < withDistance[j].distance
	})
	return withDistance
}

type pathWithDistance struct {
	path     string
	distance int
}

func (this *Marshaller) getPathsFromVariableByCriteriaWithDistance(variable model.ContentVariable, functionId string, aspectNodes []model.AspectNode, currentPath []string) (result []pathWithDistance) {
	currentPath = append(currentPath, variable.Name)
	result = []pathWithDistance{}
	aspectDistanceLevel := model.AspectMatchLevel(model.ContentVariableAspectIds(variable), aspectNodes)
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
		subResults := this.getPathsFromVariableByCriteriaWithDistance(sub, functionId, aspectNodes, currentPath)
		if len(subResults) > 0 {
			result = append(result, subResults...)
		}
	}
	return result
}
