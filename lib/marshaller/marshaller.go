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
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
)

type Converter interface {
	Cast(in interface{}, from CharacteristicId, to CharacteristicId) (out interface{}, err error)
}

type ConceptRepo interface {
	GetConceptOfCharacteristic(characteristicId string) (conceptId string, err error)
	GetCharacteristic(id CharacteristicId) (model.Characteristic, error)
	GetRootCharacteristics(ids []CharacteristicId) (result []CharacteristicId)
	GetCharacteristicsOfFunction(functionId string) (characteristicIds []string, err error)
}

type DeviceRepository interface {
	GetDeviceType(id string) (result model.DeviceType, err error, code int)
}

type Marshaller struct {
	converter   Converter
	ConceptRepo ConceptRepo
	devicerepo  DeviceRepository
}

func New(converter Converter, concepts ConceptRepo, devicerepo DeviceRepository) *Marshaller {
	return &Marshaller{
		converter:   converter,
		ConceptRepo: concepts,
		devicerepo:  devicerepo,
	}
}
