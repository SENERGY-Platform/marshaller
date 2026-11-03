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

package devicerepository

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/SENERGY-Platform/device-repository/v3/lib/client"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/service-commons/pkg/cache"
	"github.com/SENERGY-Platform/service-commons/pkg/signal"
)

// timeout bounds every request to the device-repository
const timeout = 10 * time.Second

// DeviceRepository reads through the device-repository client with the token of this service, so
// that a read continues the trace of the request that caused it; answers are cached for a minute
// and cleared by the cache invalidation signals.
type DeviceRepository struct {
	cache  *cache.Cache
	client client.Interface
	access *config.Access
}

func New(config config.Config, access *config.Access) (*DeviceRepository, error) {
	c, err := cache.New(cache.Config{
		CacheInvalidationSignalHooks: map[cache.Signal]cache.ToKey{
			signal.Known.CacheInvalidationAll:        nil,
			signal.Known.DeviceTypeCacheInvalidation: nil, //invalidate all, because services are cached and dependent on device-types
			signal.Known.AspectCacheInvalidation: func(signalValue string) (cacheKey string) {
				return "aspect-nodes." + signalValue
			},
		},
	})
	if err != nil {
		return nil, err
	}
	repo := &DeviceRepository{cache: c, access: access}
	repo.client = client.NewClient(config.DeviceRepositoryUrl, client.Options{Timeout: timeout, OptionalAuthTokenForApiGatewayRequest: repo.token})
	return repo, nil
}

func (this *DeviceRepository) token() (string, error) {
	token, err := this.access.Ensure()
	return string(token), err
}

func (this *DeviceRepository) GetProtocol(ctx context.Context, id string) (result model.Protocol, err error) {
	return cache.Use(this.cache, "protocol."+id, func() (model.Protocol, error) {
		token, err := this.token()
		if err != nil {
			return model.Protocol{}, err
		}
		protocol, err, _ := this.client.ReadProtocol(ctx, id, token)
		return protocol, err
	}, func(protocol model.Protocol) error {
		if protocol.Id == "" {
			return errors.New("invalid protocol loaded from cache")
		}
		return nil
	}, time.Minute)
}

func (this *DeviceRepository) GetDeviceType(ctx context.Context, id string) (result model.DeviceType, err error, code int) {
	code = http.StatusOK
	result, err = cache.Use(this.cache, "device-type."+id, func() (dt model.DeviceType, terr error) {
		token, terr := this.token()
		if terr != nil {
			code = http.StatusInternalServerError
			return dt, terr
		}
		dt, terr, code = this.client.ReadDeviceType(ctx, id, token)
		return dt, terr
	}, func(deviceType model.DeviceType) error {
		if deviceType.Id == "" {
			return errors.New("invalid device-type loaded from cache")
		}
		return nil
	}, time.Minute)
	return result, err, code
}

func (this *DeviceRepository) GetService(ctx context.Context, id string) (result model.Service, err error) {
	result, err, _ = this.GetServiceWithErrCode(ctx, id)
	return
}

func (this *DeviceRepository) GetServiceWithErrCode(ctx context.Context, id string) (result model.Service, err error, code int) {
	code = http.StatusOK
	result, err = cache.Use(this.cache, "service."+id, func() (service model.Service, terr error) {
		service, terr, code = this.client.GetService(ctx, id)
		return service, terr
	}, func(service model.Service) error {
		if service.Id == "" {
			return errors.New("invalid service loaded from cache")
		}
		return nil
	}, time.Minute)
	return result, err, code
}

func (this *DeviceRepository) GetAspectNode(ctx context.Context, id string) (result model.AspectNode, err error) {
	return cache.Use(this.cache, "aspect-nodes."+id, func() (aspect model.AspectNode, terr error) {
		aspect, terr, _ = this.client.GetAspectNode(ctx, id)
		return aspect, terr
	}, func(aspectNode model.AspectNode) error {
		if aspectNode.Id == "" {
			return errors.New("invalid aspect-node loaded from cache")
		}
		return nil
	}, time.Minute)
}
