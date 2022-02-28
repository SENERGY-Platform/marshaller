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
)

func RemoveVoidVariables(contents []model.Content) (result []model.Content) {
	result = []model.Content{}
	for _, content := range contents {
		if !content.ContentVariable.IsVoid {
			content.ContentVariable = removeVoidVariables(content.ContentVariable)
			result = append(result, content)
		}
	}
	return result
}

func removeVoidVariables(variable model.ContentVariable) model.ContentVariable {
	for _, sub := range variable.SubContentVariables {
		if !sub.IsVoid {
			variable.SubContentVariables = append(variable.SubContentVariables, removeVoidVariables(sub))
		}
	}
	return variable
}
