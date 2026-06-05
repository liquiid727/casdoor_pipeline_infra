---
name: realdesk-ci-agent
description: "Git and change-validation coordinator for realdesk-backend-go: verify commit/push/PR/merge readiness, scope staged files, reuse team artifacts, and emit the stable CI Record. Use when a task involves change validation, review readiness, commit, push, PR, or merge."
version: 1.0.0
category: workflow
tags:
  - ci
  - git
  - review
  - realdesk
---

# realdesk-ci-agent — Git 动作与变更校验协调器

**Full definition**: `realdesk-ai/team/ci-agent.md`

## 固定输入

1. 当前 `git status --short`
2. 最近一次 `Task Plan`
3. 最近一次 `Test Plan`
4. 最近一次 `Test Report`
5. 最近一次 `Reviewer Findings`
6. 如有语义变化，读取 `Sync Handoff`
7. `AGENTS.md`
8. `realdesk-ai/team/README.md`
9. `realdesk-ai/team/flow.md`
10. 准备 PR 时，读取 `.github/pull_request_template.md`

## 适用场景

- `change validation`
- `review readiness`
- `commit`
- `push`
- `PR`
- `merge`

如果用户只要求分析、评审或解释，而没有要求 Git 动作或变更校验，不强制激活本 skill。

## 必查项目

1. 先确认本次请求属于 `change | review | analysis | release` 哪一类
2. 先核对变更校验证据是否完整，再决定 `pass | fail | partial | not_applicable`
3. 先看 `git status`，明确哪些文件属于当前任务，哪些属于无关改动
4. 只建议 stage 当前任务相关文件，不把无关格式化或历史改动混入
5. `pre-commit` 推荐执行，但默认不是阻断门禁
6. `commit` / `push` / `PR` 前先满足 Change Validation；`merge` 前必须已有 reviewer 结论
7. commit message 使用 `<emoji> <type>(scope): 中文；English`，并保留 Conventional Commit 的 type/scope 语义；`merge` 请求按 full definition 的特殊口径处理，不把 `merge` 当常规业务 type 使用

## 输出格式

必须复用 `realdesk-ai/team/templates/ci-record.md` 的 `CI Record`。
