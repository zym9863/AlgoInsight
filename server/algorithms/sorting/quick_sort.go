package sorting

import (
	"gin/algorithms"
	"gin/models"
	"strconv"
)

// QuickSort 快速排序算法
type QuickSort struct {
	algorithms.BaseAlgorithm
}

// NewQuickSort 创建快速排序算法实例
func NewQuickSort() *QuickSort {
	return &QuickSort{
		BaseAlgorithm: algorithms.BaseAlgorithm{
			ID:              "quick_sort",
			Name:            "快速排序",
			Category:        models.CategorySorting,
			Description:     "快速排序是一种高效的排序算法，采用分治法策略。选择一个基准元素，将数组分为两部分，递归地对子数组进行排序。",
			TimeComplexity:  "O(n log n)",
			SpaceComplexity: "O(log n)",
			Parameters: []models.Parameter{
				{
					Name:         "pivot_strategy",
					Type:         "string",
					Description:  "基准选择策略",
					DefaultValue: "last",
					Required:     false,
					Options:      []string{"first", "last", "middle", "random"},
				},
			},
			Stable:   false,
			InPlace:  true,
			Adaptive: false,
		},
	}
}

// Execute 执行快速排序
func (qs *QuickSort) Execute(data interface{}, tracker models.StepTracker) (interface{}, error) {
	// 验证输入
	if err := qs.ValidateInput(data); err != nil {
		return nil, err
	}

	// 转换数据类型
	arr, ok := data.([]interface{})
	if !ok {
		return nil, algorithms.ErrInvalidInput
	}

	// 复制数组以避免修改原数据
	result := make([]interface{}, len(arr))
	copy(result, arr)

	// 执行排序
	err := qs.Sort(result, tracker)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Sort 快速排序实现
func (qs *QuickSort) Sort(data []interface{}, tracker models.StepTracker) error {
	n := len(data)
	if n <= 1 {
		tracker.AddStep("数组长度小于等于1，无需排序", data, []int{})
		return nil
	}

	tracker.SetPhase("初始化")
	tracker.AddStep("开始快速排序", data, []int{})

	// 调用递归排序
	qs.quickSortRecursive(data, 0, n-1, tracker, 0)

	tracker.SetPhase("完成")
	tracker.AddStep("快速排序完成", data, []int{})
	return nil
}

// quickSortRecursive 递归快速排序
func (qs *QuickSort) quickSortRecursive(data []interface{}, low, high int, tracker models.StepTracker, depth int) {
	if low < high {
		// 设置当前阶段
		tracker.SetPhase("递归深度 " + strconv.Itoa(depth) + " - 分区")
		
		// 显示当前处理的子数组
		highlights := make([]int, 0)
		for i := low; i <= high; i++ {
			highlights = append(highlights, i)
		}
		tracker.AddStep("处理子数组 ["+strconv.Itoa(low)+", "+strconv.Itoa(high)+"]", data, highlights)

		// 分区操作
		pivotIndex := qs.partition(data, low, high, tracker)

		// 显示分区结果
		tracker.AddStep("分区完成，基准位置: "+strconv.Itoa(pivotIndex), data, []int{pivotIndex})
		tracker.AddNote("基准元素 " + qs.toString(data[pivotIndex]) + " 已就位")

		// 递归排序左半部分
		if pivotIndex-1 > low {
			tracker.SetPhase("递归深度 " + strconv.Itoa(depth+1) + " - 左子数组")
			qs.quickSortRecursive(data, low, pivotIndex-1, tracker, depth+1)
		}

		// 递归排序右半部分
		if pivotIndex+1 < high {
			tracker.SetPhase("递归深度 " + strconv.Itoa(depth+1) + " - 右子数组")
			qs.quickSortRecursive(data, pivotIndex+1, high, tracker, depth+1)
		}
	}
}

// partition 分区操作
func (qs *QuickSort) partition(data []interface{}, low, high int, tracker models.StepTracker) int {
	// 选择最后一个元素作为基准
	pivot := data[high]
	tracker.AddStep("选择基准元素: "+qs.toString(pivot), data, []int{high})

	i := low - 1 // 小于基准的元素的索引

	for j := low; j < high; j++ {
		// 比较当前元素与基准
		tracker.AddStep("比较元素", data, []int{j, high})
		
		if qs.compare(data[j], pivot) <= 0 {
			// 当前元素小于等于基准
			tracker.AddComparison(j, high, -1)
			
			i++
			if i != j {
				// 交换元素
				data[i], data[j] = data[j], data[i]
				tracker.AddOperation(models.OpTypeSwap, []int{i, j}, 
					[]interface{}{data[i], data[j]}, "移动小于基准的元素")
				tracker.AddStep("交换元素", data, []int{i, j})
			}
		} else {
			// 当前元素大于基准
			tracker.AddComparison(j, high, 1)
		}
	}

	// 将基准元素放到正确位置
	i++
	if i != high {
		data[i], data[high] = data[high], data[i]
		tracker.AddOperation(models.OpTypeSwap, []int{i, high}, 
			[]interface{}{data[i], data[high]}, "将基准元素放到正确位置")
		tracker.AddStep("基准元素就位", data, []int{i})
	}

	return i
}

// compare 比较两个元素
func (qs *QuickSort) compare(a, b interface{}) int {
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
func (qs *QuickSort) toString(value interface{}) string {
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
func (qs *QuickSort) ValidateInput(data interface{}) error {
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

// IsStable 快速排序不是稳定的
func (qs *QuickSort) IsStable() bool {
	return false
}

// IsInPlace 快速排序是原地的
func (qs *QuickSort) IsInPlace() bool {
	return true
}

// IsAdaptive 快速排序不是自适应的
func (qs *QuickSort) IsAdaptive() bool {
	return false
}

// GetComplexity 获取复杂度信息
func (qs *QuickSort) GetComplexity() algorithms.ComplexityInfo {
	return algorithms.ComplexityInfo{
		TimeComplexity: algorithms.ComplexityCase{
			Best:    "O(n log n)",
			Average: "O(n log n)",
			Worst:   "O(n²)",
		},
		SpaceComplexity: algorithms.ComplexityCase{
			Best:    "O(log n)",
			Average: "O(log n)",
			Worst:   "O(n)",
		},
	}
}
