package sorting

import (
	"gin/algorithms"
	"gin/models"
	"strconv"
)

// InsertionSort 插入排序算法
type InsertionSort struct {
	algorithms.BaseAlgorithm
}

// NewInsertionSort 创建插入排序算法实例
func NewInsertionSort() *InsertionSort {
	return &InsertionSort{
		BaseAlgorithm: algorithms.BaseAlgorithm{
			ID:              "insertion_sort",
			Name:            "插入排序",
			Category:        models.CategorySorting,
			Description:     "插入排序是一种简单直观的排序算法。它的工作原理是通过构建有序序列，对于未排序数据，在已排序序列中从后向前扫描，找到相应位置并插入。",
			TimeComplexity:  "O(n²)",
			SpaceComplexity: "O(1)",
			Parameters:      []models.Parameter{},
			Stable:          true,
			InPlace:         true,
			Adaptive:        true,
		},
	}
}

// Execute 执行插入排序
func (is *InsertionSort) Execute(data interface{}, tracker models.StepTracker) (interface{}, error) {
	// 验证输入
	if err := is.ValidateInput(data); err != nil {
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
	err := is.Sort(result, tracker)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Sort 插入排序实现
func (is *InsertionSort) Sort(data []interface{}, tracker models.StepTracker) error {
	n := len(data)
	if n <= 1 {
		tracker.AddStep("数组长度小于等于1，无需排序", data, []int{})
		return nil
	}

	tracker.SetPhase("初始化")
	tracker.AddStep("开始插入排序，第一个元素已排序", data, []int{0})

	// 从第二个元素开始插入
	for i := 1; i < n; i++ {
		tracker.SetPhase("第 " + strconv.Itoa(i) + " 轮插入")
		
		current := data[i]
		tracker.AddStep("取出待插入元素: "+is.toString(current), data, []int{i})
		tracker.AddNote("当前元素: " + is.toString(current))

		// 在已排序部分中找到插入位置
		j := i - 1
		
		// 显示已排序部分
		sortedHighlights := make([]int, 0)
		for k := 0; k <= i-1; k++ {
			sortedHighlights = append(sortedHighlights, k)
		}
		tracker.AddStep("已排序部分: [0, "+strconv.Itoa(i-1)+"]", data, sortedHighlights)

		// 向后移动比当前元素大的元素
		for j >= 0 && is.compare(data[j], current) > 0 {
			tracker.AddStep("比较: "+is.toString(data[j])+" > "+is.toString(current), data, []int{j, j+1})
			tracker.AddComparison(j, i, 1)
			
			// 向右移动元素
			data[j+1] = data[j]
			tracker.AddOperation(models.OpTypeMove, []int{j + 1}, []interface{}{data[j+1]}, "向右移动元素")
			tracker.AddStep("向右移动: "+is.toString(data[j+1]), data, []int{j + 1})
			
			j--
		}

		// 插入当前元素到正确位置
		data[j+1] = current
		tracker.AddOperation(models.OpTypeInsert, []int{j + 1}, []interface{}{current}, "插入元素")
		tracker.AddStep("插入元素到位置 "+strconv.Itoa(j+1), data, []int{j + 1})

		// 显示当前排序状态
		currentSortedHighlights := make([]int, 0)
		for k := 0; k <= i; k++ {
			currentSortedHighlights = append(currentSortedHighlights, k)
		}
		tracker.AddStep("第"+strconv.Itoa(i)+"轮完成，已排序: [0, "+strconv.Itoa(i)+"]", data, currentSortedHighlights)
	}

	tracker.SetPhase("完成")
	tracker.AddStep("插入排序完成", data, []int{})
	return nil
}

// compare 比较两个元素
func (is *InsertionSort) compare(a, b interface{}) int {
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
func (is *InsertionSort) toString(value interface{}) string {
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
func (is *InsertionSort) ValidateInput(data interface{}) error {
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

// IsStable 插入排序是稳定的
func (is *InsertionSort) IsStable() bool {
	return true
}

// IsInPlace 插入排序是原地的
func (is *InsertionSort) IsInPlace() bool {
	return true
}

// IsAdaptive 插入排序是自适应的
func (is *InsertionSort) IsAdaptive() bool {
	return true
}

// GetComplexity 获取复杂度信息
func (is *InsertionSort) GetComplexity() algorithms.ComplexityInfo {
	return algorithms.ComplexityInfo{
		TimeComplexity: algorithms.ComplexityCase{
			Best:    "O(n)",
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