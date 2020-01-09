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
	"github.com/SENERGY-Platform/marshaller-service/lib/marshaller/model"
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
	this.mux.Lock()
	defer this.mux.Unlock()

	this.concepts = map[string]model.Concept{model.NullConcept.Id: model.NullConcept}
	this.characteristics = map[string]model.Characteristic{model.NullCharacteristic.Id: model.NullCharacteristic}
	this.conceptByCharacteristic = map[string]model.Concept{}
	this.rootCharacteristicByCharacteristic = map[string]model.Characteristic{}

	for _, element := range temp {
		this.register(element.Concept, element.Characteristics)
	}
	return nil
}

func (this *ConceptRepo) register(concept model.Concept, characteristics []model.Characteristic) {
	for _, characteristic := range characteristics {
		concept.CharacteristicIds = append(concept.CharacteristicIds, characteristic.Id)
		this.characteristics[characteristic.Id] = characteristic
		this.conceptByCharacteristic[characteristic.Id] = concept
		this.rootCharacteristicByCharacteristic[characteristic.Id] = characteristic
		for _, descendent := range getCharacteristicDescendents(characteristic) {
			this.rootCharacteristicByCharacteristic[descendent.Id] = characteristic
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
		err = token.GetJSON(this.config.PermissionsSearchUrl+"/jwt/list/concepts/r/"+strconv.Itoa(limit)+"/"+strconv.Itoa(offset)+"/name/asc", &temp)
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

func (this *ConceptRepo) loadConcept(id string) (result model.Concept, err error) {
	token, err := this.access.Ensure()
	if err != nil {
		return result, err
	}
	err = token.GetJSON(this.config.SemanticRepositoryUrl+"/concepts/"+url.PathEscape(id), &result)
	return
}

func (this *ConceptRepo) loadCharacteristic(id string) (result model.Characteristic, err error) {
	token, err := this.access.Ensure()
	if err != nil {
		return result, err
	}
	err = token.GetJSON(this.config.SemanticRepositoryUrl+"/characteristics/"+url.PathEscape(id), &result)
	return
}
