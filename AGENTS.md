# AI Agent Guide — cluster-backup-operator

ACM hub disaster recovery operator. Drives Velero to back up and restore hub resources.
Does NOT handle managed cluster app DR (that's OADP) or Velero internals.

**Repo:** https://github.com/stolostron/cluster-backup-operator
**Docs:** https://docs.redhat.com/en/documentation/red_hat_advanced_cluster_management_for_kubernetes/2.16/html/business_continuity/index

## Quick Reference (load this first, read deeper files only when needed)

### Build/Test/Lint
```bash
make build                    # build binary (includes generate + fmt + vet)
make test                     # full test suite (includes manifests + generate + fmt + vet)
make lint                     # golangci-lint v2.7.2
make manifests generate       # regenerate CRDs + deepcopy (required after api/ changes)
```

### Key Design Rules (never violate)
1. **CRD gate** — `main.go` blocks startup until Velero CRDs exist. Don't remove.
2. **One active Restore** — second Restore gets `FinishedWithErrors`.
3. **Schedule/Restore conflict** — non-paused BackupSchedule + active Restore = `FailedValidation`.
4. **Package-level slices** — never append directly to `backupManagedClusterResources` etc. Copy first.
5. **local-cluster exclusion** — must be enforced in EVERY loop touching Hive resources.
6. **Test functions ≤ 60 statements** — `funlen` linter will fail CI.
7. **Lines ≤ 120 chars** — `lll` linter.
8. **DCO sign-off** — all commits need `--signoff`.
9. **Both Dockerfiles** — new top-level packages need `COPY` in both `Dockerfile` AND `Dockerfile.rhtap`.

### File Map (what's where)
```
main.go                     — Entrypoint, scheme, TLS, CRD gate
api/v1beta1/                — CRD types, webhook validation
controllers/
  backup.go                 — Backup content selection (API groups, labels, resource lists)
  schedule_controller.go    — BackupScheduleReconciler (creates 5 Velero schedules)
  schedule.go               — Schedule phase, collision detection, cron parsing
  restore_controller.go     — RestoreReconciler (main reconcile loop)
  restore.go                — Restore phase machine, Velero restore properties, backup selection
  restore_post.go           — Post-restore: cleanup, auto-import-secret creation, activation
  pre_backup.go             — Pre-backup: MSA tokens, Hive/AI/metal secret labeling
  utils.go                  — Shared helpers (MSA secret selection, hub ID, BSL check)
  create_helper_test.go     — Test-only builder-pattern constructors
config/                     — Kustomize (CRDs, RBAC, manager, webhook, samples)
hack/crds/                  — Extra CRDs for envtest (Velero, OCM, Hive, OpenShift)
pkg/tlsconfig/              — TLS config from APIServer profile
```

### Coding Conventions
- Go 1.25, `goimports` formatting
- Structured logging: `log.FromContext(ctx).Info("msg", "key", val)`
- Error wrapping: `fmt.Errorf("...: %w", err)`
- Ginkgo v2 + Gomega + envtest for tests
- Test helpers in `create_helper_test.go` — use builder pattern, don't create raw objects
- golangci-lint v2.7.2: `funlen` (60 stmts), `lll` (120 chars), `misspell`, `unparam`, `errcheck`

### CI Checks
- `images`, `unit-tests`, `sonar`, `pr-image-mirror`, `crd-and-gen-files-check`
- Release branches add: Konflux build + Enterprise Contract
- PRs squash-merged via Tide

## Deep-Dive References (read only when needed for the specific task)

For **code changes** (adding resources, modifying restore/backup behavior, fixing bugs):
→ Read the Code Map at `~/workspace/src/github.com/sahare/ai-tools/projects/acm-backup-triage/CODE_MAP.md`
  Has function signatures, line numbers, branch difference matrix, test patterns, common operations.

For **customer issue investigation**:
→ Read Investigation Playbooks at `~/workspace/src/github.com/sahare/ai-tools/projects/acm-backup-triage/INVESTIGATION_PLAYBOOKS.md`
  Decision trees for every symptom category.
→ Read Knowledge Base at `~/workspace/src/github.com/sahare/ai-tools/projects/acm-backup-triage/KNOWLEDGE_BASE.md`
  22 indexed issues with root causes, diagnostics, resolutions.

For **CVE/security fix work**:
→ Read KB CVE sections at `~/workspace/src/github.com/sahare/ai-tools/projects/acm-backup-triage/KNOWLEDGE_BASE.md#standard-cve-fix-workflow-business-continuity--sustaining-admins`

For **past incident details**:
→ Read incident writeups at `~/workspace/src/github.com/sahare/ai-tools/projects/acm-backup-triage/incidents/`

## Related Repositories

| Repo | What it contains |
|------|-----------------|
| [multiclusterhub-operator](https://github.com/stolostron/multiclusterhub-operator) | Helm chart deploying this operator |
| [OADP operator](https://github.com/openshift/oadp-operator) | Installs Velero, manages DPA |
| [Velero](https://github.com/vmware-tanzu/velero) | Backup/restore engine |

## Common Pitfalls (quick-check before submitting)

- Package-level slice mutation without local copy → data race
- Second Restore resource instead of editing existing one → rejected
- Go default TLS → wrong; use `BuildTLSConfig` from APIServer profile
- `PartiallyFailed` Velero restores → often normal (empty backup files)
- Container name is `cluster-backup`, not `manager`
- `createAutoImportSecret()` is plain `Create`, not upsert — silently fails on `AlreadyExists`
- `docs/ARCHITECTURE.md` vs `docs/architecture.md` case collision exists on `main`
- `openshift-cherrypick-robot` stops at first branch failure
