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

package devicerepository

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/models/go/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const traceparent = "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"

// tracedContext carries the span of traceparent, with the propagator that writes it
func tracedContext(t *testing.T) context.Context {
	previous := otel.GetTextMapPropagator()
	otel.SetTextMapPropagator(propagation.TraceContext{})
	t.Cleanup(func() { otel.SetTextMapPropagator(previous) })
	traceId, _ := trace.TraceIDFromHex("0af7651916cd43dd8448eb211c80319c")
	spanId, _ := trace.SpanIDFromHex("b7ad6b7169203331")
	return trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{TraceID: traceId, SpanID: spanId, TraceFlags: trace.FlagsSampled}))
}

// repository answers like device-repository and the token endpoint, and records the
// traceparent and Authorization header of every read by its path
type repository struct {
	mux           sync.Mutex
	traceparent   map[string]string
	authorization map[string]string
}

func (this *repository) server(t *testing.T) *httptest.Server {
	this.traceparent, this.authorization = map[string]string{}, map[string]string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/openid-connect/token") {
			_ = json.NewEncoder(w).Encode(config.OpenidToken{AccessToken: "service", ExpiresIn: 300})
			return
		}
		this.mux.Lock()
		this.traceparent[r.URL.Path] = r.Header.Get("traceparent")
		this.authorization[r.URL.Path] = r.Header.Get("Authorization")
		this.mux.Unlock()
		id := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
		_ = json.NewEncoder(w).Encode(map[string]string{"id": id})
	}))
	t.Cleanup(server.Close)
	return server
}

func TestReadsCarryTheTraceAndTheTokenOfTheService(t *testing.T) {
	ctx := tracedContext(t)
	rec := &repository{}
	server := rec.server(t)
	conf := config.Config{DeviceRepositoryUrl: server.URL, AuthEndpoint: server.URL, LogLevel: "error"}
	repo, err := New(conf, config.NewAccess(conf))
	if err != nil {
		t.Fatal(err)
	}

	if _, err = repo.GetService(ctx, "s1"); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.GetProtocol(ctx, "p1"); err != nil {
		t.Fatal(err)
	}
	if _, err, _ = repo.GetDeviceType(ctx, "dt1"); err != nil {
		t.Fatal(err)
	}
	var node models.AspectNode
	if node, err = repo.GetAspectNode(ctx, "a1"); err != nil || node.Id != "a1" {
		t.Fatal(node, err)
	}

	rec.mux.Lock()
	defer rec.mux.Unlock()
	for _, path := range []string{"/services/s1", "/protocols/p1", "/device-types/dt1", "/aspect-nodes/a1"} {
		if got := rec.traceparent[path]; got != traceparent {
			t.Errorf("%v sent traceparent %q", path, got)
		}
		if got := rec.authorization[path]; got != "Bearer service" {
			t.Errorf("%v sent Authorization %q", path, got)
		}
	}
}
