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
	"github.com/SENERGY-Platform/marshaller/lib/tests/mocks"
	"io"
	"net/url"
)

func ExampleGet_characteristicsPathSame() {
	mocks.DeviceRepo.SetServiceJson(danfossTemperatureService)
	serviceId := "urn:infai:ses:service:f306de41-a55b-45ed-afc9-039bbe53db1b"               //danfos getTemperatureService
	characteristicId := "urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a" //temperature celcius

	resp, err := get(ServerUrl + "/characteristic-paths/" + url.PathEscape(serviceId) + "/" + url.PathEscape(characteristicId))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	fmt.Println(err, resp.StatusCode, string(result))

	//output:
	//<nil> 200 {"path":"temperature.level","service_characteristic_id":"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"}

}

func ExampleGet_characteristicsPathMatching() {
	mocks.DeviceRepo.SetServiceJson(danfossTemperatureService)
	serviceId := "urn:infai:ses:service:f306de41-a55b-45ed-afc9-039bbe53db1b"               //danfos getTemperatureService
	characteristicId := "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683" //temperature kelvin

	resp, err := get(ServerUrl + "/characteristic-paths/" + url.PathEscape(serviceId) + "/" + url.PathEscape(characteristicId))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	fmt.Println(err, resp.StatusCode, string(result))

	//output:
	//<nil> 200 {"path":"temperature.level","service_characteristic_id":"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a"}

}

func ExampleGet_characteristicsPathNotMatching() {
	mocks.DeviceRepo.SetServiceJson(danfossTemperatureService)
	serviceId := "urn:infai:ses:service:f306de41-a55b-45ed-afc9-039bbe53db1b"               //danfos getTemperatureService
	characteristicId := "urn:infai:ses:characteristic:5caa707d-dc08-4f3b-bd9f-f08935c8dd3c" //percentage

	resp, err := get(ServerUrl + "/characteristic-paths/" + url.PathEscape(serviceId) + "/" + url.PathEscape(characteristicId))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	fmt.Println(err, resp.StatusCode, string(result))

	//output:
	//<nil> 404 characteristic not in service

}
