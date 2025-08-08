# 算法洞察平台架构设计

## 项目概述

算法洞察平台是一个交互式算法学习和性能分析工具，支持算法可视化和性能对比功能。

### 技术栈
- **后端**: Go + Gin框架
- **前端**: Svelte + TypeScript + Vite
- **包管理**: pnpm
- **部署**: Docker单容器部署

## 系统架构

```
Frontend (Svelte) ←→ Backend (Gin) ←→ Algorithm Engine
     ↓                    ↓                ↓
  UI Components      REST APIs        Algorithm Modules
  - Visualizer       - /api/algorithms    - Sorting
  - Controls         - /api/visualize     - Searching  
  - Charts           - /api/benchmark     - Graph
  - Config           - /api/compare       - Data Structures
```

## API接口设计

### 算法管理
- `GET /api/algorithms` - 获取所有可用算法列表
- `GET /api/algorithms/{category}` - 按类别获取算法
- `GET /api/algorithms/{id}/info` - 获取算法详细信息

### 可视化
- `POST /api/visualize/execute` - 执行算法并返回步骤数据
- `GET /api/visualize/step/{sessionId}/{stepId}` - 获取特定步骤状态
- `POST /api/visualize/reset` - 重置可视化状态

### 性能测试
- `POST /api/benchmark/run` - 运行性能测试
- `GET /api/benchmark/results/{testId}` - 获取测试结果
- `POST /api/benchmark/compare` - 对比多个算法性能

### 数据生成
- `POST /api/data/generate` - 生成测试数据
- `GET /api/data/presets` - 获取预设数据集

## 数据模型

### Algorithm - 算法定义
```go
type Algorithm struct {
    ID              string      `json:"id"`
    Name            string      `json:"name"`
    Category        string      `json:"category"`
    Description     string      `json:"description"`
    TimeComplexity  string      `json:"timeComplexity"`
    SpaceComplexity string      `json:"spaceComplexity"`
    Parameters      []Parameter `json:"parameters"`
}
```

### VisualizationStep - 可视化步骤
```go
type VisualizationStep struct {
    StepID      int           `json:"stepId"`
    Description string        `json:"description"`
    Data        interface{}   `json:"data"`
    Highlights  []int         `json:"highlights"`
    Comparisons []Comparison  `json:"comparisons"`
    Operations  []Operation   `json:"operations"`
}
```

### BenchmarkResult - 性能测试结果
```go
type BenchmarkResult struct {
    AlgorithmID   string        `json:"algorithmId"`
    DataSize      int           `json:"dataSize"`
    ExecutionTime time.Duration `json:"executionTime"`
    MemoryUsage   int64         `json:"memoryUsage"`
    Operations    int64         `json:"operations"`
    Timestamp     time.Time     `json:"timestamp"`
}
```

## 前端组件架构

### 组件层次结构
```
App.svelte
├── Header.svelte
├── AlgorithmSelector.svelte
├── DataInput.svelte
├── VisualizationPanel.svelte
│   ├── Canvas.svelte
│   ├── Controls.svelte
│   └── StepInfo.svelte
├── BenchmarkPanel.svelte
│   ├── TestConfig.svelte
│   ├── ResultsChart.svelte
│   └── ComparisonTable.svelte
└── Footer.svelte
```

### 状态管理
使用Svelte stores管理全局状态：
- `selectedAlgorithm` - 当前选择的算法
- `visualizationSteps` - 可视化步骤数据
- `benchmarkResults` - 性能测试结果

## 算法模块组织

### 算法接口
```go
type Algorithm interface {
    Execute(data interface{}, tracker *StepTracker) (interface{}, error)
    GetInfo() AlgorithmInfo
    ValidateInput(data interface{}) error
}

type StepTracker interface {
    AddStep(description string, data interface{}, highlights []int)
    AddComparison(i, j int, result int)
    AddOperation(op string, indices []int)
}
```

### 模块结构
```
server/algorithms/
├── interface.go          // 算法接口定义
├── sorting/              // 排序算法
├── searching/            // 搜索算法
├── graph/                // 图算法
└── data_structures/      // 数据结构
```

## 项目目录结构

```
AlgoInsight/
├── Dockerfile                 // 统一部署配置
├── docker-compose.yml         // 开发环境配置
├── server/                    // Go后端
│   ├── main.go               // 服务器入口
│   ├── config/               // 配置管理
│   ├── handlers/             // HTTP处理器
│   ├── models/               // 数据模型
│   ├── services/             // 业务逻辑
│   ├── algorithms/           // 算法实现
│   └── utils/                // 工具函数
├── web/                      // Svelte前端
│   ├── src/
│   │   ├── components/       // UI组件
│   │   ├── stores/           // 状态管理
│   │   ├── services/         // API服务
│   │   └── types/            // TypeScript类型
│   └── public/
└── docs/                     // 文档
```

## 部署策略

### Docker多阶段构建
1. 构建前端静态文件
2. 构建Go后端二进制文件
3. 最终镜像包含后端服务和前端静态文件

### 开发流程
- 后端开发：`go run main.go`
- 前端开发：`pnpm dev`
- 生产构建：`docker build`

## 核心功能实现计划

### 功能一：交互式算法可视化
1. 算法执行步骤追踪
2. 实时动画展示
3. 步骤控制（播放/暂停/单步）
4. 数据状态高亮显示

### 功能二：算法性能评测与对比
1. 多算法并行测试
2. 性能指标收集（时间、内存、操作次数）
3. 图表化结果展示
4. 详细对比分析报告
