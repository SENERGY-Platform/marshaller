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

package marshaller

import (
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"strings"
)

func UsePathAllowList(contents []model.Content, pathAllowList []string) (result []model.Content) {
	if len(pathAllowList) == 0 {
		return contents
	}
	for _, content := range contents {
		content.ContentVariable = removeCharacteristicIdFromVariablesNotInPathAllowList(content.ContentVariable, pathAllowList, []string{})
		result = append(result, content)
	}
	return result
}

func removeCharacteristicIdFromVariablesNotInPathAllowList(variable model.ContentVariable, pathAllowList []string, previousPath []string) model.ContentVariable {
	currentPath := append(previousPath, variable.Name)
	currentPathString := strings.Join(currentPath, ".")
	if !contains(pathAllowList, currentPathString) {
		variable.CharacteristicId = ""
		//don't remove characteristic id from sub variables of matching/allowed variables --> loop in if body
		for i, sub := range variable.SubContentVariables {
			variable.SubContentVariables[i] = removeCharacteristicIdFromVariablesNotInPathAllowList(sub, pathAllowList, currentPath)
		}
	}
	return variable
}
