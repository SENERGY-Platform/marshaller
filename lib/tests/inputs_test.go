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
	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
)

func ExampleConfigurableService_Find_configurable1() {
	protocol := model.Protocol{
		Id:      "p1",
		Name:    "p1",
		Handler: "p1",
		ProtocolSegments: []model.ProtocolSegment{
			{Id: "p1.1", Name: "body"},
			{Id: "p1.2", Name: "head"},
		},
	}

	service1 := model.Service{
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
							CharacteristicId: temperature.Celsius,
						},
						{
							Id:               "c1.1.1",
							Name:             "color",
							Type:             model.Structure,
							CharacteristicId: color.Rgb,
							SubContentVariables: []model.ContentVariable{
								{
									Id:    "c2.1.4",
									Name:  "foo",
									Type:  model.String,
									Value: "bar1",
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

	service2 := model.Service{
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
							Id:    "c2.1.4",
							Name:  "foo",
							Value: "bar2",
							Type:  model.String,
						},
						{
							Id:    "c2.1.3",
							Name:  "bar",
							Value: "foo2",
							Type:  model.String,
						},
						{
							Id:               "c2.1.2",
							Name:             "temperature",
							Type:             model.Float,
							CharacteristicId: temperature.Celsius,
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

	configurblesList, err := TestFindConfigurables(temperature.Celsius, []model.Service{service1, service2})
	if err != nil {
		fmt.Println(err, configurblesList)
		return
	}
	for index, value := range configurblesList[0].Values {
		if value.Path == "r" {
			configurblesList[0].Values[index].Value = "255"
		}
		if value.Path == "g" {
			configurblesList[0].Values[index].Value = "255"
		}
		if value.Path == "b" {
			configurblesList[0].Values[index].Value = "0"
		}
	}

	fmt.Println(TestMarshalInputs(protocol, service1, 37, temperature.Celsius, nil, configurblesList...))
	fmt.Println(TestMarshalInputs(protocol, service2, 37, temperature.Celsius, nil, configurblesList...))

	//output:
	//map[body:{"color":{"blue":0,"foo":"bar1","green":255,"red":255},"temperature":37}] <nil>
	//map[body:{"bar":"foo2","color":"#ffff00","color_2":"#ffff00","foo":"bar2","temperature":37}] <nil>
}

func ExampleConfigurableService_Find_configurableUnused() {
	protocol := model.Protocol{
		Id:      "p1",
		Name:    "p1",
		Handler: "p1",
		ProtocolSegments: []model.ProtocolSegment{
			{Id: "p1.1", Name: "body"},
			{Id: "p1.2", Name: "head"},
		},
	}

	service1 := model.Service{
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
							CharacteristicId: temperature.Celsius,
						},
						{
							Id:               "c1.1.1",
							Name:             "color",
							Type:             model.Structure,
							CharacteristicId: color.Rgb,
							SubContentVariables: []model.ContentVariable{
								{
									Id:    "c2.1.4",
									Name:  "foo",
									Type:  model.String,
									Value: "bar1",
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

	fmt.Println(TestMarshalInputs(protocol, service1, 37, temperature.Celsius, nil, configurables.Configurable{
		CharacteristicId: example.Lux,
		Values: []configurables.ConfigurableCharacteristicValue{
			{
				Label: "nope",
				Path:  "",
				Value: "20",
			},
		},
	}))

	//output:
	//map[body:{"color":{"blue":0,"foo":"bar1","green":0,"red":0},"temperature":37}] <nil>
}

func ExampleMarshaller_MarshalInputs_marshalInput1() {
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
		Inputs: []model.Content{
			{
				Id: "c1",
				ContentVariable: model.ContentVariable{
					Id:   "c1.1",
					Name: "payload",
					Type: model.Structure,
					SubContentVariables: []model.ContentVariable{
						{
							Id:   "c1.1.1",
							Name: "color",
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
	input := "#ff0064"
	inputCharacteristic := example.Hex
	result, err := TestMarshalInputs(protocol, service, input, inputCharacteristic, nil)
	fmt.Println(result, err)

	//output:
	//map[body:{"color":{"blue":100,"green":0,"red":255}}] <nil>
}

func ExampleMarshaller_MarshalInputs_marshalMultiInput1() {
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
		Inputs: []model.Content{
			{
				Id: "c1",
				ContentVariable: model.ContentVariable{
					Id:   "c1.1",
					Name: "payload",
					Type: model.Structure,
					SubContentVariables: []model.ContentVariable{
						{
							Id:               "c1.1.1",
							Name:             "color",
							Type:             model.Structure,
							CharacteristicId: example.Rgb,
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
							Id:               "c1.1.2",
							Name:             "color_2",
							Type:             model.Structure,
							CharacteristicId: example.Rgb,
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
	input := "#ff0064"
	inputCharacteristic := example.Hex
	result, err := TestMarshalInputs(protocol, service, input, inputCharacteristic, nil)
	fmt.Println(result, err)

	//output:
	//map[body:{"color":{"blue":100,"green":0,"red":255},"color_2":{"blue":100,"green":0,"red":255}}] <nil>
}

func ExampleMarshaller_MarshalInputs_marshalMultiInput2() {
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
		Inputs: []model.Content{
			{
				Id: "c1",
				ContentVariable: model.ContentVariable{
					Id:   "c1.1",
					Name: "payload",
					Type: model.Structure,
					SubContentVariables: []model.ContentVariable{
						{
							Id:               "c1.1.1",
							Name:             "color",
							Type:             model.Structure,
							CharacteristicId: example.Rgb,
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
							Id:               "c1.1.2",
							Name:             "color_2",
							Type:             model.String,
							CharacteristicId: example.Hex,
						},
					},
				},
				Serialization:     "json",
				ProtocolSegmentId: "p1.1",
			},
		},
	}
	input := "#ff0064"
	inputCharacteristic := example.Hex
	result, err := TestMarshalInputs(protocol, service, input, inputCharacteristic, nil)
	fmt.Println(result, err)

	//output:
	//map[body:{"color":{"blue":100,"green":0,"red":255},"color_2":"#ff0064"}] <nil>
}

func ExampleMarshaller_MarshalInputs_marshalInput2() {
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
		Inputs: []model.Content{
			{
				Id: "c1",
				ContentVariable: model.ContentVariable{
					Id:   "c1.1",
					Name: "payload",
					Type: model.Structure,
					SubContentVariables: []model.ContentVariable{
						{
							Id:               "c1.1.1",
							Name:             "color",
							CharacteristicId: example.Hex,
							Type:             model.String,
						},
					},
				},
				Serialization:     "json",
				ProtocolSegmentId: "p1.1",
			},
		},
	}
	input := map[string]interface{}{
		"r": float64(255),
		"g": float64(0),
		"b": float64(100),
	}
	inputCharacteristic := example.Rgb
	result, err := TestMarshalInputs(protocol, service, input, inputCharacteristic, nil)
	fmt.Println(result, err)

	//output:
	//map[body:{"color":"#ff0064"}] <nil>
}

func ExampleMarshaller_MarshalInputs_marshalInputMulti() {
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
		Inputs: []model.Content{
			{
				Id: "c1",
				ContentVariable: model.ContentVariable{
					Id:   "c1.1",
					Name: "payload",
					Type: model.Structure,
					SubContentVariables: []model.ContentVariable{
						{
							Id:   "c1.1.1",
							Name: "color",
							Type: model.Structure,
							SubContentVariables: []model.ContentVariable{
								{
									Id:               "sr",
									Name:             "red",
									Type:             model.Integer,
									CharacteristicId: example.Rgb + ".r",
									Value:            float64(255),
								},
								{
									Id:               "sg",
									Name:             "green",
									Type:             model.Integer,
									CharacteristicId: example.Rgb + ".g",
									Value:            float64(255),
								},
								{
									Id:               "sb",
									Name:             "blue",
									Type:             model.Integer,
									CharacteristicId: example.Rgb + ".b",
									Value:            float64(255),
								},
							},
						},
						{
							Id:               "c1.1.2",
							Name:             "bri",
							Type:             model.Integer,
							CharacteristicId: example.Lux,
							Value:            float64(100),
						},
					},
				},
				Serialization:     "json",
				ProtocolSegmentId: "p1.1",
			},
		},
	}
	fmt.Println(TestMarshalInputs(protocol, service, "#ff0064", example.Hex, nil))
	fmt.Println(TestMarshalInputs(protocol, service, float64(25), example.Lux, nil))

	//output:
	//map[body:{"bri":100,"color":{"blue":100,"green":0,"red":255}}] <nil>
	//map[body:{"bri":25,"color":{"blue":255,"green":255,"red":255}}] <nil>
}

func ExampleMarshaller_MarshalInputs_marshalInputMultiXml() {
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
		Inputs: []model.Content{
			{
				Id: "c1",
				ContentVariable: model.ContentVariable{
					Id:   "c1.1",
					Name: "payload",
					Type: model.Structure,
					SubContentVariables: []model.ContentVariable{
						{
							Id:   "c1.1.1",
							Name: "color",
							Type: model.Structure,
							SubContentVariables: []model.ContentVariable{
								{
									Id:               "sr",
									Name:             "-red",
									Type:             model.Integer,
									CharacteristicId: example.Rgb + ".r",
									Value:            float64(255),
								},
								{
									Id:               "sg",
									Name:             "-green",
									Type:             model.Integer,
									CharacteristicId: example.Rgb + ".g",
									Value:            float64(255),
								},
								{
									Id:               "sb",
									Name:             "-blue",
									Type:             model.Integer,
									CharacteristicId: example.Rgb + ".b",
									Value:            float64(255),
								},
							},
						},
						{
							Id:               "c1.1.2",
							Name:             "bri",
							Type:             model.Integer,
							CharacteristicId: example.Lux,
							Value:            float64(100),
						},
					},
				},
				Serialization:     "xml",
				ProtocolSegmentId: "p1.1",
			},
		},
	}
	fmt.Println(TestMarshalInputs(protocol, service, "#ff0064", example.Hex, nil))
	fmt.Println(TestMarshalInputs(protocol, service, float64(25), example.Lux, nil))

	//output:
	//map[body:<payload><bri>100</bri><color blue="100" green="0" red="255"/></payload>] <nil>
	//map[body:<payload><bri>25</bri><color blue="255" green="255" red="255"/></payload>] <nil>
}

func ExampleMarshaller_MarshalInputs_marshalInputNull() {
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
		Inputs: []model.Content{
			{
				Id: "c1",
				ContentVariable: model.ContentVariable{
					Id:   "c1.1",
					Name: "payload",
					Type: model.Structure,
					SubContentVariables: []model.ContentVariable{
						{
							Id:   "c1.1.1",
							Name: "color",
							Type: model.Structure,
							SubContentVariables: []model.ContentVariable{
								{
									Id:               "sr",
									Name:             "-red",
									Type:             model.Integer,
									CharacteristicId: example.Rgb + ".r",
									Value:            float64(255),
								},
								{
									Id:               "sg",
									Name:             "-green",
									Type:             model.Integer,
									CharacteristicId: example.Rgb + ".g",
									Value:            float64(0),
								},
								{
									Id:               "sb",
									Name:             "-blue",
									Type:             model.Integer,
									CharacteristicId: example.Rgb + ".b",
									Value:            float64(100),
								},
							},
						},
						{
							Id:               "c1.1.2",
							Name:             "bri",
							Type:             model.Integer,
							CharacteristicId: example.Lux,
							Value:            float64(100),
						},
					},
				},
				Serialization:     "xml",
				ProtocolSegmentId: "p1.1",
			},
		},
	}
	fmt.Println(TestMarshalInputs(protocol, service, nil, "", nil))
	fmt.Println(TestMarshalInputs(protocol, service, nil, model.NullCharacteristic.Id, nil))
	fmt.Println(TestMarshalInputs(protocol, service, "something", "", nil))
	fmt.Println(TestMarshalInputs(protocol, service, "something", model.NullCharacteristic.Id, nil))
	fmt.Println(TestMarshalInputs(protocol, service, map[string]string{"foo": "bar"}, "", nil))
	fmt.Println(TestMarshalInputs(protocol, service, map[string]string{"foo": "bar"}, model.NullCharacteristic.Id, nil))
	fmt.Println(TestMarshalInputs(protocol, service, "#ff0064", "", nil))
	fmt.Println(TestMarshalInputs(protocol, service, "#ff0064", model.NullCharacteristic.Id, nil))

	//output:
	//map[body:<payload><bri>100</bri><color blue="100" green="0" red="255"/></payload>] <nil>
	//map[body:<payload><bri>100</bri><color blue="100" green="0" red="255"/></payload>] <nil>
	//map[body:<payload><bri>100</bri><color blue="100" green="0" red="255"/></payload>] <nil>
	//map[body:<payload><bri>100</bri><color blue="100" green="0" red="255"/></payload>] <nil>
	//map[body:<payload><bri>100</bri><color blue="100" green="0" red="255"/></payload>] <nil>
	//map[body:<payload><bri>100</bri><color blue="100" green="0" red="255"/></payload>] <nil>
	//map[body:<payload><bri>100</bri><color blue="100" green="0" red="255"/></payload>] <nil>
	//map[body:<payload><bri>100</bri><color blue="100" green="0" red="255"/></payload>] <nil>
}

func ExampleMarshaller_MarshalInputs_marshalEmptyService() {
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
	}
	fmt.Println(TestMarshalInputs(protocol, service, nil, "", nil))
	fmt.Println(TestMarshalInputs(protocol, service, nil, model.NullCharacteristic.Id, nil))
	fmt.Println(TestMarshalInputs(protocol, service, "something", "", nil))
	fmt.Println(TestMarshalInputs(protocol, service, "something", model.NullCharacteristic.Id, nil))
	fmt.Println(TestMarshalInputs(protocol, service, map[string]string{"foo": "bar"}, "", nil))
	fmt.Println(TestMarshalInputs(protocol, service, map[string]string{"foo": "bar"}, model.NullCharacteristic.Id, nil))
	fmt.Println(TestMarshalInputs(protocol, service, "#ff0064", "", nil))
	fmt.Println(TestMarshalInputs(protocol, service, "#ff0064", model.NullCharacteristic.Id, nil))

	fmt.Println(TestMarshalInputs(protocol, service, "#ff0064", example.Hex, nil))
	fmt.Println(TestMarshalInputs(protocol, service, map[string]interface{}{"r": float64(255), "g": float64(0), "b": float64(100)}, example.Rgb, nil))

	//output:
	//map[] <nil>
	//map[] <nil>
	//map[] <nil>
	//map[] <nil>
	//map[] <nil>
	//map[] <nil>
	//map[] <nil>
	//map[] <nil>
	//map[] <nil>
	//map[] <nil>
}
