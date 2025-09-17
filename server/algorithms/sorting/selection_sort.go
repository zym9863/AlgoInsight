package sorting

import (
	"gin/algorithms"
	"gin/models"
	"strconv"
)

// SelectionSort 选择排序算法
type SelectionSort struct {
	algorithms.BaseAlgorithm
}

// NewSelectionSort 创建选择排序算法实例
func NewSelectionSort() *SelectionSort {
	return &SelectionSort{
		BaseAlgorithm: algorithms.BaseAlgorithm{
			ID:              "selection_sort",
			Name:            "选择排序",
			Category:        models.CategorySorting,
			Description:     "选择排序是一种简单直观的排序算法。它的工作原理是：首先在未排序序列中找到最小（大）元素，存放到排序序列的起始位置，然后再从剩余未排序元素中继续寻找最小（大）元素，然后放到已排序序列的末尾。",
			TimeComplexity:  "O(n²)",
			SpaceComplexity: "O(1)",
			Parameters:      []models.Parameter{},
			Stable:          false,
			InPlace:         true,
			Adaptive:        false,
		},
	}
}

// Execute 执行选择排序
func (ss *SelectionSort) Execute(data interface{}, tracker models.StepTracker) (interface{}, error) {
	// 验证输入
	if err := ss.ValidateInput(data); err != nil {
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
	err := ss.Sort(result, tracker)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Sort 选择排序实现
func (ss *SelectionSort) Sort(data []interface{}, tracker models.StepTracker) error {
	n := len(data)
	if n <= 1 {
		tracker.AddStep("数组长度小于等于1，无需排序", data, []int{})
		return nil
	}

	tracker.SetPhase("初始化")
	tracker.AddStep("开始选择排序", data, []int{})

	// 外层循环：确定每个位置的元素
	for i := 0; i < n-1; i++ {
		tracker.SetPhase("第 " + strconv.Itoa(i+1) + " 轮选择")

		// 假设当前位置的元素是最小的
		minIndex := i
		tracker.AddStep("假设位置 "+strconv.Itoa(i)+" 的元素 "+ss.toString(data[i])+" 是最小值",
			data, []int{i})

		// 显示已排序部分和待排序部分
		sortedHighlights := make([]int, 0)
		for j := 0; j < i; j++ {
			sortedHighlights = append(sortedHighlights, j)
		}
		if len(sortedHighlights) > 0 {
			tracker.AddStep("已排序部分: [0, "+strconv.Itoa(i-1)+"], 寻找 ["+
				strconv.Itoa(i)+", "+strconv.Itoa(n-1)+"] 中的最小值", data, sortedHighlights)
		}

		// 内层循环：在未排序部分找到最小元素
		for j := i + 1; j < n; j++ {
			tracker.AddStep("比较 "+ss.toString(data[j])+" 与当前最小值 "+ss.toString(data[minIndex]),
				data, []int{j, minIndex})

			// 比较当前元素与当前最小值
			if ss.compare(data[j], data[minIndex]) < 0 {
				// 找到更小的元素
				tracker.AddComparison(j, minIndex, -1)
				tracker.AddStep("发现更小的元素 "+ss.toString(data[j])+" 在位置 "+strconv.Itoa(j),
					data, []int{j})
				tracker.AddNote("更新最小值索引：" + strconv.Itoa(minIndex) + " -> " + strconv.Itoa(j))
				minIndex = j
			} else {
				// 当前元素不是最小的
				tracker.AddComparison(j, minIndex, 1)
			}
		}

		// 显示本轮找到的最小元素
		if minIndex != i {
			tracker.AddStep("本轮最小元素: "+ss.toString(data[minIndex])+" 在位置 "+strconv.Itoa(minIndex),
				data, []int{minIndex})

			// 交换最小元素到正确位置
			tracker.AddStep("交换位置 "+strconv.Itoa(i)+" 和 "+strconv.Itoa(minIndex)+" 的元素",
				data, []int{i, minIndex})

			data[i], data[minIndex] = data[minIndex], data[i]
			tracker.AddOperation(models.OpTypeSwap, []int{i, minIndex},
				[]interface{}{data[i], data[minIndex]}, "将最小元素放到正确位置")

			tracker.AddStep("交换完成", data, []int{i})
		} else {
			tracker.AddStep("最小元素已在正确位置", data, []int{i})
		}

		// 显示本轮结果
		newSortedHighlights := make([]int, 0)
		for k := 0; k <= i; k++ {
			newSortedHighlights = append(newSortedHighlights, k)
		}
		tracker.AddStep("第 "+strconv.Itoa(i+1)+" 轮完成，位置 "+strconv.Itoa(i)+" 确定",
			data, newSortedHighlights)
		tracker.AddNote("前 " + strconv.Itoa(i+1) + " 个元素已排序")
	}

	tracker.SetPhase("完成")
	tracker.AddStep("选择排序完成", data, []int{})
	return nil
}

// compare 比较两个元素
func (ss *SelectionSort) compare(a, b interface{}) int {
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
func (ss *SelectionSort) toString(value interface{}) string {
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
func (ss *SelectionSort) ValidateInput(data interface{}) error {
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

// IsStable 选择排序不是稳定的
func (ss *SelectionSort) IsStable() bool {
	return false
}

// IsInPlace 选择排序是原地的
func (ss *SelectionSort) IsInPlace() bool {
	return true
}

// IsAdaptive 选择排序不是自适应的
func (ss *SelectionSort) IsAdaptive() bool {
	return false
}

// GetComplexity 获取复杂度信息
func (ss *SelectionSort) GetComplexity() algorithms.ComplexityInfo {
	return algorithms.ComplexityInfo{
		TimeComplexity: algorithms.ComplexityCase{
			Best:    "O(n²)",
			Average: "O(n²)",
			Worst:   "O(n²)",
		},
		SpaceComplexity: algorithms.ComplexityCase{
			Best:    "O(1)",
			Average: "O(1)",
			Worst:   "O(1)",
		},
	}
}