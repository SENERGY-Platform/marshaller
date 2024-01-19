/*
 * Copyright 2024 InfAI (CC SES)
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

package metrics

import (
	"github.com/SENERGY-Platform/marshaller/lib/api/messages"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"
)

func NewMetrics() *Metrics {
	reg := prometheus.NewRegistry()

	result := &Metrics{
		httphandler: promhttp.HandlerFor(
			reg,
			promhttp.HandlerOpts{
				Registry: reg,
			},
		),
		MarshallingRequestsSummary: prometheus.NewSummary(prometheus.SummaryOpts{
			Name: "marshaller_marshalling_request_handling_duration_summary",
			Help: "summary for handling duration (in μs) of marshalling request",
		}),
		MarshallingRequests: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "marshaller_marshalling_request_handling_duration",
			Help:    "histogram vec for handling duration (in μs) of marshalling request",
			Buckets: []float64{500, 600, 700, 800, 900, 1000, 2000, 3000, 4000, 5000, 10000, 50000, 100000, 1000000},
		}, []string{"call_source", "endpoint", "service_id", "function_ids"}),

		UnmarshallingRequestsSummary: prometheus.NewSummary(prometheus.SummaryOpts{
			Name: "marshaller_unmarshalling_request_handling_duration_summary",
			Help: "summary for handling duration (in μs) of unmarshalling request",
		}),
		UnmarshallingRequests: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "marshaller_unmarshalling_request_handling_duration",
			Help:    "histogram vec for handling duration (in μs) of unmarshalling request",
			Buckets: []float64{500, 600, 700, 800, 900, 1000, 2000, 3000, 4000, 5000, 10000, 50000, 100000, 1000000},
		}, []string{"call_source", "endpoint", "service_id", "function_ids"}),
	}

	reg.MustRegister(
		result.MarshallingRequestsSummary,
		result.MarshallingRequests,

		result.UnmarshallingRequestsSummary,
		result.UnmarshallingRequests,
	)

	return result
}

type Metrics struct {
	httphandler http.Handler

	MarshallingRequestsSummary prometheus.Summary
	MarshallingRequests        *prometheus.HistogramVec

	UnmarshallingRequestsSummary prometheus.Summary
	UnmarshallingRequests        *prometheus.HistogramVec
}

func (this *Metrics) LogMarshallingRequest(request *http.Request, endpoint string, msg messages.MarshallingV2Request, duration time.Duration) {
	if this == nil {
		return
	}
	dur := float64(duration.Microseconds())

	this.MarshallingRequestsSummary.Observe(dur)

	functionIds := []string{}
	for _, data := range msg.Data {
		functionIds = append(functionIds, data.FunctionId)
	}
	sort.Strings(functionIds)

	this.MarshallingRequests.WithLabelValues(getCallSource(request), endpoint, msg.Service.Id, strings.Join(functionIds, ",")).Observe(dur)
}

func (this *Metrics) LogUnmarshallingRequest(request *http.Request, endpoint string, msg messages.UnmarshallingV2Request, duration time.Duration) {
	if this == nil {
		return
	}
	dur := float64(duration.Microseconds())
	this.MarshallingRequestsSummary.Observe(dur)
	this.MarshallingRequests.WithLabelValues(getCallSource(request), endpoint, msg.Service.Id, msg.FunctionId).Observe(dur)
}

func getCallSource(req *http.Request) (result string) {
	result = req.RemoteAddr
	remoteAddr, _, err := net.SplitHostPort(result)
	if err == nil && remoteAddr != "" {
		remoteHosts, _ := net.LookupAddr(remoteAddr)
		if len(remoteHosts) > 0 {
			sort.Strings(remoteHosts)
			result = remoteHosts[0]
		}
	}
	return result
}
