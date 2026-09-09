# The api is registered by reflection

## Scope

Holds for `lib/api` since the router rewrite of 2026-09-09, which replaced
`julienschmidt/httprouter` with `net/http`. Read this before adding, renaming or removing
an endpoint.

Not about the controller behind the endpoints — that is
[the client mirrors the controller](the-client-mirrors-the-controller.md). Not about which
behavior the rewrite deliberately kept, which is
[behavior preserved through the router rewrite](behavior-preserved-through-the-router-rewrite.md).

## How a route gets registered

Each endpoint file declares a struct, registers an instance of it in `init()`, and carries
one method per route:

```go
func init() {
	endpoints = append(endpoints, &Marshalling{})
}

type Marshalling struct{}

func (this *Marshalling) Marshal(config config.Config, router *http.ServeMux, ctrl Controller, m *metrics.Metrics) {
	router.HandleFunc("POST /marshal", func(writer http.ResponseWriter, request *http.Request) { … })
}
```

`GetRouterWithoutMiddleware` walks `endpoints`, reflects over each object, and calls every
method it can type-assert to `EndpointMethod`. The shape follows the device-repository's
`lib/api`, with one addition: the alias carries the metrics collector, because the
per-request metrics need the decoded request body and therefore cannot sit in a middleware.

One method per route, rather than one per endpoint group, is deliberate — the swagger
annotations attach to the method, so a route and its documentation stay together.

## A drifted signature is silent, which is why api_test.go lists every route

The registration is a type assertion:

```go
f, ok := m.Interface().(EndpointMethod)
if ok { result[name] = f }
```

`ok == false` is the normal case, because the endpoint structs also carry methods that are
not routes. So a method whose signature no longer matches the alias — a parameter added,
removed or reordered — is skipped exactly like a helper. It compiles, it is never called,
and its route answers 404.

Nothing else catches it. The compiler is satisfied, the controller tests pass because the
operation itself is fine, and a test that drives one endpoint passes because it names its
own route.

`lib/api/api_test.go` is the guard: it builds the router against a stub controller and
asserts that every route this service serves answers something other than 404 or 405.
**A new route has to be added to that list**, and renaming one that callers use becomes a
visible decision rather than a silent 404.

Two properties of that test are load-bearing:

- It asserts "not 404 and not 405", not a status. An empty body reaches the handler; what
  the handler answers is the subject of the tests under `lib/tests`. Asserting on a status
  here would make the test break on every validation change.
- 405 is checked separately, because `http.ServeMux` answers 405 once a pattern matches the
  path but not the method — so a route registered under the wrong verb is otherwise
  indistinguishable from a working one.

The stub controller earns its lines twice: because it has to implement `Controller`, a
method added to that interface without an endpoint fails to compile in the test instead of
shipping unreachable.

## Adding an endpoint

1. Add the operation to `Controller` in `lib/api/interfaces.go` and implement it in
   `lib/controller`.
2. Add the method to the endpoint struct, with the exact `EndpointMethod` signature.
3. Add the route to the list in `lib/api/api_test.go`.
4. Add the swagger annotations and run `go generate ./...`; the spec under `docs/` is
   committed and CI fails on a difference.
5. Add the client method in `lib/client` — the client implements the same `Controller`, so
   omitting it is a compile error, but calling the wrong path is not.
