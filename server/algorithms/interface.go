package algorithms

import (
	"errors"
	"gin/models"
)

// Algorithm 算法接口
type Algorithm interface {
	// Execute 执行算法
	Execute(data interface{}, tracker models.StepTracker) (interface{}, error)

	// GetInfo 获取算法信息
	GetInfo() *models.Algorithm

	// ValidateInput 验证输入数据
	ValidateInput(data interface{}) error

	// GetCategory 获取算法类别
	GetCategory() string

	// GetComplexity 获取复杂度信息
	GetComplexity() ComplexityInfo
}

// ComplexityInfo 复杂度信息
type ComplexityInfo struct {
	TimeComplexity  ComplexityCase `json:"timeComplexity"`  // 时间复杂度
	SpaceComplexity ComplexityCase `json:"spaceComplexity"` // 空间复杂度
}

// ComplexityCase 复杂度情况
type ComplexityCase struct {
	Best    string `json:"best"`    // 最佳情况
	Average string `json:"average"` // 平均情况
	Worst   string `json:"worst"`   // 最坏情况
}

// SortingAlgorithm 排序算法接口
type SortingAlgorithm interface {
	Algorithm

	// Sort 排序方法
	Sort(data []interface{}, tracker models.StepTracker) error

	// IsStable 是否稳定排序
	IsStable() bool

	// IsInPlace 是否原地排序
	IsInPlace() bool

	// IsAdaptive 是否自适应
	IsAdaptive() bool
}

// SearchingAlgorithm 搜索算法接口
type SearchingAlgorithm interface {
	Algorithm

	// Search 搜索方法
	Search(data []interface{}, target interface{}, tracker models.StepTracker) (int, error)

	// RequiresSorted 是否需要已排序数据
	RequiresSorted() bool
}

// GraphAlgorithm 图算法接口
type GraphAlgorithm interface {
	Algorithm

	// ProcessGraph 处理图
	ProcessGraph(graph *models.GraphData, tracker models.StepTracker) (interface{}, error)

	// GetGraphType 获取支持的图类型
	GetGraphType() string // "directed", "undirected", "both"
}

// TreeAlgorithm 树算法接口
type TreeAlgorithm interface {
	Algorithm

	// ProcessTree 处理树
	ProcessTree(tree *models.TreeData, tracker models.StepTracker) (interface{}, error)

	// GetTreeType 获取支持的树类型
	GetTreeType() string // "binary", "n-ary", "both"
}

// BaseAlgorithm 基础算法结构
type BaseAlgorithm struct {
	ID              string             `json:"id"`
	Name            string             `json:"name"`
	Category        string             `json:"category"`
	Description     string             `json:"description"`
	TimeComplexity  string             `json:"timeComplexity"`
	SpaceComplexity string             `json:"spaceComplexity"`
	Parameters      []models.Parameter `json:"parameters"`
	Stable          bool               `json:"stable"`
	InPlace         bool               `json:"inPlace"`
	Adaptive        bool               `json:"adaptive"`
}

// GetInfo 获取算法信息
func (b *BaseAlgorithm) GetInfo() *models.Algorithm {
	return &models.Algorithm{
		ID:              b.ID,
		Name:            b.Name,
		Category:        b.Category,
		Description:     b.Description,
		TimeComplexity:  b.TimeComplexity,
		SpaceComplexity: b.SpaceComplexity,
		Parameters:      b.Parameters,
		Stable:          b.Stable,
		InPlace:         b.InPlace,
		Adaptive:        b.Adaptive,
	}
}

// GetCategory 获取算法类别
func (b *BaseAlgorithm) GetCategory() string {
	return b.Category
}

// ValidateInput 基础输入验证
func (b *BaseAlgorithm) ValidateInput(data interface{}) error {
	if data == nil {
		return ErrInvalidInput
	}
	return nil
}

// AlgorithmRegistry 算法注册表
type AlgorithmRegistry struct {
	algorithms map[string]Algorithm
}

// NewAlgorithmRegistry 创建算法注册表
func NewAlgorithmRegistry() *AlgorithmRegistry {
	return &AlgorithmRegistry{
		algorithms: make(map[string]Algorithm),
	}
}

// Register 注册算法
func (r *AlgorithmRegistry) Register(algorithm Algorithm) {
	r.algorithms[algorithm.GetInfo().ID] = algorithm
}

// Get 获取算法
func (r *AlgorithmRegistry) Get(id string) (Algorithm, bool) {
	algorithm, exists := r.algorithms[id]
	return algorithm, exists
}

// GetAll 获取所有算法
func (r *AlgorithmRegistry) GetAll() []Algorithm {
	algorithms := make([]Algorithm, 0, len(r.algorithms))
	for _, algorithm := range r.algorithms {
		algorithms = append(algorithms, algorithm)
	}
	return algorithms
}

// GetByCategory 按类别获取算法
func (r *AlgorithmRegistry) GetByCategory(category string) []Algorithm {
	algorithms := make([]Algorithm, 0)
	for _, algorithm := range r.algorithms {
		if algorithm.GetCategory() == category {
			algorithms = append(algorithms, algorithm)
		}
	}
	return algorithms
}

// 错误定义
var (
	ErrInvalidInput     = errors.New("输入数据无效")
	ErrUnsupportedType  = errors.New("不支持的数据类型")
	ErrExecutionTimeout = errors.New("算法执行超时")
)
