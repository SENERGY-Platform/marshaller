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

package xml

import (
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/serialization/base"
	"github.com/SENERGY-Platform/models/go/models"
	"reflect"
	"testing"
)

func TestUnmarshalSimpleInt(t *testing.T) {
	value := `<i>24</i>`

	marshaller, ok := base.Get(Format)
	if !ok {
		return
	}

	out, err := marshaller.Unmarshal(value, model.ContentVariable{Name: "i", Type: model.Integer})

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(out, float64(24)) {
		t.Fatal(out)
	}
}

func TestUnmarshalSimpleString(t *testing.T) {
	value := `<s>foobar</s>`

	marshaller, ok := base.Get(Format)
	if !ok {
		return
	}

	out, err := marshaller.Unmarshal(value, model.ContentVariable{Name: "s", Type: model.String})

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(out, "foobar") {
		t.Fatal(out)
	}
}

func TestUnmarshalSimpleMap(t *testing.T) {
	value := `<example attr="attrVal"><body>bodyVal</body></example>`

	marshaller, ok := base.Get(Format)
	if !ok {
		return
	}

	out, err := marshaller.Unmarshal(value, model.ContentVariable{
		Name: "example",
		Type: model.Structure,
		SubContentVariables: []model.ContentVariable{
			{Name: "-attr"},
			{Name: "body"},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(out, map[string]interface{}{"-attr": "attrVal", "body": "bodyVal"}) {
		t.Fatal(out)
	}
}

func TestUnmarshalSerializationOptionXmlAttribute(t *testing.T) {
	value := `<example attr="attrVal"><body>bodyVal</body></example>`

	marshaller, ok := base.Get(Format)
	if !ok {
		return
	}

	out, err := marshaller.Unmarshal(value, model.ContentVariable{
		Name: "example",
		Type: model.Structure,
		SubContentVariables: []model.ContentVariable{
			{Name: "attr", SerializationOptions: []string{models.SerializationOptionXmlAttribute}},
			{Name: "body"},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(out, map[string]interface{}{"attr": "attrVal", "body": "bodyVal"}) {
		t.Fatal(out)
	}
}

func TestUnmarshalSimpleList(t *testing.T) {
	value := `<list><element>1</element><element>2</element><element>3</element></list>`

	marshaller, ok := base.Get(Format)
	if !ok {
		return
	}

	out, err := marshaller.Unmarshal(value, model.ContentVariable{
		Name: "list",
		Type: model.List,
		SubContentVariables: []model.ContentVariable{
			{Name: "*", Type: model.Integer},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(out, map[string]interface{}{"element": []interface{}{float64(1), float64(2), float64(3)}}) {
		t.Fatal(out)
	}
}

//TODO: list tests
