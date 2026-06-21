# 🐠 水族箱智能控制系统

一个专为高档白子孔雀鱼和名贵水草设计的全栈智能匹配系统，实现光照和溶氧量的智能协调控制，防止水草因光照和氧气配合不当而腐烂。

## 🏗️ 系统架构

```
┌─────────────────────────────────────────────────┐
│                   前端 (Vue3)                   │
│  ┌──────────┬──────────┬──────────┬──────────┐  │
│  │ 仪表盘  │ 光照排程 │ 溶氧量   │ 设备控制 │  │
│  └──────────┴──────────┴──────────┴──────────┘  │
└─────────────────────────────────────────────────┘
                              ↓ HTTP API
┌─────────────────────────────────────────────────┐
│              后端 (Go + Gin)                   │
│  ┌─────────────────┬─────────────────────────┐  │
│  │ 光照排程器模块  │  溶氧量匹配模型模块     │  │
│  │  - model.go     │  - model.go             │  │
│  │  - service.go   │  - service.go (算法)    │  │
│  │  - controller.go│  - controller.go        │  │
│  └─────────────────┴─────────────────────────┘  │
│  ┌─────────────────┬─────────────────────────┐  │
│  │ 设备控制模块    │  传感器数据模块         │  │
│  └─────────────────┴─────────────────────────┘  │
└─────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────┐
│              MySQL 数据库                      │
│  - light_schedules   (光照排程表)             │
│  - oxygen_configs    (溶氧量配置表)           │
│  - device_status     (设备状态表)             │
│  - sensor_data       (传感器数据表)           │
└─────────────────────────────────────────────────┘
```

## ✨ 核心功能

### 1. 光照排程器模块
- 支持配置多时段光照排程（开灯时间、关灯时间、亮度）
- 实时计算当前光照状态和下一步动作
- 支持排程的增删改查
- 自动判断当前生效的排程

### 2. 溶氧量匹配模型模块
- 基于灯光瓦数和水温的智能匹配算法
- 3×3 匹配矩阵（3档光照 × 3档水温 = 9种配置）
- 气泵档位 1-5 档自动调节
- 无匹配配置时自动公式计算兜底
- 匹配结果可视化展示

### 3. 设备控制模块
- 手动/自动模式切换
- 灯光开关和亮度调节（0-100%）
- 气泵开关和档位调节（1-5档）
- 快捷操作（全开/全关/全部自动）
- 实时设备状态监控

### 4. 传感器数据模块
- 水温和灯光瓦数数据采集
- 24小时数据趋势图表
- 统计数据计算（平均值、当前值）

## 🛠️ 技术栈

**后端**
- Go 1.21
- Gin 框架
- GORM ORM
- MySQL 8.0+
- Viper 配置管理
- CORS 跨域支持

**前端**
- Vue 3 (Composition API)
- Vite 构建工具
- Element Plus UI 组件库
- ECharts 图表库
- Axios HTTP 客户端
- Vue Router 路由

## 📁 项目结构

```
da6/
├── backend/                    # Go 后端
│   ├── cmd/
│   │   └── main.go            # 主入口
│   ├── config/
│   │   └── config.yaml        # 配置文件
│   ├── internal/
│   │   ├── common/            # 通用响应
│   │   ├── config/            # 配置加载
│   │   ├── database/          # 数据库连接
│   │   ├── router/            # 路由注册
│   │   └── modules/
│   │       ├── lighting/      # 光照排程器模块 (多文件)
│   │       │   ├── model.go
│   │       │   ├── service.go
│   │       │   └── controller.go
│   │       ├── oxygen/        # 溶氧量匹配模型模块 (多文件)
│   │       │   ├── model.go
│   │       │   ├── service.go
│   │       │   └── controller.go
│   │       ├── device/        # 设备控制模块
│   │       └── sensor/        # 传感器模块
│   └── go.mod
├── frontend/                   # Vue3 前端
│   ├── src/
│   │   ├── views/              # 页面组件
│   │   ├── router/             # 路由配置
│   │   ├── utils/              # API 封装
│   │   ├── styles/             # 全局样式
│   │   ├── App.vue
│   │   └── main.js
│   ├── package.json
│   └── vite.config.js
├── database/
│   └── schema.sql             # 数据库脚本
└── README.md
```

## 🚀 快速开始

### 前置要求
- Go 1.21+
- Node.js 18+
- MySQL 8.0+

### 1. 初始化数据库

```bash
# 导入数据库脚本
mysql -u root -p < database/schema.sql
```

脚本会自动创建：
- 数据库 `aquarium_control`
- 4 张数据表
- 初始测试数据

### 2. 配置数据库连接

修改 `backend/config/config.yaml`：

```yaml
database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: your_password  # 修改为你的密码
  dbname: aquarium_control
```

### 3. 启动后端

```bash
cd backend

# 下载依赖
go mod tidy

# 运行
go run cmd/main.go
```

后端服务将在 `http://localhost:8080` 启动

### 4. 启动前端

```bash
cd frontend

# 安装依赖
npm install

# 开发模式运行
npm run dev
```

前端服务将在 `http://localhost:5173` 启动

### 5. 访问系统

打开浏览器访问 `http://localhost:5173`

## 📡 API 接口文档

### 光照排程器 API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/lighting/schedules` | 获取所有排程 |
| GET | `/api/v1/lighting/schedules/:id` | 获取单个排程 |
| POST | `/api/v1/lighting/schedules` | 创建排程 |
| PUT | `/api/v1/lighting/schedules/:id` | 更新排程 |
| DELETE | `/api/v1/lighting/schedules/:id` | 删除排程 |
| GET | `/api/v1/lighting/status` | 获取当前光照状态 |

### 溶氧量匹配 API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/oxygen/configs` | 获取所有配置 |
| POST | `/api/v1/oxygen/configs` | 创建配置 |
| PUT | `/api/v1/oxygen/configs/:id` | 更新配置 |
| DELETE | `/api/v1/oxygen/configs/:id` | 删除配置 |
| POST | `/api/v1/oxygen/match` | 计算匹配结果 |
| GET | `/api/v1/oxygen/matrix` | 获取匹配矩阵 |

### 设备控制 API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/device` | 获取所有设备 |
| GET | `/api/v1/device/dashboard` | 获取仪表盘状态 |
| GET | `/api/v1/device/:type/:name` | 获取单个设备 |
| PUT | `/api/v1/device/:type/:name/toggle` | 开关设备 |
| PUT | `/api/v1/device/:type/:name/value` | 调节参数 |
| PUT | `/api/v1/device/:type/:name/manual` | 切换手动/自动 |

### 传感器数据 API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/sensor/data` | 获取历史数据 |
| GET | `/api/v1/sensor/data/latest` | 获取最新数据 |
| POST | `/api/v1/sensor/data` | 上传传感器数据 |
| GET | `/api/v1/sensor/stats` | 获取统计数据 |

## 🧮 溶氧量匹配算法

### 匹配逻辑
1. 优先匹配数据库中的配置表
2. 根据 `灯光瓦数区间` + `水温区间` 查找对应档位
3. 无匹配配置时使用公式计算

### 计算公式
```
泵档位 = clamp(round((灯光瓦数/100 * 2) + (温度/40 * 3)), 1, 5)
```

### 匹配矩阵 (默认配置)

| 灯光瓦数 \ 水温 | 18-24°C | 24.1-28°C | 28.1-32°C |
|-----------------|---------|-----------|-----------|
| 0-30W (低光)    | 1档     | 2档       | 3档       |
| 31-60W (中光)   | 2档     | 3档       | 4档       |
| 61-100W (高光)  | 3档     | 4档       | 5档       |

## 🎨 界面预览

系统采用深海蓝主题玻璃拟态设计风格：
- 仪表盘：实时数据卡片 + 趋势图表 + 智能匹配计算器
- 光照排程：排程列表 + 时间线可视化 + 编辑弹窗
- 溶氧量控制：配置管理 + 匹配矩阵 + 实时计算
- 设备控制：分Tab管理灯光和气泵 + 快捷操作

## 🔧 核心代码说明

### 光照排程核心逻辑
`backend/internal/modules/lighting/service.go`
- `GetCurrentStatus()`: 实时计算当前光照状态
- `CalculateWattage()`: 亮度百分比转实际瓦数

### 溶氧量匹配核心算法
`backend/internal/modules/oxygen/service.go`
- `CalculateMatch()`: 匹配主逻辑
- `calculateFallbackLevel()`: 兜底公式计算
- `GetMatchMatrix()`: 匹配矩阵数据

### 设备控制核心
`backend/internal/modules/device/service.go`
- `SetManualMode()`: 手动/自动切换
- `GetDashboardStatus()`: 仪表盘聚合数据

## 📝 注意事项

1. **数据库密码**：首次运行请修改 `backend/config/config.yaml` 中的数据库密码
2. **时区设置**：确保 MySQL 和系统时区一致，避免时间判断错误
3. **手动模式**：切换到手动模式后，系统不会自动控制设备，需手动操作
4. **传感器数据**：系统支持传感器数据接入，可通过 API 上报水温和灯光瓦数

## 🤝 扩展建议

- 接入实际硬件控制（通过串口或 MQTT 发送控制指令）
- 添加定时任务，自动根据排程控制设备
- 添加告警功能（水温过高、溶氧量过低等）
- 接入移动端推送通知
- 添加更多传感器支持（PH值、TDS等）

## 📄 License

MIT License
