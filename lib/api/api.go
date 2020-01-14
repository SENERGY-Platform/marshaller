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
	"github.com/SENERGY-Platform/marshaller/lib/api/util"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	jwt_http_router "github.com/SmartEnergyPlatform/jwt-http-router"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

type DeviceRepository interface {
	GetService(token config.Impersonate, serviceId string) (model.Service, error)
	GetProtocol(token config.Impersonate, id string) (model.Protocol, error)
}

var endpoints = []func(router *jwt_http_router.Router, config config.Config, marshaller *marshaller.Marshaller, configurableService *configurables.ConfigurableService, deviceRepo DeviceRepository){}

func Start(ctx context.Context, config config.Config, marshaller *marshaller.Marshaller, configurableService *configurables.ConfigurableService, deviceRepo DeviceRepository) (closed context.Context) {
	log.Println("start api")
	router := GetRouter(config, marshaller, configurableService, deviceRepo)
	log.Println("add logging and cors")
	corsHandler := util.NewCors(router)
	logger := util.NewLogger(corsHandler, config.LogLevel)
	log.Println("listen on port", config.ServerPort)
	srv := &http.Server{Addr: ":" + config.ServerPort, Handler: logger}
	closed, close := context.WithCancel(context.Background())
	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Println("ERROR:", err)
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

func GetRouter(config config.Config, marshaller *marshaller.Marshaller, configurableService *configurables.ConfigurableService, deviceRepo DeviceRepository) (router *jwt_http_router.Router) {
	router = jwt_http_router.New(jwt_http_router.JwtConfig{ForceAuth: true, ForceUser: true})
	for _, e := range endpoints {
		log.Println("add endpoints: " + runtime.FuncForPC(reflect.ValueOf(e).Pointer()).Name())
		e(router, config, marshaller, configurableService, deviceRepo)
	}
	return
}
