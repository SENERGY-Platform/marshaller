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

package model

import "slices"

//ContentVariable.AspectId is deprecated in favor of ContentVariable.AspectIds, and the
//aspect fields of the marshaller requests follow it. A deprecated single value is an alias
//for a list with one element and is folded into that list at the request boundary, so that
//everything behind it evaluates the list only.

// ContentVariableAspectIds returns the aspects of a content variable. The deprecated
// AspectId is only used if AspectIds is empty: a device-repository read fills both fields,
// but a device-type written before the aspect lists carries just AspectId.
func ContentVariableAspectIds(variable ContentVariable) []string {
	if len(variable.AspectIds) > 0 {
		return variable.AspectIds
	}
	if variable.AspectId != "" {
		return []string{variable.AspectId}
	}
	return nil
}

// AspectIdsAlias folds a deprecated single aspect id into a list of aspect ids.
func AspectIdsAlias(aspectId string, aspectIds []string) []string {
	if aspectId != "" && !slices.Contains(aspectIds, aspectId) {
		return append(aspectIds, aspectId)
	}
	return aspectIds
}

// AspectNodesAlias folds a deprecated single aspect node into a list of aspect nodes. An
// unset node means no aspect, the way an empty aspect id means none.
func AspectNodesAlias(aspectNode AspectNode, aspectNodes []AspectNode) []AspectNode {
	if aspectNode.Id == "" || ContainsAspectNode(aspectNodes, aspectNode.Id) {
		return aspectNodes
	}
	return append(aspectNodes, aspectNode)
}

func ContainsAspectNode(aspectNodes []AspectNode, aspectId string) bool {
	return slices.ContainsFunc(aspectNodes, func(node AspectNode) bool {
		return node.Id == aspectId
	})
}

// AspectMatchLevel reports how closely the aspects of a content variable match the queried
// aspect nodes, as a distance in the aspect tree: 0 for the node itself, 1 for one of its
// children, 2 for a deeper descendant, -1 for no match. Every queried node has to be
// matched, because a query naming more than one aspect is an AND, like the device-type
// criteria filter of the device-repository. The level of the worst matched node is
// returned, so the value stays a tree distance: every queried aspect lies within it.
// Without a queried node everything matches at level 0.
func AspectMatchLevel(aspectIds []string, aspectNodes []AspectNode) int {
	level := 0
	for _, node := range aspectNodes {
		nodeLevel := aspectNodeMatchLevel(aspectIds, node)
		if nodeLevel < 0 {
			return -1
		}
		if nodeLevel > level {
			level = nodeLevel
		}
	}
	return level
}

// aspectNodeMatchLevel returns the smallest distance between one queried aspect node and
// the aspects of a content variable, because a variable carrying several aspects is as
// close to the query as its closest one.
func aspectNodeMatchLevel(aspectIds []string, node AspectNode) int {
	level := -1
	for _, aspectId := range aspectIds {
		candidate := -1
		switch {
		case aspectId == node.Id:
			return 0
		case slices.Contains(node.ChildIds, aspectId):
			candidate = 1
		case slices.Contains(node.DescendentIds, aspectId):
			candidate = 2
		}
		if candidate > -1 && (level < 0 || candidate < level) {
			level = candidate
		}
	}
	return level
}
