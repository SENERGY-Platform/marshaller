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

package converter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/SENERGY-Platform/marshaller/lib/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TestCastCarriesTheTraceToTheConverter(t *testing.T) {
	previous := otel.GetTextMapPropagator()
	otel.SetTextMapPropagator(propagation.TraceContext{})
	t.Cleanup(func() { otel.SetTextMapPropagator(previous) })
	traceId, _ := trace.TraceIDFromHex("0af7651916cd43dd8448eb211c80319c")
	spanId, _ := trace.SpanIDFromHex("b7ad6b7169203331")
	ctx := trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{TraceID: traceId, SpanID: spanId, TraceFlags: trace.FlagsSampled}))

	var mux sync.Mutex
	sent := map[string]string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/openid-connect/token") {
			_ = json.NewEncoder(w).Encode(config.OpenidToken{AccessToken: "service", ExpiresIn: 300})
			return
		}
		mux.Lock()
		sent[r.URL.Path] = r.Header.Get("traceparent")
		mux.Unlock()
		_, _ = w.Write([]byte(`{"output": 1}`))
	}))
	defer server.Close()
	conf := config.Config{ConverterUrl: server.URL, AuthEndpoint: server.URL, LogLevel: "error"}
	c := New(conf, config.NewAccess(conf))

	if _, err := c.Cast(ctx, 1, "a", "b"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.CastWithExtension(ctx, 1, "a", "b", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := c.TryExtension(ctx, ExtensionCall{Input: 1}); err != nil {
		t.Fatal(err)
	}

	mux.Lock()
	defer mux.Unlock()
	for _, path := range []string{"/conversions/a/b", "/extended-conversions/a/b", "/extension-call"} {
		if got := sent[path]; got != "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01" {
			t.Errorf("%v sent traceparent %q", path, got)
		}
	}
}
