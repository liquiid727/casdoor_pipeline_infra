# 认证网关 Roadmap 与并行任务拆分

## 0. 当前实现状态

截至当前代码状态，roadmap 已落地以下能力：

- P0 已完成：callback scheme、signed state、相对回跳限制、安全 cookie、token invalid 清 cookie 并重登、Inactive 主链路拦截、纯内存单测基线。
- P1 已完成第一版：新增 `gateway/config`、`gateway/core`，并让 `service/oauth.go` 的登录跳转复用 `gateway/core.OAuthFlow`。
- P2 已完成第一版：新增 `sdk/httpauth` middleware，支持 `RequireLogin`、`CurrentUser`、`RequireScope`、`RequireRole`。
- T8/P3 已完成基础能力：新增 `gateway/standalone` 身份头覆盖注入 helper，新增网关 Prometheus 指标。

仍未完成或只完成基础版本的内容：

- `service/proxy.go` 尚未完全迁移到 `gateway/core` 的 AuthDecision / Principal 模型。
- `gateway/standalone` 尚未成为完整独立网关，当前仍以 helper 形式提供身份头注入。
- `sdk/httpauth` 尚未实现 JWKS 缓存、introspection 兜底和跨语言 SDK。
- 认证审计字段、trace id 全链路透传、按 Site 灰度开关、单 Site 回滚策略和 E2E 回归仍待补齐。

## 1. 目标

本文用于指导后续多线程并行推进认证网关能力建设。目标形态是：

- 独立认证网关：保护 Web、老系统、私有化产品入口。
- 嵌入式 SDK / Middleware：服务 API、微服务、产品内部服务。
- 统一 Gateway Core：复用 OAuth 跳转、callback、state/session、token 校验、用户上下文和错误模型。

当前代码的主要入口是：

- `service/proxy.go`：Host 路由、认证拦截、规则、反代。
- `service/oauth.go`：OAuth 登录跳转和 `/caswaf-handler` callback。
- `service/util.go`：域名、证书、认证客户端、scheme 等辅助逻辑。
- `object/site.go` / `object/site_cache.go`：Site 控制面模型和运行期缓存。

## 2. 总体路线

### P0：现有 Site 认证网关闭环加固

目标：不做大重构，先让当前闭环达到内测和后续抽象的基础安全线。

必须解决：

- HTTPS / Ingress TLS 终止场景 callback scheme 正确。
- OAuth `state` 不能直接信任原始路径。
- callback 回跳不能跳到外部 URL。
- 认证 cookie 具备 `HttpOnly`、`Secure`、`SameSite` 和清理策略。
- token 解析失败不能返回 500，应清 cookie 并重新登录。
- `Site.Status == "Inactive"` 时主请求链路直接拒绝。
- 最小单元测试覆盖当前闭环。

### P1：冻结核心契约并抽 Gateway Core

目标：把认证核心能力从反代流程中抽出来，避免后续独立网关和 SDK 复制逻辑。

建议新增：

```text
gateway/
  config/
  core/
  standalone/
```

核心契约：

- `SiteAuthConfig`
- `AuthSession`
- `Principal`
- `AuthDecision`
- `AuthError`

核心能力：

- `BeginAuth`
- `HandleCallback`
- `VerifyAccessToken`
- `BuildPrincipal`
- `ClearSession`

### P2：Standalone Gateway 与 Embedded SDK 双形态

目标：同一套 core 同时支持独立部署和业务内嵌。

新增能力：

- `gateway/standalone`：Host 路由、TLS、反代、规则、身份头、审计。
- `sdk/httpauth`：Go HTTP middleware，提供 token 校验和用户上下文。

### P3：生产化运维与灰度

目标：具备企业级上线、灰度、排障和回滚能力。

新增能力：

- 认证审计字段
- Prometheus 指标
- trace id 透传
- 按 Site 灰度开关
- 单 Site 回滚策略
- 最小 E2E 回归

## 3. 并行线程总览

| 线程 | 阶段 | Owner | 主要文件 | 是否可并行 | 说明 |
| --- | --- | --- | --- | --- | --- |
| T0 | 准备 | 文档 / 集成 | `docs/identity/*` | 先做 | 先提交现有文档，清理工作区 |
| T1 | P0 | OAuth 安全 | `service/oauth.go`, `service/util.go`, `service/oauth_state.go` | 可与 T2/T3 局部并行 | 负责 scheme、state、callback |
| T2 | P0 | Cookie / token | `service/oauth.go`, `service/proxy.go`, `service/cookie.go` | 与 T1 有冲突风险 | 建议 T1 合并后再最终落地 |
| T3 | P0 | Site 路由状态 | `service/proxy.go`, `object/site_timer_health.go` | 可独立推进 | 负责 Inactive 和路由兼容 |
| T4 | P0 | 测试基线 | `service/*_test.go`, `object/site*_test.go` | 可并行 | 尽量只写测试和 helper |
| T5 | P1 | Config/Core 架构 | `gateway/config/*`, `gateway/core/*` | P0 后开始 | 冻结核心契约 |
| T6 | P1 | Service 接 core | `service/*.go`, `gateway/standalone/*` | 依赖 T5 | 让 service 变薄 |
| T7 | P2 | Embedded SDK | `sdk/httpauth/*` | 依赖 T5 | 不碰反代和 callback |
| T8 | P3 | QA / Observability | `routers/*`, `object/record.go`, docs | 可与 P2 并行 | 指标、审计、灰度 |

## 4. P0 详细任务

### T1：OAuth callback、scheme、state 安全

职责：

- 修正 `getScheme()`。
- 为 OAuth 登录跳转生成 signed state。
- callback 校验 signed state。
- 限制 callback 回跳路径。

建议文件：

- 修改：`service/oauth.go`
- 修改：`service/util.go`
- 新增：`service/oauth_state.go`
- 新增：`service/oauth_state_test.go`

实现要求：

- scheme 判断顺序：
  1. `X-Forwarded-Proto`
  2. `X-Forwarded-Scheme`
  3. `r.TLS != nil`
  4. `http`
- state 内容包含：
  - `path`
  - `nonce`
  - `iat`
  - `exp`
- state 签名使用当前 Site 绑定 Application 的 `ClientSecret`。
- callback 只允许回跳以 `/` 开头的相对路径。
- state 过期或签名失败时返回 `400` 或 `401`。

测试：

- `TestGetScheme_UsesForwardedProto`
- `TestGetScheme_UsesTLS`
- `TestOAuthState_RoundTrip`
- `TestOAuthState_RejectsTamperedPayload`
- `TestOAuthState_RejectsExpiredPayload`
- `TestOAuthState_RejectsExternalRedirect`

验收：

```bash
go test ./service -run 'TestGetScheme|TestOAuthState' -count=1
```

### T2：认证 cookie 和 token 失败重登

职责：

- 封装认证 cookie 设置和清理。
- token 解析失败时清 cookie 并重新登录。
- 统一认证失败行为，避免 500。

建议文件：

- 修改：`service/oauth.go`
- 修改：`service/proxy.go`
- 新增：`service/cookie.go`
- 新增：`service/cookie_test.go`

实现要求：

- `pipeline_auth_access_token` 必须设置：
  - `Path=/`
  - `HttpOnly=true`
  - `SameSite=Lax`
  - `Secure` 根据请求 scheme 判断
- 新增 `clearAuthCookie()`，使用 `MaxAge=-1`。
- token invalid / expired：
  - 清 cookie
  - 重新发起 OAuth 登录跳转
  - 不返回 500

测试：

- `TestSetAuthCookie_HasSecurityAttributes`
- `TestClearAuthCookie_ExpiresCookie`
- `TestHandleRequest_InvalidTokenClearsCookieAndRedirects`

验收：

```bash
go test ./service -run 'Test.*Cookie|TestHandleRequest_InvalidToken' -count=1
```

### T3：Site 状态与路由兼容

职责：

- 让 `Inactive` Site 在主请求链路中生效。
- 保留 ACME、`MP_verify_`、`www`、`needRedirect`、`HTTPS Only` 的现有行为。

建议文件：

- 修改：`service/proxy.go`
- 可选修改：`object/site_timer_health.go`
- 新增：`service/proxy_test.go`

实现要求：

- `Site.Status == "Inactive"` 时直接拒绝，建议返回 `503`。
- ACME challenge 和 `MP_verify_` 应保持可用，避免证书续期和平台验证受影响。
- 不产生重定向环。

测试：

- `TestHandleRequest_InactiveSiteReturns503`
- `TestHandleRequest_AcmeChallengeStillWorks`
- `TestHandleRequest_MpVerifyStillWorks`
- `TestHandleRequest_WwwRedirect`
- `TestHandleRequest_HTTPSOnlyRedirect`

验收：

```bash
go test ./service -run 'TestHandleRequest_(Inactive|Acme|Mp|Www|HTTPS)' -count=1
```

### T4：测试基线与 fixture

职责：

- 先补不依赖真实 DB、真实网络、真实 OAuth 服务的单测。
- 给后续 P1/P2 提供稳定 fixture。

建议文件：

- 新增：`service/test_helpers_test.go`
- 新增：`object/site_test.go`
- 新增：`object/site_cache_test.go`
- 新增：`object/test_helpers_site_test.go`
- 新增：`service/testdata/static-site/index.html`

测试范围：

- `joinPath`
- `getHostNonWww`
- `getDomainWithoutPort`
- `Site.GetChallengeMap`
- `Site.GetHost`
- `GetSiteByDomain` 大小写
- `otherDomains` 命中
- `ApplicationObj` 绑定

验收：

```bash
go test ./service ./object -run 'Test(JoinPath|GetHost|GetDomain|Site_|RefreshSiteMap|GetSiteByDomain)' -count=1
```

## 5. P1 详细任务

### T5：新增 Gateway Config/Core

职责：

- 定义运行期认证网关配置快照。
- 定义认证核心契约。
- 不依赖反向代理。

建议文件：

```text
gateway/config/site_config.go
gateway/config/resolver.go
gateway/core/session.go
gateway/core/state.go
gateway/core/token.go
gateway/core/principal.go
gateway/core/errors.go
gateway/core/decision.go
```

核心结构：

- `SiteAuthConfig`
- `UpstreamConfig`
- `AuthConfig`
- `Principal`
- `AuthSession`
- `AuthDecision`

测试：

- `authgateway/core/site_resolver_test.go`
- `authgateway/core/oauth_flow_test.go`
- `authgateway/core/decision_test.go`

验收：

```bash
go test ./gateway/config ./gateway/core -count=1
```

### T6：service 接入 Gateway Core

职责：

- 保持外部行为不变。
- 将 `service/oauth.go` 和 `service/proxy.go` 中的认证决策逐步迁移到 core。
- `service` 只保留 HTTP glue、proxy、TLS、启动装配。

建议文件：

- 修改：`service/oauth.go`
- 修改：`service/proxy.go`
- 修改：`service/util.go`
- 新增：`gateway/standalone/handler.go`

验收：

```bash
go test ./gateway/config ./gateway/core ./service -count=1
make backend
```

## 6. P2 详细任务

### T7：Embedded SDK / Middleware

职责：

- 为业务服务提供内嵌认证能力。
- 不处理浏览器 callback，不做代理。

建议文件：

```text
sdk/httpauth/options.go
sdk/httpauth/middleware.go
sdk/httpauth/context.go
sdk/httpauth/errors.go
sdk/httpauth/middleware_test.go
sdk/httpauth/example_integration_test.go
```

能力：

- `VerifyToken`
- `RequireLogin`
- `CurrentUser`
- `RequireScope`
- `RequireRole`

验收：

```bash
go test ./sdk/httpauth -count=1
```

### T8：Standalone Gateway 增强

职责：

- 身份头透传。
- 覆盖客户端伪造同名头。
- 标准化未认证 / 未授权响应。

建议身份头：

- `X-Auth-User`
- `X-Auth-Owner`
- `X-Auth-Email`
- `X-Auth-Scope`
- `X-Trace-Id`

验收：

```bash
go test ./gateway/standalone ./service -count=1
```

## 7. P3 详细任务

### T9：审计、指标、灰度

职责：

- 网关认证链路可观测。
- 支持按 Site 灰度和回滚。

指标建议：

- `gateway_login_start_total`
- `gateway_callback_success_total`
- `gateway_callback_failure_total`
- `gateway_token_invalid_total`
- `gateway_inactive_site_block_total`
- `gateway_upstream_error_total`

日志字段：

- `site`
- `domain`
- `authApplication`
- `principal`
- `action`
- `status`
- `reason`
- `upstream`
- `duration`
- `traceId`

验收：

```bash
go test ./service ./routers ./object -count=1
```

## 8. QA Gate

### P0 Gate

- HTTP、HTTPS、Ingress TLS 终止下 callback URL 正确。
- 篡改 state、过期 state、外部回跳 URL 必须失败。
- cookie 安全属性齐全。
- token invalid 不返回 500。
- Inactive Site 不进入 OAuth 和反代主链路。
- ACME / `MP_verify_` 不被误伤。

### P1 Gate

- `gateway/core` 有独立单测。
- `service` 迁移后外部行为不变。
- core 不依赖反向代理。
- `Principal`、`AuthSession`、`SiteAuthConfig` 契约稳定。

### P2 Gate

- Standalone Gateway 与 SDK 解析同一 token 得到一致 `Principal`。
- SDK 不依赖全局 `SiteMap`。
- 身份头由网关注入，并覆盖客户端伪造值。
- 401 / 403 语义稳定。

### P3 Gate

- 指标能区分 login、callback、token invalid、inactive block、upstream error。
- 日志可定位 site、user、reason、trace。
- 支持单 Site 灰度和回滚演练。

## 9. 合并顺序

推荐合并顺序：

1. T0：文档基线。
2. T4：测试 helper 和基础单测。
3. T1：OAuth state / scheme。
4. T2：cookie / token 失败重登。
5. T3：Inactive Site / 路由兼容。
6. P0 全量验收。
7. T5：Gateway Config/Core。
8. T6：service 接 core。
9. T7 / T8：SDK 与 Standalone Gateway 并行。
10. T9：审计、指标、灰度。

## 10. 不建议现在做的事

- 不要一开始就做多语言 SDK。
- 不要现在重做公司唯一 Identity 主数据平台。
- 不要先拆数据库、KMS、Secret 子系统。
- 不要在 `Principal` 契约未定前做身份头透传。
- 不要让 SDK 线程修改 `object/site_cache.go`。
- 不要在抽 core 的同时继续往 `service/proxy.go` 塞新认证逻辑。

## 11. 总体验收命令

P0：

```bash
go test ./service ./object -count=1
make backend
```

P1：

```bash
go test ./gateway/config ./gateway/core ./service -count=1
make backend
```

P2：

```bash
go test ./gateway/... ./sdk/httpauth ./service -count=1
make backend
```

P3：

```bash
go test ./gateway/... ./sdk/httpauth ./service ./routers ./object -count=1
make backend
```
