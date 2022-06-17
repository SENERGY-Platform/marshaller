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

func TestMarshalListIndexed(t *testing.T) {
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
		Inputs: []model.Content{
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
							FunctionId:       model.CONTROLLING_FUNCTION_PREFIX + "setTemperature",
							AspectId:         "inside_air",
							Value:            12,
						},
						{
							Id:               "outside",
							Name:             "1",
							Type:             model.Integer, //results like 26,85 will be rounded to 27
							CharacteristicId: characteristics.Celsius,
							FunctionId:       model.CONTROLLING_FUNCTION_PREFIX + "setTemperature",
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

	t.Run("inside 300 kelvin path", testMarshal(apiurl, api.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            300,
				CharacteristicId: characteristics.Kelvin,
				Paths:            []string{"temperature.0"},
				FunctionId:       "",
			},
		},
	}, map[string]string{"body": `[27,13]`}))

	t.Run("inside 300 kelvin path and function", testMarshal(apiurl, api.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            300,
				CharacteristicId: characteristics.Kelvin,
				Paths:            []string{"temperature.0"},
				FunctionId:       model.CONTROLLING_FUNCTION_PREFIX + "setTemperature",
			},
		},
	}, map[string]string{"body": `[27,13]`}))

	t.Run("outside 300 kelvin path", testMarshal(apiurl, api.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            300,
				CharacteristicId: characteristics.Kelvin,
				Paths:            []string{"temperature.1"},
				FunctionId:       "",
			},
		},
	}, map[string]string{"body": `[12,27]`}))

	t.Run("outside 300 kelvin path and function", testMarshal(apiurl, api.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            300,
				CharacteristicId: characteristics.Kelvin,
				Paths:            []string{"temperature.1"},
				FunctionId:       model.CONTROLLING_FUNCTION_PREFIX + "setTemperature",
			},
		},
	}, map[string]string{"body": `[12,27]`}))

	t.Run("inside and outside 300 kelvin path", testMarshal(apiurl, api.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            300,
				CharacteristicId: characteristics.Kelvin,
				Paths:            []string{"temperature.0", "temperature.1"},
				FunctionId:       "",
			},
		},
	}, map[string]string{"body": `[27,27]`}))

	t.Run("inside 400 and outside 500 kelvin path", testMarshal(apiurl, api.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            400,
				CharacteristicId: characteristics.Kelvin,
				Paths:            []string{"temperature.0"},
				FunctionId:       "",
			},
			{
				Value:            500,
				CharacteristicId: characteristics.Kelvin,
				Paths:            []string{"temperature.1"},
				FunctionId:       "",
			},
		},
	}, map[string]string{"body": `[127,227]`}))

	t.Run("inside and outside 300 kelvin functionId", testMarshal(apiurl, api.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            300,
				CharacteristicId: characteristics.Kelvin,
				Paths:            nil,
				FunctionId:       model.CONTROLLING_FUNCTION_PREFIX + "setTemperature",
			},
		},
	}, map[string]string{"body": `[27,27]`}))
}

func TestMarshalListVariable(t *testing.T) {
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
		Inputs: []model.Content{
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
							FunctionId:       model.CONTROLLING_FUNCTION_PREFIX + "setTemperature",
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

	t.Run("* 300 kelvin path", testMarshal(apiurl, api.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            300,
				CharacteristicId: characteristics.Kelvin,
				Paths:            []string{"temperature.*"},
				FunctionId:       "",
			},
		},
	}, map[string]string{"body": `[27]`}))

	t.Run("functionId 300 kelvin", testMarshal(apiurl, api.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            300,
				CharacteristicId: characteristics.Kelvin,
				Paths:            nil,
				FunctionId:       model.CONTROLLING_FUNCTION_PREFIX + "setTemperature",
			},
		},
	}, map[string]string{"body": `[27]`}))

	t.Run("indexed setting of var len list values", func(t *testing.T) {
		t.Skip("would be nice but is not supported")

		t.Run("0 300 kelvin path", testMarshal(apiurl, api.MarshallingV2Request{
			Service:  service,
			Protocol: protocol,
			Data: []model.MarshallingV2RequestData{
				{
					Value:            300,
					CharacteristicId: characteristics.Kelvin,
					Paths:            []string{"temperature.0"},
					FunctionId:       "",
				},
			},
		}, map[string]string{"body": `[27]`}))

		t.Run("0 300 kelvin path and function", testMarshal(apiurl, api.MarshallingV2Request{
			Service:  service,
			Protocol: protocol,
			Data: []model.MarshallingV2RequestData{
				{
					Value:            300,
					CharacteristicId: characteristics.Kelvin,
					Paths:            []string{"temperature.0"},
					FunctionId:       model.CONTROLLING_FUNCTION_PREFIX + "setTemperature",
				},
			},
		}, map[string]string{"body": `[27]`}))

		t.Run("0 and 1 300 kelvin path", testMarshal(apiurl, api.MarshallingV2Request{
			Service:  service,
			Protocol: protocol,
			Data: []model.MarshallingV2RequestData{
				{
					Value:            300,
					CharacteristicId: characteristics.Kelvin,
					Paths:            []string{"temperature.0", "temperature.1"},
					FunctionId:       "",
				},
			},
		}, map[string]string{"body": `[27,27]`}))

		t.Run("0 400 and 1 500 kelvin path", testMarshal(apiurl, api.MarshallingV2Request{
			Service:  service,
			Protocol: protocol,
			Data: []model.MarshallingV2RequestData{
				{
					Value:            400,
					CharacteristicId: characteristics.Kelvin,
					Paths:            []string{"temperature.0"},
					FunctionId:       "",
				},
				{
					Value:            500,
					CharacteristicId: characteristics.Kelvin,
					Paths:            []string{"temperature.1"},
					FunctionId:       "",
				},
			},
		}, map[string]string{"body": `[127,227]`}))
	})
}

func TestMarshalListCharacteristic(t *testing.T) {
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
		Inputs: []model.Content{
			{
				Id: "content",
				ContentVariable: model.ContentVariable{
					Id:               "root",
					Name:             "root",
					Type:             model.Structure,
					CharacteristicId: "urn:infai:ses:characteristic:f48d7985-7ee7-4119-a791-bc16a953f440",
					FunctionId:       "urn:infai:ses:controlling-function:ced44f01-7328-43e3-8db0-ecd12f448758",
					SubContentVariables: []model.ContentVariable{
						{
							Id:               "segment_ids",
							Name:             "segment_ids",
							Type:             model.List,
							CharacteristicId: "urn:infai:ses:characteristic:b0bf0d79-8a23-40d3-a284-2b87e38138be",
							SubContentVariables: []model.ContentVariable{
								{
									Id:               "var",
									Name:             "*",
									Type:             model.String,
									CharacteristicId: "urn:infai:ses:characteristic:802c79a9-c96b-4848-848c-bae29fb00375",
								},
							},
						},
						{
							Id:               "iterations",
							Name:             "iterations",
							Type:             model.Integer,
							CharacteristicId: "urn:infai:ses:characteristic:63769613-47bd-4bb9-9624-083669ee6261",
						},
						{
							Id:               "customOrder",
							Name:             "customOrder",
							Type:             model.Boolean,
							CharacteristicId: "urn:infai:ses:characteristic:aeacd820-8a9e-4a68-a4c6-412537bb8afb",
						},
					},
				},
				Serialization:     "json",
				ProtocolSegmentId: "p1.1",
			},
		},
	}

	t.Run("test marshal", testMarshal(apiurl, api.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            map[string]interface{}{"iterations": 1, "customOrder": true, "segment_ids": []string{"1", "2"}},
				CharacteristicId: "urn:infai:ses:characteristic:f48d7985-7ee7-4119-a791-bc16a953f440",
				Paths:            []string{"root"},
				FunctionId:       "urn:infai:ses:controlling-function:ced44f01-7328-43e3-8db0-ecd12f448758",
			},
		},
	}, map[string]string{"body": `{"customOrder":true,"iterations":1,"segment_ids":["1","2"]}`}))

}
