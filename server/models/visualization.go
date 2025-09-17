package models

import (
	"time"
)

// VisualizationSession 可视化会话
type VisualizationSession struct {
	ID          string              `json:"id"`              // 会话ID
	AlgorithmID string              `json:"algorithmId"`     // 算法ID
	InputData   interface{}         `json:"inputData"`       // 输入数据
	Parameters  interface{}         `json:"parameters"`      // 算法参数
	Steps       []VisualizationStep `json:"steps"`           // 执行步骤
	Status      string              `json:"status"`          // 会话状态 (running, completed, error)
	CreatedAt   time.Time           `json:"createdAt"`       // 创建时间
	CompletedAt *time.Time          `json:"completedAt"`     // 完成时间
	Error       string              `json:"error,omitempty"` // 错误信息
}

// VisualizationStep 可视化步骤
type VisualizationStep struct {
	StepID      int          `json:"stepId"`      // 步骤ID
	Description string       `json:"description"` // 步骤描述
	Data        interface{}  `json:"data"`        // 当前数据状态
	Highlights  []int        `json:"highlights"`  // 高亮元素索引
	Comparisons []Comparison `json:"comparisons"` // 比较操作
	Operations  []Operation  `json:"operations"`  // 执行的操作
	Metadata    StepMetadata `json:"metadata"`    // 步骤元数据
}

// Comparison 比较操作
type Comparison struct {
	Index1 int    `json:"index1"` // 比较元素1的索引
	Index2 int    `json:"index2"` // 比较元素2的索引
	Result int    `json:"result"` // 比较结果 (-1: <, 0: =, 1: >)
	Type   string `json:"type"`   // 比较类型 (value, key, etc.)
}

// Operation 操作
type Operation struct {
	Type        string        `json:"type"`        // 操作类型 (swap, move, insert, delete, etc.)
	Indices     []int         `json:"indices"`     // 涉及的索引
	Values      []interface{} `json:"values"`      // 涉及的值
	Description string        `json:"description"` // 操作描述
}

// StepMetadata 步骤元数据
type StepMetadata struct {
	ExecutionTime time.Duration `json:"executionTime"` // 执行时间
	MemoryUsage   int64         `json:"memoryUsage"`   // 内存使用
	Complexity    string        `json:"complexity"`    // 当前复杂度
	Phase         string        `json:"phase"`         // 算法阶段
	Notes         []string      `json:"notes"`         // 备注信息
}

// VisualizationResult 可视化结果
type VisualizationResult struct {
	SessionID     string              `json:"sessionId"`     // 会话ID
	AlgorithmID   string              `json:"algorithmId"`   // 算法ID
	InputData     interface{}         `json:"inputData"`     // 输入数据
	OutputData    interface{}         `json:"outputData"`    // 输出数据
	Steps         []VisualizationStep `json:"steps"`         // 所有步骤
	TotalSteps    int                 `json:"totalSteps"`    // 总步骤数
	ExecutionTime time.Duration       `json:"executionTime"` // 总执行时间
	MemoryUsage   int64               `json:"memoryUsage"`   // 内存使用峰值
	Statistics    ExecutionStats      `json:"statistics"`    // 执行统计
}

// ExecutionStats 执行统计
type ExecutionStats struct {
	Comparisons int `json:"comparisons"` // 比较次数
	Swaps       int `json:"swaps"`       // 交换次数
	Moves       int `json:"moves"`       // 移动次数
	Accesses    int `json:"accesses"`    // 访问次数
}

// SessionStatus 会话状态常量
const (
	StatusRunning   = "running"
	StatusCompleted = "completed"
	StatusError     = "error"
	StatusCancelled = "cancelled"
)

// OperationType 操作类型常量
const (
	OpTypeSwap      = "swap"
	OpTypeMove      = "move"
	OpTypeInsert    = "insert"
	OpTypeDelete    = "delete"
	OpTypeCompare   = "compare"
	OpTypeAccess    = "access"
	OpTypeUpdate    = "update"
	OpTypeMerge     = "merge"
	OpTypeSplit     = "split"
	OpTypePartition = "partition"
	OpTypeAssign    = "assign" // 赋值操作
	OpTypeCall      = "call"   // 函数调用操作
)

// StepTracker 步骤追踪器接口
type StepTracker interface {
	AddStep(description string, data interface{}, highlights []int)
	AddComparison(index1, index2 int, result int)
	AddOperation(opType string, indices []int, values []interface{}, description string)
	SetPhase(phase string)
	AddNote(note string)
	GetSteps() []VisualizationStep
	GetStats() ExecutionStats
}

// DefaultStepTracker 默认步骤追踪器实现
type DefaultStepTracker struct {
	steps        []VisualizationStep
	stats        ExecutionStats
	currentPhase string
}

// NewStepTracker 创建新的步骤追踪器
func NewStepTracker() *DefaultStepTracker {
	return &DefaultStepTracker{
		steps: make([]VisualizationStep, 0),
		stats: ExecutionStats{},
	}
}

// AddStep 添加步骤
func (t *DefaultStepTracker) AddStep(description string, data interface{}, highlights []int) {
	step := VisualizationStep{
		StepID:      len(t.steps),
		Description: description,
		Data:        data,
		Highlights:  highlights,
		Comparisons: make([]Comparison, 0),
		Operations:  make([]Operation, 0),
		Metadata: StepMetadata{
			Phase: t.currentPhase,
			Notes: make([]string, 0),
		},
	}
	t.steps = append(t.steps, step)
}

// AddComparison 添加比较操作
func (t *DefaultStepTracker) AddComparison(index1, index2 int, result int) {
	if len(t.steps) > 0 {
		comparison := Comparison{
			Index1: index1,
			Index2: index2,
			Result: result,
			Type:   "value",
		}
		lastStep := &t.steps[len(t.steps)-1]
		lastStep.Comparisons = append(lastStep.Comparisons, comparison)
		t.stats.Comparisons++
	}
}

// AddOperation 添加操作
func (t *DefaultStepTracker) AddOperation(opType string, indices []int, values []interface{}, description string) {
	if len(t.steps) > 0 {
		operation := Operation{
			Type:        opType,
			Indices:     indices,
			Values:      values,
			Description: description,
		}
		lastStep := &t.steps[len(t.steps)-1]
		lastStep.Operations = append(lastStep.Operations, operation)

		// 更新统计
		switch opType {
		case OpTypeSwap:
			t.stats.Swaps++
		case OpTypeMove:
			t.stats.Moves++
		case OpTypeAccess:
			t.stats.Accesses++
		}
	}
}

// SetPhase 设置当前阶段
func (t *DefaultStepTracker) SetPhase(phase string) {
	t.currentPhase = phase
}

// AddNote 添加备注
func (t *DefaultStepTracker) AddNote(note string) {
	if len(t.steps) > 0 {
		lastStep := &t.steps[len(t.steps)-1]
		lastStep.Metadata.Notes = append(lastStep.Metadata.Notes, note)
	}
}

// GetSteps 获取所有步骤
func (t *DefaultStepTracker) GetSteps() []VisualizationStep {
	return t.steps
}

// GetStats 获取统计信息
func (t *DefaultStepTracker) GetStats() ExecutionStats {
	return t.stats
}
