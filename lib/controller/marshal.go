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
	"net/http"

	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
)

func (this *Controller) Marshal(request messages.MarshallingRequest) (result map[string]string, err error, code int) {
	err, code = this.normalizeMarshallingRequest(&request)
	if err != nil {
		return result, err, code
	}
	result, err = this.marshaller.MarshalInputs(*request.Protocol, request.Service, request.Data, request.CharacteristicId, request.PathAllowList, request.Configurables...)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	return result, nil, http.StatusOK
}

func (this *Controller) MarshalForService(serviceId string, characteristicId string, request messages.MarshallingRequest) (result map[string]string, err error, code int) {
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
	return this.Marshal(request)
}

func (this *Controller) MarshalV2(request messages.MarshallingV2Request) (result map[string]string, err error, code int) {
	err, code = this.normalizeMarshallingV2Request(&request)
	if err != nil {
		return result, err, code
	}
	result, err = this.marshallerV2.Marshal(request.Protocol, request.Service, request.Data)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	return result, nil, http.StatusOK
}

func (this *Controller) MarshalV2ForService(serviceId string, request messages.MarshallingV2Request) (result map[string]string, err error, code int) {
	if serviceId == "" {
		return result, errors.New("expect serviceId as parameter in path"), http.StatusBadRequest
	}
	request.Service, err = this.deviceRepo.GetService(serviceId)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	return this.MarshalV2(request)
}

// normalizeMarshallingRequest fills in the protocol from the service if the caller left it
// out, and rejects a request whose service and protocol disagree.
//
// Both failures answer 400, which is what the endpoint did before this interface existed:
// it mapped every normalization error to StatusBadRequest, a failing protocol lookup
// included. That one is arguably a 500 — changing it is a behavior change and does not
// belong in a refactor, so it is preserved here rather than quietly corrected.
func (this *Controller) normalizeMarshallingRequest(request *messages.MarshallingRequest) (err error, code int) {
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

// normalizeMarshallingV2Request answers 400 on both failures, see
// normalizeMarshallingRequest.
func (this *Controller) normalizeMarshallingV2Request(request *messages.MarshallingV2Request) (err error, code int) {
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
	return nil, http.StatusOK
}
