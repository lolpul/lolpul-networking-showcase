# Networking showcase memory - v4

- Purpose: independent public Go lifecycle examples plus a conceptual experimental networking overview.
- Repository: https://github.com/lolpul/lolpul-networking-showcase; private product remains private.
- Entry points: examples/session-lifecycle/writer.go, examples/cancellation/task.go, examples/resource-cleanup/scope.go; corresponding tests.
- Architecture: serial writer owner, cooperative task with cancel/join, scoped two-resource acquisition. Standard library only; no real networking.
- Contract limits: workers cooperate; writer/Close terminate and do not reenter; non-nil interface dependencies; failed close attempted once without retry.
- Commands: gofmt -l examples; go vet ./...; go test -timeout 30s ./...; go test -race -timeout 30s ./.... CI .github/workflows/go.yml pins Go 1.26.4.
- Audit: only generic contracts, ownership, cancellation, concurrency, and failure-test concepts selected. No source/history/config transferred; private receipts and backups are ignored.
- Validation: Windows Go 1.26.4 tests/vet pass; Linux formatting/vet/tests/race pass in [Actions 37048917736](https://github.com/lolpul/lolpul-networking-showcase/actions/runs/37048917736) for implementation 83de2a619133b0f318df16bc9e5e74affafb4a12.
- Active step: public examples accepted; automated inventory and manual confidentiality review passed. Next: Python showcase, then Kotlin, profile, and portfolio code evidence.
- Documents: [scope](spec.md), [interview notes](interview-notes.md), [patch](patches/2026-10-02-go-examples.md).
- Portfolio: https://elisey.kochura.com/work/lolpul-vpn. Site code evidence follows after all three showcases are published.
