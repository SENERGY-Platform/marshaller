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

package conceptrepo

import (
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"log"
	"net/url"
	"strconv"
)

func (this *ConceptRepo) Load() error {
	conceptIds, err := this.loadConceptIds()
	if err != nil {
		return err
	}

	type Temp struct {
		Concept         model.Concept
		Characteristics []model.Characteristic
	}
	temp := []Temp{}

	for _, conceptId := range conceptIds {
		concept, err := this.loadConcept(conceptId)
		if err != nil {
			return err
		}
		element := Temp{
			Concept: concept,
		}
		for _, characteristicId := range concept.CharacteristicIds {
			characteristic, err := this.loadCharacteristic(characteristicId)
			if err != nil {
				return err
			}
			element.Characteristics = append(element.Characteristics, characteristic)
		}
		temp = append(temp, element)
	}

	functionInfos, err := this.loadFunctions()
	if err != nil {
		return err
	}

	this.mux.Lock()
	defer this.mux.Unlock()

	this.resetToDefault()

	for _, element := range temp {
		this.register(element.Concept, element.Characteristics)
	}

	for _, f := range functionInfos {
		this.registerFunction(f)
	}

	return nil
}

func (this *ConceptRepo) resetToDefault() {
	this.concepts = map[string]model.Concept{}
	this.characteristics = map[string]model.Characteristic{}
	this.conceptByCharacteristic = map[string][]model.Concept{}
	this.rootCharacteristicByCharacteristic = map[string]model.Characteristic{}
	for _, defaultElement := range this.defaults {
		this.register(defaultElement.Concept, defaultElement.Characteristics)
	}
}

func (this *ConceptRepo) register(concept model.Concept, characteristics []model.Characteristic) {
	log.Println("load concept", concept.Name, concept.Id)
	for _, characteristic := range characteristics {
		log.Println("    load characteristic", characteristic.Name, characteristic.Id)
		concept.CharacteristicIds = append(concept.CharacteristicIds, characteristic.Id)
		this.characteristics[characteristic.Id] = characteristic
		this.conceptByCharacteristic[characteristic.Id] = append(this.conceptByCharacteristic[characteristic.Id], concept)
		this.rootCharacteristicByCharacteristic[characteristic.Id] = characteristic
		for _, descendent := range getCharacteristicDescendents(characteristic) {
			this.rootCharacteristicByCharacteristic[descendent.Id] = characteristic
			this.characteristics[descendent.Id] = descendent
		}
	}
	this.concepts[concept.Id] = concept
}

type IdWrapper struct {
	Id string `json:"id"`
}

func (this *ConceptRepo) loadConceptIds() (ids []string, err error) {
	token, err := this.access.Ensure()
	if err != nil {
		return ids, err
	}
	limit := 100
	offset := 0
	temp := []IdWrapper{}
	for len(temp) == limit || offset == 0 {
		temp = []IdWrapper{}
		endpoint := this.config.PermissionsSearchUrl + "/v3/resources/concepts?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset) + "&sort=name.asc&rights=r"
		err = token.GetJSON(endpoint, &temp)
		if err != nil {
			return ids, err
		}
		for _, wrapper := range temp {
			ids = append(ids, wrapper.Id)
		}
		offset = offset + limit
	}
	return ids, err
}

type FunctionInfo struct {
	Id        string `json:"id"`
	ConceptId string `json:"concept_id"`
}

func (this *ConceptRepo) loadFunctions() (functionInfos []FunctionInfo, err error) {
	log.Println("load functions")
	token, err := this.access.Ensure()
	if err != nil {
		return functionInfos, err
	}
	limit := 100
	offset := 0
	temp := []FunctionInfo{}
	for len(temp) == limit || offset == 0 {
		temp = []FunctionInfo{}
		endpoint := this.config.PermissionsSearchUrl + "/v3/resources/functions?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset) + "&sort=name.asc&rights=r"
		err = token.GetJSON(endpoint, &temp)
		if err != nil {
			return functionInfos, err
		}
		functionInfos = append(functionInfos, temp...)
		offset = offset + limit
	}
	return functionInfos, err
}

func (this *ConceptRepo) loadConcept(id string) (result model.Concept, err error) {
	token, err := this.access.Ensure()
	if err != nil {
		return result, err
	}
	err = token.GetJSON(this.config.DeviceRepositoryUrl+"/concepts/"+url.PathEscape(id), &result)
	return
}

func (this *ConceptRepo) loadCharacteristic(id string) (result model.Characteristic, err error) {
	token, err := this.access.Ensure()
	if err != nil {
		return result, err
	}
	err = token.GetJSON(this.config.DeviceRepositoryUrl+"/characteristics/"+url.PathEscape(id), &result)
	return
}
