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

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testSecret = "s3cr3t-not-a-real-credential"

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	location := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(location, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return location
}

// TestSecretIsLoadedFromEnv guards the combination of two things that were decided
// separately: AuthClientSecret is a types.Secret, while the config is still loaded by the
// hand-rolled environment reflection in this package. That loader assigns by reflect.Kind,
// and types.Secret is a string kind — if it ever stops being one, the secret silently stays
// empty and every token request fails with an unhelpful error.
func TestSecretIsLoadedFromEnv(t *testing.T) {
	location := writeConfig(t, `{"auth_client_secret":""}`)
	t.Setenv("AUTH_CLIENT_SECRET", testSecret)

	conf, err := Load(location)
	if err != nil {
		t.Fatal(err)
	}
	if conf.AuthClientSecret.Value() != testSecret {
		t.Error("expected the secret from the environment, got", conf.AuthClientSecret.Value())
	}
}

func TestSecretIsLoadedFromFile(t *testing.T) {
	location := writeConfig(t, `{"auth_client_secret":"`+testSecret+`"}`)

	conf, err := Load(location)
	if err != nil {
		t.Fatal(err)
	}
	if conf.AuthClientSecret.Value() != testSecret {
		t.Error("expected the secret from the file, got", conf.AuthClientSecret.Value())
	}
}

// TestSecretIsMaskedWhenPrinted is the reason the type was introduced: a careless print or
// a config dump must not carry the credential.
func TestSecretIsMaskedWhenPrinted(t *testing.T) {
	location := writeConfig(t, `{"auth_client_secret":"`+testSecret+`"}`)
	conf, err := Load(location)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("fmt of the whole config", func(t *testing.T) {
		if rendered := fmt.Sprintf("%v", conf); strings.Contains(rendered, testSecret) {
			t.Error("the secret is in the rendered config")
		}
	})
	t.Run("fmt with the plus flag", func(t *testing.T) {
		if rendered := fmt.Sprintf("%+v", conf); strings.Contains(rendered, testSecret) {
			t.Error("the secret is in the rendered config")
		}
	})
	t.Run("json marshal", func(t *testing.T) {
		rendered, err := json.Marshal(conf)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(rendered), testSecret) {
			t.Error("the secret is in the marshalled config")
		}
	})
	t.Run("Value still returns it", func(t *testing.T) {
		if conf.AuthClientSecret.Value() != testSecret {
			t.Error("expected Value to return the real secret")
		}
	})
}
