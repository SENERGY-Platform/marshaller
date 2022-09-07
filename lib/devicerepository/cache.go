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
	"errors"
	"github.com/patrickmn/go-cache"
	"log"
	"time"
)

var CacheExpiration = 60      // 60sec
var L1Size = 20 * 1024 * 1024 //20MB
var Debug = false

type Cache struct {
	l1 *cache.Cache
}

type Item struct {
	Key   string
	Value []byte
}

var ErrNotFound = errors.New("key not found in cache")

func NewCache() *Cache {
	return &Cache{l1: cache.New(time.Duration(CacheExpiration)*time.Second, time.Duration(CacheExpiration)*time.Second)}
}

func (this *Cache) get(key string) (value []byte, err error) {
	temp, found := this.l1.Get(key)
	if !found {
		err = ErrNotFound
	} else {
		var ok bool
		value, ok = temp.([]byte)
		if !ok {
			err = errors.New("unable to interprete cache result")
		}
	}
	return
}

func (this *Cache) set(key string, value []byte) {
	this.l1.Set(key, value, time.Duration(CacheExpiration)*time.Second)
	return
}
func (this *Cache) Use(key string, getter func() (interface{}, error), result interface{}) (err error) {
	value, err := this.get(key)
	if err == nil {
		err = json.Unmarshal(value, result)
		return
	} else if err != ErrNotFound {
		log.Println("WARNING: err in LocalCache::l1.Get()", err)
	}
	temp, err := getter()
	if err != nil {
		return err
	}
	value, err = json.Marshal(temp)
	if err != nil {
		return err
	}
	this.set(key, value)
	return json.Unmarshal(value, &result)
}
