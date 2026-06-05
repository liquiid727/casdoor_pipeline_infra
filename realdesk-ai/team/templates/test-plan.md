# Test Plan Template

```markdown
## Test Plan
test_track: unit | api_connectivity | e2e_scenario | triage
risk_level: P0 | P1 | P2
execution_spec_path: <realdesk-ai/spec-craft/.../execution.md，没有则写 none>
spec_sources:
  - <execution.md / specs / knowledge / none>
contexts:
  - <affected context>
owner_packages:
  - package: <./internal/...>
    threshold: 95 | none
    hard_gate: true | false
    why: <为什么它是本次变更的真相 owner 包>
test_asset_outputs:
  - <Test Plan / Test Report / Bruno README / collection refresh / none>
bruno_asset_actions:
  - <update README / run make bruno-docs / none>
ci_gates:
  - <make check-realdek-ai / make check-testing-system / make p0-unit-gate / none>
case_groups:
  - name: <happy_path / invalid_state / idempotency / concurrency / timeout_or_dependency_failure / compensation_or_rollback / late_event / route_method / auth / status_body / headers / idempotency_smoke / health_smoke / replay / failure_compensation / recovery / query_callback / failure_cluster / flake_repro / env_migration_smoke / manual_gap / uncovered_risk>
    cases:
      - <用例，没有则写 none>
commands:
  - <命令，没有则写 none>
delivery_tasks:
  - <测试任务交接项，没有则写 none>
reviewer_focus:
  - <reviewer 需要重点复核的点，没有则写 none>
```
