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
	"encoding/json"
	"fmt"
	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"reflect"
	"testing"
)

func ExampleFindConfigurablesShort() {
	if !testing.Short() {
		//skip
		fmt.Println(`[{"characteristic_id":"urn:infai:ses:characteristic:5b4eea52-e8e5-4e80-9455-0382f81a1b43","values":[{"label":"RGB r","path":"r","value":"0"},{"label":"RGB g","path":"g","value":"0"},{"label":"RGB b","path":"b","value":"0"}]}]`)
	} else {
		exampleFindConfigurables()
	}

	//output:
	//[{"characteristic_id":"urn:infai:ses:characteristic:5b4eea52-e8e5-4e80-9455-0382f81a1b43","values":[{"label":"RGB b","path":"b","value":"0"},{"label":"RGB g","path":"g","value":"0"},{"label":"RGB r","path":"r","value":"0"}]}]
}

func ExampleFindConfigurablesLong() {
	if testing.Short() {
		//skip
		fmt.Println(`[{"characteristic_id":"urn:infai:ses:characteristic:5b4eea52-e8e5-4e80-9455-0382f81a1b43","values":[{"label":"RGB b","path":"b","value":"0"},{"label":"RGB g","path":"g","value":"0"},{"label":"RGB r","path":"r","value":"0"}]}]`)
	} else {
		exampleFindConfigurables()
	}

	//output:
	//[{"characteristic_id":"urn:infai:ses:characteristic:5b4eea52-e8e5-4e80-9455-0382f81a1b43","values":[{"label":"RGB b","path":"b","value":"0"},{"label":"RGB g","path":"g","value":"0"},{"label":"RGB r","path":"r","value":"0"}]}]
}

func exampleFindConfigurables() {
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

	configurablesList, err := TestFindConfigurables(temperature.Celcius, []model.Service{service1, service2})
	if err != nil {
		fmt.Println(err)
		return
	}
	temp, _ := json.Marshal(configurablesList)
	fmt.Println(string(temp))

	//output:
	//[{"characteristic_id":"urn:infai:ses:characteristic:5b4eea52-e8e5-4e80-9455-0382f81a1b43","values":[{"label":"","path":"b","value":0},{"label":"","path":"g","value":0},{"label":"","path":"r","value":0}]}]

}

func TestFindIntersectingConfigurables1(t *testing.T) {
	assert := Assertions{t}
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

	configurablesList, err := TestFindConfigurables(temperature.Celcius, []model.Service{service1, service2})
	if err != nil {
		t.Fatal(err)
	}

	if len(configurablesList) != 1 {
		t.Fatal(configurablesList)
	}

	assert.ConfigurableListContains(configurablesList, configurables.Configurable{
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

type Assertions struct {
	*testing.T
}

func (this Assertions) ListContains(list interface{}, element interface{}) {
	this.Helper()
	if !this.listContains(list, element) {
		this.Fatal("missing element in list", list, element)
	}
}

func (this Assertions) listContains(list interface{}, element interface{}) bool {
	listValue := reflect.ValueOf(list)
	for i := 0; i < listValue.Len(); i++ {
		if reflect.DeepEqual(listValue.Index(i).Interface(), element) {
			return true
		}
	}
	return false
}

func (this Assertions) ConfigurableListContains(list []configurables.Configurable, element configurables.Configurable) {
	this.Helper()
	for _, listElement := range list {
		if listElement.CharacteristicId == element.CharacteristicId && len(listElement.Values) == len(element.Values) {
			for _, value := range element.Values {
				if this.listContains(listElement.Values, value) {
					return
				}
			}
		}
	}
	this.Fatal("missing element in list", list, element)
}
