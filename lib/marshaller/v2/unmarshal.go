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
	"runtime/debug"
	"strconv"
	"strings"
)

func (this *Marshaller) Unmarshal(protocol model.Protocol, service model.Service, characteristicId string, path string, msg map[string]string) (result interface{}, err error) {
	outputObjectMap, err := serializeOutput(msg, service, protocol)
	if err != nil {
		return result, err
	}
	pathToValue := this.getPathToValueMapFromObj([]string{}, outputObjectMap)
	value, ok := pathToValue[path]
	if !ok {
		if this.config.ReturnUnknownPathAsNull {
			return nil, nil
		}
		return result, errors.New("path not found in message")
	}

	//no conversion wanted
	if characteristicId == "" {
		return value, nil
	}

	//convert service/variable characteristic to wanted characteristic
	pathToCharacteristic := this.getPathToCharacteristicFromContents(service.Outputs)
	variableCharacteristic, ok := pathToCharacteristic[path]
	if !ok {
		return result, errors.New("path not found in service")
	}
	if variableCharacteristic == characteristicId {
		return value, nil
	} else {
		return this.converter.Cast(value, variableCharacteristic, characteristicId)
	}
}

func (this Marshaller) getPathToValueMapFromObj(startPath []string, value interface{}) (result map[string]interface{}) {
	result = map[string]interface{}{
		strings.Join(startPath, "."): value,
	}
	switch v := value.(type) {
	case map[string]interface{}:
		for k, sub := range v {
			for subPath, subValue := range this.getPathToValueMapFromObj(append(startPath, k), sub) {
				result[subPath] = subValue
			}
		}
	case []interface{}:
		for k, sub := range v {
			for subPath, subValue := range this.getPathToValueMapFromObj(append(startPath, strconv.Itoa(k)), sub) {
				result[subPath] = subValue
			}
		}
	}
	return result
}

func (this *Marshaller) getPathToCharacteristicFromContents(content []model.Content) (result map[string]string) {
	result = map[string]string{}
	for _, c := range content {
		for key, value := range this.getPathToCharacteristicFromVariable([]string{}, c.ContentVariable) {
			result[key] = value
		}
	}
	return result
}

func (this *Marshaller) getPathToCharacteristicFromVariable(startPath []string, variable model.ContentVariable) (result map[string]string) {
	startPath = append(startPath, variable.Name)
	result = map[string]string{
		strings.Join(startPath, "."): variable.CharacteristicId,
	}
	for _, sub := range variable.SubContentVariables {
		for subPath, subValue := range this.getPathToCharacteristicFromVariable(startPath, sub) {
			result[subPath] = subValue
		}
	}
	return result
}

func serializeOutput(output map[string]string, service model.Service, protocol model.Protocol) (result map[string]interface{}, err error) {
	result = map[string]interface{}{}
	for _, content := range service.Outputs {
		for _, segment := range protocol.ProtocolSegments {
			if segment.Id == content.ProtocolSegmentId {
				output, ok := output[segment.Name]
				if ok {
					marshaller, ok := serialization.Get(content.Serialization)
					if !ok {
						debug.PrintStack()
						return result, errors.New("unknown serialization " + content.Serialization)
					}
					value, err := marshaller.Unmarshal(output, content.ContentVariable)
					if err != nil {
						return result, err
					}
					result[segment.Name] = value
				}
			}
		}
	}
	return
}
