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

import (
	"reflect"
	"testing"
)

var testAspectNodes = map[string]AspectNode{
	"air": {
		Id:            "air",
		ChildIds:      []string{"inside_air", "outside_air"},
		DescendentIds: []string{"inside_air", "outside_air"},
	},
	"inside_air": {
		Id:          "inside_air",
		ParentId:    "air",
		AncestorIds: []string{"air"},
	},
	"outside_air": {
		Id:          "outside_air",
		ParentId:    "air",
		AncestorIds: []string{"air"},
	},
	"electricity": {
		Id:            "electricity",
		ChildIds:      []string{"consumption"},
		DescendentIds: []string{"consumption", "today"},
	},
	"humidity": {
		Id: "humidity",
	},
}

func testNodes(ids ...string) (result []AspectNode) {
	for _, id := range ids {
		node, ok := testAspectNodes[id]
		if !ok {
			panic("unknown test aspect node " + id)
		}
		result = append(result, node)
	}
	return result
}

func TestContentVariableAspectIds(t *testing.T) {
	t.Run("returns the aspect list if it is set", func(t *testing.T) {
		actual := ContentVariableAspectIds(ContentVariable{AspectIds: []string{"air", "humidity"}})
		if !reflect.DeepEqual(actual, []string{"air", "humidity"}) {
			t.Error(actual)
		}
	})
	t.Run("reads the deprecated aspect id as a list with one element", func(t *testing.T) {
		actual := ContentVariableAspectIds(ContentVariable{AspectId: "air"})
		if !reflect.DeepEqual(actual, []string{"air"}) {
			t.Error(actual)
		}
	})
	t.Run("prefers the aspect list over the deprecated aspect id", func(t *testing.T) {
		actual := ContentVariableAspectIds(ContentVariable{AspectId: "air", AspectIds: []string{"humidity"}})
		if !reflect.DeepEqual(actual, []string{"humidity"}) {
			t.Error(actual)
		}
	})
	t.Run("returns nothing for a content variable without an aspect", func(t *testing.T) {
		actual := ContentVariableAspectIds(ContentVariable{})
		if len(actual) != 0 {
			t.Error(actual)
		}
	})
}

func TestAspectIdsAlias(t *testing.T) {
	t.Run("folds the deprecated aspect id into an empty list", func(t *testing.T) {
		actual := AspectIdsAlias("air", nil)
		if !reflect.DeepEqual(actual, []string{"air"}) {
			t.Error(actual)
		}
	})
	t.Run("appends the deprecated aspect id to a list that misses it", func(t *testing.T) {
		actual := AspectIdsAlias("humidity", []string{"air"})
		if !reflect.DeepEqual(actual, []string{"air", "humidity"}) {
			t.Error(actual)
		}
	})
	t.Run("keeps the list unchanged if it already carries the deprecated aspect id", func(t *testing.T) {
		actual := AspectIdsAlias("air", []string{"air", "humidity"})
		if !reflect.DeepEqual(actual, []string{"air", "humidity"}) {
			t.Error(actual)
		}
	})
	t.Run("keeps the list unchanged if the deprecated aspect id is unset", func(t *testing.T) {
		actual := AspectIdsAlias("", []string{"air"})
		if !reflect.DeepEqual(actual, []string{"air"}) {
			t.Error(actual)
		}
	})
}

func TestAspectNodesAlias(t *testing.T) {
	t.Run("folds the deprecated aspect node into an empty list", func(t *testing.T) {
		actual := AspectNodesAlias(testAspectNodes["air"], nil)
		if !reflect.DeepEqual(actual, testNodes("air")) {
			t.Error(actual)
		}
	})
	t.Run("keeps the list unchanged if the deprecated aspect node is unset", func(t *testing.T) {
		actual := AspectNodesAlias(AspectNode{}, testNodes("air"))
		if !reflect.DeepEqual(actual, testNodes("air")) {
			t.Error(actual)
		}
	})
	t.Run("keeps the list unchanged if it already carries the deprecated aspect node", func(t *testing.T) {
		actual := AspectNodesAlias(testAspectNodes["air"], testNodes("air", "humidity"))
		if !reflect.DeepEqual(actual, testNodes("air", "humidity")) {
			t.Error(actual)
		}
	})
}

func TestAspectMatchLevel(t *testing.T) {
	tests := []struct {
		name      string
		aspectIds []string
		queried   []AspectNode
		expected  int
	}{
		{
			name:      "matches every aspect at level 0 without a queried aspect",
			aspectIds: []string{"air"},
			queried:   nil,
			expected:  0,
		},
		{
			name:      "matches a content variable without an aspect at level 0 without a queried aspect",
			aspectIds: nil,
			queried:   nil,
			expected:  0,
		},
		{
			name:      "matches the queried aspect itself at level 0",
			aspectIds: []string{"air"},
			queried:   testNodes("air"),
			expected:  0,
		},
		{
			name:      "matches a child of the queried aspect at level 1",
			aspectIds: []string{"inside_air"},
			queried:   testNodes("air"),
			expected:  1,
		},
		{
			name:      "matches a deeper descendant of the queried aspect at level 2",
			aspectIds: []string{"today"},
			queried:   testNodes("electricity"),
			expected:  2,
		},
		{
			name:      "reports the closest aspect of a content variable that carries several",
			aspectIds: []string{"inside_air", "air"},
			queried:   testNodes("air"),
			expected:  0,
		},
		{
			name:      "rejects an ancestor of the queried aspect",
			aspectIds: []string{"air"},
			queried:   testNodes("inside_air"),
			expected:  -1,
		},
		{
			name:      "rejects a content variable without an aspect",
			aspectIds: nil,
			queried:   testNodes("air"),
			expected:  -1,
		},
		{
			name:      "rejects a content variable that carries only one of two queried aspects",
			aspectIds: []string{"inside_air"},
			queried:   testNodes("inside_air", "humidity"),
			expected:  -1,
		},
		{
			name:      "matches a content variable that carries both queried aspects",
			aspectIds: []string{"inside_air", "humidity"},
			queried:   testNodes("inside_air", "humidity"),
			expected:  0,
		},
		{
			name:      "reports the level of the worst matched aspect of several queried ones",
			aspectIds: []string{"inside_air", "today"},
			queried:   testNodes("air", "electricity"),
			expected:  2,
		},
		{
			//an aspect hierarchy admits at most one aspect per content variable, so no
			//variable can ever carry two siblings and a query for two is unanswerable
			name:      "rejects two sibling aspects, which one content variable cannot carry",
			aspectIds: []string{"inside_air"},
			queried:   testNodes("inside_air", "outside_air"),
			expected:  -1,
		},
		{
			//an ancestor queried next to its descendant adds no constraint: the
			//descendant satisfies both, so the query collapses to the descendant
			name:      "collapses an ancestor queried alongside its descendant",
			aspectIds: []string{"inside_air"},
			queried:   testNodes("air", "inside_air"),
			expected:  1,
		},
		{
			name:      "rejects the ancestor itself if a descendant is queried alongside it",
			aspectIds: []string{"air"},
			queried:   testNodes("air", "inside_air"),
			expected:  -1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := AspectMatchLevel(test.aspectIds, test.queried)
			if actual != test.expected {
				t.Error("expected", test.expected, "got", actual)
			}
		})
	}
}
