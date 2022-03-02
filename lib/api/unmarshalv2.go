/*
 * Copyright 2019 InfAI (CC SES)
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

package api

import (
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	v2 "github.com/SENERGY-Platform/marshaller/lib/marshaller/v2"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func init() {
	endpoints = append(endpoints, UnmarshallingV2)
}

func UnmarshallingV2(router *httprouter.Router, marshaller *marshaller.Marshaller, marshallerV2 *v2.Marshaller, configurableService *configurables.ConfigurableService, deviceRepo DeviceRepository) {
	resource := "/v2/unmarshal"

	normalizeRequest := func(request *UnmarshallingV2Request) error {
		if request.Protocol.Id == "" {
			protocol, err := deviceRepo.GetProtocol(request.Service.ProtocolId)
			if err != nil {
				return err
			}
			request.Protocol = protocol
		}
		if request.Service.ProtocolId != request.Protocol.Id {
			return errors.New("expect service to reference given protocol")
		}
		if request.Path == "" {
			var aspect *model.AspectNode
			if request.AspectNode.Id == "" && request.AspectNodeId != "" {
				var err error
				request.AspectNode, err = deviceRepo.GetAspectNode(request.AspectNodeId)
				if err != nil {
					return err
				}
			}
			if request.AspectNode.Id != "" {
				aspect = &request.AspectNode
			}
			paths := marshallerV2.GetOutputPaths(request.Service, request.FunctionId, aspect)
			if len(paths) > 0 {
				log.Println("WARNING: only first path found by FunctionId and AspectNode is used for Unmarshal")
			}
			if len(paths) == 0 {
				return errors.New("no output path found for criteria")
			}
			request.Path = paths[0]
		}
		return nil
	}

	unmarshal := func(request UnmarshallingV2Request) (interface{}, error) {
		return marshallerV2.Unmarshal(request.Protocol, request.Service, request.CharacteristicId, request.Path, request.Message)
	}

	router.POST(resource+"/:serviceId", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		msg := UnmarshallingV2Request{}
		serviceId := params.ByName("serviceId")
		if serviceId == "" {
			http.Error(writer, "expect serviceId as parameter in path", http.StatusBadRequest)
			return
		}
		err := json.NewDecoder(request.Body).Decode(&msg)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		msg.Service, err = deviceRepo.GetService(serviceId)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		err = normalizeRequest(&msg)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := unmarshal(msg)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(writer).Encode(result)
		if err != nil {
			log.Println("ERROR: unable to encode response", err)
		}
	})

	router.POST(resource, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		msg := UnmarshallingV2Request{}
		err := json.NewDecoder(request.Body).Decode(&msg)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		err = normalizeRequest(&msg)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := unmarshal(msg)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(writer).Encode(result)
		if err != nil {
			log.Println("ERROR: unable to encode response", err)
		}
	})

}
