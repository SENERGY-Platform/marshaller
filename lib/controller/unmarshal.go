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

package controller

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
)

func (this *Controller) Unmarshal(request messages.UnmarshallingRequest) (result interface{}, err error, code int) {
	err, code = this.normalizeUnmarshallingRequest(&request)
	if err != nil {
		return result, err, code
	}
	result, err = this.marshaller.UnmarshalOutputs(*request.Protocol, request.Service, request.Message, request.CharacteristicId, request.PathAllowList, request.ContentVariableHints...)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	return result, nil, http.StatusOK
}

func (this *Controller) UnmarshalForService(serviceId string, characteristicId string, request messages.UnmarshallingRequest) (result interface{}, err error, code int) {
	if serviceId == "" {
		return result, errors.New("expect serviceId as parameter in path"), http.StatusBadRequest
	}
	if characteristicId == "" {
		return result, errors.New("expect characteristicId as parameter in path"), http.StatusBadRequest
	}
	request.CharacteristicId = characteristicId
	request.Service, err = this.deviceRepo.GetService(serviceId)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	return this.Unmarshal(request)
}

func (this *Controller) UnmarshalV2(request messages.UnmarshallingV2Request) (result interface{}, err error, code int) {
	err, code = this.normalizeUnmarshallingV2Request(&request)
	if err != nil {
		return result, err, code
	}
	result, err = this.marshallerV2.Unmarshal(request.Protocol, request.Service, request.CharacteristicId, request.Path, request.Message, request.SerializedOutput)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	return result, nil, http.StatusOK
}

func (this *Controller) UnmarshalV2ForService(serviceId string, request messages.UnmarshallingV2Request) (result interface{}, err error, code int) {
	if serviceId == "" {
		return result, errors.New("expect serviceId as parameter in path"), http.StatusBadRequest
	}
	request.Service, err = this.deviceRepo.GetService(serviceId)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	return this.UnmarshalV2(request)
}

// normalizeUnmarshallingRequest answers 400 on both failures, see
// normalizeMarshallingRequest.
func (this *Controller) normalizeUnmarshallingRequest(request *messages.UnmarshallingRequest) (err error, code int) {
	if request.Protocol == nil {
		protocol, err := this.deviceRepo.GetProtocol(request.Service.ProtocolId)
		if err != nil {
			return err, http.StatusBadRequest
		}
		request.Protocol = &protocol
	}
	if request.Service.ProtocolId != request.Protocol.Id {
		return errors.New("expect service to reference given protocol"), http.StatusBadRequest
	}
	return nil, http.StatusOK
}

// normalizeUnmarshallingV2Request additionally determines the output path when the caller
// gave none, from the function and the requested aspects. Every failure answers 400, which
// is what the endpoint did before.
func (this *Controller) normalizeUnmarshallingV2Request(request *messages.UnmarshallingV2Request) (err error, code int) {
	this.config.GetLogger().Debug("UnmarshallingV2Request", "request", fmt.Sprintf("%#v", request))
	if request.Protocol.Id == "" {
		protocol, err := this.deviceRepo.GetProtocol(request.Service.ProtocolId)
		if err != nil {
			return err, http.StatusBadRequest
		}
		request.Protocol = protocol
	}
	if request.Service.ProtocolId != request.Protocol.Id {
		return errors.New("expect service to reference given protocol"), http.StatusBadRequest
	}
	if request.Path == "" {
		aspects, err := this.requestAspectNodes(*request)
		if err != nil {
			return err, http.StatusBadRequest
		}
		paths := this.marshallerV2.GetOutputPaths(request.Service, request.FunctionId, aspects)
		if len(paths) > 1 {
			paths, err = this.marshallerV2.SortPathsByAspectDistance(this.deviceRepo, request.Service, aspects, paths)
			if err != nil {
				this.config.GetLogger().Error("unable to sort paths by aspect distance", "error", err)
				debug.PrintStack()
				return err, http.StatusBadRequest
			}
			this.config.GetLogger().Debug("WARNING: only one path found by FunctionId and AspectNodes is used for Unmarshal", "paths", fmt.Sprintf("%#v", paths))
		}
		if len(paths) == 0 {
			return errors.New("no output path found for criteria"), http.StatusBadRequest
		}
		request.Path = paths[0]
	}
	return nil, http.StatusOK
}

// requestAspectNodes collects the aspect nodes a request asks for and resolves the aspect
// ids among them. Both deprecated single fields are aliases for a list with one element,
// folded in by the request itself.
func (this *Controller) requestAspectNodes(request messages.UnmarshallingV2Request) (result []model.AspectNode, err error) {
	result = request.GetAspectNodes()
	for _, aspectNodeId := range request.GetAspectNodeIds() {
		if aspectNodeId == "" || model.ContainsAspectNode(result, aspectNodeId) {
			continue
		}
		aspectNode, err := this.deviceRepo.GetAspectNode(aspectNodeId)
		if err != nil {
			return nil, err
		}
		result = append(result, aspectNode)
	}
	return result, nil
}
