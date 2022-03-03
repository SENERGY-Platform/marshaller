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

func TestRgbMarshal(t *testing.T) {
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
					Id:   "color",
					Name: "color",
					Type: model.Structure,
					SubContentVariables: []model.ContentVariable{
						{
							Id:               "rgb",
							Name:             "rgb",
							Type:             model.Structure,
							CharacteristicId: characteristics.Rgb,
							FunctionId:       model.CONTROLLING_FUNCTION_PREFIX + "setColor",
							AspectId:         "device",
							SubContentVariables: []model.ContentVariable{
								{
									Id:               "rot",
									Name:             "rot",
									Type:             model.Integer,
									CharacteristicId: characteristics.RgbR,
								},
								{
									Id:               "grün",
									Name:             "grün",
									Type:             model.Integer,
									CharacteristicId: characteristics.RgbG,
								},
								{
									Id:               "blau",
									Name:             "blau",
									Type:             model.Integer,
									CharacteristicId: characteristics.RgbB,
								},
							},
						},
					},
				},
				Serialization:     "json",
				ProtocolSegmentId: "p1.1",
			},
		},
	}

	t.Run("marshal hex to rgb", testMarshal(apiurl, api.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            "#ff0064",
				CharacteristicId: characteristics.Hex,
				Paths:            []string{"color.rgb"},
			},
		},
	}, map[string]string{"body": `{"rgb":{"blau":100,"grün":0,"rot":255}}`}))

	t.Run("marshal rgb to rgb", testMarshal(apiurl, api.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            map[string]interface{}{"r": 255, "g": 0, "b": 100},
				CharacteristicId: characteristics.Rgb,
				Paths:            []string{"color.rgb"},
			},
		},
	}, map[string]string{"body": `{"rgb":{"blau":100,"grün":0,"rot":255}}`}))
}

func TestRgbUnmarshal(t *testing.T) {
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
					Id:   "color",
					Name: "color",
					Type: model.Structure,
					SubContentVariables: []model.ContentVariable{
						{
							Id:               "rgb",
							Name:             "rgb",
							Type:             model.Structure,
							CharacteristicId: characteristics.Rgb,
							FunctionId:       model.CONTROLLING_FUNCTION_PREFIX + "getColor",
							AspectId:         "device",
							SubContentVariables: []model.ContentVariable{
								{
									Id:               "rot",
									Name:             "rot",
									Type:             model.Integer,
									CharacteristicId: characteristics.RgbR,
								},
								{
									Id:               "grün",
									Name:             "grün",
									Type:             model.Integer,
									CharacteristicId: characteristics.RgbG,
								},
								{
									Id:               "blau",
									Name:             "blau",
									Type:             model.Integer,
									CharacteristicId: characteristics.RgbB,
								},
							},
						},
					},
				},
				Serialization:     "json",
				ProtocolSegmentId: "p1.1",
			},
		},
	}

	output := map[string]string{"body": `{"rgb":{"blau":100,"grün":0,"rot":255}}`}

	t.Run("inside path", testUnmarshal(apiurl, api.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Hex,
		Message:          output,
		Path:             "color.rgb",
	}, "#ff0064"))
}
