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
	"context"
	"github.com/SENERGY-Platform/converter/lib/converter/characteristics"
	"github.com/SENERGY-Platform/marshaller/lib/api"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"sync"
	"testing"
)

func TestUnmarshalListIndexed(t *testing.T) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	apiurl := setup(ctx, wg)

	protocol := model.Protocol{
		Id:      "p1",
		Name:    "p1",
		Handler: "p1",
		ProtocolSegments: []model.ProtocolSegment{
			{Id: "p1.1", Name: "body"},
			{Id: "p1.2", Name: "head"},
		},
	}
	service := model.Service{
		Id:          "sid",
		LocalId:     "slid",
		Name:        "sname",
		Interaction: model.EVENT_AND_REQUEST,
		ProtocolId:  "p1",
		Outputs: []model.Content{
			{
				Id: "content",
				ContentVariable: model.ContentVariable{
					Id:   "temperature",
					Name: "temperature",
					Type: model.List,
					SubContentVariables: []model.ContentVariable{
						{
							Id:               "inside",
							Name:             "0",
							Type:             model.Integer, //results like 26,85 will be rounded to 27
							CharacteristicId: characteristics.Celsius,
							FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
							AspectId:         "inside_air",
							Value:            12,
						},
						{
							Id:               "outside",
							Name:             "1",
							Type:             model.Integer, //results like 26,85 will be rounded to 27
							CharacteristicId: characteristics.Celsius,
							FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
							AspectId:         "outside_air",
							Value:            13,
						},
					},
				},
				Serialization:     "json",
				ProtocolSegmentId: "p1.1",
			},
		},
	}

	output := map[string]string{"body": `[400,500]`}

	t.Run("0 path", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		Message:          output,
		Path:             "temperature.0",
	}, 673.15))

	t.Run("1 path", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		Message:          output,
		Path:             "temperature.1",
	}, 773.15))

	t.Run("0 criteria", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		Message:          output,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "inside_air",
	}, 673.15))

	t.Run("1 criteria", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		Message:          output,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "outside_air",
	}, 773.15))

	t.Run("air criteria", testUnmarshalAny(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		Message:          output,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "air",
	}, []interface{}{673.15, 773.15}))

	t.Run("0 no cast", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:      service,
		Protocol:     protocol,
		Message:      output,
		FunctionId:   model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId: "inside_air",
	}, 400.0))

	t.Run("inside °C", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		Message:          output,
		CharacteristicId: characteristics.Celsius,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "inside_air",
	}, 400.0))
}

func TestUnarshalListVariable(t *testing.T) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	apiurl := setup(ctx, wg)

	protocol := model.Protocol{
		Id:      "p1",
		Name:    "p1",
		Handler: "p1",
		ProtocolSegments: []model.ProtocolSegment{
			{Id: "p1.1", Name: "body"},
			{Id: "p1.2", Name: "head"},
		},
	}
	service := model.Service{
		Id:          "sid",
		LocalId:     "slid",
		Name:        "sname",
		Interaction: model.EVENT_AND_REQUEST,
		ProtocolId:  "p1",
		Outputs: []model.Content{
			{
				Id: "content",
				ContentVariable: model.ContentVariable{
					Id:   "temperature",
					Name: "temperature",
					Type: model.List,
					SubContentVariables: []model.ContentVariable{
						{
							Id:               "var",
							Name:             "*",
							Type:             model.Integer, //results like 26,85 will be rounded to 27
							CharacteristicId: characteristics.Celsius,
							FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
							AspectId:         "inside_air",
							Value:            12,
						},
					},
				},
				Serialization:     "json",
				ProtocolSegmentId: "p1.1",
			},
		},
	}

	output := map[string]string{"body": `[400,500]`}

	t.Run("temperature path", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: "",
		Message:          output,
		Path:             "temperature",
	}, []interface{}{400.0, 500.0}))

	t.Run("* path", testUnmarshalAny(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		Message:          output,
		Path:             "temperature.*",
	}, []interface{}{673.15, 773.15}))

	t.Run("0 path", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		Message:          output,
		Path:             "temperature.0",
	}, 673.15))

	t.Run("1 path", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		Message:          output,
		Path:             "temperature.1",
	}, 773.15))

	t.Run("99 path", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		Message:          output,
		Path:             "temperature.99",
	}, nil))

	t.Run("inside_air criteria", testUnmarshalAny(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		Message:          output,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "inside_air",
	}, []interface{}{673.15, 773.15}))

	t.Run("air criteria", testUnmarshalAny(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		Message:          output,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "air",
	}, []interface{}{673.15, 773.15}))

	t.Run("inside_air no cast", testUnmarshalAny(apiurl, api.UnmarshallingV2Request{
		Service:      service,
		Protocol:     protocol,
		Message:      output,
		FunctionId:   model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId: "inside_air",
	}, []interface{}{400.0, 500.0}))

	t.Run("inside_air °C", testUnmarshalAny(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		Message:          output,
		CharacteristicId: characteristics.Celsius,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "inside_air",
	}, []interface{}{400.0, 500.0}))
}
