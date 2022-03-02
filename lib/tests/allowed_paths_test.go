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

package tests

import (
	"github.com/SENERGY-Platform/converter/lib/converter/characteristics"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"reflect"
	"testing"
)

func TestMarshallingWithAllowedPaths(t *testing.T) {
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
					Type: model.Structure,
					SubContentVariables: []model.ContentVariable{
						{
							Id:               "inside",
							Name:             "inside",
							Type:             model.Float,
							CharacteristicId: characteristics.Celsius,
							FunctionId:       model.CONTROLLING_FUNCTION_PREFIX + "setTemperature",
							AspectId:         "inside_air",
							Value:            12.0,
						},
						{
							Id:               "outside",
							Name:             "outside",
							Type:             model.Float,
							CharacteristicId: characteristics.Celsius,
							FunctionId:       model.CONTROLLING_FUNCTION_PREFIX + "setTemperature",
							AspectId:         "outside_air",
							Value:            13.0,
						},
					},
				},
				Serialization:     "json",
				ProtocolSegmentId: "p1.1",
			},
		},
	}
	t.Run("k nil", testMarshalInputs(protocol, service, 300, characteristics.Kelvin, nil, map[string]string{"body": `{"inside":27,"outside":27}`}))
	t.Run("k []", testMarshalInputs(protocol, service, 300, characteristics.Kelvin, []string{}, map[string]string{"body": `{"inside":27,"outside":27}`}))
	t.Run("k inside", testMarshalInputs(protocol, service, 300, characteristics.Kelvin, []string{"temperature.inside"}, map[string]string{"body": `{"inside":27,"outside":13}`}))
	t.Run("k outside", testMarshalInputs(protocol, service, 300, characteristics.Kelvin, []string{"temperature.outside"}, map[string]string{"body": `{"inside":12,"outside":27}`}))
	t.Run("k inside and outside", testMarshalInputs(protocol, service, 300, characteristics.Kelvin, []string{"temperature.inside", "temperature.outside"}, map[string]string{"body": `{"inside":27,"outside":27}`}))

	t.Run("c nil", testMarshalInputs(protocol, service, 300, characteristics.Celsius, nil, map[string]string{"body": `{"inside":300,"outside":300}`}))
	t.Run("c []", testMarshalInputs(protocol, service, 300, characteristics.Celsius, []string{}, map[string]string{"body": `{"inside":300,"outside":300}`}))
	t.Run("c inside", testMarshalInputs(protocol, service, 300, characteristics.Celsius, []string{"temperature.inside"}, map[string]string{"body": `{"inside":300,"outside":13}`}))
	t.Run("c outside", testMarshalInputs(protocol, service, 300, characteristics.Celsius, []string{"temperature.outside"}, map[string]string{"body": `{"inside":12,"outside":300}`}))
	t.Run("c inside and outside", testMarshalInputs(protocol, service, 300, characteristics.Celsius, []string{"temperature.inside", "temperature.outside"}, map[string]string{"body": `{"inside":300,"outside":300}`}))

}

func testMarshalInputs(protocol model.Protocol, service model.Service, value interface{}, valueCharacteristic string, paths []string, expected map[string]string) func(t *testing.T) {
	return func(t *testing.T) {
		result, err := TestMarshalInputs(protocol, service, value, valueCharacteristic, paths)
		if err != nil {
			t.Error(err)
			return
		}
		if !reflect.DeepEqual(result, expected) {
			t.Error(result)
		}
	}
}

func testUnmarshalOutputs(protocol model.Protocol, service model.Service, value map[string]string, valueCharacteristic string, paths []string, expected interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		result, err := TestUnmarshalOutputs(protocol, service, value, valueCharacteristic, paths)
		if err != nil {
			t.Error(err)
			return
		}
		if !reflect.DeepEqual(result, expected) {
			t.Error(result)
		}
	}
}

func testUnmarshalOutputsOneOf(protocol model.Protocol, service model.Service, value map[string]string, valueCharacteristic string, paths []string, oneOfExpected []interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		result, err := TestUnmarshalOutputs(protocol, service, value, valueCharacteristic, paths)
		if err != nil {
			t.Error(err)
			return
		}
		matchFound := false
		for _, expected := range oneOfExpected {
			if reflect.DeepEqual(result, expected) {
				matchFound = true
				break
			}
		}
		if !matchFound {
			t.Error(result)
		}
	}
}

func TestUnmarshallingWithAllowedPaths(t *testing.T) {
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

	t.Run("k nil", testUnmarshalOutputsOneOf(protocol, service, output, characteristics.Kelvin, nil, []interface{}{673.15, 773.15}))
	t.Run("k []", testUnmarshalOutputsOneOf(protocol, service, output, characteristics.Kelvin, []string{}, []interface{}{673.15, 773.15}))
	t.Run("k inside", testUnmarshalOutputs(protocol, service, output, characteristics.Kelvin, []string{"temperature.inside"}, 673.15))
	t.Run("k outside", testUnmarshalOutputs(protocol, service, output, characteristics.Kelvin, []string{"temperature.outside"}, 773.15))
	t.Run("k inside and outside", testUnmarshalOutputsOneOf(protocol, service, output, characteristics.Kelvin, []string{"temperature.inside", "temperature.outside"}, []interface{}{673.15, 773.15}))

	t.Run("c nil", testUnmarshalOutputsOneOf(protocol, service, output, characteristics.Celsius, nil, []interface{}{400.0, 500.0}))
	t.Run("c []", testUnmarshalOutputsOneOf(protocol, service, output, characteristics.Celsius, []string{}, []interface{}{400.0, 500.0}))
	t.Run("c inside", testUnmarshalOutputs(protocol, service, output, characteristics.Celsius, []string{"temperature.inside"}, 400.0))
	t.Run("c outside", testUnmarshalOutputs(protocol, service, output, characteristics.Celsius, []string{"temperature.outside"}, 500.0))
	t.Run("c inside and outside", testUnmarshalOutputsOneOf(protocol, service, output, characteristics.Celsius, []string{"temperature.inside", "temperature.outside"}, []interface{}{400.0, 500.0}))

}

func TestMarshallingVoidToggle(t *testing.T) {
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
					Id:         "toggle",
					Name:       "toggle",
					IsVoid:     true,
					Value:      "foo",
					Type:       model.String,
					FunctionId: model.CONTROLLING_FUNCTION_PREFIX + "toggle",
					AspectId:   "",
				},
				Serialization:     "json",
				ProtocolSegmentId: "p1.1",
			},
		},
	}
	t.Run("k nil", testMarshalInputs(protocol, service, nil, "", nil, map[string]string{}))

}
