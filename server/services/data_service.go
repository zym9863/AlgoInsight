package services

import (
	"gin/models"
	"math/rand"
	"time"
)

// DataService 数据服务
type DataService struct {
	presets []models.DataPreset
}

// NewDataService 创建数据服务
func NewDataService() *DataService {
	service := &DataService{
		presets: make([]models.DataPreset, 0),
	}
	
	// 初始化预设数据
	service.initializePresets()
	
	return service
}

// initializePresets 初始化预设数据
func (s *DataService) initializePresets() {
	// 添加一些预设数据集
	s.presets = append(s.presets, models.DataPreset{
		ID:          "small_random",
		Name:        "小型随机数组",
		Description: "包含10个随机整数的数组",
		DataType:    models.DataTypeArray,
		Size:        10,
		Pattern:     models.PatternRandom,
		Data:        []interface{}{64, 34, 25, 12, 22, 11, 90, 5, 77, 30},
		Tags:        []string{"small", "random", "integer"},
		CreatedAt:   time.Now(),
	})

	s.presets = append(s.presets, models.DataPreset{
		ID:          "sorted_array",
		Name:        "已排序数组",
		Description: "包含15个已排序整数的数组",
		DataType:    models.DataTypeArray,
		Size:        15,
		Pattern:     models.PatternSorted,
		Data:        []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		Tags:        []string{"sorted", "integer", "best_case"},
		CreatedAt:   time.Now(),
	})

	s.presets = append(s.presets, models.DataPreset{
		ID:          "reverse_sorted",
		Name:        "逆序数组",
		Description: "包含15个逆序整数的数组",
		DataType:    models.DataTypeArray,
		Size:        15,
		Pattern:     models.PatternReversed,
		Data:        []interface{}{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		Tags:        []string{"reversed", "integer", "worst_case"},
		CreatedAt:   time.Now(),
	})
}

// GenerateTestData 生成测试数据
func (s *DataService) GenerateTestData(dataType string, size int, pattern string, parameters interface{}) (interface{}, error) {
	// 验证数据类型
	if !models.ValidateDataType(dataType) {
		return nil, ErrUnsupportedDataType
	}

	// 验证数据模式
	if pattern != "" && !models.ValidateDataPattern(pattern) {
		return nil, ErrInvalidPattern
	}

	// 根据数据类型生成数据
	switch dataType {
	case models.DataTypeArray:
		return s.generateArrayData(size, pattern, parameters)
	case models.DataTypeGraph:
		return s.generateGraphData(size, pattern, parameters)
	case models.DataTypeTree:
		return s.generateTreeData(size, pattern, parameters)
	default:
		return nil, ErrUnsupportedDataType
	}
}

// generateArrayData 生成数组数据
func (s *DataService) generateArrayData(size int, pattern string, parameters interface{}) (*models.ArrayData, error) {
	values := make([]interface{}, size)
	
	switch pattern {
	case models.PatternRandom:
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < size; i++ {
			values[i] = rand.Intn(100) // 0-99的随机数
		}
	case models.PatternSorted:
		for i := 0; i < size; i++ {
			values[i] = i + 1
		}
	case models.PatternReversed:
		for i := 0; i < size; i++ {
			values[i] = size - i
		}
	case models.PatternNearlySorted:
		// 生成基本有序的数组，然后随机交换少量元素
		for i := 0; i < size; i++ {
			values[i] = i + 1
		}
		// 随机交换10%的元素
		swapCount := size / 10
		if swapCount < 1 {
			swapCount = 1
		}
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < swapCount; i++ {
			idx1 := rand.Intn(size)
			idx2 := rand.Intn(size)
			values[idx1], values[idx2] = values[idx2], values[idx1]
		}
	case models.PatternFewUnique:
		// 生成只有少量唯一值的数组
		uniqueValues := []interface{}{1, 2, 3, 4, 5}
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < size; i++ {
			values[i] = uniqueValues[rand.Intn(len(uniqueValues))]
		}
	default:
		// 默认生成随机数据
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < size; i++ {
			values[i] = rand.Intn(100)
		}
	}

	return &models.ArrayData{
		Values: values,
		Type:   "int",
	}, nil
}

// generateGraphData 生成图数据
func (s *DataService) generateGraphData(size int, pattern string, parameters interface{}) (*models.GraphData, error) {
	nodes := make([]models.GraphNode, size)
	edges := make([]models.GraphEdge, 0)

	// 生成节点
	for i := 0; i < size; i++ {
		nodes[i] = models.GraphNode{
			ID:    "node_" + string(rune(i)),
			Label: "Node " + string(rune(i)),
			Value: i,
			X:     float64(i * 50),
			Y:     float64(i * 50),
		}
	}

	// 简化实现：生成一些随机边
	rand.Seed(time.Now().UnixNano())
	edgeCount := size - 1 // 至少生成一个连通图
	for i := 0; i < edgeCount; i++ {
		from := rand.Intn(size)
		to := rand.Intn(size)
		if from != to {
			edges = append(edges, models.GraphEdge{
				From:   "node_" + string(rune(from)),
				To:     "node_" + string(rune(to)),
				Weight: rand.Intn(10) + 1,
				Label:  "",
			})
		}
	}

	return &models.GraphData{
		Nodes: nodes,
		Edges: edges,
		Type:  "directed",
	}, nil
}

// generateTreeData 生成树数据
func (s *DataService) generateTreeData(size int, pattern string, parameters interface{}) (*models.TreeData, error) {
	if size == 0 {
		return &models.TreeData{
			Root: nil,
			Type: "binary",
		}, nil
	}

	// 简化实现：生成二叉搜索树
	root := &models.TreeNode{
		ID:    "node_0",
		Value: size / 2, // 根节点值
		X:     0,
		Y:     0,
		Level: 0,
	}

	// 递归添加节点（简化实现）
	s.addTreeNodes(root, 1, size, 1)

	return &models.TreeData{
		Root: root,
		Type: "binary",
	}, nil
}

// addTreeNodes 递归添加树节点
func (s *DataService) addTreeNodes(parent *models.TreeNode, currentCount int, maxSize int, level int) int {
	if currentCount >= maxSize {
		return currentCount
	}

	// 添加左子节点
	if currentCount < maxSize {
		leftChild := &models.TreeNode{
			ID:     "node_" + string(rune(currentCount)),
			Value:  currentCount,
			Parent: parent,
			X:      parent.X - float64(50/(level+1)),
			Y:      parent.Y + 50,
			Level:  level,
		}
		parent.Left = leftChild
		parent.Children = append(parent.Children, leftChild)
		currentCount++

		currentCount = s.addTreeNodes(leftChild, currentCount, maxSize, level+1)
	}

	// 添加右子节点
	if currentCount < maxSize {
		rightChild := &models.TreeNode{
			ID:     "node_" + string(rune(currentCount)),
			Value:  currentCount,
			Parent: parent,
			X:      parent.X + float64(50/(level+1)),
			Y:      parent.Y + 50,
			Level:  level,
		}
		parent.Right = rightChild
		parent.Children = append(parent.Children, rightChild)
		currentCount++

		currentCount = s.addTreeNodes(rightChild, currentCount, maxSize, level+1)
	}

	return currentCount
}

// GetDataPresets 获取预设数据
func (s *DataService) GetDataPresets(dataType string) ([]models.DataPreset, error) {
	if dataType == "" {
		return s.presets, nil
	}

	// 按数据类型过滤
	filtered := make([]models.DataPreset, 0)
	for _, preset := range s.presets {
		if preset.DataType == dataType {
			filtered = append(filtered, preset)
		}
	}

	return filtered, nil
}

// 全局数据服务实例
var dataService *DataService

// GenerateTestData 生成测试数据（全局函数）
func GenerateTestData(dataType string, size int, pattern string, parameters interface{}) (interface{}, error) {
	if dataService == nil {
		dataService = NewDataService()
	}
	return dataService.GenerateTestData(dataType, size, pattern, parameters)
}

// GetDataPresets 获取预设数据（全局函数）
func GetDataPresets(dataType string) ([]models.DataPreset, error) {
	if dataService == nil {
		dataService = NewDataService()
	}
	return dataService.GetDataPresets(dataType)
}
