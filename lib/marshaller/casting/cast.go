/*
 * Copyright 2019 InfAI (CC SES)
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

package casting

import (
	"errors"
	"github.com/SENERGY-Platform/marshaller-service/lib/marshaller/casting/base"
	_ "github.com/SENERGY-Platform/marshaller-service/lib/marshaller/casting/battery"
	_ "github.com/SENERGY-Platform/marshaller-service/lib/marshaller/casting/binary"
	_ "github.com/SENERGY-Platform/marshaller-service/lib/marshaller/casting/brightness"
	_ "github.com/SENERGY-Platform/marshaller-service/lib/marshaller/casting/color"
	_ "github.com/SENERGY-Platform/marshaller-service/lib/marshaller/casting/energy"
	_ "github.com/SENERGY-Platform/marshaller-service/lib/marshaller/casting/example"
	_ "github.com/SENERGY-Platform/marshaller-service/lib/marshaller/casting/humidity"
	_ "github.com/SENERGY-Platform/marshaller-service/lib/marshaller/casting/luminiscence"
	_ "github.com/SENERGY-Platform/marshaller-service/lib/marshaller/casting/particleamount"
	_ "github.com/SENERGY-Platform/marshaller-service/lib/marshaller/casting/power"
	_ "github.com/SENERGY-Platform/marshaller-service/lib/marshaller/casting/temperature"
	_ "github.com/SENERGY-Platform/marshaller-service/lib/marshaller/casting/time"
	_ "github.com/SENERGY-Platform/marshaller-service/lib/marshaller/casting/ultraviolet"
	"github.com/SENERGY-Platform/marshaller-service/lib/marshaller/model"
	"runtime/debug"
)

var ConceptRepo = base.ConceptRepo

func Cast(in interface{}, conceptId string, from string, to string) (out interface{}, err error) {
	if from == model.NullCharacteristic.Id || to == model.NullCharacteristic.Id || conceptId == model.NullConcept.Id {
		return in, nil
	}
	return Concepts(conceptId)(from)(in)(to)
}

func Concepts(conceptId string) base.FindCastFromCharacteristicToConceptFunction {
	result, ok := base.Concepts[conceptId]
	if !ok {
		debug.PrintStack()
		return base.GetErrorFindCastFromCharacteristicToConceptFunction(errors.New("concept not found"))
	}
	return result
}
