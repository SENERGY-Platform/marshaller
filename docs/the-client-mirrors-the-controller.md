# The client mirrors the controller

## Scope

Holds for `lib/api/interfaces.go`, `lib/controller` and `lib/client` since 2026-09-09.
Read this before adding an operation, and before assuming the shared interface guarantees
more than it does.

Not about how the endpoints reach the router — that is
[the api is registered by reflection](the-api-is-registered-by-reflection.md).

## One interface, two implementations

`api.Controller` is the operation surface of this service: one method per endpoint, no http
types, and an http status alongside every result.

```go
result, err, code := ctrl.MarshalV2(request)
```

Two things implement it. `lib/controller` does the work; `lib/client` does the same calls
over http and is declared as `type Interface interface { api.Controller }`. A consumer can
therefore hold either without knowing which — that is what makes it possible to run this
service in-process in a consumer's test and swap the client for the real controller.

The interface is declared in `lib/api` rather than in `lib/controller`, because
`lib/client` has to satisfy the same contract and a controller that imported the api
package would close the loop.

`GetRouterWithoutMiddleware` exists for the same reason: mounted in an `httptest` server it
answers over http without needing a token issuer, so a consumer can exercise its own
request path against this service's real router.

## The interface proves the shape, not the wiring

This is the trap the symmetry invites. The compiler checks that the client has a method of
the right name and signature. It says nothing about the **path** that method calls:

```go
//compiles; answers 404 forever
return get[…](c.baseUrl, "/characteristic-path/"+…)
```

A typo in a URL, a wrong verb, or a path that drifted when a route was renamed all compile
and all pass every test that does not actually perform the request.

`lib/client/client_test.go` is the guard. It runs the real api in an `httptest` server
against the real controller and calls **every** method of the interface, rejecting a 404 or
405. It deliberately does not assert on the answers — what an operation returns is the
subject of the tests under `lib/tests`; what this test establishes is that client method
and endpoint agree.

Since 2026-09-09 the tests under `lib/tests` and `lib/tests/v2` also drive the service
through the client rather than through hand-built requests, so a mismatch fails a real
behavior test too and not only the client's own.

## Where the client is deliberately asymmetric

`GET /path-options` has **no** client method, and that is not an omission.

The endpoint exists in two forms, and they do not mean the same thing. An absent
`characteristic-filter` query parameter makes the GET form answer with an empty result,
while an empty `characteristic_id_filter` in the request body means "do not filter" and
yields every characteristic of the function. Offering both through one client method would
have made the semantics depend on which fields the caller happened to leave empty, so
`GetPathOptions` uses the request-body form only.

The consequence for tests: the subtests that cover the query-parameter form build their
request by hand, and `lib/tests/aspect_ids_test.go` keeps a small `getPathOptions` helper
for it. That is the one endpoint the client cannot reach.

## Adding an operation

Add it to `Controller`, implement it in `lib/controller`, add the endpoint method, add the
route to `lib/api/api_test.go`, then add the client method. Skipping the client method is a
compile error. Getting its path wrong is not — that is what the client test is for.
