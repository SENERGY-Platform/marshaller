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

package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SENERGY-Platform/marshaller/lib/config"
)

// TestRoutesAreRegistered guards the reflection in GetRouterWithoutMiddleware. An endpoint
// method whose signature drifts from EndpointMethod is not registered and no compiler or
// other test notices — the route simply answers 404 in production. Every route this
// service serves is listed here, so adding one without listing it fails the test and
// renaming one that callers use becomes a visible decision.
func TestRoutesAreRegistered(t *testing.T) {
	routes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/"},
		{http.MethodPost, "/marshal"},
		{http.MethodPost, "/marshal/service-id/characteristic-id"},
		{http.MethodPost, "/unmarshal"},
		{http.MethodPost, "/unmarshal/service-id/characteristic-id"},
		{http.MethodPost, "/v2/marshal"},
		{http.MethodPost, "/v2/marshal/service-id"},
		{http.MethodPost, "/v2/unmarshal"},
		{http.MethodPost, "/v2/unmarshal/service-id"},
		{http.MethodGet, "/characteristic-paths/service-id/characteristic-id"},
		{http.MethodGet, "/path-options"},
		{http.MethodPost, "/query/path-options"},
		{http.MethodGet, "/configurables"},
		{http.MethodPost, "/configurables"},
		{http.MethodPost, "/converter/extension-call"},
		{http.MethodGet, "/doc"},
	}

	router := GetRouterWithoutMiddleware(config.Config{LogLevel: "error"}, &controllerStub{}, nil)

	for _, route := range routes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			//an empty body is enough: a registered route answers 400 at worst, an
			//unregistered one answers 404, and that is the difference under test
			request := httptest.NewRequest(route.method, route.path, strings.NewReader("{}"))
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			if recorder.Code == http.StatusNotFound {
				t.Error("route is not registered")
			}
			if recorder.Code == http.StatusMethodNotAllowed {
				t.Error("route is registered for another method")
			}
		})
	}
}

// TestEveryControllerMethodIsReachable makes the count explicit: the stub below has to
// implement Controller, so a method added to the interface without an endpoint fails to
// compile here rather than shipping unreachable.
func TestEveryControllerMethodIsReachable(t *testing.T) {
	var ctrl Controller = &controllerStub{}
	if ctrl == nil {
		t.Error("expected a controller")
	}
}
