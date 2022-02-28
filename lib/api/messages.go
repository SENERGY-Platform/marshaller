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

package api

import (
	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
)

type MarshallingRequest struct {
	Service          model.Service                `json:"service,omitempty"`           //semi-optional, may be determined by request path
	Protocol         *model.Protocol              `json:"protocol,omitempty"`          //semi-optional, may be determined by request path
	CharacteristicId string                       `json:"characteristic_id,omitempty"` //semi-optional, may be determined by request path
	Configurables    []configurables.Configurable `json:"configurables,omitempty"`     //optional, may be empty
	Data             interface{}                  `json:"data"`

	/*
		optional
		if len > 0: apply data only on given ContentVariable paths
		useful if 2 variables have the same characteristic assigned but are used with different functions or aspects
	*/
	PathAllowList []string `json:"path_allow_list,omitempty"`
}

type UnmarshallingRequest struct {
	Service              model.Service     `json:"service,omitempty"`           //semi-optional, may be determined by request path
	Protocol             *model.Protocol   `json:"protocol,omitempty"`          //semi-optional, may be determined by service
	CharacteristicId     string            `json:"characteristic_id,omitempty"` //semi-optional, may be determined by request path
	Message              map[string]string `json:"message"`
	ContentVariableHints []string          `json:"content_variable_hints"` //optional

	/*
		optional
		if len > 0: apply data only on given ContentVariable paths
		useful if 2 variables have the same characteristic assigned but are used with different functions or aspects
	*/
	PathAllowList []string `json:"path_allow_list,omitempty"`
}

type FindConfigurablesRequest struct {
	CharacteristicId string          `json:"characteristic_id"`
	Services         []model.Service `json:"services"`
}

type PathOptionsQuery struct {
	DeviceTypeIds          []string `json:"device_type_ids"`
	FunctionId             string   `json:"function_id"`
	AspectId               string   `json:"aspect_id"`
	CharacteristicIdFilter []string `json:"characteristic_id_filter"`
	WithoutEnvelope        bool     `json:"without_envelope"`
}
