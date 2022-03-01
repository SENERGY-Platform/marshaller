/*
 * Copyright 2021 InfAI (CC SES)
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
	"bytes"
	"encoding/json"
	"github.com/SENERGY-Platform/marshaller/lib/api"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/marshaller/lib/tests/mocks"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestPathOptionsTemperatureCelsiusArray(t *testing.T) {
	aspectId := "urn:infai:ses:aspect:a14c5efb-b0b6-46c3-982e-9fded75b5ab6"
	functionId := "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b"
	mocks.DeviceRepo.SetDeviceType(model.DeviceType{
		Id:   "TestPathOptionsTemperatureCelsiusArray",
		Name: "TestPathOptionsTemperatureCelsiusArray",
		Services: []model.Service{
			{
				Id:      "TestPathOptionsTemperatureCelsiusArray.celsius",
				LocalId: "TestPathOptionsTemperatureCelsiusArray.celsius",
				Name:    "TestPathOptionsTemperatureCelsiusArray.celsius",
				Outputs: []model.Content{
					{
						ContentVariable: model.ContentVariable{
							Name: "temperature",
							Type: model.List,
							SubContentVariables: []model.ContentVariable{
								{
									Name:             "0",
									Type:             model.Float,
									CharacteristicId: "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a",
									AspectId:         aspectId,
									FunctionId:       functionId,
								},
								{
									Name: "1",
									Type: model.String,
								},
							},
						},
					},
				},
			},
		},
	})

	t.Run("filtered to celsius with envelope", testPathOptions(
		[]string{"TestPathOptionsTemperatureCelsiusArray"},
		functionId,
		aspectId,
		[]string{"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
		false,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsTemperatureCelsiusArray": {
				{
					ServiceId:              "TestPathOptionsTemperatureCelsiusArray.celsius",
					JsonPath:               []string{"value.temperature[0]"},
					PathToCharacteristicId: map[string]string{"value.temperature[0]": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
		},
	))
	t.Run("filtered to celsius without envelope", testPathOptions(
		[]string{"TestPathOptionsTemperatureCelsiusArray"},
		functionId,
		aspectId,
		[]string{"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
		true,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsTemperatureCelsiusArray": {
				{
					ServiceId:              "TestPathOptionsTemperatureCelsiusArray.celsius",
					JsonPath:               []string{"temperature[0]"},
					PathToCharacteristicId: map[string]string{"temperature[0]": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
		},
	))

	t.Run("filtered to celsius and kelvin with envelope", testPathOptions(
		[]string{"TestPathOptionsTemperatureCelsiusArray"},
		functionId,
		aspectId,
		[]string{
			"urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683", //kelvin
			"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", //celsius
		},
		false,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsTemperatureCelsiusArray": {
				{
					ServiceId:              "TestPathOptionsTemperatureCelsiusArray.celsius",
					JsonPath:               []string{"value.temperature[0]"},
					PathToCharacteristicId: map[string]string{"value.temperature[0]": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
		},
	))
	t.Run("filtered to celsius and kelvin without envelope", testPathOptions(
		[]string{"TestPathOptionsTemperatureCelsiusArray"},
		functionId,
		aspectId,
		[]string{
			"urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683", //kelvin
			"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", //celsius
		},
		true,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsTemperatureCelsiusArray": {
				{
					ServiceId:              "TestPathOptionsTemperatureCelsiusArray.celsius",
					JsonPath:               []string{"temperature[0]"},
					PathToCharacteristicId: map[string]string{"temperature[0]": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
		},
	))

	t.Run("unfiltered with envelope", testPathOptions(
		[]string{"TestPathOptionsTemperatureCelsiusArray"},
		functionId,
		aspectId,
		nil,
		false,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsTemperatureCelsiusArray": {
				{
					ServiceId:              "TestPathOptionsTemperatureCelsiusArray.celsius",
					JsonPath:               []string{"value.temperature[0]"},
					PathToCharacteristicId: map[string]string{"value.temperature[0]": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
		},
	))
	t.Run("unfiltered with envelope", testPathOptions(
		[]string{"TestPathOptionsTemperatureCelsiusArray"},
		functionId,
		aspectId,
		nil,
		true,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsTemperatureCelsiusArray": {
				{
					ServiceId:              "TestPathOptionsTemperatureCelsiusArray.celsius",
					JsonPath:               []string{"temperature[0]"},
					PathToCharacteristicId: map[string]string{"temperature[0]": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
		},
	))
}

func TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin(t *testing.T) {
	aspectId := "urn:infai:ses:aspect:a14c5efb-b0b6-46c3-982e-9fded75b5ab6"
	functionId := "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b"
	mocks.DeviceRepo.SetDeviceType(model.DeviceType{
		Id:   "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c",
		Name: "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c",
		Services: []model.Service{
			{
				Id:      "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c.celsius",
				LocalId: "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c.celsius",
				Name:    "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c.celsius",
				Outputs: []model.Content{
					{
						ContentVariable: model.ContentVariable{
							Name: "temperature",
							Type: model.Structure,
							SubContentVariables: []model.ContentVariable{
								{
									Name:             "celsius",
									Type:             model.Float,
									CharacteristicId: "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a",
									AspectId:         aspectId,
									FunctionId:       functionId,
								},
								{
									Name: "unit",
									Type: model.String,
								},
							},
						},
					},
				},
			},
		},
	})

	mocks.DeviceRepo.SetDeviceType(model.DeviceType{
		Id:   "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k",
		Name: "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k",
		Services: []model.Service{
			{
				Id:      "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k.kelvin",
				LocalId: "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k.kelvin",
				Name:    "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k.kelvin",
				Outputs: []model.Content{
					{
						ContentVariable: model.ContentVariable{
							Name: "temperature",
							Type: model.Structure,
							SubContentVariables: []model.ContentVariable{
								{
									Name:             "kelvin",
									Type:             model.Float,
									CharacteristicId: "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683",
									AspectId:         aspectId,
									FunctionId:       functionId,
								},
								{
									Name: "unit",
									Type: model.String,
								},
							},
						},
					},
				},
			},
		},
	})

	t.Run("filtered to celsius with envelope", testPathOptions(
		[]string{"TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c", "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k"},
		functionId,
		aspectId,
		[]string{"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
		false,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c": {
				{
					ServiceId:              "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c.celsius",
					JsonPath:               []string{"value.temperature.celsius"},
					PathToCharacteristicId: map[string]string{"value.temperature.celsius": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
			"TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k": {},
		},
	))
	t.Run("filtered to celsius without envelope", testPathOptions(
		[]string{"TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c", "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k"},
		functionId,
		aspectId,
		[]string{"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
		true,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c": {
				{
					ServiceId:              "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c.celsius",
					JsonPath:               []string{"temperature.celsius"},
					PathToCharacteristicId: map[string]string{"temperature.celsius": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
			"TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k": {},
		},
	))

	t.Run("filtered to celsius and kelvin with envelope", testPathOptions(
		[]string{"TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c", "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k"},
		functionId,
		aspectId,
		[]string{
			"urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683", //kelvin
			"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", //celsius
		},
		false,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c": {
				{
					ServiceId:              "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c.celsius",
					JsonPath:               []string{"value.temperature.celsius"},
					PathToCharacteristicId: map[string]string{"value.temperature.celsius": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
			"TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k": {
				{
					ServiceId:              "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k.kelvin",
					JsonPath:               []string{"value.temperature.kelvin"},
					PathToCharacteristicId: map[string]string{"value.temperature.kelvin": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683"},
				},
			},
		},
	))
	t.Run("filtered to celsius and kelvin without envelope", testPathOptions(
		[]string{"TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c", "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k"},
		functionId,
		aspectId,
		[]string{
			"urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683", //kelvin
			"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", //celsius
		},
		true,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c": {
				{
					ServiceId:              "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c.celsius",
					JsonPath:               []string{"temperature.celsius"},
					PathToCharacteristicId: map[string]string{"temperature.celsius": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
			"TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k": {
				{
					ServiceId:              "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k.kelvin",
					JsonPath:               []string{"temperature.kelvin"},
					PathToCharacteristicId: map[string]string{"temperature.kelvin": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683"},
				},
			},
		},
	))

	t.Run("unfiltered with envelope", testPathOptions(
		[]string{"TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c", "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k"},
		functionId,
		aspectId,
		nil,
		false,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c": {
				{
					ServiceId:              "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c.celsius",
					JsonPath:               []string{"value.temperature.celsius"},
					PathToCharacteristicId: map[string]string{"value.temperature.celsius": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
			"TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k": {
				{
					ServiceId:              "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k.kelvin",
					JsonPath:               []string{"value.temperature.kelvin"},
					PathToCharacteristicId: map[string]string{"value.temperature.kelvin": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683"},
				},
			},
		},
	))
	t.Run("unfiltered with envelope", testPathOptions(
		[]string{"TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c", "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k"},
		functionId,
		aspectId,
		nil,
		true,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c": {
				{
					ServiceId:              "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_c.celsius",
					JsonPath:               []string{"temperature.celsius"},
					PathToCharacteristicId: map[string]string{"temperature.celsius": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
			"TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k": {
				{
					ServiceId:              "TestPathOptionsMultipleDeviceTypesTemperatureCelsiusAndKelvin_k.kelvin",
					JsonPath:               []string{"temperature.kelvin"},
					PathToCharacteristicId: map[string]string{"temperature.kelvin": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683"},
				},
			},
		},
	))
}

func TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin(t *testing.T) {
	aspectId := "urn:infai:ses:aspect:a14c5efb-b0b6-46c3-982e-9fded75b5ab6"
	functionId := "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b"
	mocks.DeviceRepo.SetDeviceType(model.DeviceType{
		Id:   "TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin",
		Name: "TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin",
		Services: []model.Service{
			{
				Id:      "TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin.celsius",
				LocalId: "TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin.celsius",
				Name:    "TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin.celsius",
				Outputs: []model.Content{
					{
						ContentVariable: model.ContentVariable{
							Name: "temperature",
							Type: model.Structure,
							SubContentVariables: []model.ContentVariable{
								{
									Name:             "celsius",
									Type:             model.Float,
									CharacteristicId: "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a",
									AspectId:         aspectId,
									FunctionId:       functionId,
								},
								{
									Name: "unit",
									Type: model.String,
								},
							},
						},
					},
				},
			},
			{
				Id:      "TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin.kelvin",
				LocalId: "TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin.kelvin",
				Name:    "TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin.kelvin",
				Outputs: []model.Content{
					{
						ContentVariable: model.ContentVariable{
							Name: "temperature",
							Type: model.Structure,
							SubContentVariables: []model.ContentVariable{
								{
									Name:             "kelvin",
									Type:             model.Float,
									CharacteristicId: "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683",
									AspectId:         aspectId,
									FunctionId:       functionId,
								},
								{
									Name: "unit",
									Type: model.String,
								},
							},
						},
					},
				},
			},
		},
	})

	t.Run("filtered to celsius with envelope", testPathOptions(
		[]string{"TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin"},
		functionId,
		aspectId,
		[]string{"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
		false,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin": {
				{
					ServiceId:              "TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin.celsius",
					JsonPath:               []string{"value.temperature.celsius"},
					PathToCharacteristicId: map[string]string{"value.temperature.celsius": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
		},
	))
	t.Run("filtered to celsius without envelope", testPathOptions(
		[]string{"TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin"},
		functionId,
		aspectId,
		[]string{"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
		true,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin": {
				{
					ServiceId:              "TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin.celsius",
					JsonPath:               []string{"temperature.celsius"},
					PathToCharacteristicId: map[string]string{"temperature.celsius": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
		},
	))

	t.Run("filtered to celsius and kelvin with envelope", testPathOptions(
		[]string{"TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin"},
		functionId,
		aspectId,
		[]string{
			"urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683", //kelvin
			"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", //celsius
		},
		false,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin": {
				{
					ServiceId:              "TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin.celsius",
					JsonPath:               []string{"value.temperature.celsius"},
					PathToCharacteristicId: map[string]string{"value.temperature.celsius": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
				{
					ServiceId:              "TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin.kelvin",
					JsonPath:               []string{"value.temperature.kelvin"},
					PathToCharacteristicId: map[string]string{"value.temperature.kelvin": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683"},
				},
			},
		},
	))
	t.Run("filtered to celsius and kelvin without envelope", testPathOptions(
		[]string{"TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin"},
		functionId,
		aspectId,
		[]string{
			"urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683", //kelvin
			"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", //celsius
		},
		true,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin": {
				{
					ServiceId:              "TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin.celsius",
					JsonPath:               []string{"temperature.celsius"},
					PathToCharacteristicId: map[string]string{"temperature.celsius": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
				{
					ServiceId:              "TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin.kelvin",
					JsonPath:               []string{"temperature.kelvin"},
					PathToCharacteristicId: map[string]string{"temperature.kelvin": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683"},
				},
			},
		},
	))

	t.Run("unfiltered with envelope", testPathOptions(
		[]string{"TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin"},
		functionId,
		aspectId,
		nil,
		false,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin": {
				{
					ServiceId:              "TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin.celsius",
					JsonPath:               []string{"value.temperature.celsius"},
					PathToCharacteristicId: map[string]string{"value.temperature.celsius": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
				{
					ServiceId:              "TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin.kelvin",
					JsonPath:               []string{"value.temperature.kelvin"},
					PathToCharacteristicId: map[string]string{"value.temperature.kelvin": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683"},
				},
			},
		},
	))
	t.Run("unfiltered with envelope", testPathOptions(
		[]string{"TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin"},
		functionId,
		aspectId,
		nil,
		true,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin": {
				{
					ServiceId:              "TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin.celsius",
					JsonPath:               []string{"temperature.celsius"},
					PathToCharacteristicId: map[string]string{"temperature.celsius": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
				{
					ServiceId:              "TestPathOptionsMultipleServicesTemperatureCelsiusAndKelvin.kelvin",
					JsonPath:               []string{"temperature.kelvin"},
					PathToCharacteristicId: map[string]string{"temperature.kelvin": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683"},
				},
			},
		},
	))
}

func TestPathOptionsOneServiceTemperatureCelsiusAndKelvin(t *testing.T) {
	aspectId := "urn:infai:ses:aspect:a14c5efb-b0b6-46c3-982e-9fded75b5ab6"
	functionId := "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b"
	mocks.DeviceRepo.SetDeviceType(model.DeviceType{
		Id:   "TestPathOptionsOneServiceTemperatureCelsiusAndKelvin",
		Name: "TestPathOptionsOneServiceTemperatureCelsiusAndKelvin",
		Services: []model.Service{
			{
				Id:      "TestPathOptionsOneServiceTemperatureCelsiusAndKelvin.1",
				LocalId: "TestPathOptionsOneServiceTemperatureCelsiusAndKelvin.1",
				Name:    "TestPathOptionsOneServiceTemperatureCelsiusAndKelvin.1",
				Outputs: []model.Content{
					{
						ContentVariable: model.ContentVariable{
							Name: "temperature",
							Type: model.Structure,
							//modified from legacy test but multiple assignments of same aspect is not valid
							SubContentVariables: []model.ContentVariable{
								{
									Name:             "celsius",
									Type:             model.Float,
									CharacteristicId: "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a",
									AspectId:         aspectId,
									FunctionId:       functionId,
								},
								{
									Name:             "kelvin",
									Type:             model.Float,
									CharacteristicId: "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683",
									AspectId:         aspectId,
									FunctionId:       functionId,
								},
							},
						},
					},
				},
			},
		},
	})

	t.Run("filtered to celsius with envelope", testPathOptions(
		[]string{"TestPathOptionsOneServiceTemperatureCelsiusAndKelvin"},
		functionId,
		aspectId,
		[]string{"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
		false,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsOneServiceTemperatureCelsiusAndKelvin": {
				{
					ServiceId:              "TestPathOptionsOneServiceTemperatureCelsiusAndKelvin.1",
					JsonPath:               []string{"value.temperature.celsius"},
					PathToCharacteristicId: map[string]string{"value.temperature.celsius": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
		},
	))
	t.Run("filtered to celsius without envelope", testPathOptions(
		[]string{"TestPathOptionsOneServiceTemperatureCelsiusAndKelvin"},
		functionId,
		aspectId,
		[]string{"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
		true,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsOneServiceTemperatureCelsiusAndKelvin": {
				{
					ServiceId:              "TestPathOptionsOneServiceTemperatureCelsiusAndKelvin.1",
					JsonPath:               []string{"temperature.celsius"},
					PathToCharacteristicId: map[string]string{"temperature.celsius": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
		},
	))

	t.Run("filtered to celsius and kelvin with envelope", testPathOptions(
		[]string{"TestPathOptionsOneServiceTemperatureCelsiusAndKelvin"},
		functionId,
		aspectId,
		[]string{
			"urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683", //kelvin
			"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", //celsius
		},
		false,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsOneServiceTemperatureCelsiusAndKelvin": {
				{
					ServiceId:              "TestPathOptionsOneServiceTemperatureCelsiusAndKelvin.1",
					JsonPath:               []string{"value.temperature.celsius", "value.temperature.kelvin"},
					PathToCharacteristicId: map[string]string{"value.temperature.celsius": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", "value.temperature.kelvin": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683"},
				},
			},
		},
	))
	t.Run("filtered to celsius and kelvin without envelope", testPathOptions(
		[]string{"TestPathOptionsOneServiceTemperatureCelsiusAndKelvin"},
		functionId,
		aspectId,
		[]string{
			"urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683", //kelvin
			"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", //celsius
		},
		true,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsOneServiceTemperatureCelsiusAndKelvin": {
				{
					ServiceId:              "TestPathOptionsOneServiceTemperatureCelsiusAndKelvin.1",
					JsonPath:               []string{"temperature.celsius", "temperature.kelvin"},
					PathToCharacteristicId: map[string]string{"temperature.celsius": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", "temperature.kelvin": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683"},
				},
			},
		},
	))

	t.Run("unfiltered with envelope", testPathOptions(
		[]string{"TestPathOptionsOneServiceTemperatureCelsiusAndKelvin"},
		functionId,
		aspectId,
		nil,
		false,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsOneServiceTemperatureCelsiusAndKelvin": {
				{
					ServiceId:              "TestPathOptionsOneServiceTemperatureCelsiusAndKelvin.1",
					JsonPath:               []string{"value.temperature.celsius", "value.temperature.kelvin"},
					PathToCharacteristicId: map[string]string{"value.temperature.celsius": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", "value.temperature.kelvin": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683"},
				},
			},
		},
	))
	t.Run("unfiltered with envelope", testPathOptions(
		[]string{"TestPathOptionsOneServiceTemperatureCelsiusAndKelvin"},
		functionId,
		aspectId,
		nil,
		true,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsOneServiceTemperatureCelsiusAndKelvin": {
				{
					ServiceId:              "TestPathOptionsOneServiceTemperatureCelsiusAndKelvin.1",
					JsonPath:               []string{"temperature.celsius", "temperature.kelvin"},
					PathToCharacteristicId: map[string]string{"temperature.celsius": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", "temperature.kelvin": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683"},
				},
			},
		},
	))
}

func TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside(t *testing.T) {
	aspectId := "urn:infai:ses:aspect:a14c5efb-b0b6-46c3-982e-9fded75b5ab6"
	functionId := "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b"
	mocks.DeviceRepo.SetDeviceType(model.DeviceType{
		Id:   "TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside",
		Name: "TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside",
		Services: []model.Service{
			{
				Id:      "TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside.1",
				LocalId: "TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside.1",
				Name:    "TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside.1",
				Outputs: []model.Content{
					{
						ContentVariable: model.ContentVariable{
							Name: "temperature",
							Type: model.Structure,
							//modified from legacy test but multiple assignments of same aspect is not valid
							SubContentVariables: []model.ContentVariable{
								{
									Name: "inside",
									Type: model.Structure,
									SubContentVariables: []model.ContentVariable{
										{
											Name:             "value",
											Type:             model.Float,
											CharacteristicId: "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a",
											AspectId:         aspectId,
											FunctionId:       functionId,
										},
										{
											Name: "unit",
											Type: model.String,
										},
									},
								},
								{
									Name: "outside",
									Type: model.Structure,
									SubContentVariables: []model.ContentVariable{
										{
											Name:             "value",
											Type:             model.Float,
											CharacteristicId: "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a",
											AspectId:         aspectId,
											FunctionId:       functionId,
										},
										{
											Name: "unit",
											Type: model.String,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	t.Run("filtered to celsius with envelope", testPathOptions(
		[]string{"TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside"},
		functionId,
		aspectId,
		[]string{"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
		false,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside": {
				{
					ServiceId:              "TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside.1",
					JsonPath:               []string{"value.temperature.inside.value", "value.temperature.outside.value"},
					PathToCharacteristicId: map[string]string{"value.temperature.inside.value": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", "value.temperature.outside.value": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
		},
	))
	t.Run("filtered to celsius without envelope", testPathOptions(
		[]string{"TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside"},
		functionId,
		aspectId,
		[]string{"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
		true,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside": {
				{
					ServiceId:              "TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside.1",
					JsonPath:               []string{"temperature.inside.value", "temperature.outside.value"},
					PathToCharacteristicId: map[string]string{"temperature.inside.value": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", "temperature.outside.value": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
		},
	))

	t.Run("filtered to celsius and kelvin with envelope", testPathOptions(
		[]string{"TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside"},
		functionId,
		aspectId,
		[]string{
			"urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683", //kelvin
			"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", //celsius
		},
		false,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside": {
				{
					ServiceId:              "TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside.1",
					JsonPath:               []string{"value.temperature.inside.value", "value.temperature.outside.value"},
					PathToCharacteristicId: map[string]string{"value.temperature.inside.value": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", "value.temperature.outside.value": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
		},
	))
	t.Run("filtered to celsius and kelvin without envelope", testPathOptions(
		[]string{"TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside"},
		functionId,
		aspectId,
		[]string{
			"urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683", //kelvin
			"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", //celsius
		},
		true,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside": {
				{
					ServiceId:              "TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside.1",
					JsonPath:               []string{"temperature.inside.value", "temperature.outside.value"},
					PathToCharacteristicId: map[string]string{"temperature.inside.value": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", "temperature.outside.value": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
		},
	))

	t.Run("unfiltered with envelope", testPathOptions(
		[]string{"TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside"},
		functionId,
		aspectId,
		nil,
		false,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside": {
				{
					ServiceId:              "TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside.1",
					JsonPath:               []string{"value.temperature.inside.value", "value.temperature.outside.value"},
					PathToCharacteristicId: map[string]string{"value.temperature.inside.value": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", "value.temperature.outside.value": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
		},
	))
	t.Run("unfiltered with envelope", testPathOptions(
		[]string{"TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside"},
		functionId,
		aspectId,
		nil,
		true,
		map[string][]marshaller.PathOptionsResultElement{
			"TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside": {
				{
					ServiceId:              "TestPathOptionsOneServiceTemperatureCelsiusInsideAndOutside.1",
					JsonPath:               []string{"temperature.inside.value", "temperature.outside.value"},
					PathToCharacteristicId: map[string]string{"temperature.inside.value": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a", "temperature.outside.value": "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"},
				},
			},
		},
	))
}

func testPathOptions(deviceTypes []string, functionId string, aspect string, characteristicsFilter []string, withoutEnvelope bool, expectedResult map[string][]marshaller.PathOptionsResultElement) func(t *testing.T) {
	return func(t *testing.T) {
		buff := bytes.Buffer{}
		err := json.NewEncoder(&buff).Encode(api.PathOptionsQuery{
			DeviceTypeIds:          deviceTypes,
			FunctionId:             functionId,
			AspectId:               aspect,
			CharacteristicIdFilter: characteristicsFilter,
			WithoutEnvelope:        withoutEnvelope,
		})
		if err != nil {
			t.Error(err.Error())
			return
		}
		req, err := http.NewRequest("POST", ServerUrl+"/query/path-options", &buff)
		if err != nil {
			t.Error(err.Error())
			return
		}
		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Error(err.Error())
			return
		}

		if resp.StatusCode != 200 {
			temp, _ := ioutil.ReadAll(resp.Body)
			t.Error(resp.StatusCode, string(temp))
			return
		}

		result := map[string][]marshaller.PathOptionsResultElement{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if !reflect.DeepEqual(expectedResult, result) {
			resultJson, _ := json.Marshal(result)
			expectedJson, _ := json.Marshal(expectedResult)
			t.Error(string(resultJson), "\n", string(expectedJson))
			return
		}
	}
}
