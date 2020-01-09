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
	"encoding/json"
	"github.com/SENERGY-Platform/marshaller-service/lib/config"
	"github.com/SENERGY-Platform/marshaller-service/lib/marshaller/model"
	"log"
	"net/url"
	"runtime/debug"
)

type DeviceRepository struct {
	cache         *Cache
	repoUrl       string
	permsearchUrl string
}

func New(config config.Config) *DeviceRepository {
	return &DeviceRepository{repoUrl: config.DeviceRepositoryUrl, cache: NewCache(), permsearchUrl: config.PermissionsSearchUrl}
}

func (this *DeviceRepository) GetProtocol(token config.Impersonate, id string) (result model.Protocol, err error) {
	result, err = this.getProtocolFromCache(id)
	if err != nil {
		err = token.GetJSON(this.repoUrl+"/protocols/"+url.QueryEscape(id), &result)
		if err == nil {
			this.saveProtocolToCache(result)
		}
	}
	return
}

func (this *DeviceRepository) GetService(token config.Impersonate, id string) (result model.Service, err error) {
	result, err = this.getServiceFromCache(id)
	if err != nil {
		err = token.GetJSON(this.repoUrl+"/services/"+url.QueryEscape(id), &result)
		if err != nil {
			log.Println("ERROR:", err)
			debug.PrintStack()
			return result, err
		}
		this.saveServiceToCache(result)
	}
	return
}

func (this *DeviceRepository) getServiceFromCache(id string) (service model.Service, err error) {
	item, err := this.cache.Get("service." + id)
	if err != nil {
		return service, err
	}
	err = json.Unmarshal(item.Value, &service)
	return
}

func (this *DeviceRepository) saveServiceToCache(service model.Service) {
	buffer, _ := json.Marshal(service)
	this.cache.Set("service."+service.Id, buffer)
}

func (this *DeviceRepository) saveProtocolToCache(protocol model.Protocol) {
	buffer, _ := json.Marshal(protocol)
	this.cache.Set("protocol."+protocol.Id, buffer)
}

func (this *DeviceRepository) getProtocolFromCache(id string) (protocol model.Protocol, err error) {
	item, err := this.cache.Get("protocol." + id)
	if err != nil {
		return protocol, err
	}
	err = json.Unmarshal(item.Value, &protocol)
	return
}
