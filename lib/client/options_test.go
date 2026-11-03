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

package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// headers answers every request with an empty json object and keeps the headers of the last one
type headers struct {
	mux  sync.Mutex
	last http.Header
}

func (this *headers) server(t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		this.mux.Lock()
		this.last = r.Header.Clone()
		this.mux.Unlock()
		_, _ = w.Write([]byte("{}"))
	}))
	t.Cleanup(server.Close)
	return server
}

func (this *headers) get(key string) string {
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.last.Get(key)
}

func TestContextCarriesTheTraceToTheService(t *testing.T) {
	previous := otel.GetTextMapPropagator()
	otel.SetTextMapPropagator(propagation.TraceContext{})
	t.Cleanup(func() { otel.SetTextMapPropagator(previous) })
	traceId, _ := trace.TraceIDFromHex("0af7651916cd43dd8448eb211c80319c")
	spanId, _ := trace.SpanIDFromHex("b7ad6b7169203331")
	ctx := trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{TraceID: traceId, SpanID: spanId, TraceFlags: trace.FlagsSampled}))
	rec := &headers{}
	c := NewClient(rec.server(t).URL, Options{})

	if _, err, _ := c.MarshalV2(ctx, messages.MarshallingV2Request{}); err != nil {
		t.Fatal(err)
	}
	if got := rec.get("traceparent"); got != "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01" {
		t.Errorf("post sent traceparent %q", got)
	}
	if _, err, _ := c.GetCharacteristicPaths(ctx, "s", "c"); err != nil {
		t.Fatal(err)
	}
	if got := rec.get("traceparent"); got != "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01" {
		t.Errorf("get sent traceparent %q", got)
	}
}

// hanging answers only when the test ends or the request is cancelled
func hanging(t *testing.T) *httptest.Server {
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-release:
		case <-r.Context().Done():
		}
	}))
	t.Cleanup(func() {
		close(release)
		server.Close()
	})
	return server
}

func expectBounded(t *testing.T, call func() error) {
	t.Helper()
	done := make(chan error, 1)
	go func() { done <- call() }()
	select {
	case err := <-done:
		if err == nil {
			t.Error("a request past its bound has to fail")
		}
	case <-time.After(5 * time.Second):
		t.Error("the request was not bounded")
	}
}

func TestContextCancelsTheRequest(t *testing.T) {
	c := NewClient(hanging(t).URL, Options{})
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	expectBounded(t, func() error {
		_, err, _ := c.UnmarshalV2(ctx, messages.UnmarshallingV2Request{})
		return err
	})
}

func TestOptionsTimeoutBoundsARequest(t *testing.T) {
	c := NewClient(hanging(t).URL, Options{Timeout: 50 * time.Millisecond})
	expectBounded(t, func() error {
		_, err, _ := c.UnmarshalV2(context.Background(), messages.UnmarshallingV2Request{})
		return err
	})
}

type countingTransport struct {
	mux   sync.Mutex
	count int
}

func (this *countingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	this.mux.Lock()
	this.count++
	this.mux.Unlock()
	return http.DefaultTransport.RoundTrip(req)
}

func TestOptionsHttpClientSendsTheRequestsAndKeepsItsTimeout(t *testing.T) {
	rec := &headers{}
	transport := &countingTransport{}
	given := &http.Client{Transport: transport, Timeout: time.Hour}

	if _, err, _ := NewClient(rec.server(t).URL, Options{HttpClient: given, Timeout: time.Minute}).MarshalV2(context.Background(), messages.MarshallingV2Request{}); err != nil {
		t.Fatal(err)
	}
	if transport.count != 1 {
		t.Errorf("the request has to go through the given client, got %v requests", transport.count)
	}
	if given.Timeout != time.Hour {
		t.Errorf("the timeout is set on a copy, the given client changed to %v", given.Timeout)
	}
}

func TestOptionsTokenForTheApiGatewayIsSentWithoutOwnToken(t *testing.T) {
	rec := &headers{}
	token := func() (string, error) { return "Bearer gateway", nil }

	if _, err, _ := NewClient(rec.server(t).URL, Options{OptionalAuthTokenForApiGatewayRequest: token}).MarshalV2(context.Background(), messages.MarshallingV2Request{}); err != nil {
		t.Fatal(err)
	}
	if got := rec.get("Authorization"); got != "Bearer gateway" {
		t.Errorf("Authorization %q", got)
	}
}
