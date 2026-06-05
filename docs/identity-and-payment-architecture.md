# 身份与支付架构接入说明

## 1. 目标

本项目采用统一身份与统一支付分层：

- `oauth.xx.cn` 作为 Casdoor，负责外部认证入口
- `identity.xx.cn` 作为公司唯一身份中心
- `payment.xx.cn` 作为统一支付网关
- 各业务产品只对接 `identity.xx.cn` 和 `payment.xx.cn`

相关文档：

- [文档目录](./README.md)
- [认证网关当前能力盘点](./identity/auth-gateway-capabilities.md)
- [Identity Service PRD](./identity-service-prd.md)
- [Identity Service API 设计](./identity-service-api.md)
- [Identity Service 场景、用户行为与系统行为](./identity-service-scenarios.md)

## 2. 架构

```text
第三方 Provider
   -> Casdoor(oauth.xx.cn)
   -> Identity Service(identity.xx.cn)
   -> Unified Payment Gateway(payment.xx.cn)
   -> Product A / Product B / Admin / API Gateway
```

### 2.1 Identity Service 职责

- 对接 Casdoor OIDC
- 维护公司唯一用户
- 管理第三方账号绑定
- 维护统一会话
- 签发内部 Identity JWT
- 提供用户资料和权限查询接口
- 负责产品登录策略控制和租户访问授权

### 2.2 Payment Gateway 职责

- 校验 Identity JWT
- 识别来源产品
- 统一创建订单、支付、退款、对账
- 适配微信支付、支付宝、Stripe、Apple Pay
- 处理异步回调和订单状态流转

## 3. 接入流程

### 3.1 登录

1. 用户访问业务产品
2. 产品跳转到 `identity.xx.cn/auth/login`
3. Identity Service 跳转 Casdoor 授权页
4. Casdoor 完成第三方登录
5. Casdoor 回调 `identity.xx.cn/auth/callback`
6. Identity Service 换取 token 并拉取用户信息
7. Identity Service 完成用户绑定
8. Identity Service 签发内部 JWT
9. 用户回到业务产品

### 3.2 支付

1. 业务产品携带 Identity JWT 调用支付网关
2. Payment Gateway 校验身份和权限
3. 创建统一订单
4. 调用具体支付通道
5. 支付回调进入网关
6. 网关更新订单状态并通知业务方

## 4. Identity Service 对接点

### 4.1 Casdoor 配置

在 Casdoor 中为 `identity.xx.cn` 创建一个 Application：

- Redirect URI: `https://identity.xx.cn/auth/callback`
- Grant Type: `authorization_code`
- Scope: `openid profile email`
- Token Format: `JWT`

### 4.2 Identity Service 必备接口

- `GET /auth/login`
- `GET /auth/callback`
- `POST /auth/logout`
- `GET /me`
- `GET /me/permissions`

### 4.3 建议数据表

- `users`
- `account_bindings`
- `sessions`
- `tenants`
- `roles`
- `permissions`
- `orders`
- `payments`
- `refunds`
- `billing_records`

## 5. Payment Gateway 对接点

### 5.1 必备接口

- `POST /orders`
- `GET /orders/{id}`
- `POST /orders/{id}/pay`
- `POST /refunds`
- `GET /channels`

### 5.2 核心规则

- 所有请求先验 Identity JWT
- 同一业务单号必须幂等
- 支付回调必须可重复消费
- 订单状态以状态机推进，不允许直接覆盖

## 6. 研发约定

- Casdoor 只做上游认证，不承载公司内部用户域
- Identity Service 才是公司唯一身份源
- Payment Gateway 只信任 Identity Service 发出的身份
- 业务产品不直接依赖 Casdoor provider

## 7. 落地顺序

1. 先接通 Casdoor -> Identity Service 的 OIDC 流程
2. 再落地公司唯一用户和绑定表
3. 再接统一 JWT 和会话
4. 最后接统一支付网关和订单状态机
