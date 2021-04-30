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

package tests

import (
	"fmt"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
)

func ExampleHints1() {
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
		Id:          "s1",
		LocalId:     "s1l",
		Name:        "s1n",
		Description: "s1d",
		ProtocolId:  "p1",
		Outputs: []model.Content{
			{
				Id: "c1",
				ContentVariable: model.ContentVariable{
					Id:   "c1.1",
					Name: "payload",
					Type: model.Structure,
					SubContentVariables: []model.ContentVariable{
						{
							Id:   "c1.1.1",
							Name: "color1",
							Type: model.Structure,
							SubContentVariables: []model.ContentVariable{
								{
									Id:               "sr",
									Name:             "red",
									Type:             model.Integer,
									CharacteristicId: example.Rgb + ".r",
								},
								{
									Id:               "sg",
									Name:             "green",
									Type:             model.Integer,
									CharacteristicId: example.Rgb + ".g",
								},
								{
									Id:               "sb",
									Name:             "blue",
									Type:             model.Integer,
									CharacteristicId: example.Rgb + ".b",
								},
							},
						},
						{
							Id:   "c1.1.2",
							Name: "color2",
							Type: model.Structure,
							SubContentVariables: []model.ContentVariable{
								{
									Id:               "sr",
									Name:             "red",
									Type:             model.Integer,
									CharacteristicId: example.Rgb + ".r",
								},
								{
									Id:               "sg",
									Name:             "green",
									Type:             model.Integer,
									CharacteristicId: example.Rgb + ".g",
								},
								{
									Id:               "sb",
									Name:             "blue",
									Type:             model.Integer,
									CharacteristicId: example.Rgb + ".b",
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
	output := map[string]string{"body": "{\"color1\":{\"blue\":100,\"green\":0,\"red\":255}, \"color2\":{\"blue\":100,\"green\":255,\"red\":0}}"}
	outputCharacteristic := example.Hex
	fmt.Println(TestUnmarshalOutputs(protocol, service, output, outputCharacteristic, "c1.1.1"))
	fmt.Println(TestUnmarshalOutputs(protocol, service, output, outputCharacteristic, "c1.1.1", "foo", "bar"))
	fmt.Println(TestUnmarshalOutputs(protocol, service, output, outputCharacteristic, "c1.1.2"))
	fmt.Println(TestUnmarshalOutputs(protocol, service, output, outputCharacteristic, "c1.1.2", "foo", "bar"))
	fmt.Println(TestUnmarshalOutputs(protocol, service, output, outputCharacteristic, "foo", "bar"))

	//output:
	//#ff0064 <nil>
	//#ff0064 <nil>
	//#00ff64 <nil>
	//#00ff64 <nil>
	//#00ff64 <nil>
}
