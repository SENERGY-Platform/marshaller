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

package v2

import (
	"fmt"
	"log/slog"
	"math"
	"reflect"
	"sort"
	"strings"

	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
)

type PathAspectInfo struct {
	Path     string
	Aspects  []string
	Distance int
}

type DeviceRepository interface {
	GetAspectNode(id string) (model.AspectNode, error)
}

func (this *Marshaller) SortPathsByAspectDistance(repo DeviceRepository, service model.Service, aspects []model.AspectNode, paths []string) (result []string, err error) {
	if len(aspects) == 0 {
		return paths, nil
	}
	distancesPerAspect := []map[string]int{}
	for _, aspect := range aspects {
		distances, err := getAspectDistances(repo, aspect)
		if err != nil {
			return result, err
		}
		distancesPerAspect = append(distancesPerAspect, distances)
	}
	infoList := []PathAspectInfo{}
	for _, path := range paths {
		info := this.getOutputPathAspectInfo(service, path)
		info.Distance = pathAspectDistance(info.Aspects, distancesPerAspect)
		infoList = append(infoList, info)
	}
	sort.Slice(infoList, func(i, j int) bool {
		return infoList[i].Distance < infoList[j].Distance
	})
	for _, info := range infoList {
		result = append(result, info.Path)
	}
	slog.Debug("path aspect distance info", "aspects", fmt.Sprintf("%#v", aspects), "infoList", fmt.Sprintf("%#v", infoList))
	return result, nil
}

// pathAspectDistance aggregates the aspect distances of a path the way
// model.AspectMatchLevel aggregates the match levels: the smallest distance over the
// aspects of the path per queried aspect, then the largest of those, so the value stays a
// tree distance covering every queried aspect. A queried aspect that reaches none of the
// path aspects puts the path last, the way an unknown aspect always did.
func pathAspectDistance(pathAspects []string, distancesPerAspect []map[string]int) int {
	result := 0
	for _, distances := range distancesPerAspect {
		closest := math.MaxInt
		for _, pathAspect := range pathAspects {
			if distance, ok := distances[pathAspect]; ok && distance < closest {
				closest = distance
			}
		}
		if closest == math.MaxInt {
			return math.MaxInt
		}
		if closest > result {
			result = closest
		}
	}
	return result
}

func getAspectDistances(repo DeviceRepository, aspect model.AspectNode) (result map[string]int, err error) {
	result = map[string]int{
		aspect.Id: 0,
	}
	for _, child := range aspect.ChildIds {
		result[child] = 1
	}
	if len(aspect.ChildIds) != len(aspect.DescendentIds) {
		for _, child := range aspect.ChildIds {
			childAspect, err := repo.GetAspectNode(child)
			if err != nil {
				return result, err
			}
			temp, err := getAspectDistances(repo, childAspect)
			if err != nil {
				return result, err
			}
			for id, distance := range temp {
				result[id] = distance + 1
			}
		}
	}
	return result, nil
}

func (this *Marshaller) getOutputPathAspectInfo(service model.Service, path string) (result PathAspectInfo) {
	result.Path = path
	result.Aspects = this.getPathAspects(service.Outputs, strings.Split(path, "."))
	return result
}

func (this *Marshaller) getPathAspects(contents []model.Content, path []string) (result []string) {
	for _, c := range contents {
		result, ok := this.getPathAspectsFromContentVariable(c.ContentVariable, path, []string{})
		if ok {
			return result
		}
	}
	return result
}

func (this *Marshaller) getPathAspectsFromContentVariable(variable model.ContentVariable, targetPath []string, currentPath []string) (result []string, ok bool) {
	currentPath = append(currentPath, variable.Name)
	if len(currentPath) > len(targetPath) || !reflect.DeepEqual(currentPath, targetPath[:len(currentPath)]) {
		return nil, false
	}
	if len(currentPath) == len(targetPath) {
		return model.ContentVariableAspectIds(variable), true
	}
	for _, sub := range variable.SubContentVariables {
		result, ok = this.getPathAspectsFromContentVariable(sub, targetPath, currentPath)
		if ok {
			return result, ok
		}
	}
	return nil, false
}
