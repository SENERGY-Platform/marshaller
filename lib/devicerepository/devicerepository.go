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
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/service-commons/pkg/cache"
	"github.com/SENERGY-Platform/service-commons/pkg/signal"
	"net/http"
	"net/url"
	"runtime/debug"
	"time"
)

type DeviceRepository struct {
	cache         *cache.Cache
	repoUrl       string
	permsearchUrl string
	access        *config.Access
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
	return &DeviceRepository{repoUrl: config.DeviceRepositoryUrl, cache: c, permsearchUrl: config.PermissionsSearchUrl, access: access}, nil
}

func (this *DeviceRepository) GetProtocol(id string) (result model.Protocol, err error) {
	return cache.Use(this.cache, "protocol."+id, func() (model.Protocol, error) {
		return this.getProtocol(id)
	}, time.Minute)
}

func (this *DeviceRepository) getProtocol(id string) (result model.Protocol, err error) {
	token, err := this.access.Ensure()
	if err != nil {
		return result, err
	}
	err = token.GetJSON(this.repoUrl+"/protocols/"+url.QueryEscape(id), &result)
	return
}

func (this *DeviceRepository) GetDeviceType(id string) (result model.DeviceType, err error, code int) {
	code = http.StatusOK
	result, err = cache.Use(this.cache, "device-type."+id, func() (dt model.DeviceType, terr error) {
		dt, terr, code = this.getDeviceType(id)
		return dt, terr
	}, time.Minute)
	return result, err, code
}

func (this *DeviceRepository) getDeviceType(id string) (result model.DeviceType, err error, code int) {
	token, err := this.access.Ensure()
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	req, err := http.NewRequest("GET", this.repoUrl+"/device-types/"+url.PathEscape(id), nil)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", string(token))
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	if resp.StatusCode == http.StatusNotFound {
		return result, errors.New("device-type not found"), resp.StatusCode
	}
	if resp.StatusCode >= 300 {
		return result, errors.New("unexpected status code"), resp.StatusCode
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	return
}

func (this *DeviceRepository) GetService(id string) (result model.Service, err error) {
	result, err, _ = this.GetServiceWithErrCode(id)
	return
}

func (this *DeviceRepository) GetServiceWithErrCode(id string) (result model.Service, err error, code int) {
	code = http.StatusOK
	result, err = cache.Use(this.cache, "service."+id, func() (service model.Service, terr error) {
		service, terr, code = this.getServiceWithErrCode(id)
		return service, terr
	}, time.Minute)
	return result, err, code
}

func (this *DeviceRepository) getServiceWithErrCode(id string) (result model.Service, err error, code int) {
	token, err := this.access.Ensure()
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	req, err := http.NewRequest("GET", this.repoUrl+"/services/"+url.PathEscape(id), nil)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", string(token))
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	if resp.StatusCode == http.StatusNotFound {
		return result, errors.New("service not found"), resp.StatusCode
	}
	if resp.StatusCode >= 300 {
		return result, errors.New("unexpected status code"), resp.StatusCode
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	return
}

func (this *DeviceRepository) GetAspectNode(id string) (result model.AspectNode, err error) {
	result, err = cache.Use(this.cache, "aspect-nodes."+id, func() (aspect model.AspectNode, terr error) {
		token, err := this.access.Ensure()
		if err != nil {
			return aspect, err
		}
		req, err := http.NewRequest("GET", this.repoUrl+"/aspect-nodes/"+url.PathEscape(id), nil)
		if err != nil {
			debug.PrintStack()
			return aspect, err
		}
		req.Header.Set("Authorization", string(token))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			debug.PrintStack()
			return aspect, err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			debug.PrintStack()
			return aspect, errors.New(buf.String())
		}
		err = json.NewDecoder(resp.Body).Decode(&aspect)
		return aspect, err
	}, time.Minute)
	return result, err
}
