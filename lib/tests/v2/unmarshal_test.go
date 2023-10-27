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
	"bytes"
	"context"
	"encoding/json"
	"github.com/SENERGY-Platform/converter/lib/converter/characteristics"
	"github.com/SENERGY-Platform/marshaller/lib/api"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"net/http"
	"reflect"
	"sync"
	"testing"
)

func TestUnmarshalIncompleteCharacteristic(t *testing.T) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	apiurl := setup(ctx, wg)

	functionId := "urn:infai:ses:measuring-function:bdb6a7c8-4a3d-4fe0-bab3-ce02e09b5869" //Read Color

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
					Id:   "root",
					Name: "root",
					Type: model.Structure,
					SubContentVariables: []model.ContentVariable{
						{
							Id:               "color",
							Name:             "color",
							CharacteristicId: characteristics.Hsb,
							FunctionId:       functionId,
							AspectId:         "air",
							Type:             model.Structure,
							SubContentVariables: []model.ContentVariable{
								{
									Id:               "hue",
									Name:             "hue",
									Type:             model.Integer,
									CharacteristicId: characteristics.HsbH,
								},
								{
									Id:               "sat",
									Name:             "sat",
									Type:             model.Integer,
									CharacteristicId: characteristics.HsbS,
								},
							},
						},
						{
							Id:   "brightness",
							Name: "brightness",
							Type: model.Integer,
						},
					},
				},
				Serialization:     "json",
				ProtocolSegmentId: "p1.1",
			},
		},
	}

	output := map[string]string{"body": `{"brightness":66,"color":{"hue": 219, "sat": 68}}`}

	t.Run("run", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Rgb,
		Message:          output,
		FunctionId:       functionId,
		AspectNodeId:     "air",
	}, map[string]interface{}{"b": 128.0, "g": 71.0, "r": 41.0}))
}

func TestUnmarshalDiscombobulatedCharacteristic(t *testing.T) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	apiurl := setup(ctx, wg)

	functionId := "urn:infai:ses:measuring-function:bdb6a7c8-4a3d-4fe0-bab3-ce02e09b5869" //Read Color

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
					Id:               "root",
					Name:             "root",
					Type:             model.Structure,
					CharacteristicId: characteristics.Hsb,
					FunctionId:       functionId,
					AspectId:         "air",
					SubContentVariables: []model.ContentVariable{
						{
							Id:   "color",
							Name: "color",
							Type: model.Structure,
							SubContentVariables: []model.ContentVariable{
								{
									Id:               "hue",
									Name:             "hue",
									Type:             model.Integer,
									CharacteristicId: characteristics.HsbH,
								},
								{
									Id:               "sat",
									Name:             "sat",
									Type:             model.Integer,
									CharacteristicId: characteristics.HsbS,
								},
							},
						},
						{
							Id:               "brightness",
							Name:             "brightness",
							Type:             model.Integer,
							CharacteristicId: characteristics.HsbB,
						},
					},
				},
				Serialization:     "json",
				ProtocolSegmentId: "p1.1",
			},
		},
	}

	t.Run("run_1", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Rgb,
		Message:          map[string]string{"body": `{"brightness":66,"color":{"hue": 219, "sat": 68}}`},
		FunctionId:       functionId,
		AspectNodeId:     "air",
	}, map[string]interface{}{"b": 168.0, "g": 94.0, "r": 54.0}))

	t.Run("run_2", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Rgb,
		Message:          map[string]string{"body": `{"brightness":50,"color":{"hue": 219, "sat": 68}}`},
		FunctionId:       functionId,
		AspectNodeId:     "air",
	}, map[string]interface{}{"b": 128.0, "g": 71.0, "r": 41.0}))
}

func TestUnmarshalPrioritySort(t *testing.T) {
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
					Type: model.Structure,
					SubContentVariables: []model.ContentVariable{
						{
							Id:               "today",
							Name:             "today",
							Type:             model.Float,
							CharacteristicId: characteristics.Celsius,
							FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
							AspectId:         "today",
						},
						{
							Id:               "consumption",
							Name:             "consumption",
							Type:             model.Float,
							CharacteristicId: characteristics.Celsius,
							FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
							AspectId:         "consumption",
						},
					},
				},
				Serialization:     "json",
				ProtocolSegmentId: "p1.1",
			},
		},
	}

	output := map[string]string{"body": `{"today":400,"consumption":500}`}

	t.Run("electricity", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Celsius,
		Message:          output,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "electricity",
	}, 500.0))

	t.Run("consumption", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Celsius,
		Message:          output,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "consumption",
	}, 500.0))

	t.Run("today", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Celsius,
		Message:          output,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "today",
	}, 400.0))
}

func TestUnmarshalPrioritySort2(t *testing.T) {
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
					Type: model.Structure,
					SubContentVariables: []model.ContentVariable{
						{
							Id:               "today",
							Name:             "today",
							Type:             model.Float,
							CharacteristicId: characteristics.Celsius,
							FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
							AspectId:         "today",
						},
						{
							Id:               "electricity",
							Name:             "electricity",
							Type:             model.Float,
							CharacteristicId: characteristics.Celsius,
							FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
							AspectId:         "electricity",
						},
					},
				},
				Serialization:     "json",
				ProtocolSegmentId: "p1.1",
			},
		},
	}

	output := map[string]string{"body": `{"today":400,"electricity":500}`}

	t.Run("electricity", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Celsius,
		Message:          output,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "electricity",
	}, 500.0))

	t.Run("consumption", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Celsius,
		Message:          output,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "consumption",
	}, 400.0))

	t.Run("today", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Celsius,
		Message:          output,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "today",
	}, 400.0))
}

func TestUnmarshalling(t *testing.T) {
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
					Type: model.Structure,
					SubContentVariables: []model.ContentVariable{
						{
							Id:               "inside",
							Name:             "inside",
							Type:             model.Float,
							CharacteristicId: characteristics.Celsius,
							FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
							AspectId:         "inside_air",
						},
						{
							Id:               "outside",
							Name:             "outside",
							Type:             model.Float,
							CharacteristicId: characteristics.Celsius,
							FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
							AspectId:         "outside_air",
						},
					},
				},
				Serialization:     "json",
				ProtocolSegmentId: "p1.1",
			},
		},
	}

	output := map[string]string{"body": `{"inside":400,"outside":500}`}
	serializedOutput := map[string]interface{}{"temperature": map[string]interface{}{"inside": 400, "outside": 500}}

	t.Run("inside path", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		Message:          output,
		Path:             "temperature.inside",
	}, 673.15))

	t.Run("outside path", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		Message:          output,
		Path:             "temperature.outside",
	}, 773.15))

	t.Run("inside criteria", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		Message:          output,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "inside_air",
	}, 673.15))

	t.Run("outside criteria", testUnmarshal(apiurl, api.UnmarshallingV2Request{
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

	t.Run("inside no cast", testUnmarshal(apiurl, api.UnmarshallingV2Request{
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

	t.Run("inside path serialized", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		SerializedOutput: serializedOutput,
		Path:             "temperature.inside",
	}, 673.15))

	t.Run("outside path serialized", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		SerializedOutput: serializedOutput,
		Path:             "temperature.outside",
	}, 773.15))

	t.Run("inside criteria serialized", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		SerializedOutput: serializedOutput,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "inside_air",
	}, 673.15))

	t.Run("outside criteria serialized", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		SerializedOutput: serializedOutput,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "outside_air",
	}, 773.15))

	t.Run("air criteria serialized", testUnmarshalAny(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Kelvin,
		SerializedOutput: serializedOutput,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "air",
	}, []interface{}{673.15, 773.15}))

	t.Run("inside no cast serialized", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		SerializedOutput: serializedOutput,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "inside_air",
	}, 400.0))

	t.Run("inside °C serialized", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		SerializedOutput: serializedOutput,
		CharacteristicId: characteristics.Celsius,
		FunctionId:       model.MEASURING_FUNCTION_PREFIX + "getTemperature",
		AspectNodeId:     "inside_air",
	}, 400.0))

}

func testUnmarshal(apiurl string, request api.UnmarshallingV2Request, expectedResult interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		body := new(bytes.Buffer)
		err := json.NewEncoder(body).Encode(request)
		if err != nil {
			t.Error(err)
			return
		}
		req, err := http.NewRequest("POST", apiurl+"/v2/unmarshal", body)
		if err != nil {
			t.Error(err)
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			t.Error(resp.StatusCode, buf.String())
			return
		}
		var result interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Error(err)
			return
		}
		if !reflect.DeepEqual(result, expectedResult) {
			t.Error("\n", result, "\n", expectedResult)
			return
		}
	}
}

func testUnmarshalAny(apiurl string, request api.UnmarshallingV2Request, expectedResult []interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		body := new(bytes.Buffer)
		err := json.NewEncoder(body).Encode(request)
		if err != nil {
			t.Error(err)
			return
		}
		req, err := http.NewRequest("POST", apiurl+"/v2/unmarshal", body)
		if err != nil {
			t.Error(err)
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			t.Error(resp.StatusCode, buf.String())
			return
		}
		var result interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Error(err)
			return
		}
		found := false
		for _, expected := range expectedResult {
			if reflect.DeepEqual(result, expected) {
				found = true
			}
		}
		if !found {
			t.Error("\n", result, "\n", expectedResult)
			return
		}
	}
}
