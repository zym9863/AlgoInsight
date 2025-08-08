package models

import "time"

// DataPreset 预设数据集
type DataPreset struct {
	ID          string      `json:"id"`          // 预设ID
	Name        string      `json:"name"`        // 预设名称
	Description string      `json:"description"` // 描述
	DataType    string      `json:"dataType"`    // 数据类型
	Size        int         `json:"size"`        // 数据大小
	Pattern     string      `json:"pattern"`     // 数据模式
	Data        interface{} `json:"data"`        // 实际数据
	Tags        []string    `json:"tags"`        // 标签
	CreatedAt   time.Time   `json:"createdAt"`   // 创建时间
}

// DataGenerationRequest 数据生成请求
type DataGenerationRequest struct {
	DataType   string                 `json:"dataType"`   // 数据类型
	Size       int                    `json:"size"`       // 数据大小
	Pattern    string                 `json:"pattern"`    // 数据模式
	Parameters map[string]interface{} `json:"parameters"` // 生成参数
}

// DataGenerationResult 数据生成结果
type DataGenerationResult struct {
	DataType   string      `json:"dataType"`   // 数据类型
	Size       int         `json:"size"`       // 实际大小
	Pattern    string      `json:"pattern"`    // 数据模式
	Data       interface{} `json:"data"`       // 生成的数据
	Metadata   DataMetadata `json:"metadata"`  // 数据元数据
	GeneratedAt time.Time  `json:"generatedAt"` // 生成时间
}

// DataMetadata 数据元数据
type DataMetadata struct {
	MinValue    interface{} `json:"minValue,omitempty"`    // 最小值
	MaxValue    interface{} `json:"maxValue,omitempty"`    // 最大值
	Range       interface{} `json:"range,omitempty"`       // 值域范围
	Duplicates  int         `json:"duplicates,omitempty"`  // 重复元素数量
	Sorted      bool        `json:"sorted,omitempty"`      // 是否已排序
	Reversed    bool        `json:"reversed,omitempty"`    // 是否逆序
	Unique      bool        `json:"unique,omitempty"`      // 是否唯一
	Distribution string     `json:"distribution,omitempty"` // 分布类型
}

// ArrayData 数组数据结构
type ArrayData struct {
	Values []interface{} `json:"values"` // 数组值
	Type   string        `json:"type"`   // 元素类型 (int, float, string)
}

// GraphData 图数据结构
type GraphData struct {
	Nodes []GraphNode `json:"nodes"` // 节点列表
	Edges []GraphEdge `json:"edges"` // 边列表
	Type  string      `json:"type"`  // 图类型 (directed, undirected)
}

// GraphNode 图节点
type GraphNode struct {
	ID    string      `json:"id"`    // 节点ID
	Label string      `json:"label"` // 节点标签
	Value interface{} `json:"value"` // 节点值
	X     float64     `json:"x"`     // X坐标
	Y     float64     `json:"y"`     // Y坐标
}

// GraphEdge 图边
type GraphEdge struct {
	From   string      `json:"from"`   // 起始节点ID
	To     string      `json:"to"`     // 目标节点ID
	Weight interface{} `json:"weight"` // 边权重
	Label  string      `json:"label"`  // 边标签
}

// TreeData 树数据结构
type TreeData struct {
	Root *TreeNode `json:"root"` // 根节点
	Type string    `json:"type"` // 树类型 (binary, n-ary)
}

// TreeNode 树节点
type TreeNode struct {
	ID       string      `json:"id"`       // 节点ID
	Value    interface{} `json:"value"`    // 节点值
	Children []*TreeNode `json:"children"` // 子节点
	Parent   *TreeNode   `json:"parent,omitempty"` // 父节点（序列化时忽略）
	Left     *TreeNode   `json:"left,omitempty"`   // 左子节点（二叉树）
	Right    *TreeNode   `json:"right,omitempty"`  // 右子节点（二叉树）
	X        float64     `json:"x"`        // X坐标
	Y        float64     `json:"y"`        // Y坐标
	Level    int         `json:"level"`    // 层级
}

// MatrixData 矩阵数据结构
type MatrixData struct {
	Values [][]interface{} `json:"values"` // 矩阵值
	Rows   int             `json:"rows"`   // 行数
	Cols   int             `json:"cols"`   // 列数
	Type   string          `json:"type"`   // 元素类型
}

// DataPattern 数据模式常量
const (
	PatternRandom       = "random"        // 随机数据
	PatternSorted       = "sorted"        // 已排序数据
	PatternReversed     = "reversed"      // 逆序数据
	PatternNearlySorted = "nearly_sorted" // 近似排序数据
	PatternFewUnique    = "few_unique"    // 少量唯一值
	PatternManyDuplicates = "many_duplicates" // 大量重复值
	PatternWorstCase    = "worst_case"    // 最坏情况
	PatternBestCase     = "best_case"     // 最佳情况
	PatternAverageCase  = "average_case"  // 平均情况
)

// GetDataPatterns 获取所有数据模式
func GetDataPatterns() []string {
	return []string{
		PatternRandom,
		PatternSorted,
		PatternReversed,
		PatternNearlySorted,
		PatternFewUnique,
		PatternManyDuplicates,
		PatternWorstCase,
		PatternBestCase,
		PatternAverageCase,
	}
}

// GetDataTypes 获取所有数据类型
func GetDataTypes() []string {
	return []string{
		DataTypeArray,
		DataTypeGraph,
		DataTypeTree,
		DataTypeString,
		DataTypeMatrix,
	}
}

// ValidateDataType 验证数据类型
func ValidateDataType(dataType string) bool {
	validTypes := GetDataTypes()
	for _, validType := range validTypes {
		if dataType == validType {
			return true
		}
	}
	return false
}

// ValidateDataPattern 验证数据模式
func ValidateDataPattern(pattern string) bool {
	validPatterns := GetDataPatterns()
	for _, validPattern := range validPatterns {
		if pattern == validPattern {
			return true
		}
	}
	return false
}
