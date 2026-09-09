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
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/SENERGY-Platform/marshaller/lib/api"
	"github.com/SENERGY-Platform/marshaller/lib/api/metrics"
	"github.com/SENERGY-Platform/marshaller/lib/client"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	"github.com/SENERGY-Platform/marshaller/lib/controller"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	v2 "github.com/SENERGY-Platform/marshaller/lib/marshaller/v2"
	"github.com/SENERGY-Platform/marshaller/lib/tests/mocks"
)

var ServerUrl string

// TestClient drives the service the way a consumer does. The tests use it instead of
// hand-built requests, so a client method that calls the wrong endpoint fails here too and
// not only in the client package's own test.
var TestClient client.Interface

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
	Celsius string
}{
	Celsius: "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a",
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
	setupMock(ctx, cancelFinished)
	code := m.Run()
	return code
}

func setupMock(ctx context.Context, done *sync.WaitGroup) {
	conceptRepo, err := mocks.NewMockConceptRepo(ctx)
	if err != nil {
		panic(err)
	}
	marshaller := marshaller.New(mocks.Converter{}, conceptRepo, mocks.DeviceRepo)
	marshallerv2 := v2.New(config.Config{}, mocks.Converter{}, conceptRepo)
	configurableService := configurables.New(conceptRepo)
	TestMarshalInputs = marshaller.MarshalInputs
	TestUnmarshalOutputs = marshaller.UnmarshalOutputs
	TestFindConfigurables = configurableService.Find
	done.Add(1)
	ctrl := controller.New(config.Config{Debug: true}, marshaller, marshallerv2, configurableService, mocks.DeviceRepo, nil)
	server := httptest.NewServer(api.GetRouter(config.Config{Debug: true}, ctrl, metrics.NewMetrics(config.Config{})))
	ServerUrl = server.URL
	TestClient = client.NewClient(server.URL, nil)
	go func() {
		<-ctx.Done()
		server.Close()
		done.Done()
	}()
}

var TestMarshalInputs = func(protocol model.Protocol, service model.Service, input interface{}, inputCharacteristicId string, pathAllowList []string, configurables ...configurables.Configurable) (result map[string]string, err error) {
	return nil, errors.New("todo")
}

var TestUnmarshalOutputs = func(protocol model.Protocol, service model.Service, outputMap map[string]string, outputCharacteristicId string, pathAllowList []string, hints ...string) (result interface{}, err error) {
	return nil, errors.New("todo")
}

var TestFindConfigurables = func(notCharacteristicId string, services []model.Service) (result configurables.Configurables, err error) {
	return nil, errors.New("todo")
}

func requestFromJson[T any](body string) T {
	result := new(T)
	if err := json.Unmarshal([]byte(body), result); err != nil {
		panic(err)
	}
	return *result
}

// printJson renders a client answer the way the raw response body was printed before, so
// the expected output of an example stays what it was and does not depend on how Go
// happens to print a map.
func printJson(result interface{}, err error, code int) {
	if err != nil {
		fmt.Println(code, err)
		return
	}
	body, err := json.Marshal(result)
	fmt.Println(err, string(body))
}

// printJsonWithCode additionally prints the status, for the examples that documented it.
func printJsonWithCode(result interface{}, err error, code int) {
	if err != nil {
		fmt.Println(code, err)
		return
	}
	body, err := json.Marshal(result)
	fmt.Println(err, code, string(body))
}
