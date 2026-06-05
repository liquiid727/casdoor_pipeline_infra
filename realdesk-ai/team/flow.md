# Agent Team Flow — realdesk-ai

## 基线说明

`team/` 是仓库任务执行的默认控制面。

所有仓库内任务都应先经过：

1. `orchestrator`
2. `planner`

再按任务实际需要激活其余 agent。要求是“**每次都走 team**”，不是“**每次都启用全员**”。

如果是 P2 小改动，可以走紧凑模式；如果是 P0、高风险跨 context 变更或上线任务，则必须走完整链路。

## Team Roster

| Agent | 职责定位 | 激活时机 |
|---|---|---|
| `orchestrator` | 解析请求、识别受影响的 bounded context、路由任务 | 每次请求入口 |
| `planner` | 拆分任务、评估风险等级、分配执行顺序 | 任何多步任务 |
| `architect` | 架构分析、边界设计、跨上下文影响评估 | 新特性设计 / 跨模块变更 / 高风险路径 |
| `spec-writer` | 起草方向 spec / execution spec，并整理治理产物 | 需要先做设计冻结、再交实施时 |
| `backend/connectionplane` | connectionplane 上下文实现 | 涉及 session / 资源分配 / queue / autoscale |
| `backend/iam` | iam 上下文实现 | 涉及认证 / token / 设备策略 / 限流 |
| `backend/fleet` | fleet 上下文实现 | 涉及实例 / provider / 健康 / 扩缩容 |
| `backend/commerce` | commerce 上下文实现 | 涉及订单 / 支付 / 资产 / 账单 / 钱包 |
| `backend/operations` | operations 聚合层实现 | 涉及 admin / 配置中心 / 通知 / 促销 |
| `backend/shared` | shared 跨模块契约 | 修改共享契约 / 配置模型 / 工程规范 |
| `test-unit-agent` | changed-package coverage / 分支异常 / race / benchmark | 每次实现完成后，按需激活 |
| `test-api-connectivity-agent` | HTTP/gRPC route/auth/status/body/header/idempotency/health smoke | 触达 transport 或 API contract 时 |
| `test-e2e-scenario-agent` | spec 驱动业务链路 / replay / compensation / recovery | 触达完整业务链路时 |
| `test-triage-agent` | flaky 归因 / failure clustering / env smoke / manual gaps | 测试失败或验证缺口需要归因时 |
| `reviewer` | 代码审查 + spec governance 审查 + realdesk-ai 同步核查 | 每次实现完成后，或 spec 产物完成后 |
| `sync` | 更新 realdesk-ai 文档层 | 任何改变语义的变更后 |
| `ci` | Git 动作前的变更校验、stage 范围约束、CI Record 输出 | 请求 change validation / review readiness / commit / push / PR / merge 时 |
| `deploy` | 构建 / 环境部署 / 灰度 / 私有镜像站与凭据引用 | 上线操作 |
| `observability` | 运行时指标 / 告警 / 行为漂移检测 | 部署后 / 定期巡检 |

## 主流程

```
User Request
      │
      ▼
orchestrator
      │
      ▼
planner
      │
      ├─ architect (if needed)
      │
      ▼
backend/<context>
      │
      ▼
one or more test agents
  - test-unit-agent
  - test-api-connectivity-agent
  - test-e2e-scenario-agent
  - test-triage-agent
      │
      ▼
reviewer
      │
      ├─ sync (if needed)
      │
      ▼
deploy -> observability (if explicitly requested)
```

如果请求包含 `change validation`、`review readiness`、`commit`、`push`、`PR` 或 `merge`，则在 `reviewer` 和必需的 `sync` 之后追加 `ci-agent`，由它汇总现有 team 产物并输出 `CI Record`。

---

## Spec 流程

当任务目标是设计治理而不是直接实现时，走这条链路：

`orchestrator -> planner -> architect -> spec-writer -> reviewer -> sync`

阶段职责：

- `architect`
  - 给出 bounded context、truth owner、runtime owner、artifact class
  - 明确 `stable_spec_allowed: true | false`
- `spec-writer`
  - 只在可写结论后起草 `direction.md` / `execution.md`
  - 正式产物放入 `realdesk-ai/spec-craft/`
- `reviewer`
  - 审 spec governance，而不是只审代码
  - 确认没有把方向性设计误写成 stable truth
- `sync`
  - 只更新索引、team 文档、prompt 与治理契约
  - 不把 `spec-craft` 自动提升为 `specs`

实现链路只消费 `execution.md`，不再直接回到临时输入稿取规则。

## 任务交接协议

每个 agent 完成后必须输出：

```markdown
## Handoff
context: <受影响的 bounded context 列表>
risk_level: P0 | P1 | P2
changed:
  - 代码路径 + 层级落点
  - 状态机/流程是否变化
  - 幂等键/Redis key/契约是否变化
sync_needed: true | false
sync_targets:
  - modules/<module>.md
  - specs/<context>/*.yaml
  - knowledge/*.md
next_agent: reviewer | test | sync | ci | none
```

标准模板见 `team/templates/`。

`spec-writer` 的治理交接使用：

- `team/templates/spec-brief.md`
- `team/templates/execution-spec.md`
- `team/templates/spec-handoff.md`

`execution.md` 至少要明确：

- direction path
- risk level
- truth owner / runtime owner
- stable truth context / stable_spec_allowed
- implementation phases
- affected layers
- owner packages
- target paths
- API / contract deltas
- state / flow deltas
- config deltas
- Redis / queue deltas
- failure modes
- observability
- rollout
- test scope
- test inputs
- test asset outputs
- Bruno asset impact
- coverage gate targets
- integration requirements
- verification commands
- expected outcomes
- reviewer focus
- doc generation inputs
- doc maintenance inputs

`architect` 在 spec 模式下补充：

- `team/templates/architecture-eval.md`

`reviewer` 在 spec 模式下补充：

- `team/templates/spec-acceptance.md`

测试轨道 agents 除了 `Handoff` 之外，还必须额外输出：

- `team/templates/test-plan.md` 对应的 `Test Plan`
- `team/templates/test-report.md` 对应的 `Test Report`

`ci-agent` 在 Git 动作请求下额外输出：

- `team/templates/ci-record.md` 对应的 `CI Record`

`Test Plan` 至少要写清：

- `test_track`
- `execution_spec_path`
- `spec_sources`
- 本次受影响的 owner package
- 哪些包命中 `95%` hard gate
- 本次需要维护哪些测试资产
- Bruno 资产要不要更新
- 本次命中哪些 CI gates
- role-aware `case_groups`
- 要执行哪些命令
- 测试任务如何交给 reviewer / 后续 agent

`Test Report` 至少要包含：

- `test_track`
- `execution_spec_path`
- 使用了哪些行为矩阵 / checklist / 上游 handoff
- 跑了哪些命令
- 生成了哪些测试文件
- 测试资产和 Bruno 资产的同步结果
- P0 owner package 的覆盖率结果
- 哪些 gate 失败了
- 未覆盖风险与原因
- 后续 backlog / follow-up

## 风险等级定义

| 等级 | 定义 | 要求 |
|---|---|---|
| **P0** | 会话状态机 / 支付链路 / 资源分配 / 认证核心路径 | architect 必须介入；至少启用 `test-unit-agent`，并按接口/链路风险追加 `test-api-connectivity-agent` 与 `test-e2e-scenario-agent`；命中 hard gate 的 owner package 必须过 `95%`；并发敏感路径补 `-race`；reviewer 双重审查 |
| **P1** | 有幂等/超时/补偿风险但不涉及核心状态机 | 至少启用 `test-unit-agent`；触达公开接口时追加 `test-api-connectivity-agent`；必要时启用 `test-e2e-scenario-agent` 或 `test-triage-agent`；reviewer 审查一次 |
| **P2** | 纯聚合层 / admin 展示 / 配置读取 | 标准流程；仅启用需要的最小测试轨道；reviewer 一次 |

## 测试门禁

- `95%` 是 **P0 owner package 的硬门禁目标**，不是“全仓 go test 覆盖率”口号。
- 当前 hard gate 只作用于 `team/p0-unit-packages.txt` 中列出的首批 P0 规则包。
- P0 编排型大包在未完成拆分前，先执行“场景完整性 + 覆盖率只增不减”策略，达到稳定边界后再提升为 95% hard gate。
- `P0/P1` 与语义变更任务在进入测试前，先使用 `scripts/testing/spec-test-intake.sh` 生成标准化测试输入。
- 测试执行入口统一由 `scripts/testing/p0-unit-gate.sh` 提供；四个测试轨道都不应各自发明报告格式。
- `make check-testing-system` 负责 Bruno 资产与测试工程入口脚本的一致性检查。
- Bruno 在第一阶段只承担资产与调试入口职责，不是完整自动化回归门禁。
- `reviewer-agent` 必须复核 `Test Plan` / `Test Report`、coverage gate 结果和未覆盖风险说明。

## Git 动作约定

- `commit` / `push` / `PR` 前先满足 Change Validation，再由 `ci-agent` 汇总成 `CI Record`
- 代码或行为变更默认 `review_required: true`，但 `commit` / `push` / `PR` 可以记录 `review_status: pending`
- `merge` 前必须已有 reviewer 结论，且 `blocking_items` 清空
- `ci-agent` 只消费现有 team 产物与 `git status`，不额外引入新的 records 目录或 PR 正文解析 workflow

---

## realdesk-ai 维护原则

1. **每次语义变更后必须 sync** — 参见 `tools/sync-rules.md`
2. **realdesk-ai 不是事实源** — 只反映上游变更，不在此处发明规则
3. **sync-agent 只读取代码 + doc 后写** — 禁止凭记忆或 AI 推测补文档
4. **reviewer 必须核查 sync 状态** — 审查清单第 11 项是硬性门槛
5. **`spec-craft` 不是 stable truth** — 已审批设计只能停留在 `realdesk-ai/spec-craft/`，实现落地后才允许进入 `realdesk-ai/specs/`
