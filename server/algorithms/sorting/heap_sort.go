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
			Description:     "堆排序是一种选择排序算法，利用堆这种数据结构所设计的排序算法。堆是一个近似完全二叉树的结构，并同时满足堆积的性质。",
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
	tracker.AddStep("开始构建最大堆", data, []int{})
	
	// 从最后一个非叶子节点开始向上调整
	for i := n/2 - 1; i >= 0; i-- {
		hs.heapify(data, n, i, tracker)
	}
	
	tracker.AddStep("最大堆构建完成", data, []int{0})
	tracker.AddNote("堆顶元素是最大值: " + hs.toString(data[0]))

	// 第二阶段：逐个取出最大元素
	tracker.SetPhase("堆排序")
	
	for i := n - 1; i > 0; i-- {
		tracker.AddStep("将堆顶元素移到位置 "+strconv.Itoa(i), data, []int{0, i})
		
		// 将当前最大元素（堆顶）移到数组末尾
		data[0], data[i] = data[i], data[0]
		tracker.AddOperation(models.OpTypeSwap, []int{0, i}, []interface{}{data[0], data[i]}, "将最大元素放到正确位置")
		
		// 显示已排序部分
		sortedHighlights := make([]int, 0)
		for j := i; j < n; j++ {
			sortedHighlights = append(sortedHighlights, j)
		}
		tracker.AddStep("已排序部分: ["+strconv.Itoa(i)+", "+strconv.Itoa(n-1)+"]", data, sortedHighlights)
		
		// 重新调整剩余元素为最大堆
		tracker.AddNote("重新调整堆，堆大小: " + strconv.Itoa(i))
		hs.heapify(data, i, 0, tracker)
	}

	tracker.SetPhase("完成")
	tracker.AddStep("堆排序完成", data, []int{})
	return nil
}

// heapify 调整堆，使其满足最大堆性质
func (hs *HeapSort) heapify(data []interface{}, heapSize, rootIndex int, tracker models.StepTracker) {
	largest := rootIndex  // 假设根节点是最大的
	left := 2*rootIndex + 1     // 左子节点
	right := 2*rootIndex + 2    // 右子节点

	// 显示当前处理的节点及其子节点
	nodeHighlights := []int{rootIndex}
	if left < heapSize {
		nodeHighlights = append(nodeHighlights, left)
	}
	if right < heapSize {
		nodeHighlights = append(nodeHighlights, right)
	}
	
	tracker.AddStep("调整节点 "+strconv.Itoa(rootIndex)+" (值: "+hs.toString(data[rootIndex])+")", data, nodeHighlights)

	// 如果左子节点存在且大于根节点
	if left < heapSize {
		tracker.AddComparison(left, rootIndex, hs.compare(data[left], data[largest]))
		if hs.compare(data[left], data[largest]) > 0 {
			largest = left
			tracker.AddNote("左子节点更大: " + hs.toString(data[left]))
		}
	}

	// 如果右子节点存在且大于当前最大值
	if right < heapSize {
		tracker.AddComparison(right, largest, hs.compare(data[right], data[largest]))
		if hs.compare(data[right], data[largest]) > 0 {
			largest = right
			tracker.AddNote("右子节点更大: " + hs.toString(data[right]))
		}
	}

	// 如果最大值不是根节点，则交换并继续调整
	if largest != rootIndex {
		tracker.AddStep("需要交换节点 "+strconv.Itoa(rootIndex)+" 和 "+strconv.Itoa(largest), data, []int{rootIndex, largest})
		
		data[rootIndex], data[largest] = data[largest], data[rootIndex]
		tracker.AddOperation(models.OpTypeSwap, []int{rootIndex, largest}, 
			[]interface{}{data[rootIndex], data[largest]}, "维护堆性质")
		
		tracker.AddStep("交换后的状态", data, []int{rootIndex, largest})
		
		// 递归调整受影响的子树
		hs.heapify(data, heapSize, largest, tracker)
	} else {
		tracker.AddNote("节点 " + strconv.Itoa(rootIndex) + " 已满足堆性质")
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