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
