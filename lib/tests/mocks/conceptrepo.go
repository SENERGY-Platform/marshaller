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

package mocks

import (
	"context"
	"github.com/SENERGY-Platform/converter/lib/converter/characteristics"
	"github.com/SENERGY-Platform/device-repository/lib/api"
	"github.com/SENERGY-Platform/device-repository/lib/client"
	devicerepoconfig "github.com/SENERGY-Platform/device-repository/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/conceptrepo"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/marshaller/lib/tests/testdata"
	"net/http/httptest"
)

const exampleColor = "example_color"
const exampleRgb = "example_rgb"
const exampleHex = "example_hex"
const exampleBrightness = "example_brightness"
const exampleLux = "example_lux"

func NewMockConceptRepo(ctx context.Context) (*conceptrepo.ConceptRepo, error) {
	functions, err := testdata.GetFunctions()
	if err != nil {
		return nil, err
	}
	concepts, err := testdata.GetConcepts()
	if err != nil {
		return nil, err
	}
	characteristicsList, err := testdata.GetCharacteristics()
	if err != nil {
		return nil, err
	}

	c, db, err := client.NewTestClient()
	if err != nil {
		return nil, err
	}

	for _, function := range functions {
		err = db.SetFunction(ctx, function)
		if err != nil {
			return nil, err
		}
	}
	for _, concept := range concepts {
		err = db.SetConcept(ctx, concept)
		if err != nil {
			return nil, err
		}
	}
	for _, characteristic := range characteristicsList {
		err = db.SetCharacteristic(ctx, characteristic)
		if err != nil {
			return nil, err
		}
	}

	server := httptest.NewServer(api.GetRouterWithoutMiddleware(devicerepoconfig.Config{}, c))

	go func() {
		<-ctx.Done()
		server.Close()
	}()

	config := config.Config{
		DeviceRepositoryUrl:        server.URL,
		ConceptRepoRefreshInterval: 42000,
		LogLevel:                   "DEBUG",
	}
	return conceptrepo.New(ctx, config, MockAccess{},
		conceptrepo.ConceptRepoDefault{
			Concept: model.NullConcept,
			Characteristics: []model.Characteristic{
				model.NullCharacteristic,
			},
		},
		conceptrepo.ConceptRepoDefault{
			Concept: model.Concept{Id: exampleColor, Name: "example", BaseCharacteristicId: exampleRgb},
			Characteristics: []model.Characteristic{
				{
					Id:   exampleRgb,
					Name: "rgb",
					Type: model.Structure,
					SubCharacteristics: []model.Characteristic{
						{Id: exampleRgb + ".r", Name: "r", Type: model.Integer},
						{Id: exampleRgb + ".g", Name: "g", Type: model.Integer},
						{Id: exampleRgb + ".b", Name: "b", Type: model.Integer},
					},
				},
				{
					Id:   exampleHex,
					Name: "hex",
					Type: model.String,
				},
			},
		},
		conceptrepo.ConceptRepoDefault{
			Concept: model.Concept{Id: exampleBrightness, Name: "example-bri"},
			Characteristics: []model.Characteristic{
				{
					Id:   exampleLux,
					Name: "lux",
					Type: model.Integer,
				},
			},
		},
		conceptrepo.ConceptRepoDefault{
			Concept: model.Concept{Id: "side-celsius", Name: "side-celsius"},
			Characteristics: []model.Characteristic{
				{
					Id:   characteristics.Celsius,
					Name: "celsius",
					Type: model.Integer,
				},
				{
					Id:   "side-celsius-foo",
					Name: "side-celsius-foo",
					Type: model.Integer,
				},
			},
		},

		conceptrepo.ConceptRepoDefault{
			Concept: model.Concept{Id: "side-kelvin", Name: "side-kelvin"},
			Characteristics: []model.Characteristic{
				{
					Id:   characteristics.Kelvin,
					Name: "kelvin",
					Type: model.Integer,
				},
				{
					Id:   "side-kelvin-foo",
					Name: "side-kelvin-foo",
					Type: model.Integer,
				},
			},
		},
	)
}

type MockAccess struct{}

func (m MockAccess) Ensure() (config.Impersonate, error) {
	return "", nil
}
