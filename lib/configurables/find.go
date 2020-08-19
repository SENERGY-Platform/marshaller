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

package configurables

import (
	"errors"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/mapping"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"strings"
)

type ConceptRepo interface {
	GetConceptOfCharacteristic(characteristicId string) (conceptId string, err error)
	GetCharacteristic(id string) (model.Characteristic, error)
	GetConcept(id string) (concept model.Concept, err error)
}

type ConfigurableService struct {
	conceptrepo ConceptRepo
}

func New(repo ConceptRepo) *ConfigurableService {
	return &ConfigurableService{conceptrepo: repo}
}

func (this *ConfigurableService) Find(notCharacteristicId string, services []model.Service) (result Configurables, err error) {
	return FindIntersectingConfigurables(this.conceptrepo, notCharacteristicId, services)
}

func FindIntersectingConfigurables(repo ConceptRepo, notCharacteristicId string, services []model.Service) (result Configurables, err error) {
	notConcept, err := repo.GetConceptOfCharacteristic(notCharacteristicId)
	if err != nil {
		return nil, err
	}
	invertedIndex := map[string]map[string]bool{}
	for _, service := range services {
		characteristics := []string{}
		for _, content := range service.Inputs {
			characteristics = append(characteristics, characteristicsInContentVariable(content.ContentVariable)...)
		}
		for _, characteristic := range characteristics {
			concept, err := repo.GetConceptOfCharacteristic(characteristic)
			if err != nil {
				return nil, err
			}
			if _, ok := invertedIndex[concept]; !ok {
				invertedIndex[concept] = map[string]bool{}
			}
			invertedIndex[concept][service.Id] = true
		}
	}
	serviceCount := len(services)
	for conceptId, servicesUsingConcept := range invertedIndex {
		if conceptId != notConcept && len(servicesUsingConcept) == serviceCount {
			configurable, err := createConfigurable(repo, conceptId)
			if err != nil {
				return nil, err
			}
			result = append(result, configurable)
		}
	}
	return result, nil
}

func createConfigurable(repo ConceptRepo, conceptId string) (result Configurable, err error) {
	concept, err := repo.GetConcept(conceptId)
	if err != nil {
		return result, err
	}
	if len(concept.CharacteristicIds) == 0 {
		return result, errors.New("expect at least one characteristic for concept " + concept.Name + " " + conceptId)
	}
	characteristicId := concept.BaseCharacteristicId
	if characteristicId == "" {
		characteristicId = concept.CharacteristicIds[0]
	}
	characteristic, err := repo.GetCharacteristic(characteristicId)
	if err != nil {
		return result, err
	}
	return Configurable{
		CharacteristicId: concept.BaseCharacteristicId,
		Values:           createConfigurableValues(characteristic, characteristic.Name),
	}, nil
}

func createConfigurableValues(characteristic model.Characteristic, labelPrefix string, pathSegments ...string) (result []ConfigurableCharacteristicValue) {
	switch characteristic.Type {
	case model.Integer, model.Float:
		return []ConfigurableCharacteristicValue{
			{
				Label: strings.Join(append([]string{labelPrefix}, pathSegments...), " "),
				Path:  strings.Join(pathSegments, "."),
				Value: "0",
			},
		}
	case model.Boolean:
		return []ConfigurableCharacteristicValue{
			{
				Label: strings.Join(append([]string{labelPrefix}, pathSegments...), " "),
				Path:  strings.Join(pathSegments, "."),
				Value: "false",
			},
		}
	case model.String:
		return []ConfigurableCharacteristicValue{
			{
				Label: strings.Join(append([]string{labelPrefix}, pathSegments...), " "),
				Path:  strings.Join(pathSegments, "."),
				Value: "",
			},
		}
	case model.List:
		for _, sub := range characteristic.SubCharacteristics {
			index := "0"
			if sub.Name != mapping.VAR_LEN_PLACEHOLDER {
				index = sub.Name
			}
			result = append(result, createConfigurableValues(sub, labelPrefix, append(pathSegments, index)...)...)
		}
		return result
	case model.Structure:
		for _, sub := range characteristic.SubCharacteristics {
			result = append(result, createConfigurableValues(sub, labelPrefix, append(pathSegments, sub.Name)...)...)
		}
		return result
	default:
		return
	}
}

func characteristicsInContentVariable(variable model.ContentVariable) (result []string) {
	if variable.CharacteristicId != "" {
		result = append(result, variable.CharacteristicId)
	}
	for _, sub := range variable.SubContentVariables {
		result = append(result, characteristicsInContentVariable(sub)...)
	}
	return result
}
