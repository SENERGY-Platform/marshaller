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

package api

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/SENERGY-Platform/marshaller/lib/api/metrics"
	"github.com/SENERGY-Platform/marshaller/lib/api/util"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	"github.com/SENERGY-Platform/marshaller/lib/converter"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	v2 "github.com/SENERGY-Platform/marshaller/lib/marshaller/v2"
	"github.com/SENERGY-Platform/service-commons/pkg/accesslog"
	"github.com/julienschmidt/httprouter"
)

type DeviceRepository interface {
	GetService(serviceId string) (model.Service, error)
	GetProtocol(id string) (model.Protocol, error)
	GetServiceWithErrCode(serviceId string) (model.Service, error, int)
	GetAspectNode(id string) (model.AspectNode, error)
}

var endpoints = []func(router *httprouter.Router, config config.Config, marshaller *marshaller.Marshaller, marshallerV2 *v2.Marshaller, configurableService *configurables.ConfigurableService, deviceRepo DeviceRepository, converter *converter.Converter, metrics *metrics.Metrics){}

func Start(ctx context.Context, config config.Config, marshaller *marshaller.Marshaller, marshallerV2 *v2.Marshaller, configurableService *configurables.ConfigurableService, deviceRepo DeviceRepository, converter *converter.Converter) (closed context.Context) {
	config.GetLogger().Info("start api")
	m, err := metrics.Start(ctx, config)
	if err != nil {
		config.GetLogger().Warn("unable to serve metrics", "error", err)
	}
	router := GetRouter(config, marshaller, marshallerV2, configurableService, deviceRepo, converter, m)
	config.GetLogger().Info("add logging and cors")
	corsHandler := util.NewCors(router)
	logger := accesslog.New(corsHandler)
	config.GetLogger().Info("listen on port", "port", config.ServerPort)
	srv := &http.Server{Addr: ":" + config.ServerPort, Handler: logger}
	closed, close := context.WithCancel(context.Background())
	go func() {
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			config.GetLogger().Error("unable to listen and serve http", "error", err)
		}
		close()
	}()
	go func() {
		<-ctx.Done()
		timeout, _ := context.WithTimeout(context.Background(), 2*time.Second)
		if err := srv.Shutdown(timeout); err != nil {
			srv.Close()
		}
	}()
	return closed
}

func GetRouter(config config.Config, marshaller *marshaller.Marshaller, marshallerV2 *v2.Marshaller, configurableService *configurables.ConfigurableService, deviceRepo DeviceRepository, converter *converter.Converter, metrics *metrics.Metrics) (router *httprouter.Router) {
	router = httprouter.New()
	for _, e := range endpoints {
		config.GetLogger().Info("add endpoints", "endpoint", runtime.FuncForPC(reflect.ValueOf(e).Pointer()).Name())
		e(router, config, marshaller, marshallerV2, configurableService, deviceRepo, converter, metrics)
	}
	return
}
