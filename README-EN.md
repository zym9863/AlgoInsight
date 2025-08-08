# AlgoInsight Platform

**Language**: [中文](README.md) | **English**

An interactive algorithm learning and performance analysis tool that supports algorithm visualization and performance comparison.

## 🚀 Features

### Feature 1: Interactive Algorithm Visualization
- **Algorithm Selection**: Support for various sorting, searching, graph algorithms, etc.
- **Data Input**: Custom input, random generation, preset datasets
- **Animation Display**: Real-time visualization of algorithm execution process
- **Step Control**: Play, pause, single-step execution, speed adjustment
- **Detailed Information**: Display comparison, swap, move operation details

### Feature 2: Algorithm Performance Testing & Comparison
- **Multi-Algorithm Testing**: Test multiple algorithm performances simultaneously
- **Data Scale Configuration**: Customize test data size and type
- **Performance Metrics**: Execution time, memory usage, operation count statistics
- **Comparison Analysis**: Chart visualization of performance differences
- **Detailed Reports**: Generate comprehensive performance analysis reports

## 🛠️ Technology Stack

- **Backend**: Go + Gin Framework
- **Frontend**: Svelte + TypeScript + Vite
- **Package Manager**: pnpm
- **Deployment**: Docker + Docker Compose

## 📦 Quick Start

### Using Docker (Recommended)

1. **Clone the project**
```bash
git clone https://github.com/zym9863/AlgoInsight.git
cd AlgoInsight
```

2. **Build and start**
```bash
docker-compose up --build
```

3. **Access the application**
- Open browser and visit: http://localhost:8080

### Development Environment

#### Backend Development
```bash
cd server
go mod tidy
go run main.go
```

#### Frontend Development
```bash
cd web
pnpm install
pnpm dev
```

#### Using Docker Development Environment
```bash
# Start development environment
docker-compose --profile dev up --build

# Backend: http://localhost:8080
# Frontend: http://localhost:5173
```

## 🏗️ Project Structure

```
AlgoInsight/
├── Dockerfile                 # Production Docker configuration
├── docker-compose.yml         # Docker Compose configuration
├── README.md                  # Project documentation (Chinese)
├── README-EN.md              # Project documentation (English)
├── server/                    # Go backend
│   ├── main.go               # Server entry point
│   ├── config/               # Configuration management
│   ├── handlers/             # HTTP handlers
│   ├── models/               # Data models
│   ├── services/             # Business logic
│   ├── algorithms/           # Algorithm implementations
│   │   ├── sorting/          # Sorting algorithms
│   │   ├── searching/        # Searching algorithms
│   │   └── graph/            # Graph algorithms
│   └── utils/                # Utility functions
├── web/                      # Svelte frontend
│   ├── src/
│   │   ├── components/       # UI components
│   │   │   ├── algorithm/    # Algorithm-related components
│   │   │   ├── ui/           # Common UI components
│   │   │   └── visualization/ # Visualization components
│   │   ├── stores/           # State management
│   │   ├── services/         # API services
│   │   ├── types/            # TypeScript types
│   │   └── utils/            # Utility functions
│   ├── public/
│   └── package.json
└── docs/                     # Documentation
    └── architecture.md       # Architecture design documentation
```

## 🔧 API Endpoints

### Algorithm Management
- `GET /api/algorithms` - Get all algorithms
- `GET /api/algorithms/category/{category}` - Get algorithms by category
- `GET /api/algorithms/info/{id}` - Get algorithm detailed information

### Visualization
- `POST /api/visualize/execute` - Execute algorithm visualization
- `GET /api/visualize/step/{sessionId}/{stepId}` - Get visualization step
- `POST /api/visualize/reset` - Reset visualization state

### Performance Testing
- `POST /api/benchmark/run` - Run performance test
- `GET /api/benchmark/results/{testId}` - Get test results
- `POST /api/benchmark/compare` - Compare test results

### Data Generation
- `POST /api/data/generate` - Generate test data
- `GET /api/data/presets` - Get preset data

## 🎯 Supported Algorithms

### Sorting Algorithms
- Bubble Sort
- Quick Sort
- Merge Sort
- Heap Sort
- Insertion Sort

### Searching Algorithms
- Linear Search
- Binary Search
- Hash Search

### Graph Algorithms
- Breadth-First Search (BFS)
- Depth-First Search (DFS)
- Shortest Path (Dijkstra)
- Minimum Spanning Tree (Kruskal)

## 🚀 Deployment

### Production Environment Deployment

1. **Build Docker image**
```bash
docker build -t algoinsight .
```

2. **Run container**
```bash
docker run -p 8080:8080 algoinsight
```

3. **Using Docker Compose**
```bash
docker-compose up -d
```

### Environment Variable Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | Service port |
| `ENVIRONMENT` | development | Runtime environment |
| `MAX_EXECUTION_TIME` | 30 | Maximum execution time (seconds) |
| `MAX_DATA_SIZE` | 10000 | Maximum data scale |
| `BENCHMARK_TIMEOUT` | 60 | Performance test timeout (seconds) |
| `MAX_CONCURRENT_TESTS` | 5 | Maximum concurrent tests |

## 🤝 Contributing

1. Fork the project
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Create a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - Go Web Framework
- [Svelte](https://svelte.dev/) - Frontend Framework
- [Vite](https://vitejs.dev/) - Build Tool
- [TypeScript](https://www.typescriptlang.org/) - Type System