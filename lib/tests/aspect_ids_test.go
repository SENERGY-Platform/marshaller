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

package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/marshaller/lib/tests/mocks"
)

// TestPathOptionsAspectIds checks the aspect list of a path-options query against a
// device-type with two services that differ only in the aspects their output carries. A
// query naming more than one aspect is an AND, so only the service whose content variable
// carries all of them is offered.
func TestPathOptionsAspectIds(t *testing.T) {
	functionId := "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b"
	deviceTypeId := "TestPathOptionsAspectIds"
	insideService := deviceTypeId + ".inside"
	insideTodayService := deviceTypeId + ".inside_today"

	temperatureOutput := func(aspectIds []string) []model.Content {
		return []model.Content{
			{
				ContentVariable: model.ContentVariable{
					Name: "temperature",
					Type: model.Structure,
					SubContentVariables: []model.ContentVariable{
						{
							Name:             "value",
							Type:             model.Float,
							CharacteristicId: temperature.Celsius,
							FunctionId:       functionId,
							AspectIds:        aspectIds,
						},
					},
				},
			},
		}
	}

	mocks.DeviceRepo.SetDeviceType(model.DeviceType{
		Id:   deviceTypeId,
		Name: deviceTypeId,
		Services: []model.Service{
			{
				Id:      insideService,
				LocalId: insideService,
				Name:    insideService,
				Outputs: temperatureOutput([]string{"inside_air"}),
			},
			{
				Id:      insideTodayService,
				LocalId: insideTodayService,
				Name:    insideTodayService,
				Outputs: temperatureOutput([]string{"inside_air", "today"}),
			},
		},
	})

	pathOption := func(serviceId string) marshaller.PathOptionsResultElement {
		return marshaller.PathOptionsResultElement{
			ServiceId:              serviceId,
			JsonPath:               []string{"value.temperature.value"},
			PathToCharacteristicId: map[string]string{"value.temperature.value": temperature.Celsius},
		}
	}
	onlyInsideToday := map[string][]marshaller.PathOptionsResultElement{
		deviceTypeId: {pathOption(insideTodayService)},
	}
	bothServices := map[string][]marshaller.PathOptionsResultElement{
		deviceTypeId: {pathOption(insideService), pathOption(insideTodayService)},
	}

	query := func(aspectId string, aspectIds []string) messages.PathOptionsQuery {
		return messages.PathOptionsQuery{
			DeviceTypeIds:          []string{deviceTypeId},
			FunctionId:             functionId,
			AspectId:               aspectId,
			AspectIds:              aspectIds,
			CharacteristicIdFilter: []string{temperature.Celsius},
		}
	}

	t.Run("offers only the service carrying every queried aspect", testPathOptionsQuery(
		query("", []string{"inside_air", "today"}), onlyInsideToday))

	t.Run("covers the subtree of every queried aspect", testPathOptionsQuery(
		query("", []string{"air", "electricity"}), onlyInsideToday))

	t.Run("offers every service carrying the single queried aspect", testPathOptionsQuery(
		query("", []string{"inside_air"}), bothServices))

	t.Run("reads the deprecated aspect id as a list with one element", testPathOptionsQuery(
		query("today", nil), onlyInsideToday))

	t.Run("reads the deprecated aspect id and the aspect list as one query", testPathOptionsQuery(
		query("today", []string{"inside_air"}), onlyInsideToday))

	t.Run("offers no service if none carries every queried aspect", testPathOptionsQuery(
		query("", []string{"inside_air", "outside_air"}),
		map[string][]marshaller.PathOptionsResultElement{deviceTypeId: {}}))

	t.Run("offers only the service carrying every aspect of the aspect-ids parameter", testPathOptionsGet(
		deviceTypeId, functionId, "", []string{"inside_air", "today"}, onlyInsideToday))

	t.Run("reads the deprecated aspect-id parameter as a list with one element", testPathOptionsGet(
		deviceTypeId, functionId, "today", nil, onlyInsideToday))
}

// TestPathOptionsDeprecatedContentVariableAspectId checks that a device-type delivered by a
// device-repository that predates the aspect lists is still matched: the deprecated
// AspectId of a content variable is read as a list with one element.
func TestPathOptionsDeprecatedContentVariableAspectId(t *testing.T) {
	functionId := "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b"
	deviceTypeId := "TestPathOptionsDeprecatedContentVariableAspectId"
	serviceId := deviceTypeId + ".inside"

	mocks.DeviceRepo.SetDeviceType(model.DeviceType{
		Id:   deviceTypeId,
		Name: deviceTypeId,
		Services: []model.Service{
			{
				Id:      serviceId,
				LocalId: serviceId,
				Name:    serviceId,
				Outputs: []model.Content{
					{
						ContentVariable: model.ContentVariable{
							Name:             "temperature",
							Type:             model.Float,
							CharacteristicId: temperature.Celsius,
							FunctionId:       functionId,
							AspectId:         "inside_air",
						},
					},
				},
			},
		},
	})

	expected := map[string][]marshaller.PathOptionsResultElement{
		deviceTypeId: {
			{
				ServiceId:              serviceId,
				JsonPath:               []string{"value.temperature"},
				PathToCharacteristicId: map[string]string{"value.temperature": temperature.Celsius},
			},
		},
	}

	t.Run("matches the deprecated aspect id of a content variable by aspect list", testPathOptionsQuery(
		messages.PathOptionsQuery{
			DeviceTypeIds:          []string{deviceTypeId},
			FunctionId:             functionId,
			AspectIds:              []string{"inside_air"},
			CharacteristicIdFilter: []string{temperature.Celsius},
		}, expected))

	t.Run("matches the deprecated aspect id of a content variable by parent aspect", testPathOptionsQuery(
		messages.PathOptionsQuery{
			DeviceTypeIds:          []string{deviceTypeId},
			FunctionId:             functionId,
			AspectIds:              []string{"air"},
			CharacteristicIdFilter: []string{temperature.Celsius},
		}, expected))
}

func testPathOptionsQuery(query messages.PathOptionsQuery, expectedResult map[string][]marshaller.PathOptionsResultElement) func(t *testing.T) {
	return func(t *testing.T) {
		result, err, code := TestClient.GetPathOptions(query)
		if err != nil {
			t.Error(code, err)
			return
		}
		comparePathOptions(t, result, expectedResult)
	}
}

// testPathOptionsGet drives the query-parameter form. The client has no method for it on
// purpose — it uses the request-body form, whose semantics do not depend on which
// parameters are present — so these subtests stay on a hand-built request.
func testPathOptionsGet(deviceTypeId string, functionId string, aspectId string, aspectIds []string, expectedResult map[string][]marshaller.PathOptionsResultElement) func(t *testing.T) {
	return func(t *testing.T) {
		query := url.Values{}
		query.Set("device-type-ids", deviceTypeId)
		query.Set("characteristic-filter", temperature.Celsius)
		query.Set("function-id", functionId)
		if aspectId != "" {
			query.Set("aspect-id", aspectId)
		}
		if len(aspectIds) > 0 {
			query.Set("aspect-ids", strings.Join(aspectIds, ","))
		}
		result, err := getPathOptions(query)
		if err != nil {
			t.Error(err)
			return
		}
		comparePathOptions(t, result, expectedResult)
	}
}

func comparePathOptions(t *testing.T, result map[string][]marshaller.PathOptionsResultElement, expectedResult map[string][]marshaller.PathOptionsResultElement) {
	t.Helper()
	if !reflect.DeepEqual(result, expectedResult) {
		resultJson, _ := json.Marshal(result)
		expectedJson, _ := json.Marshal(expectedResult)
		t.Error("\n", string(resultJson), "\n", string(expectedJson))
	}
}

// getPathOptions calls GET /path-options, the one endpoint with no client method.
func getPathOptions(query url.Values) (result map[string][]marshaller.PathOptionsResultElement, err error) {
	resp, err := (&http.Client{Timeout: 10 * time.Second}).Get(ServerUrl + "/path-options?" + query.Encode())
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return result, fmt.Errorf("unexpected statuscode %v: %v", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

// TestPathOptionsWithoutEnvelope covers the without-envelope query parameter of
// GET /path-options. It used to be read from the function-id parameter, so the envelope
// could not be switched off at all through the query-parameter form.
func TestPathOptionsWithoutEnvelope(t *testing.T) {
	functionId := "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b"
	deviceTypeId := "TestPathOptionsWithoutEnvelope"
	serviceId := deviceTypeId + ".temperature"

	mocks.DeviceRepo.SetDeviceType(model.DeviceType{
		Id:   deviceTypeId,
		Name: deviceTypeId,
		Services: []model.Service{
			{
				Id:      serviceId,
				LocalId: serviceId,
				Name:    serviceId,
				Outputs: []model.Content{
					{
						ContentVariable: model.ContentVariable{
							Name:             "temperature",
							Type:             model.Float,
							CharacteristicId: temperature.Celsius,
							FunctionId:       functionId,
							AspectIds:        []string{"inside_air"},
						},
					},
				},
			},
		},
	})

	expected := func(path string) map[string][]marshaller.PathOptionsResultElement {
		return map[string][]marshaller.PathOptionsResultElement{
			deviceTypeId: {
				{
					ServiceId:              serviceId,
					JsonPath:               []string{path},
					PathToCharacteristicId: map[string]string{path: temperature.Celsius},
				},
			},
		}
	}

	t.Run("keeps the envelope by default", testPathOptionsGetRaw(
		deviceTypeId, functionId, "", expected("value.temperature")))

	t.Run("drops the envelope when without-envelope is true", testPathOptionsGetRaw(
		deviceTypeId, functionId, "true", expected("temperature")))

	t.Run("keeps the envelope when without-envelope is false", testPathOptionsGetRaw(
		deviceTypeId, functionId, "false", expected("value.temperature")))

	t.Run("rejects a without-envelope that is not a bool", func(t *testing.T) {
		query := url.Values{}
		query.Set("device-type-ids", deviceTypeId)
		query.Set("characteristic-filter", temperature.Celsius)
		query.Set("function-id", functionId)
		query.Set("aspect-ids", "inside_air")
		query.Set("without-envelope", "not-a-bool")
		resp, err := (&http.Client{Timeout: 10 * time.Second}).Get(ServerUrl + "/path-options?" + query.Encode())
		if err != nil {
			t.Error(err.Error())
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Error("expected 400, got", resp.StatusCode)
		}
	})
}

func testPathOptionsGetRaw(deviceTypeId string, functionId string, withoutEnvelope string, expectedResult map[string][]marshaller.PathOptionsResultElement) func(t *testing.T) {
	return func(t *testing.T) {
		query := url.Values{}
		query.Set("device-type-ids", deviceTypeId)
		query.Set("characteristic-filter", temperature.Celsius)
		query.Set("function-id", functionId)
		//the aspect has to be named: a v1 path-options query without one matches only
		//content variables that carry no aspect either
		query.Set("aspect-ids", "inside_air")
		if withoutEnvelope != "" {
			query.Set("without-envelope", withoutEnvelope)
		}
		result, err := getPathOptions(query)
		if err != nil {
			t.Error(err)
			return
		}
		comparePathOptions(t, result, expectedResult)
	}
}
