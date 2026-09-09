/*
 * Copyright 2026 InfAI (CC SES)
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
	"net/http"
	"sync"
	"testing"

	"github.com/SENERGY-Platform/converter/lib/converter/characteristics"
	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
)

// testUnmarshalNoPath expects the request to be rejected because its criteria select no
// output path at all.
func testUnmarshalNoPath(apiurl string, request messages.UnmarshallingV2Request) func(t *testing.T) {
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
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		if resp.StatusCode != http.StatusBadRequest {
			t.Error(resp.StatusCode, buf.String())
			return
		}
		if !bytes.Contains(buf.Bytes(), []byte("no output path found for criteria")) {
			t.Error(buf.String())
			return
		}
	}
}

// TestUnmarshalAspectIds checks the aspect list of an unmarshal request against a service
// whose two content variables differ only in the aspects they carry. A request naming more
// than one aspect is an AND, so only the variable carrying all of them is a candidate.
func TestUnmarshalAspectIds(t *testing.T) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	apiurl := setup(ctx, wg)

	functionId := model.MEASURING_FUNCTION_PREFIX + "getTemperature"

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
					//inside comes first on purpose: an OR over the queried aspects would
					//match it too and, being the closer path, hand back its value
					SubContentVariables: []model.ContentVariable{
						{
							Id:               "inside",
							Name:             "inside",
							Type:             model.Float,
							CharacteristicId: characteristics.Celsius,
							FunctionId:       functionId,
							AspectIds:        []string{"inside_air"},
						},
						{
							Id:               "inside_today",
							Name:             "inside_today",
							Type:             model.Float,
							CharacteristicId: characteristics.Celsius,
							FunctionId:       functionId,
							AspectIds:        []string{"inside_air", "today"},
						},
					},
				},
				Serialization:     "json",
				ProtocolSegmentId: "p1.1",
			},
		},
	}

	output := map[string]string{"body": `{"inside":500,"inside_today":400}`}

	t.Run("matches only the content variable carrying every queried aspect", testUnmarshal(apiurl, messages.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Celsius,
		Message:          output,
		FunctionId:       functionId,
		AspectNodeIds:    []string{"inside_air", "today"},
	}, 400.0))

	t.Run("covers the subtree of every queried aspect", testUnmarshal(apiurl, messages.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Celsius,
		Message:          output,
		FunctionId:       functionId,
		AspectNodeIds:    []string{"air", "electricity"},
	}, 400.0))

	t.Run("matches every content variable carrying the single queried aspect", testUnmarshalAny(apiurl, messages.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Celsius,
		Message:          output,
		FunctionId:       functionId,
		AspectNodeIds:    []string{"inside_air"},
	}, []interface{}{400.0, 500.0}))

	t.Run("reads the deprecated aspect node id as a list with one element", testUnmarshal(apiurl, messages.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Celsius,
		Message:          output,
		FunctionId:       functionId,
		AspectNodeId:     "today",
	}, 400.0))

	t.Run("reads the deprecated aspect node as a list with one element", testUnmarshal(apiurl, messages.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Celsius,
		Message:          output,
		FunctionId:       functionId,
		AspectNode:       model.AspectNode{Id: "today", RootId: "electricity", ParentId: "consumption"},
	}, 400.0))

	t.Run("adds the deprecated aspect node id to the queried aspect list", testUnmarshal(apiurl, messages.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Celsius,
		Message:          output,
		FunctionId:       functionId,
		AspectNodeId:     "today",
		AspectNodeIds:    []string{"inside_air"},
	}, 400.0))

	t.Run("finds no path if no content variable carries every queried aspect", testUnmarshalNoPath(apiurl, messages.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Celsius,
		Message:          output,
		FunctionId:       functionId,
		AspectNodeIds:    []string{"inside_air", "outside_air"},
	}))
}

// TestUnmarshalDeprecatedContentVariableAspectId checks that a service delivered by a
// device-repository that predates the aspect lists is still matched: the deprecated
// AspectId of a content variable is read as a list with one element.
func TestUnmarshalDeprecatedContentVariableAspectId(t *testing.T) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	apiurl := setup(ctx, wg)

	functionId := model.MEASURING_FUNCTION_PREFIX + "getTemperature"

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
							FunctionId:       functionId,
							AspectId:         "inside_air",
						},
						{
							Id:               "outside",
							Name:             "outside",
							Type:             model.Float,
							CharacteristicId: characteristics.Celsius,
							FunctionId:       functionId,
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

	t.Run("matches the deprecated aspect id of a content variable by aspect list", testUnmarshal(apiurl, messages.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Celsius,
		Message:          output,
		FunctionId:       functionId,
		AspectNodeIds:    []string{"outside_air"},
	}, 500.0))

	t.Run("matches the deprecated aspect id of a content variable by deprecated aspect node id", testUnmarshal(apiurl, messages.UnmarshallingV2Request{
		Service:          service,
		Protocol:         protocol,
		CharacteristicId: characteristics.Celsius,
		Message:          output,
		FunctionId:       functionId,
		AspectNodeId:     "outside_air",
	}, 500.0))
}

// TestMarshalAspectNodes checks the aspect node list of a marshal request. The paths a
// value is written to are selected the same way the unmarshal side selects them, so a
// request naming more than one aspect is an AND here as well.
func TestMarshalAspectNodes(t *testing.T) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	apiurl := setup(ctx, wg)

	functionId := model.CONTROLLING_FUNCTION_PREFIX + "setTemperature"
	insideAir := model.AspectNode{Id: "inside_air", RootId: "air", ParentId: "air", AncestorIds: []string{"air"}}
	today := model.AspectNode{Id: "today", RootId: "electricity", ParentId: "consumption", AncestorIds: []string{"electricity", "consumption"}}

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
					//inside comes first on purpose: an OR over the queried aspects would
					//write to it as well
					SubContentVariables: []model.ContentVariable{
						{
							Id:               "inside",
							Name:             "inside",
							Type:             model.Integer,
							CharacteristicId: characteristics.Celsius,
							FunctionId:       functionId,
							AspectIds:        []string{"inside_air"},
							Value:            12,
						},
						{
							Id:               "inside_today",
							Name:             "inside_today",
							Type:             model.Integer,
							CharacteristicId: characteristics.Celsius,
							FunctionId:       functionId,
							AspectIds:        []string{"inside_air", "today"},
							Value:            13,
						},
					},
				},
				Serialization:     "json",
				ProtocolSegmentId: "p1.1",
			},
		},
	}

	t.Run("writes only to the content variable carrying every queried aspect", testMarshal(apiurl, messages.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            27,
				CharacteristicId: characteristics.Celsius,
				FunctionId:       functionId,
				AspectNodes:      []model.AspectNode{insideAir, today},
			},
		},
	}, map[string]string{"body": `{"inside":12,"inside_today":27}`}))

	t.Run("reads the deprecated aspect node as a list with one element", testMarshal(apiurl, messages.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            27,
				CharacteristicId: characteristics.Celsius,
				FunctionId:       functionId,
				AspectNode:       &today,
			},
		},
	}, map[string]string{"body": `{"inside":12,"inside_today":27}`}))

	t.Run("adds the deprecated aspect node to the queried aspect list", testMarshal(apiurl, messages.MarshallingV2Request{
		Service:  service,
		Protocol: protocol,
		Data: []model.MarshallingV2RequestData{
			{
				Value:            27,
				CharacteristicId: characteristics.Celsius,
				FunctionId:       functionId,
				AspectNode:       &insideAir,
				AspectNodes:      []model.AspectNode{today},
			},
		},
	}, map[string]string{"body": `{"inside":12,"inside_today":27}`}))
}
