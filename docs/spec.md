# Public Go examples

## Problem and scope

An architecture overview alone cannot demonstrate how cancellation and ownership behave. Add three independently written, standard-library-only examples that a reader can run without a server or credentials. Audit of a private experimental project confirmed explicit boundaries, context propagation, serialized lifecycle operations, idempotent shutdown, cleanup after partial acquisition, and concurrent state tests. Only these generic concerns inform the new examples; no private files, interfaces, comments, protocols, settings, or history are transferred.

## Stages and acceptance

1. Implement a small writer session: ownership transfers on construction, writes serialize, shutdown waits for in-flight writes, repeated close returns the same result. Verify lifecycle and concurrent close/write tests.
2. Implement a cooperative cancellable task: parent cancellation/deadline reaches the worker; shutdown cancels and joins it. Verify success, cancellation, deadline, and concurrent shutdown with channel barriers, without sleeps.
3. Implement scoped acquisition of two resources: clean partial acquisitions, unwind in reverse order, preserve operation and cleanup errors. Verify every failure boundary and close counts.
4. Document contracts, limitations, and interview discussion; add CI. Run `gofmt`, `go vet ./...`, `go test ./...`, `go test -race ./...`, confidentiality review, and `git diff --check` before commit/push. Confirm CI for the exact published revision.

No real network access, transport implementation, production deployment, or claim of private-project validation is part of these host tests. Cancellation requires worker cooperation; closing a writer must return. Backups and the private audit receipt remain ignored. Rollback is a normal revert of the public patch, never a rewritten history.
