# Personal Blog System

一个基于 Go 语言开发的个人博客系统后端，使用 Gin 框架和 GORM 库构建，实现了完整的博客文章管理、用户认证和评论功能。

## 功能特性

### 核心功能
- 用户注册、登录、JWT 认证
- 博客文章的 CRUD 操作
- 文章评论系统
- 用户权限管理
- 数据验证和错误处理
- 请求日志记录


### 后端框架
- **Gin** - Go Web 框架
- **GORM** - Go ORM 库
- **JWT** - 身份认证

### 数据库
- **MySQL** - 主数据库

### 工具库
- **bcrypt** - 密码加密
- **validator** - 数据验证
- **viper** - 配置管理

## 系统要求

- Go 1.21+
- MySQL 5.7+
- Git

##  配置环境：

# 1.复制并修改配置文件：
cp config/config.yaml config/config.yaml.local

# 2.编辑 config/config.yaml.local：
server:
  port: 8080
  mode: "debug"

database:
  host: "localhost"
  port: 3306
  user: "root"
  password: "your_mysql_password"  # 修改为你的 MySQL 密码
  dbname: "blog_system"
  charset: "utf8mb4"

jwt:
  secret: "your-super-secret-jwt-key-change-in-production"
  expire: 24
  issuer: "blog-system"

也可以通过环境变量覆盖配置：
export BLOG_DB_PASSWORD=your_password
export BLOG_JWT_SECRET=your_jwt_secret
export BLOG_SERVER_PORT=8080

# 3.安装依赖：
go mod tidy

# 4.初始化数据库：
确保 MySQL 服务正在运行，然后创建数据库：
CREATE DATABASE blog_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 5.运行项目：
开发模式
go run main.go

#或者编译后运行
go build -o blog-system
./blog-system

# 6.验证安装：
访问健康检查端点：
curl http://localhost:8080/health
应该返回：
{
  "status": "ok",
  "message": "Blog system is running",
  "version": "1.0.0"
}

##  API 文档
# 1.认证相关
用户注册
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}


用户登录
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
响应：
{
  "success": true,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "role": "user",
      "created_at": "2023-10-01T10:00:00Z"
    }
  }
}
# 2.文章相关
获取文章列表
GET /api/v1/posts?page=1&page_size=10

创建文章（需登录）
POST /api/v1/posts
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "title": "我的第一篇文章",
  "content": "这是文章内容...",
  "summary": "文章摘要",
  "status": "published"
}

更新文章（需登录，仅作者）
PUT /api/v1/posts/1
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "title": "更新后的标题",
  "content": "更新后的内容..."
}

删除文章（需登录，仅作者）
DELETE /api/v1/posts/1
Authorization: Bearer <your_jwt_token>

# 3.评论相关
获取文章评论
GET /api/v1/comments/posts/1?page=1&page_size=20

创建评论（需登录）
POST /api/v1/comments
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "content": "这是一条评论",
  "post_id": 1
}

删除评论（需登录，仅作者）
DELETE /api/v1/comments/1
Authorization: Bearer <your_jwt_token>

# 4.评论相关
获取用户信息
GET /api/v1/users/1

获取当前用户信息（需登录）
GET /api/v1/users/profile
Authorization: Bearer <your_jwt_token>

更新用户信息（需登录）
PUT /api/v1/users/profile
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "bio": "个人简介",
  "avatar": "https://example.com/avatar.jpg"
}

##  项目结构
blog-system/
├── main.go                 # 应用入口
├── config/                 # 配置管理
│   ├── config.go
│   └── config.yaml
├── database/              # 数据库连接
│   ├── connection.go
│   └── migration.go
├── models/                # 数据模型
│   ├── user.go
│   ├── post.go
│   └── comment.go
├── controllers/           # 控制器层
│   ├── auth_controller.go
│   ├── user_controller.go
│   ├── post_controller.go
│   └── comment_controller.go
├── services/              # 业务逻辑层
│   ├── auth_service.go
│   ├── user_service.go
│   ├── post_service.go
│   └── comment_service.go
├── middleware/            # 中间件
│   ├── auth_middleware.go
│   ├── cors_middleware.go
│   ├── logger_middleware.go
│   └── security_middleware.go
├── utils/                 # 工具函数
│   ├── password_utils.go
│   ├── response_utils.go
│   ├── jwt_utils.go
│   └── validator_utils.go
└── routes/                # 路由定义
    └── routes.go


##  配置说明
# 1.服务器配置
server:
  port: 8080              # 服务端口
  mode: "debug"           # 运行模式：debug, release, test
  read_timeout: 30        # 读取超时（秒）
  write_timeout: 30       # 写入超时（秒）

# 2.数据库配置
database:
  host: "localhost"       # 数据库主机
  port: 3306              # 数据库端口
  user: "root"            # 数据库用户
  password: ""            # 数据库密码
  dbname: "blog_system"   # 数据库名称
  charset: "utf8mb4"      # 字符集
  max_idle_conns: 10      # 最大空闲连接数
  max_open_conns: 100     # 最大打开连接数

# 3.数据库配置
jwt:
  secret: "secret-key"    # JWT 密钥（生产环境请修改）
  expire: 24              # Token 过期时间（小时）
  issuer: "blog-system"   # 签发者



##  测试
# 使用 Postman 测试
1.导入 Postman 集合（可参考 docs/postman_collection.json）

2.配置环境变量：

base_url: http://localhost:8080

token: 登录后获取的 JWT token

# 测试步骤
1.注册用户

 调用 /api/v1/auth/register 注册新用户

2.用户登录

 调用 /api/v1/auth/login 获取 token

3.创建文章

 使用 token 调用 /api/v1/posts 创建文章

4.发表评论

 使用 token 调用 /api/v1/comments 发表评论

5.测试权限

 尝试修改/删除其他用户的文章和评论（应该返回 403 错误）