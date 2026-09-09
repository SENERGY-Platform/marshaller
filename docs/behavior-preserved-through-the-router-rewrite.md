# Behavior preserved through the router rewrite

## Scope

Holds for the api as of the router rewrite of 2026-09-09, which replaced
`julienschmidt/httprouter` with `net/http` and introduced the `Controller` interface. It
records three oddities that were kept on purpose and two changes that were made on purpose,
so that neither gets "fixed" or repeated by accident.

Not a description of how the routes are wired — that is
[the api is registered by reflection](the-api-is-registered-by-reflection.md).

The three preserved items are **not** endorsements. Each is a candidate for a deliberate
change with its own commit; what they are not is a refactoring artifact.

## Kept: a failed protocol lookup answers 400

`normalizeMarshallingRequest` and its three siblings in `lib/controller` answer
`StatusBadRequest` for both of their failure modes: a service and protocol that disagree,
**and** a protocol lookup against the device-repository that fails.

The second one is arguably a 500 — the caller did nothing wrong if a dependency is
unreachable. It answers 400 because that is what the endpoint did before the controller
existed: it mapped every error from its `normalizeRequest` closure to `StatusBadRequest`.
Changing it is a behavior change for callers that branch on the status, and it does not
belong in a refactor.

## Kept: metric labels use the old httprouter spelling

`lib/api/marshalv2.go` and `unmarshalv2.go` pass `resource+"/:serviceId"` as the endpoint
label, while the route itself is registered as `POST /v2/marshal/{serviceId}`.

The mismatch is intentional. The label is a **value** that dashboards and alerts query by,
not a description of the route. Renaming it starts a new time series: the old one stops
receiving data and every panel built on it goes silently empty, while the service stays
healthy. A label that disagrees with the pattern is the lesser evil, because it is visible
in the code where a reader can ask.

If the labels are ever renamed, the dashboards move in the same change.

## Kept: a v1 path-options query without an aspect matches only aspect-less variables

`GET /path-options` and `POST /query/path-options` with no aspect named match **only**
content variables that carry no aspect either. This is not the same as "do not filter by
aspect", which is what the v2 path resolution means when it gets no aspect.

The rule lives in `aspectsMatch` in `lib/marshaller/pathoptions.go` and is commented there.
It is undocumented in the api and covered by no test that predates the rewrite, so whether
anyone relies on it is unknown — which is precisely why it was preserved rather than
aligned with v2. The two endpoints therefore disagree on what "no aspect" means, and they
did so before the rewrite as well.

## Changed on purpose: a trailing slash is a 404

`httprouter` had `RedirectTrailingSlash: true`, so `POST /marshal/` was answered with a 301
to `/marshal`. `http.ServeMux` does not match a trailing slash against a pattern registered
without one, so it is now a 404.

This was decided rather than discovered. Registering both patterns would have preserved the
redirect; the choice was to let the canonical path be the only one. A caller that appends a
slash today breaks, and it breaks silently from its own point of view.

## Changed on purpose: without-envelope reads its own query parameter

`GET /path-options` read the `without-envelope` flag from the `function-id` parameter and
swallowed the `strconv.ParseBool` error. The envelope could therefore not be switched off
through the query-parameter form at all, and a non-boolean value passed unnoticed. It now
reads `without-envelope` and answers 400 on a value that is not a boolean.

The request-body form was never affected — it takes `without_envelope` as a typed field.

`TestPathOptionsWithoutEnvelope` in `lib/tests/aspect_ids_test.go` pins all four cases.
Note that those subtests name an aspect explicitly: without one they would match nothing,
because of the preserved rule above.
