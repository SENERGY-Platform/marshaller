# Marshal writes the closest input paths

## Scope

Holds for `lib/marshaller/v2` as of 2026-10-06, when controlling functions started to be
combined with aspects the way measuring functions are. Applies to a marshal request entry
that names a function and no `paths`; an entry with `paths` is written to exactly those.

## What it does

`Marshal` collects the input variables matching the function and every queried aspect,
sorted by `model.AspectMatchLevel` — 0 for the aspect itself, 1 for a child, 2 for any
deeper descendant. With the default `v2.InputPathSelection = v2.ClosestInputPaths` only the
paths at the smallest distance get the value. A request for `air` against a service with
variables for `air`, `inside_air` and `outside_air` sets `air` only; without an `air`
variable it sets both children, because they tie.

This mirrors the unmarshal side, which reads only the closest output path. A request
without aspects still writes every variable of its function, because all of them lie at
distance 0 — the behavior before aspects reached the input side.

`v2.AllInputPaths` restores writing every match. It is a package variable and not a
constant on purpose: which variables a command for a parent aspect should set was decided
for the closest one, with the explicit expectation that the requirement may change.
`GetInputPaths` is unaffected by the mode and still returns every match.

## Fixed along the way

`setContentVariableValues` used to append one copy of the inputs per path, so a second
path started again from the original inputs, and an entry matching no path returned no
inputs at all. The latter emptied the message for every following entry of the same
request. Each path is now applied to the result of the one before, and an entry without a
match leaves the inputs as they are.

Not changed: an entry naming a function that matches no input still produces no error. A
configurable that does not fit the service takes the same route, so turning it into an
error is a decision of its own.
