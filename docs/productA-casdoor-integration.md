# ProductA 接入 Casdoor 接口文档

## 1. 结论

ProductA 不需要“必须走 Casdoor 页面”。

- Web/H5：可以走浏览器重定向的 OIDC 标准流程
- App/小程序/原生端：可以先拿 QQ/微信 `code`，再交给 Casdoor 后端换成 Casdoor 的 token
- 短信：由 Casdoor 负责发送和校验，ProductA 只调用 Casdoor API

## 2. Casdoor 提供的核心接口

### 2.1 OIDC / JWT

- `GET /.well-known/openid-configuration`
- `GET /.well-known/:application/openid-configuration`
- `GET /.well-known/jwks`
- `GET /.well-known/:application/jwks`
- `GET /api/userinfo`
- `POST /api/login/oauth/access_token`
- `POST /api/login/oauth/refresh_token`
- `POST /api/login/oauth/introspect`

### 2.2 登录 / 注册

- `POST /api/login`
- `POST /api/signup`
- `GET /api/get-app-login`

### 2.3 短信验证码

- `POST /api/send-verification-code`
- `POST /api/verify-code`

## 3. QQ / 微信登录

### 3.1 推荐方式：标准 OIDC

1. ProductA 把用户重定向到 Casdoor `GET /login/oauth/authorize`
2. Casdoor 完成 QQ / 微信授权
3. Casdoor 回调 ProductA 的 `redirectUri`
4. ProductA 后端调用 `POST /api/login/oauth/access_token`
5. ProductA 用 `/.well-known/jwks` 验签 token，读取 `sub` 作为 `user_id`

### 3.2 直接接 QQ / 微信 code

如果 ProductA 已经从 QQ / 微信 SDK 拿到了第三方 `code`，可直接调：

- `POST /api/login`

请求体使用 `form.AuthForm`，关键字段：

- `application`: Casdoor Application 名称
- `provider`: provider 名称
- `code`: QQ / 微信返回的 code
- `method`: `signin` 或 `signup`
- `type`: 常用 `code`

## 4. 短信注册

### 4.1 发送验证码

`POST /api/send-verification-code`

请求体使用 `form.VerificationForm`：

- `applicationId`: `owner/name`
- `type`: `email` 或 `phone`
- `dest`: 邮箱或手机号
- `method`: `signup`
- `captchaType`
- `captchaToken`
- `clientSecret`

### 4.2 注册

`POST /api/signup`

请求体使用 `form.AuthForm`：

- `application`: Application 名称
- `organization`
- `username`
- `password`
- `name`
- `email` / `phone`
- `emailCode` / `phoneCode`
- `countryCode`

## 5. 短信登录

### 5.1 发送验证码

`POST /api/send-verification-code`

关键字段：

- `applicationId`
- `type=phone`
- `dest`
- `method=login`

### 5.2 登录

`POST /api/login`

关键字段：

- `application`
- `organization`
- `username`：手机号或邮箱
- `code`：验证码
- `signinMethod`：`Verification code`

## 6. 验证码校验

### 6.1 账号改绑 / 预校验

`POST /api/verify-code`

用途：

- 验证短信/邮箱验证码
- 通过后把 `verifiedCode`、`verifiedUserId` 写入 Casdoor session

## 7. 返回值

### 7.1 `POST /api/login/oauth/access_token`

返回 `TokenWrapper`：

- `access_token`
- `id_token`
- `refresh_token`
- `token_type`
- `expires_in`
- `scope`

### 7.2 `POST /api/login` / `POST /api/signup`

返回通用 `Response`：

- `status`
- `msg`
- `data`
- `data2`
- `data3`

## 8. ProductA 侧建议

- 后端只信任 Casdoor 签发的 JWT
- 用 `sub` 当 `user_id`
- 用 `/.well-known/jwks` 做本地验签
- `userinfo` 只作为补充，不作为主校验链路
