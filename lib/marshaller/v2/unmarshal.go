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
	convertermodel "github.com/SENERGY-Platform/converter/lib/model"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/serialization"
	"runtime/debug"
	"strconv"
	"strings"
)

var PathNotFoundInMessage = errors.New("path not found in message")

func (this *Marshaller) Unmarshal(protocol model.Protocol, service model.Service, characteristicId string, path string, msg map[string]string) (result interface{}, err error) {
	path = substitudeVariableLenPlaceholderInPath(path)

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
		return result, PathNotFoundInMessage
	}

	service.Outputs, err = substituteVariableLenListsInOutputs(service.Outputs, pathToValue)

	//no conversion wanted
	if characteristicId == "" {
		return value, nil
	}

	//convert service/variable characteristic to wanted characteristic
	pathToCharacteristic := getPathToCharacteristicFromContents(service.Outputs)
	pathToFunction := getPathToFunctionFromContents(service.Outputs)
	variableCharacteristic, ok := pathToCharacteristic[path]
	if !ok {
		return result, errors.New("path not found in service")
	}

	value, err = this.variableStructureToCharacteristicsStructure(pathToCharacteristic, path, value)
	if err != nil {
		return result, err
	}

	if variableCharacteristic == characteristicId {
		return value, nil
	} else {
		castExtensions := []convertermodel.ConverterExtension{}
		functionId := pathToFunction[path]
		if functionId != "" {
			conceptId := this.characteristics.GetConceptIdOfFunction(functionId)
			if conceptId != "" {
				concept, err := this.characteristics.GetConcept(conceptId)
				if err != nil {
					return nil, err
				}
				castExtensions = concept.Conversions
			}
		}
		if len(castExtensions) == 0 {
			return this.converter.Cast(value, variableCharacteristic, characteristicId)
		} else {
			return this.converter.CastWithExtension(value, variableCharacteristic, characteristicId, castExtensions)
		}
	}
}

func substitudeVariableLenPlaceholderInPath(path string) string {
	if path == "*" {
		path = "0"
	}
	path = strings.ReplaceAll(path, ".*.", ".0.")
	if strings.HasPrefix(path, "*.") {
		path = "0." + strings.TrimPrefix(path, "*.")
	}
	if strings.HasSuffix(path, ".*") {
		path = strings.TrimSuffix(path, ".*") + ".0"
	}
	return path
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

func getPathToFunctionFromContents(content []model.Content) (result map[string]string) {
	result = map[string]string{}
	for _, c := range content {
		temp := walkPathToMap(
			[]string{},
			c.ContentVariable, func(v model.ContentVariable) string { return v.Name },
			func(v model.ContentVariable) string { return v.FunctionId },
			func(v model.ContentVariable) []model.ContentVariable { return v.SubContentVariables },
		)
		for key, value := range temp {
			result[key] = value
		}
	}
	return result
}

func getPathToCharacteristicFromContents(content []model.Content) (result map[string]string) {
	result = map[string]string{}
	for _, c := range content {
		temp := walkPathToMap(
			[]string{},
			c.ContentVariable, func(v model.ContentVariable) string { return v.Name },
			func(v model.ContentVariable) string { return v.CharacteristicId },
			func(v model.ContentVariable) []model.ContentVariable { return v.SubContentVariables },
		)
		for key, value := range temp {
			result[key] = value
		}
	}
	return result
}

func walkPathToMap[Element any, PathElement func(Element) string, Field func(Element) FieldType, FieldType any, Sub func(Element) []Element](startPath []string, element Element, pathElement PathElement, field Field, sub Sub) (result map[string]FieldType) {
	startPath = append(startPath, pathElement(element))
	result = map[string]FieldType{
		strings.Join(startPath, "."): field(element),
	}
	for _, subElement := range sub(element) {
		for subPath, subValue := range walkPathToMap(startPath, subElement, pathElement, field, sub) {
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
					result[content.ContentVariable.Name] = value
				}
			}
		}
	}
	return
}

func (this *Marshaller) variableStructureToCharacteristicsStructure(variablePathToCharacteristic map[string]string, variablePath string, value interface{}) (result interface{}, err error) {
	value, err = normalize(value)
	if err != nil {
		return result, err
	}
	switch value.(type) {
	case map[string]interface{}:
	case []interface{}:
	default:
		return value, nil
	}

	shortVariablePathToCharacteristic := map[string]string{}
	for subVariablePath, variableCharacteristic := range variablePathToCharacteristic {
		if subVariablePath != variablePath && strings.HasPrefix(subVariablePath, variablePath+".") {
			shortVariablePath := strings.Replace(subVariablePath, variablePath+".", "", 1)
			shortVariablePathToCharacteristic[shortVariablePath] = variableCharacteristic
		}
	}

	characteristicId := variablePathToCharacteristic[variablePath]
	characteristic, err := this.characteristics.GetCharacteristic(characteristicId)
	if err != nil {
		return result, err
	}

	characteristicToVariablePath := map[string]string{}
	for path, variableCharacteristic := range shortVariablePathToCharacteristic {
		characteristicToVariablePath[variableCharacteristic] = path
	}

	variablePathToValue := this.getPathToValueMapFromObj([]string{}, value)

	characteristicIdToValue := map[string]interface{}{}
	for subVariablePath, variableValue := range variablePathToValue {
		subCharacteristicId, ok := shortVariablePathToCharacteristic[subVariablePath]
		if ok && subCharacteristicId != characteristicId {
			characteristicIdToValue[subCharacteristicId] = variableValue
		}
	}

	characteristicsIdToPath := getCharacteristicIdToPath(characteristic, []string{})

	characteristicsPathToValue := map[string]interface{}{}
	for subCharacteristicId, subValue := range characteristicIdToValue {
		characteristicsPathToValue[characteristicsIdToPath[subCharacteristicId]] = subValue
	}

	pseudoVariable := characteristicToPseudoVariable(characteristic)
	for characteristicsPath, subValue := range characteristicsPathToValue {
		pseudoVariable, err = this.setContentVariableValue(pseudoVariable, []string{}, strings.Split(characteristicsPath, "."), "", subValue)
		if err != nil {
			return result, err
		}
	}

	_, result, err = contentVariableToObject(pseudoVariable)
	if err != nil {
		return result, err
	}
	return result, err
}

func characteristicToPseudoVariable(characteristic model.Characteristic) model.ContentVariable {
	subVariables := []model.ContentVariable{}
	for _, sub := range characteristic.SubCharacteristics {
		subVariables = append(subVariables, characteristicToPseudoVariable(sub))
	}
	variableType := characteristic.Type
	if variableType == "" && len(subVariables) > 0 {
		variableType = model.Structure
	}
	return model.ContentVariable{
		Name:                characteristic.Name,
		Type:                variableType,
		SubContentVariables: subVariables,
	}
}

func getCharacteristicIdToPath(characteristic model.Characteristic, currentPath []string) (result map[string]string) {
	result = map[string]string{}
	currentPath = append(currentPath, characteristic.Name)
	if characteristic.Id != "" {
		result[characteristic.Id] = strings.Join(currentPath, ".")
	}
	for _, sub := range characteristic.SubCharacteristics {
		for characteristicId, path := range getCharacteristicIdToPath(sub, currentPath) {
			result[characteristicId] = path
		}
	}
	return result
}

func substituteVariableLenListsInOutputs(outputs []model.Content, value map[string]interface{}) (result []model.Content, err error) {
	result = []model.Content{}
	for _, c := range outputs {
		c.ContentVariable, err = substituteVariableLenListsInVariable(c.ContentVariable, []string{}, value)
		if err != nil {
			return result, err
		}
		result = append(result, c)
	}
	return result, nil
}

func substituteVariableLenListsInVariable(variable model.ContentVariable, currentPath []string, pathToValue map[string]interface{}) (model.ContentVariable, error) {
	currentPath = append(currentPath, variable.Name)
	if variable.Type == model.List && len(variable.SubContentVariables) > 0 && variable.SubContentVariables[0].Name == "*" {
		currentPathStr := strings.Join(currentPath, ".")
		currentValue := pathToValue[currentPathStr]
		list, ok := currentValue.([]interface{})
		if ok {
			max := len(list)
			for i := 0; i < max; i++ {
				clonedChild := variable.SubContentVariables[0]
				clonedChild.Name = strconv.Itoa(i)
				variable.SubContentVariables = append(variable.SubContentVariables, clonedChild)
			}
		}
	}
	temp := []model.ContentVariable{}
	for _, sub := range variable.SubContentVariables {
		newSub, err := substituteVariableLenListsInVariable(sub, currentPath, pathToValue)
		if err != nil {
			return variable, err
		}
		temp = append(temp, newSub)
	}
	variable.SubContentVariables = temp
	return variable, nil
}
