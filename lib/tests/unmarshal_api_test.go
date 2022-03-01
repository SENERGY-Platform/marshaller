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
	"io/ioutil"
	"log"
	"net/url"
	"strings"
)

func ExampleUnmarshal1() {
	mocks.DeviceRepo.SetServiceJson(danfossTemperatureService).SetProtocolJson(protocolJson)
	serviceId := "urn:infai:ses:service:f306de41-a55b-45ed-afc9-039bbe53db1b"               //Danfoss Radiator Thermostat getTemperatureService
	characteristicId := "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683" //temperature kelvin

	resp, err := post(
		ServerUrl+"/unmarshal/"+url.PathEscape(serviceId)+"/"+url.PathEscape(characteristicId),
		"application/json",
		strings.NewReader(
			`{"message": {"data":"{\"level\":21,\"updateTime\":\"2020-01-15T07:20:01.000Z\"}"}}`,
		),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	fmt.Println(err, string(result))

	//output:
	//<nil> 294.15

}

func ExampleUnmarshal2() {
	mocks.DeviceRepo.SetProtocolJson(protocolJson)

	resp, err := post(
		ServerUrl+"/unmarshal",
		"application/json",
		strings.NewReader(
			`{
					"message": {"data":"{\"level\":21,\"updateTime\":\"2020-01-15T07:20:01.000Z\"}"},
					"characteristic_id": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683",
					"service": `+danfossTemperatureService+`
				}`,
		),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	fmt.Println(err, string(result))

	//output:
	//<nil> 294.15

}

func ExampleUnmarshal3() {
	resp, err := post(
		ServerUrl+"/unmarshal",
		"application/json",
		strings.NewReader(
			`{
					"message": {"data":"{\"level\":21,\"updateTime\":\"2020-01-15T07:20:01.000Z\"}"},
					"characteristic_id": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683",
					"service": `+danfossTemperatureService+`,
					"protocol": `+protocolJson+`
				}`,
		),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	fmt.Println(err, string(result))

	//output:
	//<nil> 294.15

}

func ExampleUnmarshalWithHints1() {
	resp, err := post(
		ServerUrl+"/unmarshal",
		"application/json",
		strings.NewReader(
			`{
					"message": {"data":"{\"level1\":21,\"level2\":22,\"updateTime\":\"2020-01-15T07:20:01.000Z\"}"},
					"characteristic_id": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683",
					"service": `+danfossTemperatureServiceForHints+`,
					"protocol": `+protocolJson+`,
					"content_variable_hints": ["urn:infai:ses:content-variable:c504db64-05ea-4736-89fb-8a7a04d5c468_1", "foo", "bar"]
				}`,
		),
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(err, string(result))

	//output:
	//<nil> 294.15

}

func ExampleUnmarshalWithHints2() {
	resp, err := post(
		ServerUrl+"/unmarshal",
		"application/json",
		strings.NewReader(
			`{
					"message": {"data":"{\"level1\":21,\"level2\":22,\"updateTime\":\"2020-01-15T07:20:01.000Z\"}"},
					"characteristic_id": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683",
					"service": `+danfossTemperatureServiceForHints+`,
					"protocol": `+protocolJson+`,
					"content_variable_hints": ["urn:infai:ses:content-variable:c504db64-05ea-4736-89fb-8a7a04d5c468_2", "foo", "bar"]
				}`,
		),
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(err, string(result))

	//output:
	//<nil> 295.15

}

//Danfoss Radiator Thermostat
const danfossTemperatureService = `{
   "id":"urn:infai:ses:service:f306de41-a55b-45ed-afc9-039bbe53db1b",
   "local_id":"get_level:67-1",
   "name":"getTemperatureService",
   "description":"",
   "aspects":[
      {
         "id":"urn:infai:ses:aspect:a14c5efb-b0b6-46c3-982e-9fded75b5ab6",
         "name":"Air",
         "rdf_type":"https://senergy.infai.org/ontology/Aspect"
      }
   ],
   "protocol_id":"urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
   "inputs":[

   ],
   "outputs":[
      {
         "id":"urn:infai:ses:content:1a4ebdd9-bfc5-4208-b8ed-826286792d21",
         "content_variable":{
            "id":"urn:infai:ses:content-variable:fad90f83-32f2-4e5d-9aa6-5efbe24a8cac",
            "name":"temperature",
            "type":"https://schema.org/StructuredValue",
            "sub_content_variables":[
               {
                  "id":"urn:infai:ses:content-variable:c504db64-05ea-4736-89fb-8a7a04d5c468",
                  "name":"level",
                  "type":"https://schema.org/Float",
                  "sub_content_variables":null,
                  "characteristic_id":"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a",
                  "value":0,
                  "serialization_options":null
               },
               {
                  "id":"urn:infai:ses:content-variable:ad1ebce4-31b4-47ce-ba00-42ab7b44e982",
                  "name":"updateTime",
                  "type":"https://schema.org/Text",
                  "sub_content_variables":null,
                  "characteristic_id":"",
                  "value":null,
                  "serialization_options":null
               }
            ],
            "characteristic_id":"",
            "value":null,
            "serialization_options":null
         },
         "serialization":"json",
         "protocol_segment_id":"urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65"
      }
   ],
   "functions":[
      {
         "id":"urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
         "name":"getTemperatureFunction",
         "concept_id":"urn:infai:ses:concept:0bc81398-3ed6-4e2b-a6c4-b754583aac37",
         "rdf_type":"https://senergy.infai.org/ontology/MeasuringFunction"
      }
   ],
   "rdf_type":""
}`

//Danfoss Radiator Thermostat
const danfossTemperatureServiceForHints = `{
   "id":"urn:infai:ses:service:f306de41-a55b-45ed-afc9-039bbe53db1b",
   "local_id":"get_level:67-1",
   "name":"getTemperatureService",
   "description":"",
   "aspects":[
      {
         "id":"urn:infai:ses:aspect:a14c5efb-b0b6-46c3-982e-9fded75b5ab6",
         "name":"Air",
         "rdf_type":"https://senergy.infai.org/ontology/Aspect"
      }
   ],
   "protocol_id":"urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
   "inputs":[

   ],
   "outputs":[
      {
         "id":"urn:infai:ses:content:1a4ebdd9-bfc5-4208-b8ed-826286792d21",
         "content_variable":{
            "id":"urn:infai:ses:content-variable:fad90f83-32f2-4e5d-9aa6-5efbe24a8cac",
            "name":"temperature",
            "type":"https://schema.org/StructuredValue",
            "sub_content_variables":[
				{
                  "id":"urn:infai:ses:content-variable:c504db64-05ea-4736-89fb-8a7a04d5c468_1",
                  "name":"level1",
                  "type":"https://schema.org/Float",
                  "sub_content_variables":null,
                  "characteristic_id":"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a",
                  "value":0,
                  "serialization_options":null
               },
               {
                  "id":"urn:infai:ses:content-variable:c504db64-05ea-4736-89fb-8a7a04d5c468_2",
                  "name":"level2",
                  "type":"https://schema.org/Float",
                  "sub_content_variables":null,
                  "characteristic_id":"urn:infai:ses:characteristic:5ba31623-0ccb-4488-bfb7-f73b50e03b5a",
                  "value":0,
                  "serialization_options":null
               },
               {
                  "id":"urn:infai:ses:content-variable:ad1ebce4-31b4-47ce-ba00-42ab7b44e982",
                  "name":"updateTime",
                  "type":"https://schema.org/Text",
                  "sub_content_variables":null,
                  "characteristic_id":"",
                  "value":null,
                  "serialization_options":null
               }
            ],
            "characteristic_id":"",
            "value":null,
            "serialization_options":null
         },
         "serialization":"json",
         "protocol_segment_id":"urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65"
      }
   ],
   "functions":[
      {
         "id":"urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
         "name":"getTemperatureFunction",
         "concept_id":"urn:infai:ses:concept:0bc81398-3ed6-4e2b-a6c4-b754583aac37",
         "rdf_type":"https://senergy.infai.org/ontology/MeasuringFunction"
      }
   ],
   "rdf_type":""
}`
