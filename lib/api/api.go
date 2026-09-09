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
	"context"
	"errors"
	"net/http"
	"reflect"
	"time"

	"github.com/SENERGY-Platform/marshaller/lib/api/metrics"
	"github.com/SENERGY-Platform/marshaller/lib/api/util"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/service-commons/pkg/accesslog"
)

// EndpointMethod is the shape an endpoint registration has to have to be picked up by
// GetRouterWithoutMiddleware. Unlike the device-repository's version it carries the
// metrics collector, because this service records per-request metrics that need the
// decoded request body and therefore cannot sit in a middleware.
type EndpointMethod = func(config config.Config, router *http.ServeMux, ctrl Controller, m *metrics.Metrics)

var endpoints = []interface{}{} //list of objects with EndpointMethod

func Start(ctx context.Context, config config.Config, ctrl Controller) (closed context.Context) {
	config.GetLogger().Info("start api")
	m, err := metrics.Start(ctx, config)
	if err != nil {
		config.GetLogger().Warn("unable to serve metrics", "error", err)
	}
	handler := GetRouter(config, ctrl, m)
	config.GetLogger().Info("listen on port", "port", config.ServerPort)
	srv := &http.Server{Addr: ":" + config.ServerPort, Handler: handler}
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
		timeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := srv.Shutdown(timeout); err != nil {
			srv.Close()
		}
	}()
	return closed
}

// GetRouter doc
// @title         Marshaller API
// @version       0.1
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath  /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func GetRouter(config config.Config, ctrl Controller, m *metrics.Metrics) http.Handler {
	handler := GetRouterWithoutMiddleware(config, ctrl, m)
	config.GetLogger().Info("add cors")
	corsHandler := util.NewCors(handler)
	config.GetLogger().Info("add logging")
	return accesslog.New(corsHandler)
}

// GetRouterWithoutMiddleware returns the routes without cors and access logging. It exists
// for consumers that want to run this service in-process in their tests: mounted in an
// httptest server it answers over http without needing a token issuer.
func GetRouterWithoutMiddleware(config config.Config, ctrl Controller, m *metrics.Metrics) http.Handler {
	router := http.NewServeMux()
	config.GetLogger().Info("add heart beat endpoint")
	router.HandleFunc("GET /{$}", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	for _, e := range endpoints {
		for name, call := range getEndpointMethods(e) {
			config.GetLogger().Info("add endpoint " + name)
			call(config, router, ctrl, m)
		}
	}
	return router
}

func getEndpointMethods(e interface{}) map[string]EndpointMethod {
	result := map[string]EndpointMethod{}
	objRef := reflect.ValueOf(e)
	methodCount := objRef.NumMethod()
	for i := 0; i < methodCount; i++ {
		m := objRef.Method(i)
		f, ok := m.Interface().(EndpointMethod)
		if ok {
			name := getTypeName(objRef.Type()) + "::" + objRef.Type().Method(i).Name
			result[name] = f
		}
	}
	return result
}

func getTypeName(t reflect.Type) (res string) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}
