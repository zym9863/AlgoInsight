package sorting

import (
	"gin/algorithms"
	"gin/models"
	"strconv"
)

// ShellSort 希尔排序算法
type ShellSort struct {
	algorithms.BaseAlgorithm
}

// NewShellSort 创建希尔排序算法实例
func NewShellSort() *ShellSort {
	return &ShellSort{
		BaseAlgorithm: algorithms.BaseAlgorithm{
			ID:              "shell_sort",
			Name:            "希尔排序",
			Category:        models.CategorySorting,
			Description:     "希尔排序是插入排序的一种改进版本，也叫缩小增量排序。通过将比较的全部元素分为几个区域来提升插入排序的性能，先对间隔较远的元素进行排序，然后逐步减小间隔。",
			TimeComplexity:  "O(n log² n)",
			SpaceComplexity: "O(1)",
			Parameters: []models.Parameter{
				{
					Name:         "gap_sequence",
					Type:         "string",
					Description:  "间隔序列类型",
					DefaultValue: "shell",
					Required:     false,
					Options:      []string{"shell", "knuth", "sedgewick"},
				},
			},
			Stable:   false,
			InPlace:  true,
			Adaptive: true,
		},
	}
}

// Execute 执行希尔排序
func (ss *ShellSort) Execute(data interface{}, tracker models.StepTracker) (interface{}, error) {
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

// Sort 希尔排序实现
func (ss *ShellSort) Sort(data []interface{}, tracker models.StepTracker) error {
	n := len(data)
	if n <= 1 {
		tracker.AddStep("数组长度小于等于1，无需排序", data, []int{})
		return nil
	}

	tracker.SetPhase("初始化")
	tracker.AddStep("开始希尔排序", data, []int{})

	// 使用Shell原始间隔序列：n/2, n/4, n/8, ..., 1
	gap := n / 2

	for gap > 0 {
		tracker.SetPhase("间隔 gap = " + strconv.Itoa(gap))
		tracker.AddStep("当前间隔："+strconv.Itoa(gap), data, []int{})
		tracker.AddNote("对间隔为 " + strconv.Itoa(gap) + " 的子序列进行插入排序")

		// 对每个间隔为gap的子序列进行插入排序
		for i := gap; i < n; i++ {
			// 当前要插入的元素
			key := data[i]
			j := i

			tracker.AddStep("处理元素 "+ss.toString(key)+" (位置 "+strconv.Itoa(i)+")",
				data, []int{i})

			// 显示当前比较的子序列
			subseqHighlights := make([]int, 0)
			for k := i; k >= 0; k -= gap {
				subseqHighlights = append(subseqHighlights, k)
			}
			tracker.AddStep("子序列元素位置："+ss.formatIntArray(subseqHighlights),
				data, subseqHighlights)

			// 在子序列中找到插入位置
			for j >= gap && ss.compare(data[j-gap], key) > 0 {
				tracker.AddStep("比较 "+ss.toString(data[j-gap])+" 与 "+ss.toString(key)+
					" (间隔 "+strconv.Itoa(gap)+")", data, []int{j - gap, i})
				tracker.AddComparison(j-gap, i, 1)

				// 元素向后移动
				data[j] = data[j-gap]
				tracker.AddOperation(models.OpTypeMove, []int{j - gap, j},
					[]interface{}{data[j]}, "在子序列中移动元素")
				tracker.AddStep("移动元素 "+ss.toString(data[j])+" 到位置 "+strconv.Itoa(j),
					data, []int{j})

				j -= gap
			}

			// 如果j >= gap，说明找到了插入位置
			if j >= gap {
				tracker.AddStep("比较 "+ss.toString(data[j-gap])+" 与 "+ss.toString(key)+
					" (找到插入位置)", data, []int{j - gap})
				tracker.AddComparison(j-gap, i, -1)
			}

			// 插入元素到正确位置
			if j != i {
				data[j] = key
				tracker.AddOperation(models.OpTypeInsert, []int{j},
					[]interface{}{key}, "在子序列中插入元素")
				tracker.AddStep("插入元素 "+ss.toString(key)+" 到位置 "+strconv.Itoa(j),
					data, []int{j})
			} else {
				tracker.AddStep("元素 "+ss.toString(key)+" 已在正确位置", data, []int{i})
			}
		}

		// 显示当前间隔完成后的状态
		tracker.AddStep("间隔 "+strconv.Itoa(gap)+" 的排序完成", data, []int{})
		tracker.AddNote("所有间隔为 " + strconv.Itoa(gap) + " 的子序列已排序")

		// 缩小间隔
		gap /= 2
	}

	tracker.SetPhase("完成")
	tracker.AddStep("希尔排序完成", data, []int{})
	return nil
}

// compare 比较两个元素
func (ss *ShellSort) compare(a, b interface{}) int {
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
func (ss *ShellSort) toString(value interface{}) string {
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

// formatIntArray 格式化整数数组为字符串
func (ss *ShellSort) formatIntArray(arr []int) string {
	if len(arr) == 0 {
		return "[]"
	}
	result := "["
	for i, v := range arr {
		if i > 0 {
			result += ", "
		}
		result += strconv.Itoa(v)
	}
	result += "]"
	return result
}

// ValidateInput 验证输入数据
func (ss *ShellSort) ValidateInput(data interface{}) error {
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

// IsStable 希尔排序不是稳定的
func (ss *ShellSort) IsStable() bool {
	return false
}

// IsInPlace 希尔排序是原地的
func (ss *ShellSort) IsInPlace() bool {
	return true
}

// IsAdaptive 希尔排序是自适应的
func (ss *ShellSort) IsAdaptive() bool {
	return true
}

// GetComplexity 获取复杂度信息
func (ss *ShellSort) GetComplexity() algorithms.ComplexityInfo {
	return algorithms.ComplexityInfo{
		TimeComplexity: algorithms.ComplexityCase{
			Best:    "O(n log n)",
			Average: "O(n log² n)",
			Worst:   "O(n²)",
		},
		SpaceComplexity: algorithms.ComplexityCase{
			Best:    "O(1)",
			Average: "O(1)",
			Worst:   "O(1)",
		},
	}
}