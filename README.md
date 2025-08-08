# 算法洞察平台 (AlgoInsight)

**语言**: **中文** | [English](README-EN.md)

一个交互式算法学习和性能分析工具，支持算法可视化和性能对比功能。

## 🚀 功能特性

### 功能一：交互式算法可视化
- **算法选择**: 支持多种排序、搜索、图算法等
- **数据输入**: 自定义输入、随机生成、预设数据集
- **动画展示**: 实时可视化算法执行过程
- **步骤控制**: 播放、暂停、单步执行、速度调节
- **详细信息**: 显示比较、交换、移动等操作详情

### 功能二：算法性能评测与对比
- **多算法测试**: 同时测试多个算法性能
- **数据规模配置**: 自定义测试数据大小和类型
- **性能指标**: 执行时间、内存使用、操作次数统计
- **对比分析**: 图表化展示性能差异
- **详细报告**: 生成完整的性能分析报告

## 🛠️ 技术栈

- **后端**: Go + Gin框架
- **前端**: Svelte + TypeScript + Vite
- **包管理**: pnpm
- **部署**: Docker + Docker Compose

## 📦 快速开始

### 使用Docker (推荐)

1. **克隆项目**
```bash
git clone https://github.com/zym9863/AlgoInsight.git
cd AlgoInsight
```

2. **构建并启动**
```bash
docker-compose up --build
```

3. **访问应用**
- 打开浏览器访问: http://localhost:8080

### 开发环境

#### 后端开发
```bash
cd server
go mod tidy
go run main.go
```

#### 前端开发
```bash
cd web
pnpm install
pnpm dev
```

#### 使用Docker开发环境
```bash
# 启动开发环境
docker-compose --profile dev up --build

# 后端: http://localhost:8080
# 前端: http://localhost:5173
```

## 🏗️ 项目结构

```
AlgoInsight/
├── Dockerfile                 # 生产环境Docker配置
├── docker-compose.yml         # Docker Compose配置
├── README.md                  # 项目文档
├── server/                    # Go后端
│   ├── main.go               # 服务器入口
│   ├── config/               # 配置管理
│   ├── handlers/             # HTTP处理器
│   ├── models/               # 数据模型
│   ├── services/             # 业务逻辑
│   ├── algorithms/           # 算法实现
│   │   ├── sorting/          # 排序算法
│   │   ├── searching/        # 搜索算法
│   │   └── graph/            # 图算法
│   └── utils/                # 工具函数
├── web/                      # Svelte前端
│   ├── src/
│   │   ├── components/       # UI组件
│   │   │   ├── algorithm/    # 算法相关组件
│   │   │   ├── ui/           # 通用UI组件
│   │   │   └── visualization/ # 可视化组件
│   │   ├── stores/           # 状态管理
│   │   ├── services/         # API服务
│   │   ├── types/            # TypeScript类型
│   │   └── utils/            # 工具函数
│   ├── public/
│   └── package.json
└── docs/                     # 文档
    └── architecture.md       # 架构设计文档
```

## 🔧 API接口

### 算法管理
- `GET /api/algorithms` - 获取所有算法
- `GET /api/algorithms/category/{category}` - 按类别获取算法
- `GET /api/algorithms/info/{id}` - 获取算法详细信息

### 可视化
- `POST /api/visualize/execute` - 执行算法可视化
- `GET /api/visualize/step/{sessionId}/{stepId}` - 获取可视化步骤
- `POST /api/visualize/reset` - 重置可视化状态

### 性能测试
- `POST /api/benchmark/run` - 运行性能测试
- `GET /api/benchmark/results/{testId}` - 获取测试结果
- `POST /api/benchmark/compare` - 对比测试结果

### 数据生成
- `POST /api/data/generate` - 生成测试数据
- `GET /api/data/presets` - 获取预设数据

## 🎯 支持的算法

### 排序算法
- 冒泡排序 (Bubble Sort)
- 快速排序 (Quick Sort)
- 归并排序 (Merge Sort)
- 堆排序 (Heap Sort)
- 插入排序 (Insertion Sort)

### 搜索算法
- 线性搜索 (Linear Search)
- 二分搜索 (Binary Search)
- 哈希搜索 (Hash Search)

### 图算法
- 广度优先搜索 (BFS)
- 深度优先搜索 (DFS)
- 最短路径 (Dijkstra)
- 最小生成树 (Kruskal)

## 🚀 部署

### 生产环境部署

1. **构建Docker镜像**
```bash
docker build -t algoinsight .
```

2. **运行容器**
```bash
docker run -p 8080:8080 algoinsight
```

3. **使用Docker Compose**
```bash
docker-compose up -d
```

### 环境变量配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `PORT` | 8080 | 服务端口 |
| `ENVIRONMENT` | development | 运行环境 |
| `MAX_EXECUTION_TIME` | 30 | 最大执行时间(秒) |
| `MAX_DATA_SIZE` | 10000 | 最大数据规模 |
| `BENCHMARK_TIMEOUT` | 60 | 性能测试超时(秒) |
| `MAX_CONCURRENT_TESTS` | 5 | 最大并发测试数 |

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📝 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - Go Web框架
- [Svelte](https://svelte.dev/) - 前端框架
- [Vite](https://vitejs.dev/) - 构建工具
- [TypeScript](https://www.typescriptlang.org/) - 类型系统