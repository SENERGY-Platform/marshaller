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

func TestMarshalling(t *testing.T) {
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
					Type: model.Structure,
					SubContentVariables: []model.ContentVariable{
						{
							Id:               "inside",
							Name:             "inside",
							Type:             model.Integer, //results like 26,85 will be rounded to 27
							CharacteristicId: characteristics.Celsius,
							FunctionId:       model.CONTROLLING_FUNCTION_PREFIX + "setTemperature",
							AspectId:         "inside_air",
							Value:            12,
						},
						{
							Id:               "outside",
							Name:             "outside",
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
				Paths:            []string{"temperature.inside"},
				FunctionId:       "",
			},
		},
	}, map[string]string{"body": `{"inside":27,"outside":13}`}))

	t.Run("inside 300 kelvin path and function", testMarshal(apiurl, api.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            300,
				CharacteristicId: characteristics.Kelvin,
				Paths:            []string{"temperature.inside"},
				FunctionId:       model.CONTROLLING_FUNCTION_PREFIX + "setTemperature",
			},
		},
	}, map[string]string{"body": `{"inside":27,"outside":13}`}))

	t.Run("outside 300 kelvin path", testMarshal(apiurl, api.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            300,
				CharacteristicId: characteristics.Kelvin,
				Paths:            []string{"temperature.outside"},
				FunctionId:       "",
			},
		},
	}, map[string]string{"body": `{"inside":12,"outside":27}`}))

	t.Run("outside 300 kelvin path and function", testMarshal(apiurl, api.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            300,
				CharacteristicId: characteristics.Kelvin,
				Paths:            []string{"temperature.outside"},
				FunctionId:       model.CONTROLLING_FUNCTION_PREFIX + "setTemperature",
			},
		},
	}, map[string]string{"body": `{"inside":12,"outside":27}`}))

	t.Run("inside and outside 300 kelvin path", testMarshal(apiurl, api.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            300,
				CharacteristicId: characteristics.Kelvin,
				Paths:            []string{"temperature.inside", "temperature.outside"},
				FunctionId:       "",
			},
		},
	}, map[string]string{"body": `{"inside":27,"outside":27}`}))

	t.Run("inside 400 and outside 500 kelvin path", testMarshal(apiurl, api.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            400,
				CharacteristicId: characteristics.Kelvin,
				Paths:            []string{"temperature.inside"},
				FunctionId:       "",
			},
			{
				Value:            500,
				CharacteristicId: characteristics.Kelvin,
				Paths:            []string{"temperature.outside"},
				FunctionId:       "",
			},
		},
	}, map[string]string{"body": `{"inside":127,"outside":227}`}))

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
	}, map[string]string{"body": `{"inside":27,"outside":27}`}))
}

func testMarshal(apiurl string, request api.MarshallingV2Request, expectedResult map[string]string) func(t *testing.T) {
	return func(t *testing.T) {
		body := new(bytes.Buffer)
		err := json.NewEncoder(body).Encode(request)
		if err != nil {
			t.Error(err)
			return
		}
		req, err := http.NewRequest("POST", apiurl+"/v2/marshal", body)
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
		var result map[string]string
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
