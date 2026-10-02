# Networking showcase memory - v4

- Purpose: independent public Go lifecycle examples plus a conceptual experimental networking overview.
- Repository: https://github.com/lolpul/lolpul-networking-showcase; private product remains private.
- Entry points: examples/session-lifecycle/writer.go, examples/cancellation/task.go, examples/resource-cleanup/scope.go; corresponding tests.
- Architecture: serial writer owner, cooperative task with cancel/join, scoped two-resource acquisition. Standard library only; no real networking.
- Contract limits: workers cooperate; writer/Close terminate and do not reenter; non-nil interface dependencies; failed close attempted once without retry.
- Commands: gofmt -l examples; go vet ./...; go test -timeout 30s ./...; go test -race -timeout 30s ./.... CI .github/workflows/go.yml pins Go 1.26.4.
- Audit: only generic contracts, ownership, cancellation, concurrency, and failure-test concepts selected. No source/history/config transferred; private receipts and backups are ignored.
- Validation: Windows Go host tests and vet passed; Linux race verification and publication acceptance are pending.
- Active step: finish race, automatic/manual confidentiality review, commit/push, and confirm exact-revision Actions before starting the next showcase.
- Documents: [scope](spec.md), [interview notes](interview-notes.md), [patch](patches/2026-10-02-go-examples.md).
- Portfolio: https://elisey.kochura.com/work/lolpul-vpn. Site code evidence follows after all three showcases are published.
