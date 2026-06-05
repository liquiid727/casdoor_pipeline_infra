## Task Plan
total_tasks: 3
risk_summary: P0:0 P1:2 P2:1
parallel_groups:
  group_1: [T001, T002]
  group_2: [T003]
tasks:
  - task_id: T001
    context: team-governance
    layer: governance
    description: 引入 ci-agent/team 模板到当前仓库，补齐 change-validation 最小控制面
    truth_owner: realdesk-ai/team
    runtime_owner: repository workflow
    dependencies: []
    can_parallel: true
    risk_level: P2
    risk_reason: 仅新增流程资产，不改变运行时行为
  - task_id: T002
    context: iam, deployment, i18n, web-ui
    layer: infrastructure
    description: 运行构建与测试，修复阻断后端构建和 i18n 去重的确定性问题
    truth_owner: deployment, i18n, web/src
    runtime_owner: backend and frontend build pipeline
    dependencies: []
    can_parallel: true
    risk_level: P1
    risk_reason: 涉及仓库级 build/test 门禁与多模块翻译引用
  - task_id: T003
    context: git-delivery
    layer: governance
    description: 在专用分支上整理 stage 范围、提交 CI 记录并推送远端
    truth_owner: git history
    runtime_owner: origin
    dependencies: [T001, T002]
    can_parallel: false
    risk_level: P1
    risk_reason: 当前 master 落后 origin/master 1 个提交，直接推 master 有冲突风险
