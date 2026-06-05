# realdesk-ai Team Baseline

`realdesk-ai/team/` 不是附属说明区，而是仓库内 AI 执行任务的默认控制面。

目标只有一个：**每次执行都先走 team，再按任务大小激活需要的 agent，而不是一上来直接跳进某个 backend agent。**

## 执行基线

对仓库内的分析、设计、实现、Review、同步、发布类任务，统一采用这条基线：

1. `orchestrator` 先做入口路由
2. `planner` 再做任务拆分与风险分级
3. 按需激活 `architect`、`spec-writer`、`backend/*`、`test`、`reviewer`、`sync`、`ci`、`deploy`、`observability`

这条规则的关键点是：

- **每次都走 team**
- **不是每次都激活全员**
- **backend agent 不能作为直接入口**
- **测试阶段不再是旧单轨测试角色，而是四条明确轨道**

## 四种执行模式

### 1. Consult

适用于：

- 问题定位
- 架构解释
- 代码阅读
- 只给方案、不落代码

最小路径：

`orchestrator -> planner -> (architect / backend specialist if needed) -> answer`

说明：

- 不改文件时，可以不进入测试轨道 / `reviewer` / `sync`
- 但仍然要先明确 context、风险和是否涉及语义判断

### 2. Spec

适用于：

- 方向性 spec 起草
- execution spec 扩展
- 先做架构评估、再冻结设计治理产物

标准路径：

`orchestrator -> planner -> architect -> spec-writer -> reviewer -> sync`

说明：

- `spec-writer` 只是 drafting 角色，不是新的 truth owner
- 临时输入不进入正式 spec 执行链
- 已审批、待实施 spec 进入 `realdesk-ai/spec-craft/`
- 稳定 truth 仍然只在实现落地后进入 `realdesk-ai/specs/`
- `P0` / `P1` 新功能和任何语义变化必须先进入 `spec-craft`，再进入测试链路
- `direction.md` 必须保留 constraints / residual risks，`execution.md` 必须保留 direction linkage / risk / owner / target paths / state / config / Redis / test scope / test asset outputs / Bruno impact / gate / verification / reviewer focus，保证 backend / 四个测试轨道 / reviewer 可直接消费

### 3. Change

适用于：

- 新增功能
- 修改功能
- 修复 bug
- 改文档且影响语义表达

标准路径：

`orchestrator -> planner -> architect(if needed) -> backend/* -> test -> reviewer -> sync(if semantic change) -> ci(if git action requested)`

说明：

- 这是默认模式
- 任何代码改动至少要经过 `reviewer`
- 任何语义变更都必须经过 `sync`
- 如用户请求 `change validation`、`review readiness`、`commit`、`push`、`PR` 或 `merge`，由 `ci` 输出 `CI Record`

### 4. Release

适用于：

- 构建
- 发版
- 灰度
- 上线后巡检

标准路径：

`orchestrator -> planner -> deploy -> observability`

如果发布前还包含代码或配置变更，必须先完成 `Change` 路径，再进入发布路径。

补充约定：

- `deploy` 负责当前私有镜像站 `registry.qianpc.com`、镜像 tag、登录要求和凭据引用路径口径
- 旧私库 `183.134.217.150:5000` 已停止对外访问；涉及镜像来源时必须以上游 `doc/deploy.md` 与 `team/deploy-agent.md` 为准

## 紧凑执行原则

为了避免流程过重，小任务允许走紧凑模式，但只能做“压缩”，不能做“绕过”：

- P2 单模块小改动可以只生成 1 个 `Task Plan`
- 小型 spec 治理任务可以压缩为 1 个 `Spec Brief` + 1 个 `Execution Spec`
- 单 context 变更可以只激活 1 个 backend agent
- 纯 changed-package 补测可以只激活 `test-unit-agent`
- 无语义变化的纯格式或注释修正可以不触发 `sync`

以下步骤不得省略：

- `orchestrator`
- `planner`
- 代码改动后的所需测试轨道
- `reviewer`
- 语义变更后的 `sync`
- Git 动作请求下的 `ci`

## 标准产物

为了让流程真正可执行，所有任务统一复用 `team/templates/` 下的固定输出格式：

- `templates/routing.md`
- `templates/task-plan.md`
- `templates/test-plan.md`
- `templates/handoff.md`
- `templates/test-report.md`
- `templates/findings.md`
- `templates/sync-handoff.md`
- `templates/ci-record.md`
- `templates/spec-brief.md`
- `templates/architecture-eval.md`
- `templates/spec-acceptance.md`
- `templates/execution-spec.md`
- `templates/spec-handoff.md`
- `testing-system.md`
- `p0-unit-packages.txt`

这些产物可以出现在：

- agent 输出
- 任务记录
- PR 描述
- Review 记录

但 section 名称和字段应保持稳定，避免每次重新发明格式。

## Team-Skill 映射

- `realdesk-ai/team/*` 与 `realdesk-ai/team/templates/*` 继续作为流程真相源
- Spec 链路共享治理 skill：`.agents/skills/realdesk-spec-governance-core/SKILL.md`
- 测试共享治理 skill：`.agents/skills/realdesk-test-governance-core/SKILL.md`
- Spec 链路角色 skill：`realdesk-orchestrator`、`realdesk-planner`、`realdesk-architect`、`realdesk-spec-writer`
- 测试轨道角色 skill：`realdesk-test-unit`、`realdesk-test-api-connectivity`、`realdesk-test-e2e-scenario`、`realdesk-test-triage`
- 已有映射继续保留：`realdesk-reviewer`、`realdesk-sync`、`realdesk-deploy-agent` 与各 backend domain skill
- `observability` 本轮仍保持 team 文档入口，不强制要求 skill 映射
- 一致性检查入口：`realdesk-ai/tools/check-team-skill-alignment.sh` 与 `make check-team-skill-alignment`

## 测试基线

- `95%` 指的是 `team/p0-unit-packages.txt` 中 P0 truth-owner package 的 statement coverage hard gate，不是“全仓 95%”
- Go 原生 coverage 不提供 branch coverage；P0 分支完整性必须由行为矩阵、P0 checklist 和 reviewer 审查共同兜底
- `execution.md` 是 `P0/P1` 与语义变更任务的正式测试输入；所有测试轨道都先消费它或 `spec-test-intake` 输出
- `test-unit-agent` 负责 changed-package coverage、关键分支、异常、mock、table-driven、race、必要 benchmark
- `test-api-connectivity-agent` 负责 HTTP `httptest` / gRPC `bufconn` 的 route/method/auth/status/body/headers/idempotency/health smoke
- `test-e2e-scenario-agent` 负责 spec 驱动的完整业务场景链路、replay、失败补偿与恢复
- `test-triage-agent` 负责 flaky 归因、失败聚类、迁移/环境 smoke、人工验证缺口与未覆盖风险
- 每个测试轨道都必须输出 `Test Plan` 和 `Test Report`，并设置 `test_track`
- Bruno 第一阶段只承担资产与调试入口职责，CI 只检查资产和口径一致性，不把它写成完整自动化编排
- P0 hard gate 包清单按优先级逐批扩充；新增 P0 truth-owner package 时，必须同步更新 `team/p0-unit-packages.txt`

## 强约束

- 不允许跳过 `orchestrator`，直接从 `backend/*` 开始
- 不允许绕过 `architect` 直接让 `spec-writer` 产出正式 spec
- 不允许跳过 `planner`，直接分配实现任务
- 不允许把 `operations` 当成业务真相 owner
- 不允许把 `shared` 变成业务逻辑兜底模块
- 不允许把 `sync` 留成后续待办
- 不允许继续把活跃 baseline 写成单一旧单轨测试角色

## 推荐落地顺序

仓库内建议按五步推进：

1. 先把 `team/README.md`、`flow.md`、模板和校验脚本作为默认规范
2. 再把 `AGENTS.md`、PR 模板、Review 习惯统一到这套产物格式
3. 再补仓库内轻量自动化入口：`make check-testing-system`、`make spec-test-intake`
4. 再补 Git 动作门禁：`ci-agent` + `CI Record`
5. 最后再考虑更重的 agent wrapper 或外部编排

当前仓库已完成到第 4 步：在不新增 records 目录和 PR 正文解析 workflow 的前提下，引入 `ci-agent` 承接 Git 动作与 change validation。
