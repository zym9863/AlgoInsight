package searching

import (
	"gin/algorithms"
	"gin/models"
	"strconv"
)

// LinearSearch 线性搜索算法
type LinearSearch struct {
	algorithms.BaseAlgorithm
}

// NewLinearSearch 创建线性搜索算法实例
func NewLinearSearch() *LinearSearch {
	return &LinearSearch{
		BaseAlgorithm: algorithms.BaseAlgorithm{
			ID:              "linear_search",
			Name:            "线性搜索",
			Category:        models.CategorySearching,
			Description:     "线性搜索是最简单的搜索算法。它按顺序检查列表中的每个元素，直到找到目标元素或检查完所有元素。",
			TimeComplexity:  "O(n)",
			SpaceComplexity: "O(1)",
			Parameters: []models.Parameter{
				{
					Name:         "target",
					Type:         "interface{}",
					Description:  "要搜索的目标值",
					DefaultValue: nil,
					Required:     true,
				},
			},
			Stable:   true,
			InPlace:  true,
			Adaptive: false,
		},
	}
}

// Execute 执行线性搜索
func (ls *LinearSearch) Execute(data interface{}, tracker models.StepTracker) (interface{}, error) {
	// 验证输入
	if err := ls.ValidateInput(data); err != nil {
		return nil, err
	}

	// 这里需要从参数中获取目标值，暂时使用默认实现
	// 在实际使用中，target应该从parameters中获取
	target := 5 // 示例目标值

	// 转换数据类型
	arr, ok := data.([]interface{})
	if !ok {
		return nil, algorithms.ErrInvalidInput
	}

	// 执行搜索
	index, err := ls.Search(arr, target, tracker)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"index":  index,
		"found":  index != -1,
		"target": target,
	}, nil
}

// Search 线性搜索实现
func (ls *LinearSearch) Search(data []interface{}, target interface{}, tracker models.StepTracker) (int, error) {
	n := len(data)
	if n == 0 {
		tracker.AddStep("数组为空，搜索结束", data, []int{})
		return -1, nil
	}

	tracker.SetPhase("初始化")
	tracker.AddStep("开始线性搜索，目标值: "+ls.toString(target), data, []int{})
	tracker.AddNote("将从索引 0 开始逐个检查元素")

	// 逐个检查数组中的每个元素
	for i := 0; i < n; i++ {
		tracker.SetPhase("检查位置 " + strconv.Itoa(i))

		// 显示当前检查的元素
		tracker.AddStep("检查位置 "+strconv.Itoa(i)+" 的元素: "+ls.toString(data[i]),
			data, []int{i})

		// 比较当前元素与目标值
		cmp := ls.compare(data[i], target)
		tracker.AddComparison(i, -1, cmp) // -1表示与目标值比较

		if cmp == 0 {
			// 找到目标元素
			tracker.AddStep("找到目标元素! 位置: "+strconv.Itoa(i), data, []int{i})
			tracker.AddNote("搜索成功，共检查了 " + strconv.Itoa(i+1) + " 个元素")
			return i, nil
		} else {
			// 当前元素不匹配
			tracker.AddStep("元素 "+ls.toString(data[i])+" 不匹配，继续搜索", data, []int{})
		}
	}

	// 未找到目标元素
	tracker.SetPhase("完成")
	tracker.AddStep("搜索完毕，未找到目标元素", data, []int{})
	tracker.AddNote("搜索失败，已检查所有 " + strconv.Itoa(n) + " 个元素")
	return -1, nil
}

// compare 比较两个元素
func (ls *LinearSearch) compare(a, b interface{}) int {
	switch va := a.(type) {
	case int:
		if vb, ok := b.(int); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case float64:
		if vb, ok := b.(float64); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case string:
		if vb, ok := b.(string); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	}
	return 0
}

// toString 将元素转换为字符串
func (ls *LinearSearch) toString(value interface{}) string {
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		return v
	default:
		return "unknown"
	}
}

// ValidateInput 验证输入数据
func (ls *LinearSearch) ValidateInput(data interface{}) error {
	if data == nil {
		return algorithms.ErrInvalidInput
	}

	arr, ok := data.([]interface{})
	if !ok {
		return algorithms.ErrInvalidInput
	}

	// 检查数组大小限制
	if len(arr) > 10000 {
		return algorithms.ErrInvalidInput
	}

	return nil
}

// RequiresSorted 线性搜索不需要已排序的数据
func (ls *LinearSearch) RequiresSorted() bool {
	return false
}

// GetComplexity 获取复杂度信息
func (ls *LinearSearch) GetComplexity() algorithms.ComplexityInfo {
	return algorithms.ComplexityInfo{
		TimeComplexity: algorithms.ComplexityCase{
			Best:    "O(1)",
			Average: "O(n)",
			Worst:   "O(n)",
		},
		SpaceComplexity: algorithms.ComplexityCase{
			Best:    "O(1)",
			Average: "O(1)",
			Worst:   "O(1)",
		},
	}
}