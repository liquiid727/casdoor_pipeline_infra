## Test Report
test_track: triage
scope: mixed
risk_level: P1
execution_spec_path: none
coverage_gate:
  command: none
  result: not_run
matrix_sources:
  - git diff
  - realdesk-ai/team/test-plan-20260605-pipeline-auth-ci.md
executed_commands:
  - go test ./i18n ./deployment -> pass
  - make backend -> pass
  - cd web && yarn build -> pass
  - go test ./... -> fail
generated_tests:
  - none
asset_sync_result:
  - Added ci-agent/team governance assets and validation records under realdesk-ai/team
bruno_doc_status:
  - skipped: repository changes do not touch Bruno assets
coverage_summary:
  - package: github.com/liquiid727/pipeline-auth/deployment
    threshold: none
    actual: 0
    status: pass
  - package: github.com/liquiid727/pipeline-auth/i18n
    threshold: none
    actual: 0
    status: pass
gate_failures:
  - go test ./... fails in certificate because account_test.go requires local acme_account.key fixture
  - go test ./... fails in object because syncer_user_test.go assumes seeded users and panics on empty result
  - go test ./... fails in radius because localhost:1812 is not serving during test run
  - go test ./... fails in sync because localhost:3306 is unavailable
  - go test ./... fails in sync_v2 because test-db.v2tl.com is unreachable
  - go test ./... fails in util because system tests assert unsupported network/version behavior
  - go test ./... fails in xlsx because ../../tmpFiles/example fixture is missing
residual_risks:
  - object/mfa_totp_test.go still reports default issuer expectation drift: expected Casdoor but actual issuer is Pipeline Auth
  - Full-suite health is not green enough for merge readiness even though backend/frontend builds now pass
reviewer_focus:
  - Confirm whether remaining full-suite failures are acceptable legacy/environment debt for branch push
  - Decide whether TOTP issuer expectation should be updated or implementation reverted for merge
doc_impact:
  - realdesk-ai/team/README.md
  - realdesk-ai/team/flow.md
  - realdesk-ai/team/ci-agent.md
backlog_followups:
  - Normalize environment-dependent tests behind fixtures or t.Skip guards
  - Reconcile Pipeline Auth branding changes with tests that still assert Casdoor strings
next_agents:
  - reviewer
