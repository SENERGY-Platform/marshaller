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

package lib

import (
	"context"
	"github.com/SENERGY-Platform/marshaller/lib/api"
	"github.com/SENERGY-Platform/marshaller/lib/conceptrepo"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/configurables"
	"github.com/SENERGY-Platform/marshaller/lib/converter"
	"github.com/SENERGY-Platform/marshaller/lib/devicerepository"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	v2 "github.com/SENERGY-Platform/marshaller/lib/marshaller/v2"
	"github.com/SENERGY-Platform/service-commons/pkg/cache/invalidator"
	"github.com/SENERGY-Platform/service-commons/pkg/kafka"
	"log"
	"runtime/debug"
	"time"
)

func Start(ctx context.Context, conf config.Config) (closed context.Context, err error) {
	childCtx, cancel := context.WithCancel(ctx)
	access := config.NewAccess(conf)
	conceptRepo, err := conceptrepo.New(
		childCtx,
		conf,
		access,
		conceptrepo.ConceptRepoDefault{
			Concept: model.NullConcept,
			Characteristics: []model.Characteristic{
				model.NullCharacteristic,
			},
		},
	)
	if err != nil {
		cancel()
		return nil, err
	}

	devicerepo, err := devicerepository.New(conf, access)
	if err != nil {
		cancel()
		return nil, err
	}
	converter := converter.New(conf, access)
	marshaller := marshaller.New(converter, conceptRepo, devicerepo)
	configurableService := configurables.New(conceptRepo)

	marshallerV2 := v2.New(conf, converter, conceptRepo)

	closed = api.Start(childCtx, conf, marshaller, marshallerV2, configurableService, devicerepo, converter)
	go func() {
		<-closed.Done()
		cancel()
	}()
	return closed, nil
}

func StartCacheInvalidator(ctx context.Context, conf config.Config) error {
	if conf.KafkaUrl == "" || conf.KafkaUrl == "-" {
		return nil
	}
	return invalidator.StartCacheInvalidatorAll(ctx, kafka.Config{
		KafkaUrl:               conf.KafkaUrl,
		StartOffset:            kafka.LastOffset,
		Debug:                  conf.Debug,
		PartitionWatchInterval: time.Minute,
		InitTopic:              conf.InitTopics,
		OnError: func(err error) {
			log.Println("ERROR:", err)
			debug.PrintStack()
		},
	}, conf.CacheInvalidationKafkaTopics, nil)
}
