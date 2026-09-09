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

// Package controller holds the operations of this service, free of http types. It
// implements api.Controller; the interface is declared there because lib/client has to
// satisfy the same contract, and a controller that imported the api package would close
// that loop.
package controller

import (
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	"github.com/SENERGY-Platform/marshaller/lib/converter"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	v2 "github.com/SENERGY-Platform/marshaller/lib/marshaller/v2"
)

// DeviceRepository is the part of the device-repository this service reads.
type DeviceRepository interface {
	GetService(serviceId string) (model.Service, error)
	GetProtocol(id string) (model.Protocol, error)
	GetServiceWithErrCode(serviceId string) (model.Service, error, int)
	GetAspectNode(id string) (model.AspectNode, error)
}

type Controller struct {
	config              config.Config
	marshaller          *marshaller.Marshaller
	marshallerV2        *v2.Marshaller
	configurableService *configurables.ConfigurableService
	deviceRepo          DeviceRepository
	converter           *converter.Converter
}

// New bundles the services the endpoints used to receive one by one. The converter may be
// nil, the way the api tolerated it before: only TryConverterExtension needs it.
func New(config config.Config, m *marshaller.Marshaller, mV2 *v2.Marshaller, configurableService *configurables.ConfigurableService, deviceRepo DeviceRepository, c *converter.Converter) *Controller {
	return &Controller{
		config:              config,
		marshaller:          m,
		marshallerV2:        mV2,
		configurableService: configurableService,
		deviceRepo:          deviceRepo,
		converter:           c,
	}
}
