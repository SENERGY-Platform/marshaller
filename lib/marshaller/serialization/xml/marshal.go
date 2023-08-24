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

package xml

import (
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/clbanning/mxj"
	"slices"
	"strings"
)

func (Marshaller) Marshal(in interface{}, variable model.ContentVariable) (out string, err error) {
	in, err = rewriteFieldNamesFromSerializationOptions(in, variable, func(name string, options []string) string {
		if !strings.HasPrefix(name, "-") && slices.Contains(options, models.SerializationOptionXmlAttribute) {
			name = "-" + name
		}
		return name
	})
	mv, ok := in.(map[string]interface{})
	if !ok {
		mv = map[string]interface{}{variable.Name: in}
		temp, err := mxj.Map(mv).Xml()
		return string(temp), err
	} else {
		temp, err := mxj.Map(mv).Xml(variable.Name)
		return string(temp), err
	}
}
