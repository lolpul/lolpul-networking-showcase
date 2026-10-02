# Independent Go examples

## Intent and scope

Replace the documentation-only entry point with three runnable public examples of generic lifecycle concerns. A read-only private-source audit confirmed ownership, cancellation, partial failure cleanup, and concurrent contract tests. All example code and tests are newly written; no private file, interface, internal comment, protocol, configuration, or history was imported.

## Changes and decisions

- `examples/session-lifecycle`: 69-line serial writer owner; writes and shutdown share a mutex, close result retained.
- `examples/cancellation`: 51-line cooperative worker owner; cancel plus join using result publication through channel closure.
- `examples/resource-cleanup`: 51-line scoped acquisition; register cleanup immediately, unwind in reverse order, join errors.
- Root Go module, Linux Actions, README, interview notes, scope specification, and compact memory. Repository metadata now describes actual public examples.

Counts include comments and blank lines; executable non-test code is 130 nonblank, non-comment lines across the three files. Keeping cancellation and cleanup small is deliberate rather than adding abstractions to reach a line target.

## Verification

Windows Go 1.26.4: formatting, `go vet ./...`, and `go test -timeout 30s ./...` passed. Fourteen top-level tests plus six failure-stage subtests cover lifecycle, cancellation, virtual-clock deadline expiration, concurrent shutdown, partial acquisition, ownership counts, and aggregated errors. Tests do not sleep. Linux `go test -race` and exact-revision Actions acceptance follow publication and will be recorded after they finish.

Automatic inventory scan and manual review found no confidential values, private URLs, addresses, emails, internal source paths, or original source imports in the publication set. Matches were public portfolio/GitHub links and SVG namespace metadata. Scan results are not a guarantee. Ignored local audit receipts contain source revision evidence and are excluded from Git.

## Rollback and limits

Original main: `4c3bb6873b9f43cee9d0df4df70ebfc893fc11d0`. A timestamped ignored backup includes a target manifest for the README, memory, metadata, and Obsidian note. Revert the new public commit for rollback; do not rewrite history. Cooperative workers and terminating closers are assumptions. Host tests do not validate the private network stack, external cleanup guarantees, or production behavior. No hardware, service, network configuration, or private repository was changed.
