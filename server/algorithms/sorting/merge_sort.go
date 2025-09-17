package sorting

import (
	"gin/algorithms"
	"gin/models"
	"strconv"
)

// MergeSort 归并排序算法
type MergeSort struct {
	algorithms.BaseAlgorithm
}

// NewMergeSort 创建归并排序算法实例
func NewMergeSort() *MergeSort {
	return &MergeSort{
		BaseAlgorithm: algorithms.BaseAlgorithm{
			ID:              "merge_sort",
			Name:            "归并排序",
			Category:        models.CategorySorting,
			Description:     "归并排序是一种稳定的排序算法，采用分治法策略。将数组分为两半，递归地对子数组进行排序，然后合并已排序的子数组。",
			TimeComplexity:  "O(n log n)",
			SpaceComplexity: "O(n)",
			Parameters:      []models.Parameter{},
			Stable:          true,
			InPlace:         false,
			Adaptive:        false,
		},
	}
}

// Execute 执行归并排序
func (ms *MergeSort) Execute(data interface{}, tracker models.StepTracker) (interface{}, error) {
	// 验证输入
	if err := ms.ValidateInput(data); err != nil {
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
	err := ms.Sort(result, tracker)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Sort 归并排序实现
func (ms *MergeSort) Sort(data []interface{}, tracker models.StepTracker) error {
	n := len(data)
	if n <= 1 {
		tracker.AddStep("数组长度小于等于1，无需排序", data, []int{})
		return nil
	}

	tracker.SetPhase("初始化")
	tracker.AddStep("开始归并排序", data, []int{})

	// 调用递归排序
	ms.mergeSortRecursive(data, 0, n-1, tracker, 0)

	tracker.SetPhase("完成")
	tracker.AddStep("归并排序完成", data, []int{})
	return nil
}

// mergeSortRecursive 递归归并排序
func (ms *MergeSort) mergeSortRecursive(data []interface{}, left, right int, tracker models.StepTracker, depth int) {
	if left < right {
		// 设置当前阶段
		tracker.SetPhase("递归深度 " + strconv.Itoa(depth) + " - 分解")

		// 显示当前处理的子数组
		highlights := make([]int, 0)
		for i := left; i <= right; i++ {
			highlights = append(highlights, i)
		}
		tracker.AddStep("分解子数组 ["+strconv.Itoa(left)+", "+strconv.Itoa(right)+"]", data, highlights)

		// 找到中点
		mid := left + (right-left)/2
		tracker.AddStep("找到中点位置: "+strconv.Itoa(mid), data, []int{mid})

		// 递归排序左半部分
		tracker.SetPhase("递归深度 " + strconv.Itoa(depth+1) + " - 左子数组")
		ms.mergeSortRecursive(data, left, mid, tracker, depth+1)

		// 递归排序右半部分
		tracker.SetPhase("递归深度 " + strconv.Itoa(depth+1) + " - 右子数组")
		ms.mergeSortRecursive(data, mid+1, right, tracker, depth+1)

		// 合并已排序的两个子数组
		tracker.SetPhase("递归深度 " + strconv.Itoa(depth) + " - 合并")
		ms.merge(data, left, mid, right, tracker)
	}
}

// merge 合并两个已排序的子数组
func (ms *MergeSort) merge(data []interface{}, left, mid, right int, tracker models.StepTracker) {
	// 创建临时数组存储左右子数组
	leftArr := make([]interface{}, mid-left+1)
	rightArr := make([]interface{}, right-mid)

	// 复制数据到临时数组
	for i := 0; i < len(leftArr); i++ {
		leftArr[i] = data[left+i]
	}
	for i := 0; i < len(rightArr); i++ {
		rightArr[i] = data[mid+1+i]
	}

	// 显示准备合并的两个子数组
	leftHighlights := make([]int, 0)
	rightHighlights := make([]int, 0)
	for i := left; i <= mid; i++ {
		leftHighlights = append(leftHighlights, i)
	}
	for i := mid + 1; i <= right; i++ {
		rightHighlights = append(rightHighlights, i)
	}

	allHighlights := append(leftHighlights, rightHighlights...)
	tracker.AddStep("准备合并子数组 ["+strconv.Itoa(left)+", "+strconv.Itoa(mid)+"] 和 ["+
		strconv.Itoa(mid+1)+", "+strconv.Itoa(right)+"]", data, allHighlights)

	// 合并过程
	i, j, k := 0, 0, left

	for i < len(leftArr) && j < len(rightArr) {
		// 比较左右子数组的当前元素
		tracker.AddStep("比较元素 "+ms.toString(leftArr[i])+" 和 "+ms.toString(rightArr[j]),
			data, []int{k})

		if ms.compare(leftArr[i], rightArr[j]) <= 0 {
			// 左边元素较小或相等，选择左边元素
			data[k] = leftArr[i]
			tracker.AddComparison(left+i, mid+1+j, -1)
			tracker.AddOperation(models.OpTypeAssign, []int{k},
				[]interface{}{leftArr[i]}, "选择左子数组元素")
			i++
		} else {
			// 右边元素较小，选择右边元素
			data[k] = rightArr[j]
			tracker.AddComparison(left+i, mid+1+j, 1)
			tracker.AddOperation(models.OpTypeAssign, []int{k},
				[]interface{}{rightArr[j]}, "选择右子数组元素")
			j++
		}

		tracker.AddStep("放置元素 "+ms.toString(data[k])+" 到位置 "+strconv.Itoa(k),
			data, []int{k})
		k++
	}

	// 复制左子数组的剩余元素
	for i < len(leftArr) {
		data[k] = leftArr[i]
		tracker.AddOperation(models.OpTypeAssign, []int{k},
			[]interface{}{leftArr[i]}, "复制左子数组剩余元素")
		tracker.AddStep("复制左子数组剩余元素 "+ms.toString(leftArr[i]),
			data, []int{k})
		i++
		k++
	}

	// 复制右子数组的剩余元素
	for j < len(rightArr) {
		data[k] = rightArr[j]
		tracker.AddOperation(models.OpTypeAssign, []int{k},
			[]interface{}{rightArr[j]}, "复制右子数组剩余元素")
		tracker.AddStep("复制右子数组剩余元素 "+ms.toString(rightArr[j]),
			data, []int{k})
		j++
		k++
	}

	// 显示合并结果
	mergedHighlights := make([]int, 0)
	for i := left; i <= right; i++ {
		mergedHighlights = append(mergedHighlights, i)
	}
	tracker.AddStep("合并完成 ["+strconv.Itoa(left)+", "+strconv.Itoa(right)+"]",
		data, mergedHighlights)
	tracker.AddNote("子数组 [" + strconv.Itoa(left) + ", " + strconv.Itoa(right) + "] 已排序")
}

// compare 比较两个元素
func (ms *MergeSort) compare(a, b interface{}) int {
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
func (ms *MergeSort) toString(value interface{}) string {
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
func (ms *MergeSort) ValidateInput(data interface{}) error {
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

// IsStable 归并排序是稳定的
func (ms *MergeSort) IsStable() bool {
	return true
}

// IsInPlace 归并排序不是原地的
func (ms *MergeSort) IsInPlace() bool {
	return false
}

// IsAdaptive 归并排序不是自适应的
func (ms *MergeSort) IsAdaptive() bool {
	return false
}

// GetComplexity 获取复杂度信息
func (ms *MergeSort) GetComplexity() algorithms.ComplexityInfo {
	return algorithms.ComplexityInfo{
		TimeComplexity: algorithms.ComplexityCase{
			Best:    "O(n log n)",
			Average: "O(n log n)",
			Worst:   "O(n log n)",
		},
		SpaceComplexity: algorithms.ComplexityCase{
			Best:    "O(n)",
			Average: "O(n)",
			Worst:   "O(n)",
		},
	}
}