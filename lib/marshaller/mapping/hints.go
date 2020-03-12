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

package mapping

import (
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
)

func applyHints(content map[string]model.ContentVariable, hints []string) (result map[string]model.ContentVariable, err error) {
	if len(hints) == 0 {
		return content, nil
	}
	for _, hint := range hints {
		var used bool
		for _, variable := range content {
			used = used || hintIsUsed(variable, hint)
		}
		if used {
			return ignoreCharacteristicsFromContentExceptHint(content, hint), nil
		}
	}
	return content, nil
}

func hintIsUsed(variable model.ContentVariable, hint string) bool {
	if variable.Id == hint {
		return true
	}
	for _, sub := range variable.SubContentVariables {
		if hintIsUsed(sub, hint) {
			return true
		}
	}
	return false
}

func ignoreCharacteristicsFromContentExceptHint(content map[string]model.ContentVariable, hint string) (result map[string]model.ContentVariable) {
	result = map[string]model.ContentVariable{}
	for key, variable := range content {
		result[key] = ignoreCharacteristicsFromVariableExceptHint(variable, hint)
	}
	return result
}

func ignoreCharacteristicsFromVariableExceptHint(variable model.ContentVariable, hint string) model.ContentVariable {
	if variable.Id == hint {
		return variable
	}
	variable.CharacteristicId = ""
	temp := []model.ContentVariable{}
	for _, sub := range variable.SubContentVariables {
		temp = append(temp, ignoreCharacteristicsFromVariableExceptHint(sub, hint))
	}
	variable.SubContentVariables = temp
	return variable
}
