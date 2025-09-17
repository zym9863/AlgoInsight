package sorting

import (
	"gin/algorithms"
	"gin/models"
	"strconv"
)

// HeapSort 堆排序算法
type HeapSort struct {
	algorithms.BaseAlgorithm
}

// NewHeapSort 创建堆排序算法实例
func NewHeapSort() *HeapSort {
	return &HeapSort{
		BaseAlgorithm: algorithms.BaseAlgorithm{
			ID:              "heap_sort",
			Name:            "堆排序",
			Category:        models.CategorySorting,
			Description:     "堆排序是一种基于堆数据结构的排序算法。它首先构建最大堆，然后重复提取最大元素并重新堆化。",
			TimeComplexity:  "O(n log n)",
			SpaceComplexity: "O(1)",
			Parameters:      []models.Parameter{},
			Stable:          false,
			InPlace:         true,
			Adaptive:        false,
		},
	}
}

// Execute 执行堆排序
func (hs *HeapSort) Execute(data interface{}, tracker models.StepTracker) (interface{}, error) {
	// 验证输入
	if err := hs.ValidateInput(data); err != nil {
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
	err := hs.Sort(result, tracker)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Sort 堆排序实现
func (hs *HeapSort) Sort(data []interface{}, tracker models.StepTracker) error {
	n := len(data)
	if n <= 1 {
		tracker.AddStep("数组长度小于等于1，无需排序", data, []int{})
		return nil
	}

	tracker.SetPhase("初始化")
	tracker.AddStep("开始堆排序", data, []int{})

	// 第一阶段：构建最大堆
	tracker.SetPhase("构建最大堆")
	hs.buildMaxHeap(data, tracker)

	// 第二阶段：排序
	tracker.SetPhase("堆排序过程")
	for i := n - 1; i > 0; i-- {
		// 将堆顶（最大值）与末尾元素交换
		tracker.AddStep("交换堆顶与位置 "+strconv.Itoa(i), data, []int{0, i})
		data[0], data[i] = data[i], data[0]
		tracker.AddOperation(models.OpTypeSwap, []int{0, i},
			[]interface{}{data[0], data[i]}, "将最大元素移到正确位置")

		// 减少堆的大小并重新堆化
		tracker.AddStep("减少堆大小，重新堆化 [0, "+strconv.Itoa(i-1)+"]", data, []int{})
		tracker.AddNote("位置 " + strconv.Itoa(i) + " 的元素已确定")
		hs.heapify(data, i, 0, tracker)
	}

	tracker.SetPhase("完成")
	tracker.AddStep("堆排序完成", data, []int{})
	return nil
}

// buildMaxHeap 构建最大堆
func (hs *HeapSort) buildMaxHeap(data []interface{}, tracker models.StepTracker) {
	n := len(data)

	tracker.AddStep("开始构建最大堆", data, []int{})
	tracker.AddNote("从最后一个非叶子节点开始向上堆化")

	// 从最后一个非叶子节点开始，向上进行堆化
	for i := n/2 - 1; i >= 0; i-- {
		tracker.AddStep("对节点 "+strconv.Itoa(i)+" 进行堆化", data, []int{i})
		hs.heapify(data, n, i, tracker)
	}

	tracker.AddStep("最大堆构建完成", data, []int{})
	tracker.AddNote("堆顶元素为最大值：" + hs.toString(data[0]))
}

// heapify 堆化操作，维护最大堆性质
func (hs *HeapSort) heapify(data []interface{}, heapSize, rootIndex int, tracker models.StepTracker) {
	largest := rootIndex
	leftChild := 2*rootIndex + 1
	rightChild := 2*rootIndex + 2

	// 显示当前处理的节点及其子节点
	highlights := []int{rootIndex}
	if leftChild < heapSize {
		highlights = append(highlights, leftChild)
	}
	if rightChild < heapSize {
		highlights = append(highlights, rightChild)
	}

	tracker.AddStep("检查节点 "+strconv.Itoa(rootIndex)+" 及其子节点", data, highlights)

	// 与左子节点比较
	if leftChild < heapSize {
		tracker.AddStep("比较父节点 "+hs.toString(data[largest])+" 与左子节点 "+
			hs.toString(data[leftChild]), data, []int{largest, leftChild})

		if hs.compare(data[leftChild], data[largest]) > 0 {
			largest = leftChild
			tracker.AddComparison(leftChild, largest, 1)
			tracker.AddNote("左子节点更大")
		} else {
			tracker.AddComparison(leftChild, largest, -1)
		}
	}

	// 与右子节点比较
	if rightChild < heapSize {
		tracker.AddStep("比较当前最大值 "+hs.toString(data[largest])+" 与右子节点 "+
			hs.toString(data[rightChild]), data, []int{largest, rightChild})

		if hs.compare(data[rightChild], data[largest]) > 0 {
			largest = rightChild
			tracker.AddComparison(rightChild, largest, 1)
			tracker.AddNote("右子节点更大")
		} else {
			tracker.AddComparison(rightChild, largest, -1)
		}
	}

	// 如果最大值不是根节点，则交换并继续堆化
	if largest != rootIndex {
		tracker.AddStep("最大值在位置 "+strconv.Itoa(largest)+"，与根节点交换",
			data, []int{rootIndex, largest})

		data[rootIndex], data[largest] = data[largest], data[rootIndex]
		tracker.AddOperation(models.OpTypeSwap, []int{rootIndex, largest},
			[]interface{}{data[rootIndex], data[largest]}, "维护堆的性质")

		tracker.AddStep("交换完成，继续向下堆化", data, []int{largest})

		// 递归堆化受影响的子树
		hs.heapify(data, heapSize, largest, tracker)
	} else {
		tracker.AddStep("堆性质已满足，无需交换", data, []int{rootIndex})
	}
}

// compare 比较两个元素
func (hs *HeapSort) compare(a, b interface{}) int {
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
func (hs *HeapSort) toString(value interface{}) string {
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
func (hs *HeapSort) ValidateInput(data interface{}) error {
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

// IsStable 堆排序不是稳定的
func (hs *HeapSort) IsStable() bool {
	return false
}

// IsInPlace 堆排序是原地的
func (hs *HeapSort) IsInPlace() bool {
	return true
}

// IsAdaptive 堆排序不是自适应的
func (hs *HeapSort) IsAdaptive() bool {
	return false
}

// GetComplexity 获取复杂度信息
func (hs *HeapSort) GetComplexity() algorithms.ComplexityInfo {
	return algorithms.ComplexityInfo{
		TimeComplexity: algorithms.ComplexityCase{
			Best:    "O(n log n)",
			Average: "O(n log n)",
			Worst:   "O(n log n)",
		},
		SpaceComplexity: algorithms.ComplexityCase{
			Best:    "O(1)",
			Average: "O(1)",
			Worst:   "O(1)",
		},
	}
}