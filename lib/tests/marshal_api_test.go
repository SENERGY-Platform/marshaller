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
	"net/url"
	"strings"
	"testing"
)

func ExampleMarshalEmpty() {
	resp, err := post(
		ServerUrl+"/marshal",
		"application/json",
		strings.NewReader(
			`{
					"data": null,
					"characteristic_id": "",
					"service": `+offServiceStr+`,
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
	//<nil> {"data":"{\"power\":false}"}

}

func ExampleMarshal1() {
	if testing.Short() {
		mocks.DeviceRepo.SetServiceJson(philipsHueServiceStr).SetProtocolJson(protocolJson)
	}
	serviceId := "urn:infai:ses:service:1b0ef253-16f7-4b65-8a15-fe79fccf7e70"               //Philips-Extended-Color-Light setColorService
	characteristicId := "urn:infai:ses:characteristic:0fc343ce-4627-4c88-b1e0-d3ed29754af8" //color hex

	resp, err := post(ServerUrl+"/marshal/"+url.PathEscape(serviceId)+"/"+url.PathEscape(characteristicId), "application/json", strings.NewReader(`{"data": "#ff00ff"}`))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	fmt.Println(err, string(result))

	//output:
	//<nil> {"data":"{\"brightness\":100,\"duration\":1,\"hue\":300,\"saturation\":100}"}

}

func ExampleMarshal2() {
	if testing.Short() {
		mocks.DeviceRepo.SetServiceJson(philipsHueServiceStr).SetProtocolJson(protocolJson)
	}
	serviceId := "urn:infai:ses:service:1b0ef253-16f7-4b65-8a15-fe79fccf7e70"               //Philips-Extended-Color-Light setColorService
	characteristicId := "urn:infai:ses:characteristic:0fc343ce-4627-4c88-b1e0-d3ed29754af8" //color hex

	resp, err := post(
		ServerUrl+"/marshal/"+url.PathEscape(serviceId)+"/"+url.PathEscape(characteristicId),
		"application/json",
		strings.NewReader(
			`{
					"data": "#ff00ff",
					"configurables": [{
						"characteristic_id": "urn:infai:ses:characteristic:9e1024da-3b60-4531-9f29-464addccb13c",
						"values": [{
							"path": "",
							"value": "3",
							"value_type": "https://schema.org/Float"
						}]
					}]
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
	//<nil> {"data":"{\"brightness\":100,\"duration\":3,\"hue\":300,\"saturation\":100}"}

}

func ExampleMarshal3() {
	if testing.Short() {
		mocks.DeviceRepo.SetProtocolJson(protocolJson)
	}

	resp, err := post(
		ServerUrl+"/marshal",
		"application/json",
		strings.NewReader(
			`{
					"data": "#ff00ff",
					"characteristic_id": "urn:infai:ses:characteristic:0fc343ce-4627-4c88-b1e0-d3ed29754af8",
					"service": `+philipsHueServiceStr+`,
					"configurables": [{
						"characteristic_id": "urn:infai:ses:characteristic:9e1024da-3b60-4531-9f29-464addccb13c",
						"values": [{
							"path": "",
							"value": "3",
							"value_type": "https://schema.org/Float"
						}]
					}]
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
	//<nil> {"data":"{\"brightness\":100,\"duration\":3,\"hue\":300,\"saturation\":100}"}

}

func ExampleMarshal4() {
	resp, err := post(
		ServerUrl+"/marshal",
		"application/json",
		strings.NewReader(
			`{
					"data": "#ff00ff",
					"characteristic_id": "urn:infai:ses:characteristic:0fc343ce-4627-4c88-b1e0-d3ed29754af8",
					"service": `+philipsHueServiceStr+`,
					"protocol": `+protocolJson+`,
					"configurables": [{
						"characteristic_id": "urn:infai:ses:characteristic:9e1024da-3b60-4531-9f29-464addccb13c",
						"values": [{
							"path": "",
							"value": "3",
							"value_type": "https://schema.org/Float"
						}]
					}]
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
	//<nil> {"data":"{\"brightness\":100,\"duration\":3,\"hue\":300,\"saturation\":100}"}

}

func ExampleMarshalWithUnusedConfigurable() {
	resp, err := post(
		ServerUrl+"/marshal",
		"application/json",
		strings.NewReader(
			`{
					"data": "#ff00ff",
					"characteristic_id": "urn:infai:ses:characteristic:0fc343ce-4627-4c88-b1e0-d3ed29754af8",
					"service": `+philipsHueServiceStr+`,
					"protocol": `+protocolJson+`,
					"configurables": [{
						"characteristic_id": "urn:infai:ses:characteristic:9e1024da-3b60-4531-9f29-464addccb13c",
						"values": [{
							"path": "",
							"value": "3",
							"value_type": "https://schema.org/Float"
						}]
					},{
						"characteristic_id": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683",
						"values": [{
							"path": "",
							"value": "42",
							"value_type": "https://schema.org/Float"
						}]
					}]
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
	//<nil> {"data":"{\"brightness\":100,\"duration\":3,\"hue\":300,\"saturation\":100}"}

}

const offServiceStr = `{
         "id":"urn:infai:ses:service:59dd05fc-cd67-4f66-98de-bbed8257a868",
         "local_id":"setPower",
         "name":"setPowerOffService",
         "description":"",
         "interaction":"request",
         "aspect_ids":[
            "urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6"
         ],
         "protocol_id":"urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
         "inputs":[
            {
               "id":"urn:infai:ses:content:9f6ac32e-5f26-423b-947a-8a769773239a",
               "content_variable":{
                  "id":"urn:infai:ses:content-variable:b3cd09ff-8d5b-4b35-abad-d998bd46a05f",
                  "name":"struct",
                  "type":"https://schema.org/StructuredValue",
                  "sub_content_variables":[
                     {
                        "id":"urn:infai:ses:content-variable:71fbdd5d-294f-4161-bce5-0ad878d7d14f",
                        "name":"power",
                        "type":"https://schema.org/Boolean",
                        "sub_content_variables":null,
                        "characteristic_id":"urn:infai:ses:characteristic:7dc1bb7e-b256-408a-a6f9-044dc60fdcf5",
                        "value": false,
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
         "outputs":[
            {
               "id":"urn:infai:ses:content:77e18f9c-27f7-4755-82b5-9baba3196666",
               "content_variable":{
                  "id":"urn:infai:ses:content-variable:cbab719c-5105-49f0-9419-e3733ffae1b9",
                  "name":"struct",
                  "type":"https://schema.org/StructuredValue",
                  "sub_content_variables":[
                     {
                        "id":"urn:infai:ses:content-variable:82d51640-9907-4303-b97c-514e8de7c7ad",
                        "name":"status",
                        "type":"https://schema.org/Integer",
                        "sub_content_variables":null,
                        "characteristic_id":"urn:infai:ses:characteristic:c0353532-a8fb-4553-a00b-418cb8a80a65",
                        "value":null,
                        "serialization_options":null
                     }
                  ],
                  "characteristic_id":"",
                  "value":null,
                  "serialization_options":null
               },
               "serialization":"json",
               "protocol_segment_id":"urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a6"
            }
         ],
         "function_ids":[
            "urn:infai:ses:controlling-function:2f35150b-9df7-4cad-95bc-165fa00219fd"
         ],
         "attributes":[
            
         ],
         "service_group_key":"",
         "rdf_type":""
      }`

const philipsHueServiceStr = `{
   "id":"urn:infai:ses:service:1b0ef253-16f7-4b65-8a15-fe79fccf7e70",
   "local_id":"setColor",
   "name":"setColorService",
   "description":"",
   "aspects":[
      {
         "id":"urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6",
         "name":"Lighting",
         "rdf_type":"https://senergy.infai.org/ontology/Aspect"
      }
   ],
   "protocol_id":"urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
   "inputs":[
      {
         "id":"urn:infai:ses:content:5df35b7e-fe01-4c85-89f3-aea176d9455f",
         "content_variable":{
            "id":"urn:infai:ses:content-variable:7538f412-b574-4c42-a9ab-94ca564729cf",
            "name":"struct",
            "type":"https://schema.org/StructuredValue",
            "sub_content_variables":[
               {
                  "id":"urn:infai:ses:content-variable:4096e384-c0ab-475e-9707-f740b379bc62",
                  "name":"brightness",
                  "type":"https://schema.org/Integer",
                  "sub_content_variables":null,
                  "characteristic_id":"urn:infai:ses:characteristic:d840607c-c8f9-45d6-b9bd-2c2d444e2899",
                  "value":null,
                  "serialization_options":null
               },
               {
                  "id":"urn:infai:ses:content-variable:56b643b7-7050-4d72-a261-aca458b8a1ac",
                  "name":"duration",
                  "type":"https://schema.org/Float",
                  "sub_content_variables":null,
                  "characteristic_id":"urn:infai:ses:characteristic:9e1024da-3b60-4531-9f29-464addccb13c",
                  "value":1,
                  "serialization_options":null
               },
               {
                  "id":"urn:infai:ses:content-variable:d5ee69b3-abc3-4974-a5b0-b3bc5f9a4289",
                  "name":"hue",
                  "type":"https://schema.org/Integer",
                  "sub_content_variables":null,
                  "characteristic_id":"urn:infai:ses:characteristic:6ec70e99-8c6a-4909-8d5a-7cc12af76b9a",
                  "value":null,
                  "serialization_options":null
               },
               {
                  "id":"urn:infai:ses:content-variable:22fe620a-0978-4903-85a0-73939a7227c0",
                  "name":"saturation",
                  "type":"https://schema.org/Integer",
                  "sub_content_variables":null,
                  "characteristic_id":"urn:infai:ses:characteristic:a66dc568-c0e0-420f-b513-18e8df405538",
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
   "outputs":[

   ],
   "functions":[
      {
         "id":"urn:infai:ses:controlling-function:c54e2a89-1fb8-4ecb-8993-a7b40b355599",
         "name":"setColorFunction",
         "concept_id":"urn:infai:ses:concept:8b1161d5-7878-4dd2-a36c-6f98f6b94bf8",
         "rdf_type":"https://senergy.infai.org/ontology/ControllingFunction"
      }
   ],
   "rdf_type":""
}`

const protocolJson = `{
   "id":"urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
   "name":"standard-connector",
   "handler":"connector",
   "protocol_segments":[
      {
         "id":"urn:infai:ses:protocol-segment:9956d8b5-46fa-4381-a227-c1df69808997",
         "name":"metadata"
      },
      {
         "id":"urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
         "name":"data"
      }
   ]
}`
