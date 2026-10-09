# Verification — 2026-10-10

Go 1.26.4. Windows `go fmt ./...` produced no code changes; `go vet ./...` and `go test -count=1 -timeout=30s -v ./...` passed. Linux `go test -race -count=1 -timeout=30s ./...` passed for all three packages. The suite contains 14 top-level tests and six subcases.

There is no executable server: `go test` is the reproducible demonstration. Cooperative workers and terminating adapters are required; passing race checks does not prove progress under arbitrary blocking code. No real transport, protocol, throughput benchmark or production service was exercised.

[Actions](https://github.com/lolpul/lolpul-networking-showcase/actions/workflows/go.yml) provides the hosted formatting/vet/unit/race checks.
