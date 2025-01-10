### 使用方法

- 使用方法: python script.py <新模块名>
- 例如: python script.py github.com/YXRRXY/autoAuth

### 状态码

- 200: 请求成功
- 201: 创建成功
- 400: 请求参数错误
- 401: 未认证或认证失败
- 403: 权限不足
- 404: 资源不存在
- 500: 服务器内部错误

### 路由详情

#### 1. 用户认证相关 (/api/v1/auth)

##### 1.1 注册
- 路径: `POST /auth/register`
- 权限: 公开
- 描述: 新用户注册
- 请求体:
  ```json
  {
      "username": "string",   // 3-50字符
      "password": "string",   // 6-32字符
      "email": "string"      // 有效邮箱
  }
  ```
- 响应: `200 OK`
  ```json
  {
      "message": "注册成功"
  }
  ```

##### 1.2 登录
- 路径: `POST /auth/login`
- 权限: 公开
- 描述: 用户登录
- 请求体:
  ```json
  {
      "username": "string",
      "password": "string"
  }
  ```
- 响应: `200 OK`
  ```json
  {
      "token": "string",
      "user": {
          "id": 1,
          "username": "string",
          "email": "string",
          "status": "active",
          "last_login_at": "2024-01-20T10:00:00Z",
          "created_at": "2024-01-20T10:00:00Z",
          "updated_at": "2024-01-20T10:00:00Z"
      }
  }
  ```

##### 1.3 刷新令牌
- 路径: `POST /auth/refresh`
- 权限: 需要认证
- 描述: 刷新访问令牌
- 请求头: 
  ```
  Refresh-Token: <refresh_token>
  ```
- 响应: `200 OK`
  ```json
  {
      "access_token": "string",
      "refresh_token": "string",
      "expires_in": 7200
  }
  ```