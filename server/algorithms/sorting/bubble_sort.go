package sorting

import (
	"gin/algorithms"
	"gin/models"
	"strconv"
)

// BubbleSort 冒泡排序算法
type BubbleSort struct {
	algorithms.BaseAlgorithm
}

// NewBubbleSort 创建冒泡排序算法实例
func NewBubbleSort() *BubbleSort {
	return &BubbleSort{
		BaseAlgorithm: algorithms.BaseAlgorithm{
			ID:              "bubble_sort",
			Name:            "冒泡排序",
			Category:        models.CategorySorting,
			Description:     "冒泡排序是一种简单的排序算法，通过重复遍历要排序的数列，一次比较两个元素，如果它们的顺序错误就把它们交换过来。",
			TimeComplexity:  "O(n²)",
			SpaceComplexity: "O(1)",
			Parameters:      []models.Parameter{},
			Stable:          true,
			InPlace:         true,
			Adaptive:        true,
		},
	}
}

// Execute 执行冒泡排序
func (bs *BubbleSort) Execute(data interface{}, tracker models.StepTracker) (interface{}, error) {
	// 验证输入
	if err := bs.ValidateInput(data); err != nil {
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
	err := bs.Sort(result, tracker)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Sort 冒泡排序实现
func (bs *BubbleSort) Sort(data []interface{}, tracker models.StepTracker) error {
	n := len(data)
	if n <= 1 {
		tracker.AddStep("数组长度小于等于1，无需排序", data, []int{})
		return nil
	}

	tracker.SetPhase("初始化")
	tracker.AddStep("开始冒泡排序", data, []int{})

	// 外层循环：控制排序轮数
	for i := 0; i < n-1; i++ {
		tracker.SetPhase("第" + strconv.Itoa(i+1) + "轮冒泡")
		swapped := false

		// 内层循环：进行比较和交换
		for j := 0; j < n-i-1; j++ {
			// 添加比较步骤
			tracker.AddStep("比较元素", data, []int{j, j + 1})
			
			// 比较相邻元素
			if bs.compare(data[j], data[j+1]) > 0 {
				// 记录比较结果
				tracker.AddComparison(j, j+1, 1)
				
				// 交换元素
				data[j], data[j+1] = data[j+1], data[j]
				swapped = true
				
				// 记录交换操作
				tracker.AddOperation(models.OpTypeSwap, []int{j, j + 1}, 
					[]interface{}{data[j], data[j+1]}, "交换元素")
				
				// 添加交换后的步骤
				tracker.AddStep("交换元素", data, []int{j, j + 1})
			} else {
				// 记录比较结果（不需要交换）
				tracker.AddComparison(j, j+1, -1)
			}
		}

		// 一轮结束，最大元素已经"冒泡"到正确位置
		tracker.AddStep("第"+strconv.Itoa(i+1)+"轮结束，最大元素已就位", data, []int{n - i - 1})
		tracker.AddNote("位置 " + strconv.Itoa(n-i-1) + " 的元素已确定")

		// 如果这一轮没有发生交换，说明数组已经有序
		if !swapped {
			tracker.AddStep("未发生交换，排序完成", data, []int{})
			tracker.AddNote("数组已经有序，提前结束")
			break
		}
	}

	tracker.SetPhase("完成")
	tracker.AddStep("冒泡排序完成", data, []int{})
	return nil
}

// compare 比较两个元素
func (bs *BubbleSort) compare(a, b interface{}) int {
	// 这里简化处理，假设都是可比较的数值类型
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

// ValidateInput 验证输入数据
func (bs *BubbleSort) ValidateInput(data interface{}) error {
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

// IsStable 冒泡排序是稳定的
func (bs *BubbleSort) IsStable() bool {
	return true
}

// IsInPlace 冒泡排序是原地的
func (bs *BubbleSort) IsInPlace() bool {
	return true
}

// IsAdaptive 冒泡排序是自适应的
func (bs *BubbleSort) IsAdaptive() bool {
	return true
}

// GetComplexity 获取复杂度信息
func (bs *BubbleSort) GetComplexity() algorithms.ComplexityInfo {
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
