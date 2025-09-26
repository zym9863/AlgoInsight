package services

import (
	"gin/models"
	"math"
	"math/rand"
	"strconv"
	"time"
	"fmt"
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

	// 添加图数据预设
	s.presets = append(s.presets, models.DataPreset{
		ID:          "simple_graph",
		Name:        "简单无向图",
		Description: "包含4个节点的简单无向图",
		DataType:    models.DataTypeGraph,
		Size:        4,
		Pattern:     "simple",
		Data: models.GraphData{
			Type: "undirected",
			Nodes: []models.GraphNode{
				{ID: "A", Label: "A", Value: 0},
				{ID: "B", Label: "B", Value: 1},
				{ID: "C", Label: "C", Value: 2},
				{ID: "D", Label: "D", Value: 3},
			},
			Edges: []models.GraphEdge{
				{From: "A", To: "B", Weight: 1},
				{From: "B", To: "C", Weight: 1},
				{From: "C", To: "D", Weight: 1},
				{From: "D", To: "A", Weight: 1},
			},
		},
		Tags:        []string{"graph", "undirected", "simple"},
		CreatedAt:   time.Now(),
	})

	s.presets = append(s.presets, models.DataPreset{
		ID:          "weighted_graph",
		Name:        "加权有向图",
		Description: "包含权重的有向图示例",
		DataType:    models.DataTypeGraph,
		Size:        5,
		Pattern:     "weighted",
		Data: models.GraphData{
			Type: "directed",
			Nodes: []models.GraphNode{
				{ID: "S", Label: "起点", Value: 0},
				{ID: "A", Label: "A", Value: 1},
				{ID: "B", Label: "B", Value: 2},
				{ID: "C", Label: "C", Value: 3},
				{ID: "T", Label: "终点", Value: 4},
			},
			Edges: []models.GraphEdge{
				{From: "S", To: "A", Weight: 4},
				{From: "S", To: "B", Weight: 2},
				{From: "A", To: "C", Weight: 3},
				{From: "B", To: "A", Weight: 1},
				{From: "B", To: "C", Weight: 5},
				{From: "C", To: "T", Weight: 2},
			},
		},
		Tags:        []string{"graph", "directed", "weighted"},
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
	if size <= 0 {
		return &models.GraphData{
			Nodes: []models.GraphNode{},
			Edges: []models.GraphEdge{},
			Type:  "directed",
		}, nil
	}

	nodes := make([]models.GraphNode, size)
	edges := make([]models.GraphEdge, 0)

	// 生成节点
	for i := 0; i < size; i++ {
		nodes[i] = models.GraphNode{
			ID:    "node_" + strconv.Itoa(i),
			Label: "节点" + strconv.Itoa(i),
			Value: i,
		}
	}

	// 根据不同模式生成不同的图结构
	rand.Seed(time.Now().UnixNano())

	switch pattern {
	case models.PatternRandom:
		// 随机图：随机连接节点
		edgeCount := size + rand.Intn(size*size/4) // 边数量在 size 到 size*(size/4) 之间
		if edgeCount > size*(size-1)/2 {
			edgeCount = size * (size - 1) / 2 // 最大边数（完全图）
		}

		for i := 0; i < edgeCount; i++ {
			from := rand.Intn(size)
			to := rand.Intn(size)
			if from != to {
				// 检查边是否已存在
				exists := false
				for _, existing := range edges {
					if existing.From == nodes[from].ID && existing.To == nodes[to].ID {
						exists = true
						break
					}
				}
				if !exists {
					edges = append(edges, models.GraphEdge{
						From:   nodes[from].ID,
						To:     nodes[to].ID,
						Weight: rand.Intn(10) + 1,
						Label:  "",
					})
				}
			}
		}

	case "complete":
		// 完全图：每个节点都与其他所有节点连接
		for i := 0; i < size; i++ {
			for j := i + 1; j < size; j++ {
				edges = append(edges, models.GraphEdge{
					From:   nodes[i].ID,
					To:     nodes[j].ID,
					Weight: rand.Intn(10) + 1,
					Label:  "",
				})
			}
		}

	case "chain":
		// 链状图：节点呈链状连接
		for i := 0; i < size-1; i++ {
			edges = append(edges, models.GraphEdge{
				From:   nodes[i].ID,
				To:     nodes[i+1].ID,
				Weight: rand.Intn(10) + 1,
				Label:  "",
			})
		}

	case "star":
		// 星状图：一个中心节点连接到所有其他节点
		centerIndex := 0
		for i := 1; i < size; i++ {
			edges = append(edges, models.GraphEdge{
				From:   nodes[centerIndex].ID,
				To:     nodes[i].ID,
				Weight: rand.Intn(10) + 1,
				Label:  "",
			})
		}

	case "tree":
		// 树结构：每个节点（除根节点外）有一个父节点
		for i := 1; i < size; i++ {
			parent := (i - 1) / 2 // 类似二叉树的父节点计算
			edges = append(edges, models.GraphEdge{
				From:   nodes[parent].ID,
				To:     nodes[i].ID,
				Weight: rand.Intn(10) + 1,
				Label:  "",
			})
		}

	case "cycle":
		// 环状图：节点形成一个环
		for i := 0; i < size; i++ {
			nextIndex := (i + 1) % size
			edges = append(edges, models.GraphEdge{
				From:   nodes[i].ID,
				To:     nodes[nextIndex].ID,
				Weight: rand.Intn(10) + 1,
				Label:  "",
			})
		}

	case "grid":
		// 网格图：节点排列成网格，相邻节点相连
		if size >= 4 { // 至少需要4个节点才能形成有意义的网格
			gridSize := int(math.Sqrt(float64(size)))
			actualSize := gridSize * gridSize

			// 调整节点数量以适应网格
			if actualSize != size {
				nodes = nodes[:actualSize]
				size = actualSize
			}

			// 重新生成节点标签以反映网格位置
			for i := 0; i < size; i++ {
				row := i / gridSize
				col := i % gridSize
				nodes[i].Label = fmt.Sprintf("(%d,%d)", row, col)
			}

			// 连接相邻的网格节点
			for i := 0; i < size; i++ {
				row := i / gridSize
				col := i % gridSize

				// 右边的节点
				if col < gridSize-1 {
					right := i + 1
					edges = append(edges, models.GraphEdge{
						From:   nodes[i].ID,
						To:     nodes[right].ID,
						Weight: 1,
						Label:  "",
					})
				}

				// 下边的节点
				if row < gridSize-1 {
					down := i + gridSize
					edges = append(edges, models.GraphEdge{
						From:   nodes[i].ID,
						To:     nodes[down].ID,
						Weight: 1,
						Label:  "",
					})
				}
			}
		} else {
			// 节点太少，退化为链状图
			for i := 0; i < size-1; i++ {
				edges = append(edges, models.GraphEdge{
					From:   nodes[i].ID,
					To:     nodes[i+1].ID,
					Weight: 1,
					Label:  "",
				})
			}
		}

	default:
		// 默认生成连通的随机图
		// 先生成一个生成树确保连通性
		for i := 1; i < size; i++ {
			parent := rand.Intn(i) // 连接到之前的任意节点
			edges = append(edges, models.GraphEdge{
				From:   nodes[parent].ID,
				To:     nodes[i].ID,
				Weight: rand.Intn(10) + 1,
				Label:  "",
			})
		}

		// 然后添加一些额外的随机边
		extraEdges := rand.Intn(size / 2)
		for i := 0; i < extraEdges; i++ {
			from := rand.Intn(size)
			to := rand.Intn(size)
			if from != to {
				// 检查边是否已存在
				exists := false
				for _, existing := range edges {
					if existing.From == nodes[from].ID && existing.To == nodes[to].ID {
						exists = true
						break
					}
				}
				if !exists {
					edges = append(edges, models.GraphEdge{
						From:   nodes[from].ID,
						To:     nodes[to].ID,
						Weight: rand.Intn(10) + 1,
						Label:  "",
					})
				}
			}
		}
	}

	// 确定图的类型（有向图或无向图）
	graphType := "directed"
	if pattern == "complete" || pattern == "grid" {
		graphType = "undirected"
	}

	return &models.GraphData{
		Nodes: nodes,
		Edges: edges,
		Type:  graphType,
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
			ID:     "node_" + strconv.Itoa(currentCount),
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
			ID:     "node_" + strconv.Itoa(currentCount),
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
