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

package testdata

import (
	_ "embed"
	"encoding/json"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
)

//go:embed functions.json
var FunctionsJson []byte

//go:embed concepts.json
var ConceptsJson []byte

//go:embed characteristics.json
var CharacteristicsJson []byte

//go:embed aspectnodes.json
var AspectNodesJson []byte

func GetFunctions() (result []model.Function, err error) {
	err = json.Unmarshal(FunctionsJson, &result)
	return
}

func GetConcepts() (result []model.Concept, err error) {
	err = json.Unmarshal(ConceptsJson, &result)
	return
}

func GetCharacteristics() (result []model.Characteristic, err error) {
	err = json.Unmarshal(CharacteristicsJson, &result)
	return
}

func GetAspectNodes() (result []model.AspectNode, err error) {
	err = json.Unmarshal(AspectNodesJson, &result)
	return
}
