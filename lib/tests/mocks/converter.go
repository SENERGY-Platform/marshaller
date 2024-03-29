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

package mocks

import (
	converterService "github.com/SENERGY-Platform/converter/lib/converter"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	v2 "github.com/SENERGY-Platform/marshaller/lib/marshaller/v2"
	"github.com/SENERGY-Platform/models/go/models"
)

type Converter struct{}

func (this Converter) CastWithExtension(in interface{}, from v2.CharacteristicId, to v2.CharacteristicId, extensions []models.ConverterExtension) (out interface{}, err error) {
	converter, err := converterService.New()
	if err != nil {
		return nil, err
	}
	return converter.CastWithExtension(in, from, to, extensions)
}

func (this Converter) Cast(in interface{}, from marshaller.CharacteristicId, to marshaller.CharacteristicId) (out interface{}, err error) {
	converter, err := converterService.New()
	if err != nil {
		return nil, err
	}
	return converter.Cast(in, from, to)
}
