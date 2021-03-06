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

package conceptrepo

import (
	"context"
	"errors"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

type ConceptRepo struct {
	config config.Config
	access *config.Access

	defaults                           []ConceptRepoDefault
	concepts                           map[string]model.Concept
	characteristics                    map[string]model.Characteristic
	conceptByCharacteristic            map[string]model.Concept
	rootCharacteristicByCharacteristic map[string]model.Characteristic

	characteristicsOfFunction map[string][]string

	mux sync.Mutex
}

type ConceptRepoDefault struct {
	Concept         model.Concept
	Characteristics []model.Characteristic
}

func New(ctx context.Context, conf config.Config, access *config.Access, defaults ...ConceptRepoDefault) (result *ConceptRepo, err error) {
	result = &ConceptRepo{
		config:                             conf,
		access:                             access,
		defaults:                           defaults,
		concepts:                           map[string]model.Concept{},
		characteristics:                    map[string]model.Characteristic{},
		conceptByCharacteristic:            map[string]model.Concept{},
		rootCharacteristicByCharacteristic: map[string]model.Characteristic{},
		characteristicsOfFunction:          map[string][]string{},
	}
	err = result.Load()
	if err != nil {
		return result, err
	}
	ticker := time.NewTicker(time.Duration(conf.ConceptRepoRefreshInterval) * time.Second)
	go func() {
		defer ticker.Stop()
		<-ctx.Done()
	}()
	go func() {
		for range ticker.C {
			err = result.Load()
			if err != nil {
				log.Println("WARNING: unable to update concept repository", err)
			}
		}
	}()
	return result, nil
}

func (this *ConceptRepo) GetCharacteristicsOfFunction(functionId string) (characteristicIds []string, err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	var ok bool
	characteristicIds, ok = this.characteristicsOfFunction[functionId]
	if !ok {
		err = errors.New("unknown function-id")
	}
	return
}

func (this *ConceptRepo) GetConcept(id string) (concept model.Concept, err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	concept, ok := this.concepts[id]
	if !ok {
		debug.PrintStack()
		return concept, errors.New("no concept found for id " + id)
	}
	return concept, nil
}

func getCharacteristicDescendents(characteristic model.Characteristic) (result []model.Characteristic) {
	result = []model.Characteristic{characteristic}
	for _, child := range characteristic.SubCharacteristics {
		result = append(result, getCharacteristicDescendents(child)...)
	}
	return result
}

func (this *ConceptRepo) GetConceptOfCharacteristic(characteristicId string) (conceptId string, err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	concept, ok := this.conceptByCharacteristic[this.rootCharacteristicByCharacteristic[characteristicId].Id]
	if !ok {
		debug.PrintStack()
		return conceptId, errors.New("no concept found for characteristic id " + characteristicId)
	}
	return concept.Id, nil
}

func (this *ConceptRepo) GetCharacteristic(id string) (characteristic model.Characteristic, err error) {
	if id == "" {
		return model.NullCharacteristic, nil
	}
	this.mux.Lock()
	defer this.mux.Unlock()
	characteristic, ok := this.characteristics[id]
	if !ok {
		debug.PrintStack()
		return characteristic, errors.New("no characteristic found for id " + id)
	}
	return characteristic, nil
}

func (this *ConceptRepo) GetRootCharacteristics(ids []string) (result []string) {
	this.mux.Lock()
	defer this.mux.Unlock()
	for _, id := range ids {
		root, ok := this.rootCharacteristicByCharacteristic[id]
		if ok {
			result = append(result, root.Id)
		}
	}
	return
}

func (this *ConceptRepo) registerFunction(f FunctionInfo) {
	concept, ok := this.concepts[f.ConceptId]
	if !ok {
		log.Println("WARNING: unable to register function with unknown concept", f)
		return
	}
	this.characteristicsOfFunction[f.Id] = concept.CharacteristicIds
}
