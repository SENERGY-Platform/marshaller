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
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"net/http"
)

var DeviceRepo = (&DeviceRepoStruct{}).Init()

type DeviceRepoStruct struct {
	deviceTypes map[string]model.DeviceType
	services    map[string]model.Service
	protocols   map[string]model.Protocol
}

func (this *DeviceRepoStruct) Init() *DeviceRepoStruct {
	this.services = map[string]model.Service{}
	this.protocols = map[string]model.Protocol{}
	this.deviceTypes = map[string]model.DeviceType{}
	return this
}

func (this *DeviceRepoStruct) GetDeviceType(id string) (result model.DeviceType, err error, code int) {
	if dt, ok := this.deviceTypes[id]; ok {
		return dt, nil, 200
	} else {
		return dt, errors.New("not found"), 404
	}
}

func (this *DeviceRepoStruct) SetDeviceType(dt model.DeviceType) *DeviceRepoStruct {
	this.deviceTypes[dt.Id] = dt
	return this
}

func (this *DeviceRepoStruct) GetService(serviceId string) (model.Service, error) {
	if service, ok := this.services[serviceId]; ok {
		return service, nil
	} else {
		return model.Service{}, errors.New("not found")
	}
}

func (this *DeviceRepoStruct) GetProtocol(id string) (model.Protocol, error) {
	if protocol, ok := this.protocols[id]; ok {
		return protocol, nil
	} else {
		return model.Protocol{}, errors.New("not found")
	}
}

func (this *DeviceRepoStruct) SetProtocol(protocol model.Protocol) *DeviceRepoStruct {
	this.protocols[protocol.Id] = protocol
	return this
}

func (this *DeviceRepoStruct) SetService(service model.Service) *DeviceRepoStruct {
	this.services[service.Id] = service
	return this
}

func (this *DeviceRepoStruct) SetServiceJson(serviceStr string) *DeviceRepoStruct {
	service := model.Service{}
	json.Unmarshal([]byte(serviceStr), &service)
	return this.SetService(service)
}

func (this *DeviceRepoStruct) SetProtocolJson(protocolStr string) *DeviceRepoStruct {
	protocol := model.Protocol{}
	json.Unmarshal([]byte(protocolStr), &protocol)
	return this.SetProtocol(protocol)
}

func (this *DeviceRepoStruct) GetServiceWithErrCode(serviceId string) (result model.Service, err error, code int) {
	result, err = this.GetService(serviceId)
	if err != nil {
		code = http.StatusInternalServerError
	} else {
		code = 200
	}
	return
}

func (this *DeviceRepoStruct) GetAspectNode(id string) (result model.AspectNode, err error) {
	return model.AspectNode{
		Id:            id,
		Name:          id,
		RootId:        id,
		ParentId:      "",
		ChildIds:      []string{},
		AncestorIds:   []string{},
		DescendentIds: []string{},
	}, nil
}
