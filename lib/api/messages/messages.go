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

package messages

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

type MarshallingV2Request struct {
	Service  model.Service                    `json:"service"`  //semi-optional, may be determined by request path
	Protocol model.Protocol                   `json:"protocol"` //semi-optional, may be determined by service
	Data     []model.MarshallingV2RequestData `json:"data"`
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

type UnmarshallingV2Request struct {
	Service          model.Service  `json:"service"`           //semi-optional, may be determined by request path
	Protocol         model.Protocol `json:"protocol"`          //semi-optional, may be determined by service
	CharacteristicId string         `json:"characteristic_id"` //semi-optional, may be determined by request path

	Message          map[string]string      `json:"message"`           //semi-optional; may be needed to create serialized_output
	SerializedOutput map[string]interface{} `json:"serialized_output"` //semi-optional; may be created from message

	Path          string             `json:"path"`                      //semi-optional, may be determent by FunctionId and AspectNodes
	FunctionId    string             `json:"function_id"`               //semi-optional, to determine Path if not set
	AspectNode    model.AspectNode   `json:"aspect_node"`               //deprecated: please use AspectNodes
	AspectNodeId  string             `json:"aspect_node_id"`            //deprecated: please use AspectNodeIds
	AspectNodes   []model.AspectNode `json:"aspect_nodes,omitempty"`    //semi-optional, to determine Path if not set, may themselves be determent by AspectNodeIds
	AspectNodeIds []string           `json:"aspect_node_ids,omitempty"` //semi-optional, to determine AspectNodes if not set
}

// GetAspectNodeIds returns the aspect ids that still have to be resolved to aspect nodes.
// The deprecated AspectNodeId is an alias for a list with one element and is only used if
// the deprecated AspectNode is unset, the way it was before the lists existed.
func (this UnmarshallingV2Request) GetAspectNodeIds() []string {
	if this.AspectNode.Id != "" {
		return this.AspectNodeIds
	}
	return model.AspectIdsAlias(this.AspectNodeId, this.AspectNodeIds)
}

// GetAspectNodes returns the aspect nodes the request carries. The deprecated AspectNode is
// an alias for a list with one element.
func (this UnmarshallingV2Request) GetAspectNodes() []model.AspectNode {
	return model.AspectNodesAlias(this.AspectNode, this.AspectNodes)
}

type FindConfigurablesRequest struct {
	CharacteristicId string          `json:"characteristic_id"`
	Services         []model.Service `json:"services"`
}

type PathOptionsQuery struct {
	DeviceTypeIds          []string `json:"device_type_ids"`
	FunctionId             string   `json:"function_id"`
	AspectId               string   `json:"aspect_id"` //deprecated: please use AspectIds
	AspectIds              []string `json:"aspect_ids,omitempty"`
	CharacteristicIdFilter []string `json:"characteristic_id_filter"`
	WithoutEnvelope        bool     `json:"without_envelope"`
}

// GetAspectIds returns the aspects the query asks for. The deprecated AspectId is an alias
// for a list with one element.
func (this PathOptionsQuery) GetAspectIds() []string {
	return model.AspectIdsAlias(this.AspectId, this.AspectIds)
}
