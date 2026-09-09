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

package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	"github.com/SENERGY-Platform/marshaller/lib/api/metrics"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
)

func init() {
	endpoints = append(endpoints, &PathOptions{})
}

type PathOptions struct{}

// GetPathOptions godoc
// @Summary      list the paths of device-types that match a function and aspects
// @Description  returns the matching paths per device-type. Naming more than one aspect is an AND on one content variable, and each aspect covers its own subtree.
// @Tags         path-options
// @Produce      json
// @Security     Bearer
// @Param        device-type-ids query string true "comma separated device-type ids; an absent parameter yields an empty result"
// @Param        characteristic-filter query string true "comma separated characteristic ids; an absent parameter yields an empty result"
// @Param        function-id query string false "id of the function the path has to serve"
// @Param        aspect-ids query string false "comma separated aspect ids"
// @Param        aspect-id query string false "deprecated: use aspect-ids; an alias for a list with one element"
// @Param        without-envelope query bool false "omit the value. prefix from the returned paths"
// @Success      200 {object} map[string][]marshaller.PathOptionsResultElement
// @Failure      400 {string} string "without-envelope is not a bool"
// @Failure      500 {string} string
// @Router       /path-options [GET]
func (this *PathOptions) GetPathOptions(config config.Config, router *http.ServeMux, ctrl Controller, m *metrics.Metrics) {
	router.HandleFunc("GET /path-options", func(writer http.ResponseWriter, request *http.Request) {
		//an absent device-type or characteristic filter is an empty query here and is
		//answered with an empty result, unlike the request-body form where an absent
		//characteristic filter means "do not filter"
		if request.URL.Query().Get("device-type-ids") == "" || request.URL.Query().Get("characteristic-filter") == "" {
			writeJson(config, writer, map[string]interface{}{})
			return
		}
		query := messages.PathOptionsQuery{
			DeviceTypeIds:          splitList(request.URL.Query().Get("device-type-ids")),
			CharacteristicIdFilter: splitList(request.URL.Query().Get("characteristic-filter")),
			FunctionId:             strings.TrimSpace(request.URL.Query().Get("function-id")),
			//aspect-id is deprecated in favor of aspect-ids and is an alias for a list with one element
			AspectIds: model.AspectIdsAlias(strings.TrimSpace(request.URL.Query().Get("aspect-id")), splitList(request.URL.Query().Get("aspect-ids"))),
		}
		withoutEnvelopeStr := strings.TrimSpace(request.URL.Query().Get("without-envelope"))
		if withoutEnvelopeStr != "" {
			withoutEnvelope, err := strconv.ParseBool(withoutEnvelopeStr)
			if err != nil {
				http.Error(writer, "expect bool in without-envelope: "+err.Error(), http.StatusBadRequest)
				return
			}
			query.WithoutEnvelope = withoutEnvelope
		}
		result, err, code := ctrl.GetPathOptions(query)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writeJson(config, writer, result)
	})
}

// QueryPathOptions godoc
// @Summary      list the paths of device-types that match a function and aspects
// @Description  the request-body form of GET /path-options. An empty characteristic_id_filter means "do not filter" here, unlike the absent query parameter.
// @Tags         path-options
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        query body messages.PathOptionsQuery true "the device-types, function and aspects to match"
// @Success      200 {object} map[string][]marshaller.PathOptionsResultElement
// @Failure      400 {string} string "the body is not valid json"
// @Failure      500 {string} string
// @Router       /query/path-options [POST]
func (this *PathOptions) QueryPathOptions(config config.Config, router *http.ServeMux, ctrl Controller, m *metrics.Metrics) {
	router.HandleFunc("POST /query/path-options", func(writer http.ResponseWriter, request *http.Request) {
		query, ok := decodeBody[messages.PathOptionsQuery](writer, request)
		if !ok {
			return
		}
		result, err, code := ctrl.GetPathOptions(query)
		if err != nil {
			http.Error(writer, err.Error(), code)
			return
		}
		writeJson(config, writer, result)
	})
}

// splitList reads a comma separated query parameter. An empty parameter is no filter
// rather than a filter for the empty string.
func splitList(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return strings.Split(strings.ReplaceAll(value, " ", ""), ",")
}
