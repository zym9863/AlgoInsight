package searching

import (
	"gin/algorithms"
	"gin/models"
	"strconv"
)

// HashSearch 哈希搜索算法
type HashSearch struct {
	algorithms.BaseAlgorithm
}

// NewHashSearch 创建哈希搜索算法实例
func NewHashSearch() *HashSearch {
	return &HashSearch{
		BaseAlgorithm: algorithms.BaseAlgorithm{
			ID:              "hash_search",
			Name:            "哈希搜索",
			Category:        models.CategorySearching,
			Description:     "哈希搜索使用哈希表数据结构来实现快速搜索。通过将元素存储在哈希表中，可以在平均情况下实现O(1)的搜索时间复杂度。",
			TimeComplexity:  "O(1)",
			SpaceComplexity: "O(n)",
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
			InPlace:  false,
			Adaptive: false,
		},
	}
}

// hashTable 简单的哈希表实现
type hashTable struct {
	buckets [][]hashEntry // 使用链表法处理冲突
	size    int
}

// hashEntry 哈希表条目
type hashEntry struct {
	key   interface{}
	value int // 存储原数组中的索引
}

// Execute 执行哈希搜索
func (hs *HashSearch) Execute(data interface{}, tracker models.StepTracker) (interface{}, error) {
	// 验证输入
	if err := hs.ValidateInput(data); err != nil {
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
	index, err := hs.Search(arr, target, tracker)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"index":  index,
		"found":  index != -1,
		"target": target,
	}, nil
}

// Search 哈希搜索实现
func (hs *HashSearch) Search(data []interface{}, target interface{}, tracker models.StepTracker) (int, error) {
	n := len(data)
	if n == 0 {
		tracker.AddStep("数组为空，搜索结束", data, []int{})
		return -1, nil
	}

	tracker.SetPhase("初始化")
	tracker.AddStep("开始哈希搜索，目标值: "+hs.toString(target), data, []int{})

	// 第一阶段：构建哈希表
	tracker.SetPhase("构建哈希表")
	tracker.AddNote("为了快速搜索，先将所有元素加入哈希表")

	// 创建哈希表，大小为数组长度的2倍以减少冲突
	hashSize := max(n*2, 16)
	table := hs.createHashTable(hashSize)

	tracker.AddStep("创建哈希表，大小: "+strconv.Itoa(hashSize), data, []int{})

	// 将所有元素插入哈希表
	for i, element := range data {
		tracker.AddStep("插入元素 "+hs.toString(element)+" 到哈希表", data, []int{i})

		hash := hs.hash(element, hashSize)
		tracker.AddNote("元素 " + hs.toString(element) + " 的哈希值: " + strconv.Itoa(hash))

		// 检查是否有冲突
		if len(table.buckets[hash]) > 0 {
			tracker.AddNote("哈希冲突！桶 " + strconv.Itoa(hash) + " 已有元素")
		}

		// 插入到哈希表
		table.buckets[hash] = append(table.buckets[hash], hashEntry{
			key:   element,
			value: i,
		})

		tracker.AddOperation(models.OpTypeInsert, []int{i},
			[]interface{}{element}, "插入元素到哈希表桶 "+strconv.Itoa(hash))
	}

	tracker.AddStep("哈希表构建完成", data, []int{})

	// 第二阶段：搜索目标值
	tracker.SetPhase("哈希搜索")
	tracker.AddStep("在哈希表中搜索目标值: "+hs.toString(target), data, []int{})

	// 计算目标值的哈希值
	targetHash := hs.hash(target, hashSize)
	tracker.AddStep("目标值的哈希值: "+strconv.Itoa(targetHash), data, []int{})
	tracker.AddNote("直接访问桶 " + strconv.Itoa(targetHash))

	// 在对应的桶中搜索
	bucket := table.buckets[targetHash]
	if len(bucket) == 0 {
		// 桶为空，目标不存在
		tracker.AddStep("桶 "+strconv.Itoa(targetHash)+" 为空，目标值不存在", data, []int{})
		tracker.AddNote("搜索失败")
		return -1, nil
	}

	// 在桶中线性搜索
	tracker.AddStep("在桶 "+strconv.Itoa(targetHash)+" 中搜索，桶大小: "+strconv.Itoa(len(bucket)),
		data, []int{})

	for i, entry := range bucket {
		tracker.AddStep("检查桶中第 "+strconv.Itoa(i+1)+" 个元素: "+hs.toString(entry.key),
			data, []int{entry.value})

		if hs.compare(entry.key, target) == 0 {
			// 找到目标元素
			tracker.AddStep("找到目标元素! 位置: "+strconv.Itoa(entry.value), data, []int{entry.value})
			tracker.AddNote("搜索成功，哈希表提供了O(1)平均时间复杂度")
			return entry.value, nil
		}
	}

	// 桶中没有找到目标
	tracker.AddStep("桶 "+strconv.Itoa(targetHash)+" 中未找到目标值", data, []int{})
	tracker.AddNote("搜索失败")
	return -1, nil
}

// createHashTable 创建哈希表
func (hs *HashSearch) createHashTable(size int) *hashTable {
	return &hashTable{
		buckets: make([][]hashEntry, size),
		size:    size,
	}
}

// hash 简单的哈希函数
func (hs *HashSearch) hash(value interface{}, tableSize int) int {
	var hashValue int

	switch v := value.(type) {
	case int:
		hashValue = v
	case float64:
		hashValue = int(v)
	case string:
		// 简单的字符串哈希
		hashValue = 0
		for _, char := range v {
			hashValue = hashValue*31 + int(char)
		}
	default:
		hashValue = 0
	}

	// 确保哈希值为正数并在表大小范围内
	if hashValue < 0 {
		hashValue = -hashValue
	}
	return hashValue % tableSize
}

// max 返回两个整数中的较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// compare 比较两个元素
func (hs *HashSearch) compare(a, b interface{}) int {
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
func (hs *HashSearch) toString(value interface{}) string {
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
func (hs *HashSearch) ValidateInput(data interface{}) error {
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

// RequiresSorted 哈希搜索不需要已排序的数据
func (hs *HashSearch) RequiresSorted() bool {
	return false
}

// GetComplexity 获取复杂度信息
func (hs *HashSearch) GetComplexity() algorithms.ComplexityInfo {
	return algorithms.ComplexityInfo{
		TimeComplexity: algorithms.ComplexityCase{
			Best:    "O(1)",
			Average: "O(1)",
			Worst:   "O(n)", // 所有元素哈希到同一个桶的情况
		},
		SpaceComplexity: algorithms.ComplexityCase{
			Best:    "O(n)",
			Average: "O(n)",
			Worst:   "O(n)",
		},
	}
}