package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"gin/models"
	"sync"
	"time"
)

// BenchmarkService 性能测试服务
type BenchmarkService struct {
	tests            map[string]*models.BenchmarkTest
	mutex            sync.RWMutex
	algorithmService *AlgorithmService
}

// NewBenchmarkService 创建性能测试服务
func NewBenchmarkService() *BenchmarkService {
	return &BenchmarkService{
		tests:            make(map[string]*models.BenchmarkTest),
		algorithmService: NewAlgorithmService(),
	}
}

// RunBenchmarkTest 运行性能测试
func (s *BenchmarkService) RunBenchmarkTest(algorithmIDs []string, dataSizes []int, dataType string, testCount int, parameters interface{}) (string, error) {
	// 验证算法ID
	for _, algorithmID := range algorithmIDs {
		_, err := s.algorithmService.GetAlgorithm(algorithmID)
		if err != nil {
			return "", ErrInvalidAlgorithm
		}
	}

	// 验证数据规模
	for _, size := range dataSizes {
		if size > 10000 { // 可以从配置中读取
			return "", ErrDataSizeTooLarge
		}
	}

	// 创建测试
	testID := s.generateTestID()
	test := &models.BenchmarkTest{
		ID:           testID,
		Name:         "性能测试 " + testID,
		AlgorithmIDs: algorithmIDs,
		DataSizes:    dataSizes,
		DataType:     dataType,
		TestCount:    testCount,
		Parameters:   make(map[string]interface{}),
		Status:       models.TestStatusPending,
		CreatedAt:    time.Now(),
		Results:      make([]models.BenchmarkResult, 0),
	}

	// 保存测试
	s.mutex.Lock()
	s.tests[testID] = test
	s.mutex.Unlock()

	// 异步执行测试
	go s.executeBenchmarkTest(test)

	return testID, nil
}

// executeBenchmarkTest 执行性能测试
func (s *BenchmarkService) executeBenchmarkTest(test *models.BenchmarkTest) {
	// 更新状态为运行中
	s.mutex.Lock()
	test.Status = models.TestStatusRunning
	startTime := time.Now()
	test.StartedAt = &startTime
	s.mutex.Unlock()

	defer func() {
		// 测试完成
		s.mutex.Lock()
		test.Status = models.TestStatusCompleted
		completedTime := time.Now()
		test.CompletedAt = &completedTime
		s.mutex.Unlock()
	}()

	// 为每个算法和数据规模组合运行测试
	for _, algorithmID := range test.AlgorithmIDs {
		algorithm, err := s.algorithmService.GetAlgorithm(algorithmID)
		if err != nil {
			continue
		}

		for _, dataSize := range test.DataSizes {
			// 生成测试数据
			testData := s.generateTestData(test.DataType, dataSize)

			// 运行多次测试
			for i := 0; i < test.TestCount; i++ {
				result := s.runSingleTest(test.ID, algorithm, testData, test.DataType, dataSize, i)

				s.mutex.Lock()
				test.Results = append(test.Results, result)
				s.mutex.Unlock()
			}
		}
	}
}

// runSingleTest 运行单次测试
func (s *BenchmarkService) runSingleTest(testID string, algorithm interface{}, data interface{}, dataType string, dataSize int, runIndex int) models.BenchmarkResult {
	// 创建步骤追踪器（用于统计）
	tracker := models.NewStepTracker()

	// 记录开始时间
	startTime := time.Now()

	// 执行算法
	var err error
	var algorithmInfo *models.Algorithm

	// 类型断言获取算法接口
	if alg, ok := algorithm.(interface {
		Execute(data interface{}, tracker models.StepTracker) (interface{}, error)
		GetInfo() *models.Algorithm
	}); ok {
		_, err = alg.Execute(data, tracker)
		algorithmInfo = alg.GetInfo()
	} else {
		err = ErrInvalidAlgorithm
		algorithmInfo = &models.Algorithm{ID: "unknown", Name: "Unknown"}
	}

	// 记录结束时间
	executionTime := time.Since(startTime)

	// 获取统计信息
	stats := tracker.GetStats()

	// 构建结果
	result := models.BenchmarkResult{
		TestID:        testID,
		AlgorithmID:   algorithmInfo.ID,
		AlgorithmName: algorithmInfo.Name,
		DataSize:      dataSize,
		DataType:      dataType,
		RunIndex:      runIndex,
		ExecutionTime: executionTime,
		Operations:    int64(stats.Comparisons + stats.Swaps + stats.Moves),
		Comparisons:   int64(stats.Comparisons),
		Swaps:         int64(stats.Swaps),
		Success:       err == nil,
		Timestamp:     time.Now(),
		Metadata:      models.ResultMetadata{},
	}

	if err != nil {
		result.Error = err.Error()
	}

	return result
}

// generateTestData 生成测试数据
func (s *BenchmarkService) generateTestData(dataType string, size int) interface{} {
	// 根据数据类型生成简单的数据
	switch dataType {
	case models.DataTypeGraph:
		nodes := make([]models.GraphNode, size)
		edges := make([]models.GraphEdge, 0)
		for i := 0; i < size; i++ {
			id := fmt.Sprintf("node_%d", i)
			nodes[i] = models.GraphNode{ID: id, Label: id, Value: i}
			if i > 0 {
				prev := fmt.Sprintf("node_%d", i-1)
				edges = append(edges, models.GraphEdge{From: prev, To: id, Weight: 1})
			}
		}
		g := &models.GraphData{Nodes: nodes, Edges: edges, Type: "directed"}
		return g
	default:
		data := make([]interface{}, size)
		for i := 0; i < size; i++ {
			data[i] = i // 简单的递增序列
		}
		return data
	}
}

// GetBenchmarkResults 获取测试结果
func (s *BenchmarkService) GetBenchmarkResults(testID string) (*models.BenchmarkTest, error) {
	s.mutex.RLock()
	test, exists := s.tests[testID]
	s.mutex.RUnlock()

	if !exists {
		return nil, ErrTestNotFound
	}

	return test, nil
}

// CompareBenchmarkResults 对比测试结果
func (s *BenchmarkService) CompareBenchmarkResults(testIDs []string) (*models.ComparisonResult, error) {
	// 获取所有测试
	tests := make([]*models.BenchmarkTest, len(testIDs))
	for i, testID := range testIDs {
		test, err := s.GetBenchmarkResults(testID)
		if err != nil {
			return nil, err
		}
		tests[i] = test
	}

	// 构建对比结果
	comparison := &models.ComparisonResult{
		TestIDs:     testIDs,
		Algorithms:  make([]string, 0),
		DataSizes:   make([]int, 0),
		Comparisons: make([]models.AlgorithmComparison, 0),
		Summary:     models.ComparisonSummary{},
		Charts:      make(map[string]interface{}),
		CreatedAt:   time.Now(),
	}

	// 简化实现，返回基本对比结果
	return comparison, nil
}

// generateTestID 生成测试ID
func (s *BenchmarkService) generateTestID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// 全局性能测试服务实例
var benchmarkService *BenchmarkService

// RunBenchmarkTest 运行性能测试（全局函数）
func RunBenchmarkTest(algorithmIDs []string, dataSizes []int, dataType string, testCount int, parameters interface{}) (string, error) {
	if benchmarkService == nil {
		benchmarkService = NewBenchmarkService()
	}
	return benchmarkService.RunBenchmarkTest(algorithmIDs, dataSizes, dataType, testCount, parameters)
}

// GetBenchmarkResults 获取测试结果（全局函数）
func GetBenchmarkResults(testID string) (*models.BenchmarkTest, error) {
	if benchmarkService == nil {
		benchmarkService = NewBenchmarkService()
	}
	return benchmarkService.GetBenchmarkResults(testID)
}

// CompareBenchmarkResults 对比测试结果（全局函数）
func CompareBenchmarkResults(testIDs []string) (*models.ComparisonResult, error) {
	if benchmarkService == nil {
		benchmarkService = NewBenchmarkService()
	}
	return benchmarkService.CompareBenchmarkResults(testIDs)
}
