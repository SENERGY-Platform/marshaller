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
	"github.com/SENERGY-Platform/marshaller/lib/api"
	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/marshaller/lib/tests/mocks"
	"net/url"
	"testing"
)

func TestConfigurablesPostRequest(t *testing.T) {
	message := api.FindConfigurablesRequest{
		CharacteristicId: temperature.Celcius,
		Services:         []model.Service{serviceExample1, serviceExample2},
	}

	result := []configurables.Configurable{}

	err := postJSON(ServerUrl+"/configurables", message, &result)
	if err != nil {
		t.Error(err)
		return
	}

	assert := Assertions{t}
	assert.ConfigurableListContains(result, configurables.Configurable{
		CharacteristicId: color.Rgb,
		Values: []configurables.ConfigurableCharacteristicValue{
			{
				Label: "RGB r",
				Path:  "r",
				Value: "0",
			},
			{
				Label: "RGB g",
				Path:  "g",
				Value: "0",
			},
			{
				Label: "RGB b",
				Path:  "b",
				Value: "0",
			},
		},
	})
}

func TestConfigurablesMockGetRequest(t *testing.T) {
	if !testing.Short() {
		t.Skip("only with mocks (use '-short' test argument)")
	}
	mocks.DeviceRepo.Init().SetService(serviceExample1).SetService(serviceExample2)

	result := []configurables.Configurable{}

	err := getJSON(ServerUrl+"/configurables?characteristicId="+url.QueryEscape(temperature.Celcius)+"&serviceIds="+url.QueryEscape(serviceExample1.Id)+","+url.QueryEscape(serviceExample2.Id), &result)
	if err != nil {
		t.Error(err)
		return
	}

	assert := Assertions{t}
	assert.ConfigurableListContains(result, configurables.Configurable{
		CharacteristicId: color.Rgb,
		Values: []configurables.ConfigurableCharacteristicValue{
			{
				Label: "RGB r",
				Path:  "r",
				Value: "0",
			},
			{
				Label: "RGB g",
				Path:  "g",
				Value: "0",
			},
			{
				Label: "RGB b",
				Path:  "b",
				Value: "0",
			},
		},
	})
}

var serviceExample1 = model.Service{
	Id:          "s1",
	LocalId:     "s1l",
	Name:        "s1n",
	Description: "s1d",
	ProtocolId:  "p1",
	Inputs: []model.Content{
		{
			Id: "c1",
			ContentVariable: model.ContentVariable{
				Id:   "c1.1",
				Name: "payload",
				Type: model.Structure,
				SubContentVariables: []model.ContentVariable{
					{
						Id:               "c1.1.2",
						Name:             "temperature",
						Type:             model.Float,
						CharacteristicId: temperature.Celcius,
					},
					{
						Id:               "c1.1.1",
						Name:             "color",
						Type:             model.Structure,
						CharacteristicId: color.Rgb,
						SubContentVariables: []model.ContentVariable{
							{
								Id:   "c2.1.4",
								Name: "foo",
								Type: model.String,
							},
							{
								Id:               "sr",
								Name:             "red",
								Type:             model.Integer,
								CharacteristicId: color.RgbR,
							},
							{
								Id:               "sg",
								Name:             "green",
								Type:             model.Integer,
								CharacteristicId: color.RgbG,
							},
							{
								Id:               "sb",
								Name:             "blue",
								Type:             model.Integer,
								CharacteristicId: color.RgbB,
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

var serviceExample2 = model.Service{
	Id:          "s2",
	LocalId:     "s2l",
	Name:        "s2n",
	Description: "s2d",
	ProtocolId:  "p1",
	Inputs: []model.Content{
		{
			Id: "c2",
			ContentVariable: model.ContentVariable{
				Id:   "c2.1",
				Name: "payload",
				Type: model.Structure,
				SubContentVariables: []model.ContentVariable{
					{
						Id:   "c2.1.4",
						Name: "foo",
						Type: model.String,
					},
					{
						Id:   "c2.1.3",
						Name: "bar",
						Type: model.String,
					},
					{
						Id:               "c2.1.2",
						Name:             "temperature",
						Type:             model.Float,
						CharacteristicId: temperature.Celcius,
					},
					{
						Id:               "c2.1.1",
						Name:             "color",
						Type:             model.String,
						CharacteristicId: color.Hex,
					},
					{
						Id:               "c2.1.5",
						Name:             "color_2",
						Type:             model.String,
						CharacteristicId: color.Hex,
					},
				},
			},
			Serialization:     "json",
			ProtocolSegmentId: "p1.1",
		},
	},
}
