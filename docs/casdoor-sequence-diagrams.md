# Casdoor 接入时序图

## 1. 浏览器登录

```mermaid
sequenceDiagram
    autonumber
    participant U as "用户浏览器"
    participant F as "业务前端"
    participant B as "业务后端"
    participant C as "Casdoor"
    participant P as "QQ / 微信"

    U->>F: 访问受保护页面
    F->>B: 请求登录地址
    B-->>F: 302 /auth/login -> Casdoor authorize
    F->>C: 浏览器跳转到 Casdoor 登录页
    U->>C: 选择 QQ / 微信登录
    C->>P: 跳转第三方授权页
    U->>P: 完成授权 / 扫码
    P-->>C: 回调 code
    C-->>B: 回调业务后端 /auth/callback?code=...
    B->>C: 后端 POST /api/login/oauth/access_token
    C-->>B: 返回 access_token / jwt
    B-->>F: 设置业务 session 或返回业务 JWT
```

## 2. 短信验证码登录

```mermaid
sequenceDiagram
    autonumber
    participant U as "用户浏览器"
    participant F as "业务前端"
    participant B as "业务后端"
    participant C as "Casdoor"
    participant SMS as "阿里云短信"

    U->>F: 输入手机号
    F->>B: 请求发送验证码
    B->>C: 调用 Casdoor 登录/验证码接口
    C->>SMS: 发送短信验证码
    SMS-->>U: 收到验证码
    U->>F: 输入验证码
    F->>B: 提交验证码
    B->>C: 校验验证码并换 token
    C-->>B: 返回 access_token / jwt
    B-->>F: 设置业务 session 或返回业务 JWT
```

## 3. 业务侧验 JWT

```mermaid
sequenceDiagram
    autonumber
    participant A as "业务 App"
    participant J as "Casdoor JWKS"

    A->>J: 获取公钥
    J-->>A: 返回 JWKS
    A->>A: 验签 JWT
    A->>A: 读取 sub 作为 user_id
    A->>A: 执行业务逻辑
```

## 4. 关键说明

- `登录页跳转` 是浏览器前向重定向，不是纯服务端 API 调用
- `token 交换 / userinfo / jwks` 才是后端 API 调用
- `QQ / 微信` 只会被 Casdoor 直接对接，业务 App 不直接接它们
- `阿里云短信` 也是 Casdoor 内部的发送通道，业务 App 不直接发短信
