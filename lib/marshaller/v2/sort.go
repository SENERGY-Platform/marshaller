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
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"math"
	"reflect"
	"sort"
	"strings"
)

type pathAspectInfo struct {
	path     string
	aspect   string
	distance int
}

type DeviceRepository interface {
	GetAspectNode(id string) (model.AspectNode, error)
}

func (this *Marshaller) SortPathsByAspectDistance(repo DeviceRepository, service model.Service, aspect *model.AspectNode, paths []string) (result []string, err error) {
	if aspect == nil {
		return paths, nil
	}
	distances, err := getAspectDistances(repo, *aspect)
	if err != nil {
		return result, err
	}
	infoList := []pathAspectInfo{}
	for _, path := range paths {
		info := this.getOutputPathAspectInfo(service, path)
		var ok bool
		info.distance, ok = distances[info.aspect]
		if !ok {
			info.distance = math.MaxInt
		}
		infoList = append(infoList, info)
	}
	sort.Slice(infoList, func(i, j int) bool {
		return infoList[i].distance < infoList[j].distance
	})
	for _, info := range infoList {
		result = append(result, info.path)
	}
	return result, nil
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

func (this *Marshaller) getOutputPathAspectInfo(service model.Service, path string) (result pathAspectInfo) {
	result.path = path
	result.aspect = this.getPathAspect(service.Outputs, strings.Split(path, "."))
	return result
}

func (this *Marshaller) getPathAspect(contents []model.Content, path []string) (result string) {
	for _, c := range contents {
		result, ok := this.getPathAspectFromContentVariable(c.ContentVariable, path, []string{})
		if ok {
			return result
		}
	}
	return result
}

func (this *Marshaller) getPathAspectFromContentVariable(variable model.ContentVariable, targetPath []string, currentPath []string) (result string, ok bool) {
	currentPath = append(currentPath, variable.Name)
	if len(currentPath) > len(targetPath) || !reflect.DeepEqual(currentPath, targetPath[:len(currentPath)]) {
		return "", false
	}
	if len(currentPath) == len(targetPath) {
		return variable.AspectId, true
	}
	for _, sub := range variable.SubContentVariables {
		result, ok = this.getPathAspectFromContentVariable(sub, targetPath, currentPath)
		if ok {
			return result, ok
		}
	}
	return "", false
}
