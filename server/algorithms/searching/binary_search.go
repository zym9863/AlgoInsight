package searching

import (
	"gin/algorithms"
	"gin/models"
	"strconv"
)

// BinarySearch 二分搜索算法
type BinarySearch struct {
	algorithms.BaseAlgorithm
}

// NewBinarySearch 创建二分搜索算法实例
func NewBinarySearch() *BinarySearch {
	return &BinarySearch{
		BaseAlgorithm: algorithms.BaseAlgorithm{
			ID:              "binary_search",
			Name:            "二分搜索",
			Category:        models.CategorySearching,
			Description:     "二分搜索是一种在有序数组中查找特定元素的搜索算法。通过重复将搜索区间分半来工作。",
			TimeComplexity:  "O(log n)",
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

// Execute 执行二分搜索
func (bs *BinarySearch) Execute(data interface{}, tracker models.StepTracker) (interface{}, error) {
	// 验证输入
	if err := bs.ValidateInput(data); err != nil {
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
	index, err := bs.Search(arr, target, tracker)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"index":  index,
		"found":  index != -1,
		"target": target,
	}, nil
}

// Search 二分搜索实现
func (bs *BinarySearch) Search(data []interface{}, target interface{}, tracker models.StepTracker) (int, error) {
	n := len(data)
	if n == 0 {
		tracker.AddStep("数组为空，搜索结束", data, []int{})
		return -1, nil
	}

	tracker.SetPhase("初始化")
	tracker.AddStep("开始二分搜索，目标值: "+bs.toString(target), data, []int{})

	left := 0
	right := n - 1

	for left <= right {
		// 计算中间位置
		mid := left + (right-left)/2
		
		// 显示当前搜索区间
		highlights := []int{left, mid, right}
		tracker.AddStep("搜索区间 ["+strconv.Itoa(left)+", "+strconv.Itoa(right)+"], 中间位置: "+strconv.Itoa(mid), 
			data, highlights)

		// 比较中间元素与目标值
		cmp := bs.compare(data[mid], target)
		tracker.AddComparison(mid, -1, cmp) // -1表示与目标值比较

		if cmp == 0 {
			// 找到目标元素
			tracker.AddStep("找到目标元素，位置: "+strconv.Itoa(mid), data, []int{mid})
			tracker.AddNote("搜索成功")
			return mid, nil
		} else if cmp < 0 {
			// 中间元素小于目标值，搜索右半部分
			tracker.AddStep("中间元素小于目标值，搜索右半部分", data, []int{mid})
			tracker.AddNote("更新搜索区间: [" + strconv.Itoa(mid+1) + ", " + strconv.Itoa(right) + "]")
			left = mid + 1
		} else {
			// 中间元素大于目标值，搜索左半部分
			tracker.AddStep("中间元素大于目标值，搜索左半部分", data, []int{mid})
			tracker.AddNote("更新搜索区间: [" + strconv.Itoa(left) + ", " + strconv.Itoa(mid-1) + "]")
			right = mid - 1
		}
	}

	// 未找到目标元素
	tracker.SetPhase("完成")
	tracker.AddStep("搜索区间为空，未找到目标元素", data, []int{})
	tracker.AddNote("搜索失败")
	return -1, nil
}

// compare 比较两个元素
func (bs *BinarySearch) compare(a, b interface{}) int {
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
func (bs *BinarySearch) toString(value interface{}) string {
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
func (bs *BinarySearch) ValidateInput(data interface{}) error {
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

	// 二分搜索要求数组已排序，这里简化处理
	// 在实际应用中应该检查数组是否已排序

	return nil
}

// RequiresSorted 二分搜索需要已排序的数据
func (bs *BinarySearch) RequiresSorted() bool {
	return true
}

// GetComplexity 获取复杂度信息
func (bs *BinarySearch) GetComplexity() algorithms.ComplexityInfo {
	return algorithms.ComplexityInfo{
		TimeComplexity: algorithms.ComplexityCase{
			Best:    "O(1)",
			Average: "O(log n)",
			Worst:   "O(log n)",
		},
		SpaceComplexity: algorithms.ComplexityCase{
			Best:    "O(1)",
			Average: "O(1)",
			Worst:   "O(1)",
		},
	}
}
