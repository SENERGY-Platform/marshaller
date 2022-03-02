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
	"strings"
)

func New(config config.Config, converter Converter) *Marshaller {
	return &Marshaller{
		config:    config,
		converter: converter,
	}
}

type Marshaller struct {
	config    config.Config
	converter Converter
}

type CharacteristicId = string

type Converter interface {
	Cast(in interface{}, from CharacteristicId, to CharacteristicId) (out interface{}, err error)
}

func (this *Marshaller) GetInputPaths(service model.Service, functionId string, aspectNode *model.AspectNode) (result []string) {
	return this.getPathsFromContentsByCriteria(service.Inputs, functionId, aspectNode)
}

func (this *Marshaller) GetOutputPaths(service model.Service, functionId string, aspectNode *model.AspectNode) (result []string) {
	return this.getPathsFromContentsByCriteria(service.Outputs, functionId, aspectNode)
}

func (this *Marshaller) getPathsFromContentsByCriteria(contents []model.Content, functionId string, aspectNode *model.AspectNode) (result []string) {
	result = []string{}
	for _, c := range contents {
		subResults := this.getPathsFromVariableByCriteria(c.ContentVariable, functionId, aspectNode, []string{})
		if len(subResults) > 0 {
			result = append(result, subResults...)
		}
	}
	return result
}

func (this *Marshaller) getPathsFromVariableByCriteria(variable model.ContentVariable, functionId string, aspectNode *model.AspectNode, currentPath []string) (result []string) {
	currentPath = append(currentPath, variable.Name)
	result = []string{}
	aspectMatches := false
	if aspectNode == nil || variable.AspectId == aspectNode.Id || contains(aspectNode.DescendentIds, variable.AspectId) {
		aspectMatches = true
	}
	functionMatches := false
	if functionId == "" || variable.FunctionId == functionId {
		functionMatches = true
	}
	if aspectMatches && functionMatches {
		return []string{strings.Join(currentPath, ".")}
	}
	for _, sub := range variable.SubContentVariables {
		subResults := this.getPathsFromVariableByCriteria(sub, functionId, aspectNode, currentPath)
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
