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
	tracker.AddStep("开始插入排序", data, []int{})
	tracker.AddNote("第一个元素默认为已排序")

	// 从第二个元素开始，逐个插入到已排序序列中
	for i := 1; i < n; i++ {
		tracker.SetPhase("插入第 " + strconv.Itoa(i+1) + " 个元素")

		// 当前要插入的元素
		key := data[i]
		tracker.AddStep("选择待插入元素: "+is.toString(key)+" (位置 "+strconv.Itoa(i)+")",
			data, []int{i})

		// 显示已排序部分和未排序部分
		sortedHighlights := make([]int, 0)
		for j := 0; j < i; j++ {
			sortedHighlights = append(sortedHighlights, j)
		}
		tracker.AddStep("已排序部分: [0, "+strconv.Itoa(i-1)+"], 未排序部分: ["+
			strconv.Itoa(i)+", "+strconv.Itoa(n-1)+"]", data, sortedHighlights)

		// 在已排序序列中找到插入位置
		j := i - 1

		// 向后移动大于key的元素
		for j >= 0 && is.compare(data[j], key) > 0 {
			tracker.AddStep("比较 "+is.toString(data[j])+" 与 "+is.toString(key),
				data, []int{j, i})
			tracker.AddComparison(j, i, 1)

			// 元素向后移动
			data[j+1] = data[j]
			tracker.AddOperation(models.OpTypeMove, []int{j, j + 1},
				[]interface{}{data[j+1]}, "向后移动较大元素")
			tracker.AddStep("移动元素 "+is.toString(data[j+1])+" 到位置 "+strconv.Itoa(j+1),
				data, []int{j + 1})

			j--
		}

		// 如果j >= 0，说明找到了合适的位置或者到达了数组开头
		if j >= 0 {
			tracker.AddStep("比较 "+is.toString(data[j])+" 与 "+is.toString(key)+
				" (找到插入位置)", data, []int{j})
			tracker.AddComparison(j, i, -1)
		}

		// 插入key到正确位置
		insertPos := j + 1
		if insertPos != i {
			data[insertPos] = key
			tracker.AddOperation(models.OpTypeInsert, []int{insertPos},
				[]interface{}{key}, "插入元素到正确位置")
			tracker.AddStep("插入元素 "+is.toString(key)+" 到位置 "+strconv.Itoa(insertPos),
				data, []int{insertPos})
		} else {
			tracker.AddStep("元素 "+is.toString(key)+" 已在正确位置", data, []int{i})
		}

		// 显示当前状态
		newSortedHighlights := make([]int, 0)
		for k := 0; k <= i; k++ {
			newSortedHighlights = append(newSortedHighlights, k)
		}
		tracker.AddStep("第 "+strconv.Itoa(i+1)+" 轮完成，已排序部分: [0, "+strconv.Itoa(i)+"]",
			data, newSortedHighlights)
		tracker.AddNote("前 " + strconv.Itoa(i+1) + " 个元素已排序")
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