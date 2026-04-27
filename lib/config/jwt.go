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

package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"net/url"
)

type Impersonate string

func (this Impersonate) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", string(this))
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err = client.Do(req)
	if err == nil && resp.StatusCode == 401 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		resp.Body.Close()
		slog.Error("http access denied response", "url", url, "status", resp.StatusCode, "error", buf.String())
		err = errors.New("access denied")
	}
	return
}

func (this Impersonate) PostJSON(url string, body interface{}, result interface{}) (err error) {
	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(body)
	if err != nil {
		return err
	}
	resp, err := this.Post(url, "application/json", b)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		temp, _ := io.ReadAll(resp.Body)
		err = errors.New(string(temp))
		return err
	}
	defer resp.Body.Close()
	if result != nil {
		err = json.NewDecoder(resp.Body).Decode(result)
	}
	return err
}

func (this Impersonate) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", string(this))
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err = client.Do(req)
	if err == nil && resp.StatusCode == 401 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		log.Println(buf.String())
		err = errors.New("access denied")
	}
	return
}

func (this Impersonate) GetJSON(url string, result interface{}) (err error) {
	resp, err := this.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		temp, _ := io.ReadAll(resp.Body)
		err = errors.New(string(temp))
		slog.Error("http error response", "url", url, "status", resp.StatusCode, "error", err)
		debug.PrintStack()
		return err
	}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		slog.Error("unable to decode http response", "url", url, "status", resp.StatusCode, "error", err)
		debug.PrintStack()
		return err
	}
	return err
}

type OpenidToken struct {
	AccessToken      string    `json:"access_token"`
	ExpiresIn        float64   `json:"expires_in"`
	RefreshExpiresIn float64   `json:"refresh_expires_in"`
	RefreshToken     string    `json:"refresh_token"`
	TokenType        string    `json:"token_type"`
	RequestTime      time.Time `json:"-"`
}

type Access struct {
	openid *OpenidToken
	config Config
}

func NewAccess(config Config) *Access {
	return &Access{config: config}
}

func (this *Access) Ensure() (token Impersonate, err error) {
	if this.openid == nil {
		this.openid = &OpenidToken{}
	}
	duration := TimeNow().Sub(this.openid.RequestTime).Seconds()

	if this.openid.AccessToken != "" && this.openid.ExpiresIn-this.config.AuthExpirationTimeBuffer > duration {
		token = Impersonate("Bearer " + this.openid.AccessToken)
		return
	}

	if this.openid.RefreshToken != "" && this.openid.RefreshExpiresIn-this.config.AuthExpirationTimeBuffer > duration {
		this.config.GetLogger().Info("refresh token", "expires_in", this.openid.RefreshExpiresIn, "duration", duration)
		err = refreshOpenidToken(this.openid, this.config)
		if err != nil {
			this.config.GetLogger().Warn("unable to use refreshtoken", "error", err)
		} else {
			token = Impersonate("Bearer " + this.openid.AccessToken)
			return
		}
	}

	this.config.GetLogger().Info("get new access token")
	err = getOpenidToken(this.openid, this.config)
	if err != nil {
		this.config.GetLogger().Warn("unable to get new access token", "error", err)
		this.openid = &OpenidToken{}
	}
	token = Impersonate("Bearer " + this.openid.AccessToken)
	return
}

func getOpenidToken(token *OpenidToken, config Config) (err error) {
	requesttime := TimeNow()
	resp, err := http.PostForm(config.AuthEndpoint+"/auth/realms/master/protocol/openid-connect/token", url.Values{
		"client_id":     {config.AuthClientId},
		"client_secret": {config.AuthClientSecret},
		"grant_type":    {"client_credentials"},
	})

	if err != nil {
		config.GetLogger().Error("ERROR: getOpenidToken::PostForm()", "error", err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		config.GetLogger().Error("ERROR: getOpenidToken()", "status", resp.StatusCode, "error", string(body))
		err = errors.New("access denied")
		resp.Body.Close()
		return
	}
	err = json.NewDecoder(resp.Body).Decode(token)
	token.RequestTime = requesttime
	return
}

func refreshOpenidToken(token *OpenidToken, config Config) (err error) {
	requesttime := TimeNow()
	resp, err := http.PostForm(config.AuthEndpoint+"/auth/realms/master/protocol/openid-connect/token", url.Values{
		"client_id":     {config.AuthClientId},
		"client_secret": {config.AuthClientSecret},
		"refresh_token": {token.RefreshToken},
		"grant_type":    {"refresh_token"},
	})

	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		config.GetLogger().Error("ERROR: refreshOpenidToken()", "status", resp.StatusCode, "error", string(body))
		err = errors.New("access denied")
		resp.Body.Close()
		return
	}
	err = json.NewDecoder(resp.Body).Decode(token)
	token.RequestTime = requesttime
	return
}
