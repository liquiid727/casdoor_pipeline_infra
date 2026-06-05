# Casdoor 剩余真实依赖面盘点

这份清单只统计还会阻塞“从 Casdoor 运行时抽离到自有 provider 运行时”的真实依赖面，不把版权头、历史注释、纯品牌文案噪音混进来。

当前结论按四块收敛：`配置`、`OAuth`、`短信`、`文档/外显命名`。

## 1. 配置

### 1.1 Site 站点代理配置仍保留 Casdoor 语义

- `object/site.go`
  - `AuthApplication` 的数据库列名仍是 `casdoor_application`
  - JSON 兼容入口仍保留 `casdoorApplication`
  - 这是站点级“受保护应用”配置的主入口，后续如果要改成自有 provider/runtime 语义，这里是第一刀
- `web/src/SiteEditPage.js`
  - 站点编辑页字段标题仍是 `Casdoor app`
- `web/src/SiteListPage.js`
  - 旧列表列已被注释，但字段名仍是 `authApplication`

建议替换入口：

- 先保留存量兼容，新增中性语义字段名，例如 `providerApplication` 或 `authRuntimeApplication`
- 在 `object/site.go` 里统一做双向兼容：
  - DB: 先兼容旧列 `casdoor_application`
  - JSON: 继续兼容 `casdoorApplication`
  - 内部运行时: 逐步只使用新字段名

### 1.2 自举会话常量仍绑定 Casdoor built-in 语义

- `object/session.go`
  - `AuthApplication = "app-built-in"`
  - `CasdoorOrganization = "built-in"`
  - 删除 session、清理 Beego session 时仍用这两个常量判断“系统自身会话”
- `controllers/account.go`
  - 注销与踢 session 仍用 `object.AuthApplication`
- `object/user.go`
  - 用户删除 session 也仍用 `AuthApplication`

建议替换入口：

- 先把“系统内建 runtime 会话”从 Casdoor 命名中抽成独立常量组
- 例如集中到 `object/session.go`：
  - `BuiltinRuntimeApplication`
  - `BuiltinRuntimeOrganization`

### 1.3 Site OAuth 代理 client 仍直接从本机 Casdoor 运行时配置构造

- `service/util.go`
  - `getAuthServerClientFromSite()` 仍直接返回 `*casdoorsdk.Client`
  - endpoint 仍从 `conf.GetConfigString("origin")` 取值
  - 空值 fallback 仍是 `http://localhost:8000`
- `service/proxy.go`
  - `site.AuthApplication != ""` 就进入 OAuth 代理保护逻辑
  - cookie 读取先看 `pipeline_auth_access_token`，再兼容 `casdoor_access_token`

建议替换入口：

- 先在 `service/util.go` 外提一层中性 facade，例如：
  - `getProviderRuntimeClientFromSite()`
  - `ProviderRuntimeClient.ExchangeCode(...)`
  - `ProviderRuntimeClient.ParseAccessToken(...)`
- 这样后续可以先替 SDK，再替字段命名，不必一次动完

### 1.4 built-in 运行时默认值仍写死在配置层

- `conf/web_config.go`
  - `DefaultApplication` 为空时 fallback 仍是 `app-built-in`
- `conf/app.conf`
  - `defaultApplication = "app-built-in"`
  - `radiusDefaultOrganization = "built-in"`
- `object/site_cache.go`
  - `getCasdoorCertMap()` 函数名和局部变量名仍保留 Casdoor 语义

建议替换入口：

- 这批值最好和 `object/session.go` 的 built-in runtime 常量一起收敛
- 优先把“系统默认应用 / 默认组织 / 默认证书映射”抽成中性 runtime 常量，避免配置层继续反向灌入 Casdoor 语义

### 1.5 Provider 配置目录里仍保留 Casdoor 具体类型

这部分不属于文案，而是配置模型本身仍把 `Casdoor` 当成一个具体 provider type。

- `storage/storage.go`
  - `case "Casdoor"` 仍是存储 provider 工厂分支
- `storage/casdoor.go`
  - 仍存在 `NewCasdoorStorageProvider(...)`
  - 直接依赖 `github.com/casdoor/oss/casdoor`
- `object/storage.go`
  - `provider.Type == "Casdoor"` 时走特殊证书加载和上传 URL 逻辑
- `deployment/deploy.go`
  - 部署逻辑仍针对 `provider.Type == "Casdoor"` 做存储分支
- `web/src/Setting.js`
  - Storage provider 枚举里仍有 `Casdoor`
  - Log provider 枚举里仍有 `Casdoor Permission Log`
- `web/src/provider/StorageProviderFields.js`
  - 存储表单字段按 `provider.type === "Casdoor"` 做专门分支
- `object/provider.go`
  - `provider.Type == "Casdoor Permission Log"` 仍有日志 provider 特殊实现

建议替换入口：

- 这块可以归到“配置面”优先级中高的位置
- 如果你的目标是把系统抽成“自有 provider runtime”，那这些具体 provider type 至少要先变成 legacy alias，而不是继续作为主命名暴露在配置模型和 UI 里

### 1.6 旧字段兼容还散落在同步器映射配置

- `object/syncer.go`
  - `TableColumn.UnmarshalJSON()` 仍兼容旧字段 `casdoorName`
- `web/src/table/SyncerTableColumnTable.js`
  - UI 标题仍是 `Casdoor column`
- `object/syncer_user.go`
  - 内部辅助函数名仍是 `getCasdoorColumns()`

建议替换入口：

- 保留 `casdoorName` 只做反序列化兼容
- 内部和 UI 全部切到 `targetField` / `userField` 一类中性命名

### 1.7 运行时默认展示名仍回退到 Casdoor

- `object/mfa_totp.go`
  - TOTP 初始化时，未显式传入 issuer 会 fallback 为 `Casdoor`

这不是单纯文案，而是会直接出现在用户的认证器 App 里。

建议替换入口：

- 把默认 issuer 提升成 runtime 常量或配置项
- 最好跟 built-in runtime 命名一起收敛，避免迁移后 MFA 仍对外显示 Casdoor

## 2. OAuth

### 2.1 站点代理登录链仍直接以 Casdoor OAuth server 为假设

- `service/oauth.go`
  - `getSignInURL()` 直接拼 `/login/oauth/authorize`
  - callback 固定为 `/caswaf-handler`
  - token 交换直接调用 `object.GetAuthorizationCodeToken(...)`
  - 成功后写入 `pipeline_auth_access_token`
- `service/proxy.go`
  - 受保护站点的准入判断依赖 `ParseJwtToken()`
- `service/util.go`
  - 这里仍把 SDK client 构造细节暴露给上层逻辑

建议替换入口：

- 第一批不要先改流程，先把“授权 URL 生成 / code 交换 / token 校验”包进抽象层
- 最合适的切口就是：
  - `service/oauth.go`
  - `service/proxy.go`
  - `service/util.go`

### 2.2 代码库里仍保留一个“Casdoor 作为上游 IdP”的真实实现

- `idp/casdoor.go`
  - `CasdoorIdProvider`
  - `NewCasdoorIdProvider(...)`
  - token URL 固定 `/api/login/oauth/access_token`
  - userinfo 固定 `/api/userinfo`
- `idp/provider.go`
  - factory dispatch 里仍有 `case "Casdoor"`
- `object/provider.go`
  - OAuth provider 归一化逻辑仍把 `provider.Type == "Casdoor"` 当成一级分支
- `object/user.go`
  - 用户账号绑定 schema 里仍保留 `Casdoor string 'json:"casdoor"'`

这不是单纯文案，而是一个真实的上游 provider 实现。

建议替换入口：

- 如果后续“自有 provider runtime”不再把 Casdoor 当作独立上游类型，这个文件要么被泛化成通用 OIDC provider，要么被收敛成 legacy adapter

### 2.3 前端 OAuth provider 类型和按钮入口仍有 Casdoor 专名

- `web/src/auth/Provider.js`
  - provider registry 里仍有 `Casdoor`
  - 授权 URL 构造仍单独对 `provider.type === "Casdoor"` 分支
- `web/src/auth/ProviderButton.js`
  - `provider.type === "Casdoor"` 时走 `CasdoorLoginButton`
- `web/src/auth/CasdoorLoginButton.js`
  - 按钮组件名、文案、图标资源都还是 Casdoor
- `web/public/ProviderHintRedirect.js`
  - 静态脚本也仍保留 `Casdoor` provider registry
  - 并且会为 `provider.type === "Casdoor"` 直接拼 `/login/oauth/authorize`

建议替换入口：

- 这是前端第一批最适合 rename 的入口
- 可以先改成中性组件，例如 `UpstreamOidcLoginButton` 或 `ProviderRuntimeLoginButton`
- 然后把 `provider.type === "Casdoor"` 收敛成 legacy alias

### 2.4 前端 callback 兼容层仍带 Casdoor 旧协议痕迹

- `web/src/auth/AuthCallback.js`
  - 仍兼容 `__casdoor_callback_react`
  - 仍兼容 `casdoor_callback_react_fallback`
  - 注释和成功提示里仍有 Casdoor 自身语义
- `web/public/AuthCallbackHandler.js`
  - 静态 callback handler 成功提示里仍是 “apps protected by Casdoor”
- `routers/lightweight_auth_filter.go`
  - 服务端 callback HTML 注入逻辑仍同时识别
    - `__pipeline_auth_callback_react`
    - `__casdoor_callback_react`

建议替换入口：

- 先保留旧 query key 的兼容读取
- 但把内部主路径收敛到 `pipeline_auth_*` 或新的 provider runtime key
- 这会让后续彻底删 legacy key 更容易

### 2.5 MCP / OAuth 元数据 header 仍暴露 Casdoor realm

- `routers/base.go`
  - `WWW-Authenticate: Bearer realm="casdoor", resource_metadata="..."`

这不是 UI 文案，而是对外协议面。

建议替换入口：

- 把 realm 提升成可配置值或 runtime 常量
- 这里最好和 OAuth facade 一起处理，避免外部客户端继续把系统识别成 Casdoor

## 3. 短信

### 3.1 短信发送主干仍直接绑定 `github.com/casdoor/go-sms-sender`

- `go.mod`
  - 仍直接依赖 `github.com/casdoor/go-sms-sender`
- `object/sms.go`
  - `getSmsClient()` / `SendSms()` 的标准 provider 分发全部建立在 `go-sms-sender` 的 type 常量和 client 接口上
  - 只有 `Custom HTTP SMS` 和 `Alibaba Cloud PNVS SMS` 是本仓自定义分支

这块是“真实依赖面”，因为不只是 import 名称，还决定了 provider type、参数位序和发送协议。

建议替换入口：

- 第一批不要先替所有短信 provider
- 先抽一个中性网关，例如：
  - `type SmsGateway interface { Send(provider *Provider, content string, phoneNumbers ...string) error }`
- 让 `object.SendSms()` 只做 runtime 入口，现有 `go-sms-sender` 逻辑先退到 adapter 层

### 3.2 验证码和 MFA 都直接压在当前短信入口上

- `object/verification.go`
  - `SendVerificationCodeToPhone()` 直接调用 `SendSms(provider, code, dest)`
- `controllers/verification.go`
  - 应用级 provider 选择后直接进入短信验证码发送
- `controllers/service.go`
  - `/api/send-sms` 测试入口直接调用 `object.SendSms(...)`
- `object/mfa.go` / `web/src/auth/MfaSetupPage.js`
  - `sms` 仍是 MFA 方式之一

建议替换入口：

- 短信抽象的单一总入口就是 `object.SendSms()`
- 只要这个函数先被网关化，验证码、MFA、测试发送三条链都会一起受益

### 3.3 前端 provider 类型目录仍和短信库实现细节强绑定

- `web/src/Setting.js`
  - SMS provider 枚举和展示入口都在这里
- `web/src/ProviderEditPage.js`
  - 表单字段差异按 provider type 硬编码分支
- `web/src/provider/SmsProviderFields.js`
  - 哪些 provider 没有 `signName` / `templateCode` 也按具体 type 常量写死

建议替换入口：

- 如果后续只是“去 Casdoor 运行时化”，这层可以先不大改
- 如果目标是“自有 provider runtime”，那短信能力描述最好逐步从“按 type 分支”改为“按 capability 分支”

## 4. 文档 / 外显命名

### 4.1 架构文档仍把 Casdoor 当成正式运行时角色

- `docs/identity-and-payment-architecture.md`
  - 架构图、流程、职责和落地顺序都仍以 `Casdoor -> Identity Service` 为主叙述
  - 这份文档不是噪音，它会直接影响后续实现入口命名和团队认知

建议替换入口：

- 这份文档应先从“Casdoor 是正式运行时角色”改成“上游认证 provider / legacy auth runtime”
- 否则代码已抽离、文档仍把 Casdoor 当核心角色，后续命名会继续回流

### 4.2 README 仍混有大量 Casdoor 对外叙述

- `README.md`
  - 已经加入 “independent fork / self-contained” 的新说明
  - 但主体仍保留大量 Casdoor 产品叙述、链接、badge 和版权落款

这块不一定要第一批就改实现，但它属于“文档依赖面”，因为会持续给新入口命名施加反向约束。

### 4.3 前端外显命名仍有成体系的 Casdoor 文案

高频入口包括：

- `web/src/SiteEditPage.js`: `Casdoor app`
- `web/src/table/SyncerTableColumnTable.js`: `Casdoor column`
- `web/src/TourConfig.js`: `About Casdoor`
- `web/src/auth/AuthCallback.js`
  - CAS 登录成功提示仍是 `Now you can visit apps protected by Casdoor.`
- `web/src/auth/LoginPage.js`
  - CAS 登录成功提示仍是 `Now you can visit apps protected by Casdoor.`
- `web/public/AuthCallbackHandler.js`
  - 轻量 callback handler 成功提示仍是 `Now you can visit apps protected by Casdoor.`
- `web/src/ProviderEditPage.js`
  - 新建 Email provider 默认 title 仍是 `Casdoor Verification Code`
  - 新建 Log provider 默认 type 仍是 `Casdoor Permission Log`
- `web/src/provider/EmailProviderFields.js`
  - 默认验证邮件文案仍是 `You have requested a verification code at Casdoor...`
  - 默认邀请邮件文案仍是 `You have invited to join Casdoor...`
- `web/src/auth/Web3Auth.js`
  - Web3 typed-data domain name 仍是 `Casdoor`
  - app metadata name / description 仍是 `Casdoor` / `Connect a wallet using Casdoor`
- `web/public/index.html`
  - 默认 `title` 仍是 `Casdoor`
  - description meta 仍把产品定义为 Casdoor
- `web/src/locales/*/data.json`
  - 上面这些 key 的多语言副本都还在
  - 以及一批“Use same DB as Casdoor”“default Casdoor login page”之类的提示文案

建议替换入口：

- 第一批先改英文源 key 对应的主入口页面
- 然后统一回刷 `web/src/locales/*/data.json` 和 `i18n/locales/*/data.json`
- 否则你会在 UI rename 完成后继续被旧 i18n key 拉回去

### 4.4 OpenAPI / Swagger 元数据仍把系统定义为 Casdoor

- `routers/router.go`
  - Beego swagger 注解仍是
    - `@Title Casdoor RESTful API`
    - `@Description Swagger Docs of Casdoor Backend API`
    - `@ExternalDocs Find out more about Casdoor`
- `swagger/swagger.yml`
  - `info.title` / `info.description` 仍是 Casdoor
- `swagger/swagger.json`
  - 生成后的 swagger 元数据同样仍是 Casdoor

这属于对外文档面，不只是内部注释，因为 API explorer、SDK 生成和第三方接入时都会看到。

建议替换入口：

- 如果你准备先改命名而不立刻重生全部文档，至少先把源注解和 swagger 产物一起换成中性 runtime 名称

### 4.5 监控指标命名仍暴露 Casdoor 运行时身份

- `object/prometheus.go`
  - metric name 仍是
    - `casdoor_api_throughput`
    - `casdoor_api_latency`
    - `casdoor_cpu_usage`
    - `casdoor_memory_usage`
    - `casdoor_total_throughput`
  - help text 也仍直接写 Casdoor

这不是实现细节，而是外部监控、告警、Dashboard 和指标兼容面。

建议替换入口：

- 这块适合在“文档 / 外显命名”阶段一起处理
- 如果担心破坏现有监控，可以先加新指标名或做迁移窗口，再移除旧指标名

### 4.6 模块名和导出符号里仍保留 Casdoor 入口名

这类点未必改变运行时行为，但它们本身就是“第一批命名替换入口”。

- `web/src/common/CasdoorAppConnector.js`
  - 文件名仍保留 `Casdoor`
  - 但导出实现已经是 `generatePipelineAuthAppUrl` / `PipelineAuthAppQrCode` / `PipelineAuthAppUrl`
- `web/src/table/MfaAccountTable.js`
  - 仍从 `../common/CasdoorAppConnector` 导入
- `web/src/auth/CasdoorLoginButton.js`
  - 文件名和默认导出名都仍是 Casdoor

建议替换入口：

- 这类点很适合优先做纯命名层 rename
- 因为行为已经基本中性化，改名成本低、收益高，能明显减少后续实现层继续回流到 Casdoor 命名

## 审计说明

这份清单的完成性校验采用了两类证据：

- 先按四块做目标检索：
  - 配置：`casdoor_application`、`casdoorApplication`、`casdoorName`、`app-built-in`、`built-in`
  - OAuth：`CasdoorIdProvider`、`provider.type === "Casdoor"`、`casdoor_access_token`、`__casdoor_callback_react`
  - 短信：`go-sms-sender`、`SendSms`、短信 provider 分支
  - 文档/外显命名：`Casdoor app`、`About Casdoor`、`Casdoor RESTful API`、`casdoor_*` metrics
- 再对剩余 `Casdoor/casdoor` 命中做收口分类：
  - 已纳入四块清单的真实依赖面
  - 明确排除的测试、版权头、纯注释、以及 CAS 协议能力本身

当前审计后，新增补录的高价值命中主要集中在：

- Email / Log provider 默认文案与默认 type
- CAS 登录成功提示文案
- Web3Auth 展示名与 typed-data domain
- 文件名级 rename 入口，如 `CasdoorAppConnector`

## 替换矩阵

下面这张矩阵把“第一批最适合开刀的入口”压缩成更可执行的视图。

| 分块 | 当前入口 | 依赖类型 | 第一批建议动作 |
| --- | --- | --- | --- |
| 配置 | `object/site.go` `AuthApplication` / `casdoor_application` | 站点配置模型 | 新增中性字段名，旧 DB/JSON 字段只保留兼容 |
| 配置 | `object/session.go` `AuthApplication` / `CasdoorOrganization` | built-in runtime 常量 | 抽成中性 built-in runtime 常量 |
| 配置 | `conf/app.conf` / `conf/web_config.go` 中 `app-built-in` / `built-in` fallback | 默认运行时配置 | 改为引用中性常量，不在配置层硬编码 Casdoor 语义 |
| 配置 | `storage/casdoor.go` / `storage/storage.go` / `object/storage.go` | Storage provider type | 把 `Casdoor` 存储类型收敛成 legacy alias 或中性 provider |
| 配置 | `object/provider.go` `Casdoor Permission Log` | Log provider type | 改成中性日志 provider 名，旧 type 做兼容别名 |
| OAuth | `service/util.go` | SDK client 构造 | 抽 `ProviderRuntimeClient` facade |
| OAuth | `service/oauth.go` | 授权 URL / code exchange | 收进 facade，去掉直接 Casdoor OAuth server 假设 |
| OAuth | `service/proxy.go` | access token 读取与校验 | 通过 facade 处理，并保留旧 cookie 兼容窗口 |
| OAuth | `idp/casdoor.go` / `idp/provider.go` | 上游 Casdoor IdP | 泛化成通用 OIDC adapter 或收敛为 legacy adapter |
| OAuth | `web/src/auth/CasdoorLoginButton.js` / `web/src/auth/ProviderButton.js` / `web/src/auth/Provider.js` | 前端 OAuth provider 入口 | 组件和 type 改中性名，`Casdoor` 留作 legacy alias |
| OAuth | `web/public/ProviderHintRedirect.js` / `web/public/AuthCallbackHandler.js` | 静态 OAuth 辅助链路 | 和 React 主链一起去 Casdoor 化 |
| 短信 | `object/sms.go` | 短信总入口 | 抽 `SmsGateway`，把 `go-sms-sender` 下沉为 adapter |
| 短信 | `object/verification.go` / `controllers/service.go` | 验证码和测试短信 | 继续通过 `object.SendSms()` 单一入口走网关 |
| 文档 | `docs/identity-and-payment-architecture.md` | 架构叙事 | 从 Casdoor 角色叙事改为 provider/runtime 叙事 |
| 文档 | `routers/router.go` / `swagger/swagger.yml` / `swagger/swagger.json` | OpenAPI 元数据 | 对外 API 标题/描述改中性 runtime 名 |
| 文档 | `web/public/index.html` / `web/src/TourConfig.js` / locales | 页面元数据和 UI 文案 | 去 Casdoor 化并统一回刷多语言 |
| 文档 | `object/prometheus.go` | 监控指标名 | 规划指标迁移窗口，逐步从 `casdoor_*` 切到新命名 |

## 第一批最值得直接替换的命名 / 实现入口

如果目标是继续把系统从“Casdoor 运行时”抽成“自有 provider 运行时”，我建议第一批就动下面这些切口：

1. `service/util.go`
   - 把 `getAuthServerClientFromSite()` 抽成中性 runtime client facade
2. `service/oauth.go`
   - 把授权 URL / code exchange / cookie 写入收进 facade
3. `service/proxy.go`
   - 把 token 校验、旧 cookie fallback 收进 facade
4. `object/site.go`
   - 给 `AuthApplication` 建立中性命名，旧 `casdoor_application` / `casdoorApplication` 只保留兼容
5. `idp/casdoor.go`
   - 先决定是泛化成通用 OIDC adapter，还是收敛成 legacy adapter
6. `object/sms.go`
   - 把 `go-sms-sender` 退到 adapter 层，保住 `SendSms()` 作为唯一业务入口
7. `storage/casdoor.go`
   - 把 `Casdoor` 存储 provider 从主类型收敛成 legacy alias 或中性 runtime provider
8. `web/src/auth/CasdoorLoginButton.js`
   - 这是最清晰的前端 rename 起点
9. `web/src/auth/ProviderButton.js`
   - 把 `provider.type === "Casdoor"` 收敛成 legacy alias
10. `web/public/ProviderHintRedirect.js`
   - 前端静态 OAuth 辅助链路要和 React 主链一起改，不然会漏掉无 React fallback 的入口
11. `object/mfa_totp.go`
   - 把默认 issuer 从 Casdoor 抽成中性 runtime 默认值
12. `web/src/common/CasdoorAppConnector.js`
   - 文件名级入口已中性化实现，适合先 rename
13. `web/src/SiteEditPage.js`
   - 把 `Casdoor app` 改成新的运行时术语
14. `routers/router.go` / `swagger/swagger.yml`
   - 对外 API 元数据要一起去 Casdoor 化
15. `docs/identity-and-payment-architecture.md`
   - 把架构叙事从 “Casdoor 是角色本体” 改成 “provider/runtime 是角色本体”

## 这轮刻意没有算进“剩余真实依赖面”的内容

下面这些我这轮没有计入四块核心清单：

- CAS 协议本身的能力与路由
  - 例如 `object/token_cas.go`
  - `web/src/auth/CasLogout.js`
  - `web/src/EntryPage.js` 里的 `/cas/:owner/:casApplicationName/...`
  - 这些是 CAS protocol 支持，不是 Casdoor 品牌依赖
- 版权头里的 `The Casdoor Authors`
- 大量测试数据或示例字符串中的 `casdoor`
- 不影响当前抽离路径的第三方仓库 import，例如 `notify2`、`gomail`、`oss`
- 单纯历史注释且不构成运行时入口的品牌词

这些后面可以继续清理，但不是“从 Casdoor 运行时抽成自有 provider 运行时”的第一阻塞面。
