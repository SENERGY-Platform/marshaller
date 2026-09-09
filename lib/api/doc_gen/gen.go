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

// This package exists only to carry the go:generate directive for the swagger spec, the
// way the device-repository does it. It sits next to the api package rather than inside
// it so that the generator's own -d and -o paths stay relative to one known directory.
//
// The generated files under docs/ are committed on purpose, so the spec is readable
// without a build. Regenerate with `go generate ./...` from the repository root.
package main

//go:generate go tool swag init --instanceName marshaller -o ../../../docs --parseDependency -d .. -g api.go

func main() {}
