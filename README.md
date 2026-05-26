# NOAMS - 网络运维自动化管理系统

Network Operations Automation Management System

基于 Go + Vue 3 的全栈网络设备运维管理平台，实现对 H3C、华为、思科、锐捷等主流厂商网络设备的统一纳管、自动化巡检、配置备份与告警管理。

## 系统架构

```
┌──────────────────────────────────────────────────┐
│                  前端层 (Vue 3)                    │
│  设备管理 │ 设备巡检 │ 配置备份 │ 状态监控 │ 告警管理 │
└──────────────────────┬───────────────────────────┘
                       │ HTTP REST API
┌──────────────────────▼───────────────────────────┐
│              后端服务层 (Go + Gin)                 │
│  设备管理 │ 巡检调度 │ 配置管理 │ 告警引擎 │ 任务调度 │
└──────────────────────┬───────────────────────────┘
                       │
┌──────────────────────▼───────────────────────────┐
│                自动化执行层                         │
│           Netmiko (Python SSH 微服务)              │
└──────────────────────┬───────────────────────────┘
                       │ SSH
┌──────────────────────▼───────────────────────────┐
│              被管网络设备层                         │
│  H3C │ 华为 │ 思科 │ 锐捷 │ 等                     │
└──────────────────────────────────────────────────┘
```

## 技术栈

| 层级 | 技术 |
|---|---|
| 前端 | Vue 3 + TypeScript + Vite + Element Plus + Pinia |
| 后端 | Go + Gin + GORM |
| 数据库 | MySQL 8.0 / SQLite（开发模式） |
| SSH 自动化 | Netmiko (Python 微服务) |
| 容器化 | Docker + Docker Compose |

## 功能模块

### 设备管理
- 设备全生命周期管理（添加/编辑/删除）
- 支持 H3C Comware、华为 VRP、思科 IOS、锐捷 OS
- 凭据 AES-256 加密存储
- 支持设备分组管理
- 多维度搜索过滤
- JSON 导入/导出

### 设备巡检
- 定时/手动触发设备巡检
- 自动采集 CPU、内存、接口状态
- 并发 SSH 执行（goroutine 池）
- 异常指标自动检测与告警
- 巡检记录查看与 CSV 导出

### 配置备份
- 一键触发配置备份
- 配置内容查看与下载
- 多版本历史追溯

### 状态监控
- Dashboard 大屏概览
- 设备在线率统计
- CPU/内存使用率 TOP 排行
- 最近巡检记录列表

### 告警管理
- 多级告警（严重/警告/提示）
- 设备离线/上线自动告警
- 巡检异常自动告警
- 告警确认与处理流程

### 定时任务
- Cron 表达式调度管理
- 预置任务模板（巡检、备份）
- 执行日志查看

## 快速开始

### 环境要求

- Go 1.22+
- Node.js 18+
- Docker + Docker Compose（可选）

### 本地开发

```bash
# 1. 启动后端
cd backend
export GOPROXY=https://goproxy.cn,direct   # 国内加速
go mod tidy
go run main.go

# 2. 启动前端（新终端）
cd frontend
npm install
npm run dev

# 3. 启动 Netmiko 微服务（新终端）
cd netmiko-worker
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python app.py
```

### Docker 部署

```bash
cp .env.example .env
docker compose up -d
```

访问 `http://localhost:3000`

默认管理员账号：`admin` / `admin123`

## 项目结构

```
NOAMS/
├── backend/               # Go 后端
│   ├── main.go            # 入口：服务启动、数据库迁移
│   ├── config/            # 配置管理
│   ├── models/            # 数据模型（8 个核心表）
│   ├── handlers/          # API 处理器
│   ├── services/          # 业务逻辑层
│   ├── middleware/         # 认证与 CORS 中间件
│   ├── routes/            # 路由注册
│   └── Dockerfile
├── frontend/              # Vue 3 前端
│   ├── src/
│   │   ├── views/         # 页面组件
│   │   ├── api/           # API 接口层
│   │   ├── router/        # 路由配置
│   │   ├── stores/        # Pinia 状态管理
│   │   └── layouts/       # 布局组件
│   └── Dockerfile
├── netmiko-worker/        # SSH 自动化微服务
│   ├── app.py             # Flask + Netmiko
│   └── Dockerfile
├── nginx/                 # 反向代理配置
├── docker-compose.yml     # 服务编排
└── .env.example           # 环境变量模板
```

## API 概览

| 模块 | 端点 | 说明 |
|---|---|---|
| 认证 | `POST /api/v1/auth/login` | 用户登录 |
| 设备 | `GET/POST /api/v1/devices` | 设备列表/创建 |
| 设备 | `GET/PUT/DELETE /api/v1/devices/:id` | 设备详情/编辑/删除 |
| 分组 | `GET/POST /api/v1/groups` | 分组列表/创建 |
| 巡检 | `POST /api/v1/devices/:id/inspect` | 单设备巡检 |
| 巡检 | `POST /api/v1/inspections/batch` | 批量巡检 |
| 巡检 | `GET /api/v1/inspections/export` | 巡检记录 CSV 导出 |
| 配置 | `POST /api/v1/configs/backup` | 配置备份 |
| 配置 | `GET /api/v1/configs/export/:id` | 配置内容下载 |
| 告警 | `GET /api/v1/alerts` | 告警列表 |
| 监控 | `GET /api/v1/monitor/dashboard` | Dashboard 数据 |

完整接口列表见 [实施文档](./NOAMS实施文档.md)

## 测试

```bash
# 后端测试
cd backend && go test ./... -count=1

# 前端检查
cd frontend && npx vue-tsc --noEmit

# 前端构建
cd frontend && npx vite build
```

## 实施计划

| 阶段 | 周期 | 内容 |
|---|---|---|
| 一期：基础平台 | 6 周 | 环境搭建、设备管理、用户认证、Netmiko 微服务、Dashboard |
| 二期：自动化能力 | 4 周 | 自动巡检、配置备份、告警引擎 |
| 三期：增强优化 | 3 周 | AP 管理、定时任务、配置回滚、权限细化 |

## License

MIT
