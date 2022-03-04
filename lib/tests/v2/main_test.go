/*
 * Copyright 2022 InfAI (CC SES)
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

package v2

import (
	"context"
	"github.com/SENERGY-Platform/marshaller/lib/api"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	v2 "github.com/SENERGY-Platform/marshaller/lib/marshaller/v2"
	"github.com/SENERGY-Platform/marshaller/lib/tests/mocks"
	"net/http/httptest"
	"sync"
)

func setup(ctx context.Context, done *sync.WaitGroup) (serverUrl string) {
	conceptRepo, err := mocks.NewMockConceptRepo(ctx)
	if err != nil {
		panic(err)
	}
	marshaller := marshaller.New(mocks.Converter{}, conceptRepo, mocks.DeviceRepo)
	marshallerv2 := v2.New(config.Config{ReturnUnknownPathAsNull: true}, mocks.Converter{}, conceptRepo)
	configurableService := configurables.New(conceptRepo)
	done.Add(1)
	server := httptest.NewServer(api.GetRouter(marshaller, marshallerv2, configurableService, mocks.DeviceRepo))
	serverUrl = server.URL
	go func() {
		<-ctx.Done()
		server.Close()
		done.Done()
	}()
	return
}
