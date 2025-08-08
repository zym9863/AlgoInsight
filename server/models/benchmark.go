package models

import (
	"time"
)

// BenchmarkTest 性能测试定义
type BenchmarkTest struct {
	ID           string                 `json:"id"`           // 测试ID
	Name         string                 `json:"name"`         // 测试名称
	AlgorithmIDs []string               `json:"algorithmIds"` // 参与测试的算法ID列表
	DataSizes    []int                  `json:"dataSizes"`    // 测试数据规模
	DataType     string                 `json:"dataType"`     // 数据类型
	TestCount    int                    `json:"testCount"`    // 每个配置的测试次数
	Parameters   map[string]interface{} `json:"parameters"`   // 测试参数
	Status       string                 `json:"status"`       // 测试状态
	CreatedAt    time.Time              `json:"createdAt"`    // 创建时间
	StartedAt    *time.Time             `json:"startedAt"`    // 开始时间
	CompletedAt  *time.Time             `json:"completedAt"`  // 完成时间
	Results      []BenchmarkResult      `json:"results"`      // 测试结果
	Error        string                 `json:"error,omitempty"` // 错误信息
}

// BenchmarkResult 单次性能测试结果
type BenchmarkResult struct {
	TestID        string        `json:"testId"`        // 测试ID
	AlgorithmID   string        `json:"algorithmId"`   // 算法ID
	AlgorithmName string        `json:"algorithmName"` // 算法名称
	DataSize      int           `json:"dataSize"`      // 数据规模
	DataType      string        `json:"dataType"`      // 数据类型
	RunIndex      int           `json:"runIndex"`      // 运行索引（第几次运行）
	ExecutionTime time.Duration `json:"executionTime"` // 执行时间
	MemoryUsage   int64         `json:"memoryUsage"`   // 内存使用（字节）
	Operations    int64         `json:"operations"`    // 操作次数
	Comparisons   int64         `json:"comparisons"`   // 比较次数
	Swaps         int64         `json:"swaps"`         // 交换次数
	Success       bool          `json:"success"`       // 是否成功
	Error         string        `json:"error,omitempty"` // 错误信息
	Timestamp     time.Time     `json:"timestamp"`     // 时间戳
	Metadata      ResultMetadata `json:"metadata"`     // 结果元数据
}

// ResultMetadata 结果元数据
type ResultMetadata struct {
	CPUUsage      float64 `json:"cpuUsage"`      // CPU使用率
	GCCount       int     `json:"gcCount"`       // GC次数
	GCTime        time.Duration `json:"gcTime"`  // GC时间
	AllocatedMem  int64   `json:"allocatedMem"`  // 分配的内存
	SystemMem     int64   `json:"systemMem"`     // 系统内存
	GoRoutines    int     `json:"goRoutines"`    // 协程数量
}

// BenchmarkSummary 性能测试汇总
type BenchmarkSummary struct {
	TestID       string                    `json:"testId"`       // 测试ID
	AlgorithmID  string                    `json:"algorithmId"`  // 算法ID
	DataSize     int                       `json:"dataSize"`     // 数据规模
	RunCount     int                       `json:"runCount"`     // 运行次数
	AvgTime      time.Duration             `json:"avgTime"`      // 平均执行时间
	MinTime      time.Duration             `json:"minTime"`      // 最短执行时间
	MaxTime      time.Duration             `json:"maxTime"`      // 最长执行时间
	StdDev       time.Duration             `json:"stdDev"`       // 标准差
	AvgMemory    int64                     `json:"avgMemory"`    // 平均内存使用
	AvgOps       int64                     `json:"avgOps"`       // 平均操作次数
	SuccessRate  float64                   `json:"successRate"`  // 成功率
	Performance  PerformanceMetrics        `json:"performance"`  // 性能指标
}

// PerformanceMetrics 性能指标
type PerformanceMetrics struct {
	TimeComplexity  string  `json:"timeComplexity"`  // 实际时间复杂度
	SpaceComplexity string  `json:"spaceComplexity"` // 实际空间复杂度
	Efficiency      float64 `json:"efficiency"`      // 效率评分 (0-100)
	Stability       float64 `json:"stability"`       // 稳定性评分 (0-100)
	Scalability     float64 `json:"scalability"`     // 可扩展性评分 (0-100)
}

// ComparisonResult 算法对比结果
type ComparisonResult struct {
	TestIDs     []string                   `json:"testIds"`     // 对比的测试ID
	Algorithms  []string                   `json:"algorithms"`  // 对比的算法
	DataSizes   []int                      `json:"dataSizes"`   // 数据规模
	Comparisons []AlgorithmComparison      `json:"comparisons"` // 详细对比
	Summary     ComparisonSummary          `json:"summary"`     // 对比汇总
	Charts      map[string]interface{}     `json:"charts"`      // 图表数据
	CreatedAt   time.Time                  `json:"createdAt"`   // 创建时间
}

// AlgorithmComparison 算法对比
type AlgorithmComparison struct {
	AlgorithmID   string                 `json:"algorithmId"`   // 算法ID
	AlgorithmName string                 `json:"algorithmName"` // 算法名称
	Results       map[int]BenchmarkSummary `json:"results"`     // 按数据规模分组的结果
	OverallScore  float64                `json:"overallScore"`  // 综合评分
	Strengths     []string               `json:"strengths"`     // 优势
	Weaknesses    []string               `json:"weaknesses"`    // 劣势
}

// ComparisonSummary 对比汇总
type ComparisonSummary struct {
	BestTime       string  `json:"bestTime"`       // 时间最优算法
	BestMemory     string  `json:"bestMemory"`     // 内存最优算法
	BestStability  string  `json:"bestStability"`  // 稳定性最优算法
	BestScalability string `json:"bestScalability"` // 可扩展性最优算法
	Recommendation string  `json:"recommendation"` // 推荐算法
	Reason         string  `json:"reason"`         // 推荐理由
}

// TestStatus 测试状态常量
const (
	TestStatusPending   = "pending"
	TestStatusRunning   = "running"
	TestStatusCompleted = "completed"
	TestStatusFailed    = "failed"
	TestStatusCancelled = "cancelled"
)

// DataType 数据类型常量
const (
	DataTypeArray  = "array"
	DataTypeGraph  = "graph"
	DataTypeTree   = "tree"
	DataTypeString = "string"
	DataTypeMatrix = "matrix"
)

// BenchmarkConfig 性能测试配置
type BenchmarkConfig struct {
	MaxExecutionTime time.Duration `json:"maxExecutionTime"` // 最大执行时间
	MaxMemoryUsage   int64         `json:"maxMemoryUsage"`   // 最大内存使用
	MaxDataSize      int           `json:"maxDataSize"`      // 最大数据规模
	TimeoutDuration  time.Duration `json:"timeoutDuration"`  // 超时时间
	WarmupRuns       int           `json:"warmupRuns"`       // 预热运行次数
	EnableProfiling  bool          `json:"enableProfiling"`  // 是否启用性能分析
}

// GetDefaultBenchmarkConfig 获取默认性能测试配置
func GetDefaultBenchmarkConfig() *BenchmarkConfig {
	return &BenchmarkConfig{
		MaxExecutionTime: 30 * time.Second,
		MaxMemoryUsage:   100 * 1024 * 1024, // 100MB
		MaxDataSize:      10000,
		TimeoutDuration:  60 * time.Second,
		WarmupRuns:       3,
		EnableProfiling:  false,
	}
}
