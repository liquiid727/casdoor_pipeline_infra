## Test Plan
test_track: triage
risk_level: P1
execution_spec_path: none
spec_sources:
  - git diff
  - realdesk-ai/team/ci-agent.md
contexts:
  - iam
  - deployment
  - i18n
  - web-ui
  - git-delivery
owner_packages:
  - package: github.com/liquiid727/pipeline-auth/deployment
    threshold: none
    hard_gate: false
    why: 本次构建阻断直接出现在 deployment 包编译与测试
  - package: github.com/liquiid727/pipeline-auth/i18n
    threshold: none
    hard_gate: false
    why: locale 重复键会直接阻断当前仓库的去重测试
test_asset_outputs:
  - realdesk-ai/team/test-plan-20260605-pipeline-auth-ci.md
  - realdesk-ai/team/test-report-20260605-pipeline-auth-ci.md
  - realdesk-ai/team/findings-20260605-pipeline-auth-ci.md
  - realdesk-ai/team/ci-record-20260605-pipeline-auth-ci.md
bruno_asset_actions:
  - none
ci_gates:
  - go test ./i18n ./deployment
  - make backend
  - cd web && yarn build
  - go test ./...
case_groups:
  - name: env_migration_smoke
    cases:
      - 验证 ci-agent 引入后，后端 build 与前端 build 仍然可执行
      - 验证 deployment 与 i18n 的定向门禁从 fail 恢复为 pass
  - name: uncovered_risk
    cases:
      - 识别全仓 go test 中仍然存在的环境依赖测试与行为漂移
commands:
  - go test ./i18n ./deployment
  - make backend
  - cd web && yarn build
  - go test ./...
delivery_tasks:
  - 将通过与失败的命令结果落盘到 CI Record
  - 提交前仅排除 .cache/，其余当前改动按 repo-wide fork cleanup 一并交付
reviewer_focus:
  - deployment nil provider 测试策略是否合理
  - i18n key 去重是否引入跨命名空间文案回退风险
  - 全仓 go test 失败是否属于已知环境依赖还是新语义回归
