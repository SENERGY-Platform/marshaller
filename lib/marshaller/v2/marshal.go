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
	"errors"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/serialization"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func (this *Marshaller) Marshal(protocol model.Protocol, service model.Service, data []model.MarshallingV2RequestData) (result map[string]string, err error) {
	for _, value := range data {
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

	if variable.CharacteristicId != characteristic {
		variable.Value, err = this.converter.Cast(value, characteristic, variable.CharacteristicId)
		if err != nil {
			return variable, err
		}
	} else {
		variable.Value = value
	}
	return variable, nil
}

func (this *Marshaller) contentsToMessage(protocol model.Protocol, inputs []model.Content) (result map[string]string, err error) {
	result = map[string]string{}
	for _, input := range inputs {
		obj, _, err := contentVariableToObject(input.ContentVariable)
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
	return result, nil
}

func contentVariableToObject(variable model.ContentVariable) (name string, obj interface{}, err error) {
	name = variable.Name
	switch variable.Type {
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
			subName, subObj, err := contentVariableToObject(sub)
			if err != nil {
				return name, obj, err
			}
			temp[subName] = subObj
		}
		return name, temp, nil
	case model.List:
		temp := make([]interface{}, len(variable.SubContentVariables))
		for _, sub := range variable.SubContentVariables {
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
		return name, temp, nil
	default:
		return name, obj, errors.New("unknown variable type:" + variable.Type)
	}
}
