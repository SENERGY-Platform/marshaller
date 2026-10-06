<a href="https://github.com/SENERGY-Platform/marshaller/actions/workflows/tests.yml" rel="nofollow">
    <img src="https://github.com/SENERGY-Platform/marshaller/actions/workflows/tests.yml/badge.svg?branch=master" alt="Tests" />
</a>

# Documentation

`docs/` holds the generated OpenAPI spec — `marshaller_swagger.json`,
`marshaller_swagger.yaml` and `marshaller_docs.go`, regenerated with `go generate ./...`
and committed, with CI failing on a difference — next to hand-written notes on this
service's own particularities:

- [the api is registered by reflection](docs/the-api-is-registered-by-reflection.md) — how a
  route gets registered, why a drifted method signature is a silent 404, and what to touch
  when adding an endpoint
- [the client mirrors the controller](docs/the-client-mirrors-the-controller.md) — one
  interface for the api and the go client, why that does not prove the client calls the
  right path, and the one endpoint the client deliberately does not offer
- [behavior preserved through the router rewrite](docs/behavior-preserved-through-the-router-rewrite.md)
  — three oddities kept on purpose and two changes made on purpose, so neither gets
  "fixed" or repeated by accident
- [marshal writes the closest input paths](docs/marshal-writes-the-closest-input-paths.md)
  — which input variables a function with aspects sets, and the switch to set all of them

The spec itself is served under `/doc`; set `enable_swagger_ui` to also serve the swagger
ui under `/swagger`.

# Tests
- use the `-short` flag as test argument to test with mocks
- if the `-short` flag is not used the tests will try to call services defined in `lib/tests/testdata/config.json`
