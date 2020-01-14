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
	"context"
	"errors"
	"flag"
	"github.com/SENERGY-Platform/marshaller/lib/api"
	"github.com/SENERGY-Platform/marshaller/lib/conceptrepo"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	"github.com/SENERGY-Platform/marshaller/lib/converter"
	"github.com/SENERGY-Platform/marshaller/lib/devicerepository"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/marshaller/lib/tests/mocks"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

var ServerUrl string

var example = struct {
	Brightness string
	Lux        string
	Color      string
	Rgb        string
	Hex        string
}{
	Brightness: "example_brightness",
	Lux:        "example_lux",
	Color:      "example_color",
	Rgb:        "example_rgb",
	Hex:        "example_hex",
}

var temperature = struct {
	Celcius string
}{
	Celcius: "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a",
}

var color = struct {
	Rgb  string
	RgbR string
	RgbG string
	RgbB string
	Hex  string
}{
	Rgb:  "urn:infai:ses:characteristic:5b4eea52-e8e5-4e80-9455-0382f81a1b43",
	RgbR: "urn:infai:ses:characteristic:dfe6be4a-650c-4411-8d87-062916b48951",
	RgbG: "urn:infai:ses:characteristic:5ef27837-4aca-43ad-b8f6-4d95cf9ed99e",
	RgbB: "urn:infai:ses:characteristic:590af9ef-3a5e-4edb-abab-177cb1320b17",
	Hex:  "urn:infai:ses:characteristic:0fc343ce-4627-4c88-b1e0-d3ed29754af8",
}

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	flag.Parse()
	cancelFinished := &sync.WaitGroup{}
	defer func() {
		cancelFinished.Wait()
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if testing.Short() {
		setupMock(ctx, cancelFinished)
	} else {
		setupExternal(ctx, cancelFinished)
	}
	code := m.Run()
	return code
}

func setupMock(ctx context.Context, done *sync.WaitGroup) {
	marshaller := marshaller.New(mocks.Converter{}, mocks.ConceptRepo{})
	configurableService := configurables.New(mocks.ConceptRepo{})
	TestMarshalInputs = marshaller.MarshalInputs
	TestUnmarshalOutputs = marshaller.UnmarshalOutputs
	TestFindConfigurables = configurableService.Find
}

func setupExternal(ctx context.Context, done *sync.WaitGroup) {
	conf, err := config.Load("testdata/config.json")
	if err != nil {
		panic(err)
	}
	access := config.NewAccess(conf)
	conceptRepo, err := conceptrepo.New(
		ctx,
		conf,
		access,
		conceptrepo.ConceptRepoDefault{
			Concept: model.NullConcept,
			Characteristics: []model.Characteristic{
				model.NullCharacteristic,
			},
		},
		conceptrepo.ConceptRepoDefault{
			Concept: model.Concept{Id: example.Color, Name: "example", BaseCharacteristicId: example.Rgb},
			Characteristics: []model.Characteristic{
				{
					Id:   example.Rgb,
					Name: "rgb",
					Type: model.Structure,
					SubCharacteristics: []model.Characteristic{
						{Id: example.Rgb + ".r", Name: "r", Type: model.Integer},
						{Id: example.Rgb + ".g", Name: "g", Type: model.Integer},
						{Id: example.Rgb + ".b", Name: "b", Type: model.Integer},
					},
				},
				{
					Id:   example.Hex,
					Name: "hex",
					Type: model.String,
				},
			},
		},
		conceptrepo.ConceptRepoDefault{
			Concept: model.Concept{Id: example.Brightness, Name: "example-bri"},
			Characteristics: []model.Characteristic{
				{
					Id:   example.Lux,
					Name: "lux",
					Type: model.Integer,
				},
			},
		},
	)
	if err != nil {
		panic(err)
	}
	marshaller := marshaller.New(converter.New(conf, access), conceptRepo)
	configurableService := configurables.New(conceptRepo)
	devicerepo := devicerepository.New(conf)

	TestMarshalInputs = marshaller.MarshalInputs
	TestUnmarshalOutputs = marshaller.UnmarshalOutputs
	TestFindConfigurables = configurableService.Find

	done.Add(1)
	server := httptest.NewServer(api.GetRouter(conf, marshaller, configurableService, devicerepo))
	ServerUrl = server.URL
	go func() {
		<-ctx.Done()
		server.Close()
		done.Done()
	}()
}

var TestMarshalInputs = func(protocol model.Protocol, service model.Service, input interface{}, inputCharacteristicId string, configurables ...configurables.Configurable) (result map[string]string, err error) {
	return nil, errors.New("todo")
}

var TestUnmarshalOutputs = func(protocol model.Protocol, service model.Service, outputMap map[string]string, outputCharacteristicId string) (result interface{}, err error) {
	return nil, errors.New("todo")
}

var TestFindConfigurables = func(notCharacteristicId string, services []model.Service) (result configurables.Configurables, err error) {
	return nil, errors.New("todo")
}
