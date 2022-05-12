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
	convertermodel "github.com/SENERGY-Platform/converter/lib/model"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	v2 "github.com/SENERGY-Platform/marshaller/lib/marshaller/v2"
)

type Converter struct{}

func (this Converter) CastWithExtension(in interface{}, from v2.CharacteristicId, to v2.CharacteristicId, extensions []model.ConverterExtensions) (out interface{}, err error) {
	converter, err := converterService.New()
	if err != nil {
		return nil, err
	}
	targetExtension := []convertermodel.ConverterExtension{}
	for _, e := range extensions {
		targetExtension = append(targetExtension, convertermodel.ConverterExtension{
			From:            e.F,
			To:              e.To,
			Distance:        e.Distance,
			F:               e.F,
			PlaceholderName: e.PlaceholderName,
		})
	}
	return converter.CastWithExtension(in, from, to, targetExtension)
}

func (this Converter) Cast(in interface{}, from marshaller.CharacteristicId, to marshaller.CharacteristicId) (out interface{}, err error) {
	converter, err := converterService.New()
	if err != nil {
		return nil, err
	}
	return converter.Cast(in, from, to)
}
