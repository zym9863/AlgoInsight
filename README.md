# 算法洞察平台 (AlgoInsight)

**语言**: **中文** | [English](README-EN.md)

<p align="left">
    <a href="https://github.com/zym9863/AlgoInsight/actions" target="_blank"><img alt="CI" src="https://img.shields.io/badge/CI-pending-lightgrey" /></a>
    <a href="LICENSE" target="_blank"><img alt="License" src="https://img.shields.io/badge/License-MIT-blue.svg" /></a>
    <img alt="Go Version" src="https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go" />
    <img alt="Svelte" src="https://img.shields.io/badge/Svelte-5-orange?logo=svelte" />
    <img alt="Docker" src="https://img.shields.io/badge/Docker-ready-2496ED?logo=docker" />
    <img alt="Status" src="https://img.shields.io/badge/status-active-success" />
</p>

一个面向学习与教学的**交互式算法可视化与性能分析平台**。通过统一的接口与步进追踪机制（Step Tracker），可在前端实时呈现排序 / 搜索 / 图等算法的执行过程，并支持多算法基准测试、数据生成、性能指标对比与结果可视化。

> 设计目标：降低算法学习的“抽象门槛”，用结构化数据 + 可视化动画，帮助理解时间/空间复杂度及不同算法策略差异。

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

### 功能三：数据集生成与管理
- **多模式**: 随机、近乎有序、逆序、重复值、特定分布
- **预设数据集**: 内置多种规模与特征的预设数据
- **统一结构**: 后端统一包装数据与元信息，便于复现实验

### 功能四：统一算法抽象接口
- 接口定义：`Execute(data, tracker)`
- 步骤追踪：比较、交换、访问等操作统一编码，前端可重构渲染
- 易扩展：新增算法只需实现接口并注册元信息

### 功能五：可视化交互控制
- 支持播放 / 暂停 / 单步 / 快速跳转
- 调速：实时调整动画速度
- 高亮：当前操作元素、比较对、已排序区间

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

> 更详细的字段说明与扩展示例可参考：`test_api.py`（含调用示例与简单性能测量）。

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

> 下列列表基于当前代码库实际实现（`server/algorithms`）。README 早期占位的 Dijkstra / Kruskal 尚未落地，避免误导已临时移除，可在 Roadmap 中查看计划。

### 排序算法
- 冒泡排序 (Bubble Sort)
- 快速排序 (Quick Sort)
- 归并排序 (Merge Sort)
- 堆排序 (Heap Sort)
- 插入排序 (Insertion Sort)
- 选择排序 (Selection Sort)
- 希尔排序 (Shell Sort)

### 搜索算法
- 线性搜索 (Linear Search)
- 二分搜索 (Binary Search)
- 哈希搜索 (Hash Search)

### 图算法
- 广度优先搜索 (BFS)
- 深度优先搜索 (DFS)

## 🧪 本地 API 快速测试

使用自带脚本：

```bash
python test_api.py
```

或手动调用示例（Windows CMD）：

```cmd
curl -X GET http://localhost:8080/api/algorithms
curl -X POST http://localhost:8080/api/data/generate -H "Content-Type: application/json" -d "{\"dataType\":\"array\",\"size\":10,\"pattern\":\"random\"}"
```

典型可视化请求：

```bash
curl -X POST http://localhost:8080/api/visualize/execute \
    -H "Content-Type: application/json" \
    -d '{"algorithmId":"bubble_sort","data":[5,3,8,2],"parameters":{}}'
```

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

## 🧭 Roadmap

- [ ] 增加图算法：Dijkstra、Kruskal、Prim、Topological Sort
- [ ] 增加更多性能指标：内存峰值、GC 次数
- [ ] 前端：步骤滑块快速跳转 / 断点标记
- [ ] 算法步骤导出（JSON / GIF / 视频）
- [ ] 国际化 (i18n) 自动切换
- [ ] Online Playground Demo 部署
- [ ] CI：单元测试 + 构建发布流水线

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📝 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## ❓ FAQ

**Q: 为什么某些算法只显示名称没有复杂度?**  
A: 复杂度信息来自算法元数据，若未填写或算法初版实现尚未补充会为空，欢迎补充 PR。

**Q: 如何新增一个算法?**  
A: 在 `server/algorithms/<category>/` 下创建文件，实现接口（参考现有排序算法）并在注册逻辑中加入元信息；前端会自动通过 `/api/algorithms` 获取。

**Q: 是否支持自定义数据结构?**  
A: 当前聚焦数组与基本图结构，后续会在 Roadmap 中扩展。

**Q: Docker 镜像是否可用于生产?**  
A: 目前镜像为教学/演示优化，生产可自行加反向代理、日志、鉴权。

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - Go Web框架
- [Svelte](https://svelte.dev/) - 前端框架
- [Vite](https://vitejs.dev/) - 构建工具
- [TypeScript](https://www.typescriptlang.org/) - 类型系统