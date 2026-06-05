# 认证网关当前能力盘点

## 1. 文档目标

本文用于盘点当前这套认证网关代码库已经具备的能力边界，帮助研发、架构和接入方回答三个问题：

- 它现在能承担什么角色
- 业务系统可以怎样接入它
- 哪些能力应该由下游 `identity.xx.cn` 再承接一层

本文基于当前仓库实现整理，重点描述“已经存在的功能”，不展开未来规划。

## 2. 当前定位

从当前实现看，这套系统更接近：

- 上游统一认证入口
- 多协议认证服务器
- 用户、组织、权限管理后台
- 目录与预配接口提供方
- MCP Server 的访问与工具授权网关

如果放到统一身份体系里，它更适合放在 `oauth.xx.cn` 这一层，承担外部认证协议和上游账号接入，不建议直接作为公司内部唯一身份源。

## 3. 面向终端用户的能力

### 3.1 账号基础能力

当前已经具备以下用户侧基础流程：

- 注册
- 登录
- 登出
- SSO 登出
- 账号资料查询
- 邀请码加入
- 验证码校验
- 忘记密码与重置密码

对应接口入口主要集中在：

- `POST /api/signup`
- `POST /api/login`
- `GET,POST /api/logout`
- `GET,POST /api/sso-logout`
- `GET /api/get-account`
- `POST /api/send-verification-code`
- `POST /api/verify-code`
- `POST /api/set-password`

## 3.2 登录方式

当前支持的登录方式比较完整，既覆盖现代 Web/移动端，也兼容企业老系统：

- 用户名密码登录
- OAuth 2.0 / OIDC 第三方登录
- SAML 登录
- CAS 登录
- Device Authorization Grant 设备码登录
- Native SSO 流程
- WebAuthn / Passkey 无密码登录

说明：

- 登录完成后的返回结果可以是会话登录、授权码、隐式 token、SAML 响应或 CAS ticket
- 这意味着它既能作为登录门户，也能作为标准协议服务器

## 3.3 多因子认证

当前已确认支持的 MFA 方式有：

- TOTP 认证器应用
- 短信验证码
- 邮件验证码
- RADIUS 二次校验
- Push 审批式验证

配套能力还包括：

- Recovery Code 恢复码
- 首选 MFA 类型设置
- 组织级 MFA remember 时长控制

## 3.4 登录策略与风险控制

应用与组织层已经具备一批细粒度登录控制项：

- 应用级禁用登录
- 组织级禁用登录
- 自动登录开关
- 独占登录开关
- 允许访客登录
- 允许注册
- 失败登录次数限制
- 冻结时长
- 验证码重发间隔
- 标签过滤登录
- IP 白名单与 IP 限制
- 按入口规则控制登录方式

这说明它不只是“发 token”，而是已经具备登录策略网关属性。

## 3.5 登录页与注册页可配置

应用对象中已经支持大量界面级配置，可用于不同业务线接入时做登录体验隔离：

- 登录项配置
- 注册项配置
- 页头和页脚 HTML
- 页面自定义 HTML
- 登录表单 CSS
- 移动端 CSS
- 主题配置
- 背景图
- 表单偏移与侧边内容
- Terms of Use
- 自定义跳转链接

这意味着当前网关不仅能做协议接入，也能直接承载面向终端用户的登录页。

## 4. 面向应用接入方的能力

### 4.1 OAuth 2.0 / OIDC 授权服务器能力

当前已经具备较完整的 OAuth 2.0 / OIDC 服务端能力，支持：

- Authorization Code
- PKCE
- Resource Owner Password Credentials
- Client Credentials
- Implicit
- Refresh Token
- Device Code
- JWT Bearer Grant
- Token Exchange

同时已经具备配套标准端点：

- Token
- Refresh
- Introspection
- OpenID Discovery
- OAuth Authorization Server Metadata
- JWKS
- WebFinger
- OAuth Protected Resource Metadata

## 4.2 OAuth 扩展能力

在标准协议基础上，还实现了几项比较实用的扩展能力：

- Dynamic Client Registration
- Client Assertion
- DPoP 绑定
- Consent Grant / Revoke
- `resource` 参数校验
- 应用级 Scope 配置
- 自定义 Token Attribute

这使它更适合作为统一认证接入层，而不仅是最小化 OIDC 登录服务。

## 4.3 SAML 与 CAS 兼容

对于历史系统或企业软件接入，当前也提供：

- SAML metadata
- ACS
- SAML redirect
- CAS validate
- CAS serviceValidate
- CAS proxyValidate
- CAS proxy

因此下游旧系统不一定必须先升级成 OIDC 才能接这套网关。

## 4.4 WebAuthn / Passkey

当前已经支持完整的 WebAuthn 两阶段流程：

- 注册开始
- 注册完成
- 登录开始
- 登录完成

并且支持 discoverable login，因此可以覆盖无用户名的 passkey 登录场景。

## 5. 面向管理员的能力

### 5.1 用户与组织管理

后台已经提供较完整的 IAM 管理能力：

- Organization 管理
- User 管理
- Group 管理
- Invitation 管理
- Session 管理
- Token 管理

这说明它已经不是“单应用登录模块”，而是多租户风格的身份管理后台。

## 5.2 应用与 Provider 管理

管理员可以配置：

- Application
- Provider
- Resource
- Cert
- Key
- Site
- Form
- Rule

其中 `Application` 是接入核心对象，可配置：

- Client ID / Secret
- Redirect URI
- Grant Type
- Scope
- Token 格式
- Token 签名方式
- Token 字段
- SAML 属性
- 登录方法
- 注册项与登录项
- WebAuthn 开关

其中 `Provider` 已不只是社交登录，还覆盖多种外部能力源：

- OAuth / OIDC 身份源
- SAML 身份源
- 短信服务商
- 邮件服务商
- 通知推送服务
- 存储服务
- 实名认证服务
- 日志服务

## 5.3 权限与授权管理

当前内置 Casbin 风格授权中心能力，已可管理：

- Role
- Permission
- Model
- Adapter
- Enforcer
- Policy

并且提供：

- `/api/enforce`
- `/api/batch-enforce`

这意味着它不仅能做身份认证，也能直接承担统一授权判断入口的一部分职责。

## 5.4 运营、审计与运维面

当前后台还提供了以下运维与审计相关能力：

- Dashboard 指标
- 登录热力图
- MFA 覆盖率
- Provider 分布
- Record 日志
- Webhook 与事件重放
- Ticket
- System Info
- Health
- Metrics
- Prometheus 信息

这部分能力使其具备一定平台化运维基础。

## 6. 面向企业基础设施的能力

### 6.1 LDAP / LDAPS

系统内置 LDAP 服务端，支持：

- Bind
- Search
- 用户目录映射
- 组目录映射
- 组织目录映射

因此它可以对接依赖 LDAP 的办公软件、VPN、网络设备或老系统。

## 6.2 SCIM 2.0

系统提供 `/scim/*` 接口，可用于：

- 用户创建
- 用户读取
- 用户更新
- Patch 更新
- 用户列表查询

这让它具备作为预配接口提供方的基础能力。

## 6.3 RADIUS

系统内置 RADIUS Server，支持：

- Access-Request
- Accounting-Request
- Access-Challenge

并支持与 TOTP MFA 组合使用，适合网络接入或传统认证设备场景。

## 6.4 上游身份同步能力

从同步器实现看，当前已覆盖多种上游目录或 IAM 系统的同步能力，包括：

- 数据库
- LDAP / Active Directory
- SCIM
- Google Workspace
- Azure AD
- Okta
- Keycloak
- Lark
- WeCom
- DingTalk

因此它既能做认证网关，也能承担一部分身份汇聚职责。

## 7. 第三方身份源接入能力

当前已确认接入或显式支持的身份源包括：

- Google
- GitHub
- GitLab
- Okta
- OIDC
- Azure AD
- Azure AD B2C
- WeChat
- WeCom
- DingTalk
- Lark
- Telegram
- ADFS
- MetaMask / Web3

另外通过 Goth provider 还扩展了更多 OAuth 平台，因此整体接入面已经比较广。

### 7.1 常用 OAuth 2.0 / OIDC 第三方登录支持矩阵

下表只列出当前代码里已经明确支持、且在实际接入中最常见的一批 Provider。

| Provider 类型 | 协议 / 类型 | 代码侧支持情况 | 前端专属按钮 |
| --- | --- | --- | --- |
| `OIDC` | 通用 OIDC | 已内建 | 有 |
| `Okta` | OIDC | 已内建 | 有 |
| `ADFS` | OIDC / OAuth 风格企业登录 | 已内建 | 有 |
| `AzureAD` | OIDC | 已内建 | 有 |
| `AzureADB2C` | OIDC | 已内建 | 有 |
| `Google` | OAuth 2.0 / OIDC | 已内建 | 有 |
| `GitHub` | OAuth 2.0 | 已内建 | 有 |
| `GitLab` | OAuth 2.0 / OIDC 风格 | 已内建 | 有 |
| `Facebook` | OAuth 2.0 | 已内建 | 有 |
| `LinkedIn` | OAuth 2.0 | 已内建 | 有 |
| `QQ` | OAuth 2.0 | 已内建 | 有 |
| `WeChat` | OAuth 2.0 | 已内建，支持 PC 与 `Mobile` 子类型 | 有 |
| `WeCom` | OAuth 2.0 | 已内建，支持 `Internal` 与 `Third-party` 子类型 | 有 |
| `Lark` | OAuth 2.0 | 已内建 | 有 |
| `DingTalk` | OAuth 2.0 / OpenID 风格 | 已内建 | 有 |
| `Weibo` | OAuth 2.0 | 已内建 | 有 |
| `Gitee` | OAuth 2.0 | 已内建 | 有 |
| `Baidu` | OAuth 2.0 | 已内建 | 有 |
| `Alipay` | OAuth 2.0 | 已内建 | 有 |
| `Infoflow` | OAuth 2.0 | 已内建，支持 `Internal` 与 `Third-party` 子类型 | 有 |
| `Apple` | OAuth 2.0 / OIDC 风格 | 通过 Goth 支持 | 有 |
| `Slack` | OAuth 2.0 | 通过 Goth 支持 | 有 |
| `Steam` | OpenID / OAuth 兼容接入 | 通过 Goth 支持 | 有 |
| `Bilibili` | OAuth 2.0 | 已内建 | 有 |
| `Douyin` | OAuth 2.0 | 已内建 | 有 |
| `Kwai` | OAuth 2.0 | 已内建 | 有 |
| `Telegram` | Telegram Login | 已内建 | 通用按钮 |
| `Twitter` | OAuth 2.0 | 已内建 | 通用按钮 |
| `Custom` / `Custom Flexible` | 自定义 OAuth / OIDC | 已内建 | 通用按钮 |

补充说明：

- 你刚刚提到的 `微信`、`QQ`、`企业微信`、`飞书`、`钉钉`，当前代码里都是明确支持的，不是只能通过通用 OIDC 间接接。
- `WeChat` 还区分了 PC 扫码和 `Mobile` 子类型。
- `WeCom` 明确区分 `Internal` 和 `Third-party` 两种模式。
- `OIDC`、`Custom`、`Custom Flexible` 这三类更适合接企业内部或自建身份源。

### 7.2 通过 Goth 扩展支持的更多 Provider

除上表外，当前还通过 Goth 适配层支持更多 OAuth Provider，包括：

- `Amazon`
- `Auth0`
- `BattleNet`
- `Bitbucket`
- `Box`
- `CloudFoundry`
- `Dailymotion`
- `Deezer`
- `DigitalOcean`
- `Discord`
- `Dropbox`
- `EveOnline`
- `Fitbit`
- `Gitea`
- `Heroku`
- `InfluxCloud`
- `Instagram`
- `Intercom`
- `Kakao`
- `Lastfm`
- `Line`
- `Mailru`
- `Meetup`
- `MicrosoftOnline`
- `Naver`
- `Nextcloud`
- `OneDrive`
- `Oura`
- `Patreon`
- `Paypal`
- `SalesForce`
- `Shopify`
- `Soundcloud`
- `Spotify`
- `Strava`
- `Stripe`
- `TikTok`
- `Tumblr`
- `Twitch`
- `Typetalk`
- `Uber`
- `VK`
- `Wepay`
- `Xero`
- `Yahoo`
- `Yammer`
- `Yandex`
- `Zoom`

这批 Provider 大多可以直接作为 OAuth 登录源接入，但前端不一定都提供独立样式按钮；没有专属按钮的场景，仍可通过通用按钮和统一授权跳转逻辑使用。

## 8. MCP Gateway 能力

这是当前代码库比较特别的一层能力。

### 8.1 MCP Server 管理

后台已经提供 MCP 相关对象管理能力：

- Server 列表与详情
- Server 新增、更新、删除
- 在线 Server 查询
- MCP Tool 同步
- MCP Access Token 获取

## 8.2 MCP 代理转发

当前支持通过统一入口将请求代理到上游 MCP Server：

- `GET /api/server/:owner/:name`
- `POST /api/server/:owner/:name`

代理时具备以下特性：

- 校验上游 Server URL
- 可自动注入 Bearer Token
- 对 `tools/call` 执行工具白名单检查
- 工具不在 allowlist 时会直接拒绝

这说明 MCP 这一层不是简单反代，而是带了工具级授权控制。

## 8.3 MCP Tool 与 Scope 的关系

应用模型里的 Scope 已经支持配置允许访问的 MCP Tools，因此从设计上看，当前系统具备把：

- OAuth Scope
- 应用授权
- MCP Tool 访问权限

三者挂到一起的基础模型。

## 8.4 自带 MCP 入口

系统本身还暴露了：

- `POST /api/mcp`

说明它不仅能代理第三方 MCP Server，也已经预留了自身作为 MCP 入口的能力。

## 9. 在统一身份体系中的推荐边界

如果按 `oauth.xx.cn -> identity.xx.cn -> product` 的分层来放置，建议边界如下：

### 9.1 适合由当前认证网关承担的职责

- 第三方身份源接入
- OIDC / OAuth / SAML / CAS 协议适配
- 终端用户登录页
- MFA 与 Passkey
- 上游认证结果签发
- MCP Server 访问网关
- LDAP / SCIM / RADIUS 兼容接口

### 9.2 更适合由下游 Identity Service 承担的职责

- 公司唯一用户主数据
- 第三方账号绑定后的统一账号归并
- 内部租户、产品、员工、客户统一身份视图
- 公司内部统一 JWT
- 产品级权限聚合
- 内部会话生命周期
- 面向业务域的资料与权限查询 API

## 10. 总结

当前这套系统已经具备“统一登录入口 + 多协议认证服务器 + 多因子认证 + 用户与权限后台 + 目录与预配接口 + MCP 工具网关”的完整雏形。

如果仅从“现在能不能用”来看，它已经足以承担：

- 外部认证入口
- 协议兼容层
- 应用接入登录中心
- MCP 访问与工具授权入口

如果从“公司统一身份底座”来看，它更适合作为上游认证层，而不是最终唯一身份源。
