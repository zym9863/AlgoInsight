package graph

import (
	"gin/algorithms"
	"gin/models"
	"strconv"
)

// BFS 广度优先搜索算法
type BFS struct {
	algorithms.BaseAlgorithm
}

// NewBFS 创建广度优先搜索算法实例
func NewBFS() *BFS {
	return &BFS{
		BaseAlgorithm: algorithms.BaseAlgorithm{
			ID:              "bfs",
			Name:            "广度优先搜索",
			Category:        models.CategoryGraph,
			Description:     "广度优先搜索（BFS）是一种图遍历算法，从起始节点开始，优先访问距离起始节点最近的节点。使用队列来实现。",
			TimeComplexity:  "O(V + E)",
			SpaceComplexity: "O(V)",
			Parameters: []models.Parameter{
				{
					Name:         "start_node",
					Type:         "int",
					Description:  "起始节点",
					DefaultValue: 0,
					Required:     true,
				},
			},
			Stable:   true,
			InPlace:  false,
			Adaptive: false,
		},
	}
}

// Graph 图的表示结构
type Graph struct {
	Vertices int                 `json:"vertices"` // 顶点数
	AdjList  map[int][]int      `json:"adjList"`  // 邻接表
	Nodes    map[int]interface{} `json:"nodes"`    // 节点数据
}

// Execute 执行广度优先搜索
func (bfs *BFS) Execute(data interface{}, tracker models.StepTracker) (interface{}, error) {
	// 验证输入
	if err := bfs.ValidateInput(data); err != nil {
		return nil, err
	}

	// 转换数据类型
	graph, ok := data.(*Graph)
	if !ok {
		return nil, algorithms.ErrInvalidInput
	}

	// 默认从节点0开始搜索
	startNode := 0

	// 执行搜索
	result, err := bfs.Search(graph, startNode, tracker)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Search 广度优先搜索实现
func (bfs *BFS) Search(graph *Graph, startNode int, tracker models.StepTracker) (map[string]interface{}, error) {
	if graph.Vertices == 0 {
		tracker.AddStep("图为空，搜索结束", nil, []int{})
		return map[string]interface{}{
			"visited": []int{},
			"order":   []int{},
		}, nil
	}

	tracker.SetPhase("初始化")
	tracker.AddStep("开始广度优先搜索，起始节点: "+strconv.Itoa(startNode), nil, []int{startNode})

	// 初始化数据结构
	visited := make(map[int]bool)
	queue := []int{startNode}
	visitOrder := []int{}
	
	tracker.AddStep("初始化队列，添加起始节点", nil, []int{startNode})
	tracker.AddNote("队列: [" + strconv.Itoa(startNode) + "]")

	// BFS主循环
	for len(queue) > 0 {
		tracker.SetPhase("BFS遍历")
		
		// 从队列头部取出节点
		current := queue[0]
		queue = queue[1:]
		
		// 显示队列状态
		queueStr := "["
		for i, node := range queue {
			if i > 0 {
				queueStr += ", "
			}
			queueStr += strconv.Itoa(node)
		}
		queueStr += "]"
		
		tracker.AddStep("从队列取出节点: "+strconv.Itoa(current), nil, []int{current})
		tracker.AddNote("当前队列: " + queueStr)

		// 如果节点未被访问过
		if !visited[current] {
			// 标记为已访问
			visited[current] = true
			visitOrder = append(visitOrder, current)
			
			tracker.AddStep("访问节点: "+strconv.Itoa(current), nil, []int{current})
			tracker.AddOperation(models.OpTypeAccess, []int{current}, []interface{}{current}, "访问节点")
			
			// 获取邻接节点
			neighbors := graph.AdjList[current]
			if len(neighbors) > 0 {
				tracker.AddStep("检查节点 "+strconv.Itoa(current)+" 的邻接节点", nil, neighbors)
				
				// 将未访问的邻接节点加入队列
				newNodes := []int{}
				for _, neighbor := range neighbors {
					if !visited[neighbor] && !bfs.inQueue(queue, neighbor) {
						queue = append(queue, neighbor)
						newNodes = append(newNodes, neighbor)
					}
				}
				
				if len(newNodes) > 0 {
					tracker.AddStep("将未访问的邻接节点加入队列", nil, newNodes)
					
					// 更新队列显示
					newQueueStr := "["
					for i, node := range queue {
						if i > 0 {
							newQueueStr += ", "
						}
						newQueueStr += strconv.Itoa(node)
					}
					newQueueStr += "]"
					tracker.AddNote("更新后队列: " + newQueueStr)
				}
			} else {
				tracker.AddNote("节点 " + strconv.Itoa(current) + " 没有邻接节点")
			}
		}
	}

	tracker.SetPhase("完成")
	tracker.AddStep("BFS搜索完成", nil, visitOrder)
	
	// 显示访问顺序
	orderStr := "["
	for i, node := range visitOrder {
		if i > 0 {
			orderStr += ", "
		}
		orderStr += strconv.Itoa(node)
	}
	orderStr += "]"
	tracker.AddNote("访问顺序: " + orderStr)

	return map[string]interface{}{
		"visited":    visited,
		"order":      visitOrder,
		"startNode":  startNode,
		"totalNodes": len(visitOrder),
	}, nil
}

// inQueue 检查节点是否在队列中
func (bfs *BFS) inQueue(queue []int, node int) bool {
	for _, n := range queue {
		if n == node {
			return true
		}
	}
	return false
}

// ValidateInput 验证输入数据
func (bfs *BFS) ValidateInput(data interface{}) error {
	if data == nil {
		return algorithms.ErrInvalidInput
	}

	graph, ok := data.(*Graph)
	if !ok {
		return algorithms.ErrInvalidInput
	}

	// 检查图的大小限制
	if graph.Vertices > 1000 {
		return algorithms.ErrInvalidInput
	}

	return nil
}

// GetComplexity 获取复杂度信息
func (bfs *BFS) GetComplexity() algorithms.ComplexityInfo {
	return algorithms.ComplexityInfo{
		TimeComplexity: algorithms.ComplexityCase{
			Best:    "O(V + E)",
			Average: "O(V + E)",
			Worst:   "O(V + E)",
		},
		SpaceComplexity: algorithms.ComplexityCase{
			Best:    "O(V)",
			Average: "O(V)",
			Worst:   "O(V)",
		},
	}
}