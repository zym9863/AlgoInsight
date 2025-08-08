package models

// Algorithm 算法定义结构
type Algorithm struct {
	ID              string      `json:"id"`              // 算法唯一标识
	Name            string      `json:"name"`            // 算法名称
	Category        string      `json:"category"`        // 算法类别 (sorting, searching, graph, etc.)
	Description     string      `json:"description"`     // 算法描述
	TimeComplexity  string      `json:"timeComplexity"`  // 时间复杂度
	SpaceComplexity string      `json:"spaceComplexity"` // 空间复杂度
	Parameters      []Parameter `json:"parameters"`      // 算法参数
	Stable          bool        `json:"stable"`          // 是否稳定（排序算法）
	InPlace         bool        `json:"inPlace"`         // 是否原地算法
	Adaptive        bool        `json:"adaptive"`        // 是否自适应
}

// Parameter 算法参数定义
type Parameter struct {
	Name         string      `json:"name"`              // 参数名称
	Type         string      `json:"type"`              // 参数类型 (int, float, string, bool)
	Description  string      `json:"description"`       // 参数描述
	DefaultValue interface{} `json:"defaultValue"`      // 默认值
	Required     bool        `json:"required"`          // 是否必需
	Min          interface{} `json:"min,omitempty"`     // 最小值（数值类型）
	Max          interface{} `json:"max,omitempty"`     // 最大值（数值类型）
	Options      []string    `json:"options,omitempty"` // 可选值（枚举类型）
}

// AlgorithmCategory 算法类别常量
const (
	CategorySorting       = "sorting"
	CategorySearching     = "searching"
	CategoryGraph         = "graph"
	CategoryTree          = "tree"
	CategoryDynamicProg   = "dynamic_programming"
	CategoryGreedy        = "greedy"
	CategoryBacktracking  = "backtracking"
	CategoryDivideConquer = "divide_conquer"
)

// GetAlgorithmCategories 获取所有算法类别
func GetAlgorithmCategories() []string {
	return []string{
		CategorySorting,
		CategorySearching,
		CategoryGraph,
		CategoryTree,
		CategoryDynamicProg,
		CategoryGreedy,
		CategoryBacktracking,
		CategoryDivideConquer,
	}
}

// AlgorithmRegistry 算法注册表
type AlgorithmRegistry struct {
	algorithms map[string]*Algorithm
}

// NewAlgorithmRegistry 创建新的算法注册表
func NewAlgorithmRegistry() *AlgorithmRegistry {
	return &AlgorithmRegistry{
		algorithms: make(map[string]*Algorithm),
	}
}

// Register 注册算法
func (r *AlgorithmRegistry) Register(algorithm *Algorithm) {
	r.algorithms[algorithm.ID] = algorithm
}

// Get 获取算法
func (r *AlgorithmRegistry) Get(id string) (*Algorithm, bool) {
	algorithm, exists := r.algorithms[id]
	return algorithm, exists
}

// GetAll 获取所有算法
func (r *AlgorithmRegistry) GetAll() []*Algorithm {
	algorithms := make([]*Algorithm, 0, len(r.algorithms))
	for _, algorithm := range r.algorithms {
		algorithms = append(algorithms, algorithm)
	}
	return algorithms
}

// GetByCategory 按类别获取算法
func (r *AlgorithmRegistry) GetByCategory(category string) []*Algorithm {
	algorithms := make([]*Algorithm, 0)
	for _, algorithm := range r.algorithms {
		if algorithm.Category == category {
			algorithms = append(algorithms, algorithm)
		}
	}
	return algorithms
}
