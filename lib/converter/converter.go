/*
 * Copyright 2020 InfAI (CC SES)
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
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	"github.com/SENERGY-Platform/models/go/models"
	"net/url"
)

type Converter struct {
	config config.Config
	access *config.Access
}

func New(config config.Config, access *config.Access) *Converter {
	return &Converter{config: config, access: access}
}

func (this *Converter) Cast(in interface{}, from marshaller.CharacteristicId, to marshaller.CharacteristicId) (out interface{}, err error) {
	if from == to {
		return in, nil
	}
	token, err := this.access.Ensure()
	if err != nil {
		return out, err
	}
	err = token.PostJSON(this.config.ConverterUrl+"/conversions/"+url.PathEscape(from)+"/"+url.PathEscape(to), in, &out)
	return out, err
}

func (this *Converter) CastWithExtension(in interface{}, from marshaller.CharacteristicId, to marshaller.CharacteristicId, extensions []models.ConverterExtension) (out interface{}, err error) {
	if from == to {
		return in, nil
	}
	token, err := this.access.Ensure()
	if err != nil {
		return out, err
	}
	err = token.PostJSON(this.config.ConverterUrl+"/extended-conversions/"+url.PathEscape(from)+"/"+url.PathEscape(to), map[string]interface{}{
		"input":      in,
		"extensions": extensions,
	}, &out)
	return out, err
}

func (this *Converter) TryExtension(call ExtensionCall) (resp ExtensionCallResponse, err error) {
	token, err := this.access.Ensure()
	if err != nil {
		return resp, err
	}
	err = token.PostJSON(this.config.ConverterUrl+"/extension-call", call, &resp)
	return resp, err
}

type ExtensionCall struct {
	Input     interface{}               `json:"input"`
	Extension models.ConverterExtension `json:"extension"`
}

type ExtensionCallResponse struct {
	Output interface{} `json:"output"`
	Error  string      `json:"error,omitempty"`
}
