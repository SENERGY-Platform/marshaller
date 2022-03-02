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
	"encoding/json"
	"github.com/SENERGY-Platform/converter/lib/converter/characteristics"
	"github.com/SENERGY-Platform/marshaller/lib/conceptrepo"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/marshaller/lib/tests/testdata"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
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

	searchMockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		//endpoint := this.config.PermissionsSearchUrl + "/v3/resources/concepts?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset) + "&sort=name.asc&rights=r"
		//endpoint := this.config.PermissionsSearchUrl + "/v3/resources/functions?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset) + "&sort=name.asc&rights=r"

		log.Println("TEST-DEBUG: searchMockServer:", request.Method, request.URL.String(), request.URL.Path)

		limitStr := request.URL.Query().Get("limit")
		limit := 100
		if limitStr != "" {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}

		offsetStr := request.URL.Query().Get("offset")
		offset := 0
		if offsetStr != "" {
			offset, err = strconv.Atoi(offsetStr)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}
		if strings.Contains(request.URL.String(), "/v3/resources/functions") {
			end := limit + offset
			if end > len(functions) {
				end = len(functions)
			}
			json.NewEncoder(writer).Encode(functions[offset:end])
			return
		}

		if strings.Contains(request.URL.String(), "/v3/resources/concepts") {
			end := limit + offset
			if end > len(concepts) {
				end = len(concepts)
			}
			json.NewEncoder(writer).Encode(concepts[offset:end])
			return
		}

		http.Error(writer, request.URL.String(), http.StatusNotFound)
	}))

	devicerepoMockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		//this.config.DeviceRepositoryUrl+"/concepts/"+url.PathEscape(id)
		//this.config.DeviceRepositoryUrl+"/characteristics/"+url.PathEscape(id)

		log.Println("TEST-DEBUG: devicerepoMockServer:", request.Method, request.URL.String(), request.URL.Path)

		parts := strings.Split(request.URL.Path, "/")
		id := parts[len(parts)-1]

		if strings.Contains(request.URL.String(), "/concepts/") {
			for _, element := range concepts {
				if element.Id == id {
					json.NewEncoder(writer).Encode(element)
					return
				}
			}
		}

		if strings.Contains(request.URL.String(), "/characteristics/") {
			for _, element := range characteristicsList {
				if element.Id == id {
					json.NewEncoder(writer).Encode(element)
					return
				}
			}
		}

		log.Println("TEST_ERROR: no match found", request.URL.Path, request.URL.String())
		http.Error(writer, request.URL.Path, http.StatusNotFound)
	}))

	go func() {
		<-ctx.Done()
		searchMockServer.Close()
		devicerepoMockServer.Close()
	}()

	config := config.Config{
		PermissionsSearchUrl:       searchMockServer.URL,
		DeviceRepositoryUrl:        devicerepoMockServer.URL,
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
