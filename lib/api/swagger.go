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

package api

import (
	"net/http"
	"strings"

	_ "github.com/SENERGY-Platform/marshaller/docs"
	"github.com/SENERGY-Platform/marshaller/lib/api/metrics"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"
)

func init() {
	endpoints = append(endpoints, &Swagger{})
}

type Swagger struct{}

func (this *Swagger) Swagger(config config.Config, router *http.ServeMux, ctrl Controller, m *metrics.Metrics) {
	if config.EnableSwaggerUi {
		router.HandleFunc("GET /swagger/{pathname...}", func(res http.ResponseWriter, req *http.Request) {
			httpSwagger.WrapHandler(res, req)
		})
	}

	router.HandleFunc("GET /doc", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		doc, err := swag.ReadDoc("marshaller")
		if err != nil {
			config.GetLogger().Error("unable to read swagger doc", "error", err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		//remove empty host to enable developer-swagger-api service to replace it; can not use cleaner delete on json object, because developer-swagger-api is sensible to formatting; better alternative is refactoring of developer-swagger-api/apis/db/db.py
		doc = strings.Replace(doc, `"host": "",`, "", 1)
		_, _ = writer.Write([]byte(doc))
	})
}
