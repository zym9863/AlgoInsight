# AlgoInsight Platform

**Language**: [中文](README.md) | **English**

<p align="left">
    <a href="https://github.com/zym9863/AlgoInsight/actions" target="_blank"><img alt="CI" src="https://img.shields.io/badge/CI-pending-lightgrey" /></a>
    <a href="LICENSE" target="_blank"><img alt="License" src="https://img.shields.io/badge/License-MIT-blue.svg" /></a>
    <img alt="Go Version" src="https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go" />
    <img alt="Svelte" src="https://img.shields.io/badge/Svelte-5-orange?logo=svelte" />
    <img alt="Docker" src="https://img.shields.io/badge/Docker-ready-2496ED?logo=docker" />
    <img alt="Status" src="https://img.shields.io/badge/status-active-success" />
</p>

An interactive platform focusing on algorithm learning, visualization and performance analysis. With a unified step-tracking abstraction, it renders execution flow of sorting / searching / graph algorithms in real time, supports multi-algorithm benchmarking, dataset generation, metric comparison and visual reporting.

> Goal: Lower the cognitive barrier of classic algorithms by turning abstract steps into structured trace + visual animation.

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

### Feature 3: Dataset Generation & Management
- **Modes**: random, nearly-sorted, reversed, duplicated, custom patterns
- **Preset datasets**: Bundled multi-size presets for reproducibility
- **Unified format**: Standard metadata wrapper for experiments

### Feature 4: Unified Algorithm Abstraction
- Interface: `Execute(data, tracker)`
- Step tracking: comparisons / swaps / accesses are normalized for consistent rendering
- Extensibility: Implement interface + register metadata

### Feature 5: Visualization Interaction
- Play / pause / single-step / jump
- Speed control (dynamic)
- Highlight current operation, compare pair, sorted region

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

> For more concrete examples see: `test_api.py` (includes sample calls & simple performance measurements).

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
- Selection Sort
- Shell Sort

### Searching Algorithms
- Linear Search
- Binary Search
- Hash Search

### Graph Algorithms
- Breadth-First Search (BFS)
- Depth-First Search (DFS)
- Shortest Path Algorithm (Dijkstra)
- Minimum Spanning Tree (Kruskal)
- Minimum Spanning Tree (Prim)
- Topological Sort

## 🧪 Local API Quick Test

Using bundled script:

```bash
python test_api.py
```

Manual calls (Windows CMD):

```cmd
curl -X GET http://localhost:8080/api/algorithms
curl -X POST http://localhost:8080/api/data/generate -H "Content-Type: application/json" -d "{\"dataType\":\"array\",\"size\":10,\"pattern\":\"random\"}"
```

Visualization example:

```bash
curl -X POST http://localhost:8080/api/visualize/execute \
    -H "Content-Type: application/json" \
    -d '{"algorithmId":"bubble_sort","data":[5,3,8,2],"parameters":{}}'
```

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

## 🧭 Roadmap

- [x] Graph algorithms: Dijkstra, Kruskal, Prim, Topological Sort
- [ ] More metrics: peak memory, GC stats
- [ ] Frontend: step slider & breakpoint markers
- [ ] Export steps (JSON / GIF / Video)
- [ ] i18n auto switch
- [ ] Public online playground demo
- [ ] CI: unit tests + build pipeline

## 🤝 Contributing

1. Fork the project
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Create a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ❓ FAQ

**Q: Why some algorithms miss complexity info?**  
A: Metadata not filled yet; contributions welcome.

**Q: How do I add a new algorithm?**  
A: Create file under `server/algorithms/<category>/`, implement interface (see sorting examples), register metadata; frontend auto fetches via `/api/algorithms`.

**Q: Custom data structures supported?**  
A: Currently arrays & basic graph structures; extensions planned (see Roadmap).

**Q: Is the Docker image production ready?**  
A: Geared toward learning/demo; add reverse proxy, logging, auth for production.

## 🙏 Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - Go Web Framework
- [Svelte](https://svelte.dev/) - Frontend Framework
- [Vite](https://vitejs.dev/) - Build Tool
- [TypeScript](https://www.typescriptlang.org/) - Type System