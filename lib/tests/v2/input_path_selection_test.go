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
	"context"
	"sync"
	"testing"

	"github.com/SENERGY-Platform/converter/lib/converter/characteristics"
	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	v2 "github.com/SENERGY-Platform/marshaller/lib/marshaller/v2"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
)

// TestMarshalInputPathSelection checks which of the input paths matching a controlling
// function and its aspects receive the value: by default only the closest ones, with
// v2.AllInputPaths every match.
func TestMarshalInputPathSelection(t *testing.T) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := setup(ctx, wg)

	functionId := model.CONTROLLING_FUNCTION_PREFIX + "setTemperature"
	air := model.AspectNode{Id: "air", RootId: "air", ChildIds: []string{"inside_air", "outside_air"}, DescendentIds: []string{"inside_air", "outside_air"}}
	insideAir := model.AspectNode{Id: "inside_air", RootId: "air", ParentId: "air", AncestorIds: []string{"air"}}

	protocol := model.Protocol{
		Id:               "p1",
		Name:             "p1",
		Handler:          "p1",
		ProtocolSegments: []model.ProtocolSegment{{Id: "p1.1", Name: "body"}},
	}
	variable := func(name string, aspectId string, value int) model.ContentVariable {
		return model.ContentVariable{
			Id:               name,
			Name:             name,
			Type:             model.Integer,
			CharacteristicId: characteristics.Celsius,
			FunctionId:       functionId,
			AspectIds:        []string{aspectId},
			Value:            value,
		}
	}
	service := model.Service{
		Id:          "sid",
		LocalId:     "slid",
		Name:        "sname",
		Interaction: model.REQUEST,
		ProtocolId:  "p1",
		Inputs: []model.Content{{
			Id: "content",
			ContentVariable: model.ContentVariable{
				Id:   "temperature",
				Name: "temperature",
				Type: model.Structure,
				SubContentVariables: []model.ContentVariable{
					variable("inside", "inside_air", 11),
					variable("outside", "outside_air", 12),
					variable("air", "air", 10),
				},
			},
			Serialization:     "json",
			ProtocolSegmentId: "p1.1",
		}},
	}
	request := func(data ...model.MarshallingV2RequestData) messages.MarshallingV2Request {
		return messages.MarshallingV2Request{Service: service, Protocol: protocol, Data: data}
	}
	setTemperature := func(aspectNodes ...model.AspectNode) model.MarshallingV2RequestData {
		return model.MarshallingV2RequestData{
			Value:            27,
			CharacteristicId: characteristics.Celsius,
			FunctionId:       functionId,
			AspectNodes:      aspectNodes,
		}
	}

	t.Run("a parent aspect sets only the variable carrying it", testMarshal(c, request(setTemperature(air)),
		map[string]string{"body": `{"air":27,"inside":11,"outside":12}`}))

	t.Run("a child aspect sets only its variable", testMarshal(c, request(setTemperature(insideAir)),
		map[string]string{"body": `{"air":10,"inside":27,"outside":12}`}))

	//all matches lie at the same distance, so all of them are the closest; this is also the
	//case with more than one path, which must write every one of them into the same input
	t.Run("without aspect every variable of the function is set", testMarshal(c, request(setTemperature()),
		map[string]string{"body": `{"air":27,"inside":27,"outside":27}`}))

	t.Run("a request matching no path leaves the inputs to the next one", testMarshal(c, request(
		model.MarshallingV2RequestData{
			Value:            30,
			CharacteristicId: characteristics.Celsius,
			FunctionId:       model.CONTROLLING_FUNCTION_PREFIX + "unknown",
		},
		setTemperature(insideAir),
	), map[string]string{"body": `{"air":10,"inside":27,"outside":12}`}))

	t.Run("all input paths", func(t *testing.T) {
		v2.InputPathSelection = v2.AllInputPaths
		defer func() { v2.InputPathSelection = v2.ClosestInputPaths }()
		testMarshal(c, request(setTemperature(air)),
			map[string]string{"body": `{"air":27,"inside":27,"outside":27}`})(t)
	})
}
