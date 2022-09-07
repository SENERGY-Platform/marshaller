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

package devicerepository

import (
	"errors"
	"testing"
)

func TestCache(t *testing.T) {
	cache := NewCache()
	temp := ""
	err := cache.Use("foo", func() (interface{}, error) {
		return "bar", nil
	}, &temp)
	if err != nil {
		t.Error(err)
		return
	}
	if temp != "bar" {
		t.Error(temp)
		return
	}

	temp = ""
	err = cache.Use("foo", func() (interface{}, error) {
		return "bar", nil
	}, &temp)
	if err != nil {
		t.Error(err)
		return
	}
	if temp != "bar" {
		t.Error(temp)
		return
	}

	temp = ""
	err = cache.Use("batz", func() (interface{}, error) {
		return "", errors.New("should fail")
	}, &temp)
	if err == nil {
		t.Error(err, "should have failed")
		return
	}

}
