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

	"github.com/SENERGY-Platform/converter/lib/converter/characteristics"
	"github.com/SENERGY-Platform/marshaller/lib/api"
	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	"github.com/SENERGY-Platform/marshaller/lib/controller"
	"github.com/SENERGY-Platform/marshaller/lib/converter"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	v2 "github.com/SENERGY-Platform/marshaller/lib/marshaller/v2"
	"github.com/SENERGY-Platform/marshaller/lib/tests/mocks"
	"github.com/SENERGY-Platform/models/go/models"
)

// setup runs the real api over http, backed by the real controller and the test mocks, and
// returns a client pointed at it. That is the pairing under test: the client has to reach
// the endpoint the controller method is served under, and the shared interface does not
// prove that on its own — a client method calling the wrong path still compiles.
func setup(t *testing.T, ctx context.Context, wg *sync.WaitGroup) Interface {
	t.Helper()
	conceptRepo, err := mocks.NewMockConceptRepo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	conf := config.Config{ReturnUnknownPathAsNull: true, LogLevel: "error"}
	m := marshaller.New(mocks.Converter{}, conceptRepo, mocks.DeviceRepo)
	mV2 := v2.New(conf, mocks.Converter{}, conceptRepo)
	ctrl := controller.New(conf, m, mV2, configurables.New(conceptRepo), mocks.DeviceRepo, nil)

	wg.Add(1)
	server := httptest.NewServer(api.GetRouter(conf, ctrl, nil))
	go func() {
		<-ctx.Done()
		server.Close()
		wg.Done()
	}()
	return NewClient(server.URL, nil)
}

// TestClientReachesEveryEndpoint calls every method of the interface and rejects a 404 or
// 405. A wrong path or verb in the client is invisible to the compiler and to the api's own
// tests, and it is the mistake this whole package can make.
func TestClientReachesEveryEndpoint(t *testing.T) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := setup(t, ctx, wg)

	protocol := model.Protocol{
		Id:               "p1",
		Name:             "p1",
		Handler:          "p1",
		ProtocolSegments: []model.ProtocolSegment{{Id: "p1.1", Name: "body"}},
	}
	service := model.Service{
		Id:         "sid",
		LocalId:    "slid",
		Name:       "sname",
		ProtocolId: "p1",
		Outputs: []model.Content{{
			Id:                "content",
			ContentVariable:   model.ContentVariable{Id: "value", Name: "value", Type: model.Float, CharacteristicId: characteristics.Celsius},
			Serialization:     "json",
			ProtocolSegmentId: "p1.1",
		}},
		Inputs: []model.Content{{
			Id:                "content",
			ContentVariable:   model.ContentVariable{Id: "value", Name: "value", Type: model.Float, CharacteristicId: characteristics.Celsius},
			Serialization:     "json",
			ProtocolSegmentId: "p1.1",
		}},
	}
	mocks.DeviceRepo.SetService(service)
	mocks.DeviceRepo.SetProtocol(protocol)

	//every call is expected to be answered by its endpoint; what the answer says is the
	//subject of the api and controller tests, not of this one
	calls := map[string]func() (error, int){
		"Marshal": func() (error, int) {
			_, err, code := c.Marshal(messages.MarshallingRequest{Service: service, Protocol: &protocol, CharacteristicId: characteristics.Celsius, Data: 1})
			return err, code
		},
		"MarshalForService": func() (error, int) {
			_, err, code := c.MarshalForService("sid", characteristics.Celsius, messages.MarshallingRequest{Data: 1})
			return err, code
		},
		"Unmarshal": func() (error, int) {
			_, err, code := c.Unmarshal(messages.UnmarshallingRequest{Service: service, Protocol: &protocol, CharacteristicId: characteristics.Celsius, Message: map[string]string{"body": `{"value":1}`}})
			return err, code
		},
		"UnmarshalForService": func() (error, int) {
			_, err, code := c.UnmarshalForService("sid", characteristics.Celsius, messages.UnmarshallingRequest{Message: map[string]string{"body": `{"value":1}`}})
			return err, code
		},
		"MarshalV2": func() (error, int) {
			_, err, code := c.MarshalV2(messages.MarshallingV2Request{Service: service, Protocol: protocol, Data: []model.MarshallingV2RequestData{{Value: 1, CharacteristicId: characteristics.Celsius, Paths: []string{"value"}}}})
			return err, code
		},
		"MarshalV2ForService": func() (error, int) {
			_, err, code := c.MarshalV2ForService("sid", messages.MarshallingV2Request{Data: []model.MarshallingV2RequestData{{Value: 1, CharacteristicId: characteristics.Celsius, Paths: []string{"value"}}}})
			return err, code
		},
		"UnmarshalV2": func() (error, int) {
			_, err, code := c.UnmarshalV2(messages.UnmarshallingV2Request{Service: service, Protocol: protocol, CharacteristicId: characteristics.Celsius, Path: "value", Message: map[string]string{"body": `{"value":1}`}})
			return err, code
		},
		"UnmarshalV2ForService": func() (error, int) {
			_, err, code := c.UnmarshalV2ForService("sid", messages.UnmarshallingV2Request{CharacteristicId: characteristics.Celsius, Path: "value", Message: map[string]string{"body": `{"value":1}`}})
			return err, code
		},
		"GetCharacteristicPaths": func() (error, int) {
			_, err, code := c.GetCharacteristicPaths("sid", characteristics.Celsius)
			return err, code
		},
		"GetPathOptions": func() (error, int) {
			_, err, code := c.GetPathOptions(messages.PathOptionsQuery{})
			return err, code
		},
		"FindConfigurables": func() (error, int) {
			_, err, code := c.FindConfigurables(characteristics.Celsius, []model.Service{service})
			return err, code
		},
		"FindConfigurablesForServiceIds": func() (error, int) {
			_, err, code := c.FindConfigurablesForServiceIds(characteristics.Celsius, []string{"sid"})
			return err, code
		},
	}

	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			err, code := call()
			if code == http.StatusNotFound {
				t.Error("client called a path the api does not serve:", err)
			}
			if code == http.StatusMethodNotAllowed {
				t.Error("client used a verb the api does not serve on that path:", err)
			}
		})
	}
}

// TestTryConverterExtensionReachesEndpoint is separate because the test setup runs without
// a converter, so the endpoint answers 500 by design — which still proves it was reached.
func TestTryConverterExtensionReachesEndpoint(t *testing.T) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := setup(t, ctx, wg)

	_, _, code := c.TryConverterExtension(converterExtensionCall())
	if code == http.StatusNotFound || code == http.StatusMethodNotAllowed {
		t.Error("client did not reach the converter extension endpoint, got", code)
	}
	if code != http.StatusInternalServerError {
		t.Error("expected 500 from a service without a converter, got", code)
	}
}

func converterExtensionCall() converter.ExtensionCall {
	return converter.ExtensionCall{Input: 1, Extension: models.ConverterExtension{}}
}
