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
	"encoding/json"
	"errors"
	convertermodel "github.com/SENERGY-Platform/converter/lib/model"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/serialization"
	"log"
	"math"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
)

func (this *Marshaller) Marshal(protocol model.Protocol, service model.Service, data []model.MarshallingV2RequestData) (result map[string]string, err error) {
	for _, value := range data {
		if len(value.Paths) == 0 && value.FunctionId != "" {
			value.Paths = this.GetInputPaths(service, value.FunctionId, value.AspectNode)
		}
		service.Inputs, err = this.setContentVariableValues(service.Inputs, value.Paths, value.CharacteristicId, value.Value)
		if err != nil {
			return result, err
		}
	}
	return this.contentsToMessage(protocol, service.Inputs)
}

func (this *Marshaller) setContentVariableValues(inputs []model.Content, paths []string, characteristic string, value interface{}) (result []model.Content, err error) {
	for _, path := range paths {
		pathParts := strings.Split(path, ".")
		for _, input := range inputs {
			input.ContentVariable, err = this.setContentVariableValue(input.ContentVariable, []string{}, pathParts, characteristic, value)
			if err != nil {
				return result, err
			}
			result = append(result, input)
		}
	}
	return result, nil
}

func (this *Marshaller) setContentVariableValue(variable model.ContentVariable, currentPath []string, pathParts []string, characteristic string, value interface{}) (model.ContentVariable, error) {
	currentPath = append(currentPath, variable.Name)
	index := len(currentPath) - 1
	//searched path is shorter than current path
	if len(currentPath) > len(pathParts) {
		return variable, nil
	}
	//current path segment is different from searched
	if currentPath[index] != pathParts[index] {
		return variable, nil
	}
	var err error

	//not at the end of the correct path
	if len(currentPath) < len(pathParts) {
		for i, sub := range variable.SubContentVariables {
			variable.SubContentVariables[i], err = this.setContentVariableValue(sub, currentPath, pathParts, characteristic, value)
			if err != nil {
				return variable, err
			}
		}
		return variable, nil
	}

	//should never happen but check to be sure
	if !reflect.DeepEqual(currentPath, pathParts) {
		return variable, errors.New("wtf")
	}

	if characteristic != "" && variable.CharacteristicId != characteristic {
		castExtensions := []convertermodel.ConverterExtension{}
		if variable.FunctionId != "" {
			conceptId := this.characteristics.GetConceptIdOfFunction(variable.FunctionId)
			if conceptId != "" {
				concept, err := this.characteristics.GetConcept(conceptId)
				if err != nil {
					return variable, err
				}
				castExtensions = concept.Conversions
			}
		}
		if len(castExtensions) == 0 {
			variable.Value, err = this.converter.Cast(value, characteristic, variable.CharacteristicId)
		} else {
			variable.Value, err = this.converter.CastWithExtension(value, characteristic, variable.CharacteristicId, castExtensions)
		}
		if err != nil {
			return variable, err
		}
		if variable.Type == model.Integer {
			switch v := variable.Value.(type) {
			case float64:
				variable.Value = int64(math.Round(v))
			case float32:
				variable.Value = int64(math.Round(float64(v)))
			}
		}
	} else {
		variable.Value = value
	}

	//handle complex variables/data-structures like rgb
	if variable.CharacteristicId != "" && (variable.Type == model.Structure || variable.Type == model.List) {
		characteristicToVariablePath := getCharacteristicToPath(variable, []string{})
		targetCharacteristic, err := this.characteristics.GetCharacteristic(variable.CharacteristicId)
		if err != nil {
			return variable, err
		}
		normalizedValue, err := normalize(variable.Value)
		if err != nil {
			return variable, err
		}
		characteristicPathToValue := this.getPathToValueMapFromObj([]string{targetCharacteristic.Name}, normalizedValue)
		variablePathToCharacteristicsValue := getVariablePathToCharacteristicsValue(targetCharacteristic, []string{}, characteristicToVariablePath, characteristicPathToValue)
		for subPath, subValue := range variablePathToCharacteristicsValue {
			subPathParts := strings.Split(subPath, ".")
			if !reflect.DeepEqual(currentPath, subPathParts) {
				variable, err = this.setContentVariableValue(variable, []string{}, subPathParts, "", subValue)
				if err != nil {
					return variable, err
				}
			}
		}
		variable.Value = nil
	}

	return variable, nil
}

func normalize(value interface{}) (result interface{}, err error) {
	temp, err := json.Marshal(value)
	if err != nil {
		debug.PrintStack()
		return result, err
	}
	err = json.Unmarshal(temp, &result)
	if err != nil {
		debug.PrintStack()
		return result, err
	}
	return
}

func getVariablePathToCharacteristicsValue(characteristic model.Characteristic, currentPath []string, characteristicToVariablePath map[string]string, characteristicPathToValue map[string]interface{}) (result map[string]interface{}) {
	result = map[string]interface{}{}
	currentPath = append(currentPath, characteristic.Name)
	if path, ok := characteristicToVariablePath[characteristic.Id]; ok {
		if value, ok := characteristicPathToValue[strings.Join(currentPath, ".")]; ok {
			result[path] = value
		}
	}
	for _, sub := range characteristic.SubCharacteristics {
		for subPath, subValue := range getVariablePathToCharacteristicsValue(sub, currentPath, characteristicToVariablePath, characteristicPathToValue) {
			result[subPath] = subValue
		}
	}
	return result
}

func getCharacteristicToPath(variable model.ContentVariable, currentPath []string) (result map[string]string) {
	result = map[string]string{}
	currentPath = append(currentPath, variable.Name)
	if variable.CharacteristicId != "" {
		result[variable.CharacteristicId] = strings.Join(currentPath, ".")
	}
	for _, sub := range variable.SubContentVariables {
		for characteristic, path := range getCharacteristicToPath(sub, currentPath) {
			result[characteristic] = path
		}
	}
	return result
}

func (this *Marshaller) contentsToMessage(protocol model.Protocol, inputs []model.Content) (result map[string]string, err error) {
	result = map[string]string{}
	for _, input := range inputs {
		if !input.ContentVariable.IsVoid {
			_, obj, err := contentVariableToObject(input.ContentVariable)
			if err != nil {
				return result, err
			}
			s, ok := serialization.Get(input.Serialization)
			if !ok {
				return result, errors.New("unknown serialization " + input.Serialization)
			}
			segmentName := ""
			for _, segment := range protocol.ProtocolSegments {
				if segment.Id == input.ProtocolSegmentId {
					segmentName = segment.Name
					break
				}
			}
			if segmentName != "" {
				result[segmentName], err = s.Marshal(obj, input.ContentVariable)
				if err != nil {
					return result, err
				}
			} else {
				log.Println("WARNING: protocol-segment not found " + input.ProtocolSegmentId)
			}
		}
	}
	return result, nil
}

func contentVariableToObject(variable model.ContentVariable) (name string, obj interface{}, err error) {
	name = variable.Name
	switch variable.Type {
	case "":
		return name, variable.Value, nil
	case model.String:
		return name, variable.Value, nil
	case model.Boolean:
		return name, variable.Value, nil
	case model.Integer:
		return name, variable.Value, nil
	case model.Float:
		return name, variable.Value, nil
	case model.Structure:
		temp := map[string]interface{}{}
		for _, sub := range variable.SubContentVariables {
			if !sub.IsVoid {
				subName, subObj, err := contentVariableToObject(sub)
				if err != nil {
					return name, obj, err
				}
				temp[subName] = subObj
			}
		}
		return name, temp, nil
	case model.List:
		temp := make([]interface{}, len(variable.SubContentVariables))
		for _, sub := range variable.SubContentVariables {
			if !sub.IsVoid {
				subName, subObj, err := contentVariableToObject(sub)
				if err != nil {
					return name, obj, err
				}
				if subName == "*" {
					if len(temp) != 1 {
						return name, obj, errors.New("expect * only on list with one element")
					}
					temp[0] = subObj
					return name, temp, nil
				}
				index, err := strconv.Atoi(subName)
				if err != nil {
					return name, obj, errors.New("unable to marshal list with index " + subName + " in " + variable.Name + " " + variable.Id + ": " + err.Error())
				}
				temp[index] = subObj
			}
		}
		return name, temp, nil
	default:
		return name, obj, errors.New("unknown variable type:" + string(variable.Type))
	}
}
