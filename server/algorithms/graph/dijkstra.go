package graph

import (
	"container/heap"
	"fmt"
	"gin/algorithms"
	"gin/models"
	"math"
)

// Dijkstra Dijkstra最短路径算法
type Dijkstra struct {
	algorithms.BaseAlgorithm
}

// NewDijkstra 创建Dijkstra实例
func NewDijkstra() *Dijkstra {
	return &Dijkstra{
		BaseAlgorithm: algorithms.BaseAlgorithm{
			ID:              "graph_dijkstra",
			Name:            "Dijkstra最短路径算法",
			Category:        models.CategoryGraph,
			Description:     "计算从起始节点到所有其他节点的最短路径，适用于带权重的有向图和无向图。",
			TimeComplexity:  "O((V+E)logV)",
			SpaceComplexity: "O(V)",
			Parameters: []models.Parameter{
				{
					Name:         "start",
					Type:         "string",
					Description:  "起始节点ID",
					DefaultValue: "node_0",
					Required:     false,
				},
			},
			Stable:   false,
			InPlace:  false,
			Adaptive: false,
		},
	}
}

// PriorityQueue 优先队列，用于存储节点和距离
type PriorityQueue []*Item

// Item 优先队列中的元素
type Item struct {
	nodeID   string
	distance float64
	index    int
}

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// Execute 执行Dijkstra算法
func (d *Dijkstra) Execute(data interface{}, tracker models.StepTracker) (interface{}, error) {
	if err := d.ValidateInput(data); err != nil {
		return nil, err
	}

	graph, ok := data.(*models.GraphData)
	if !ok {
		if g2, ok2 := data.(models.GraphData); ok2 {
			graph = &g2
		} else {
			return nil, algorithms.ErrInvalidInput
		}
	}

	tracker.SetPhase("初始化")
	tracker.AddStep("开始Dijkstra最短路径算法", graph, []int{})

	// 构建节点索引和邻接表
	idx := make(map[string]int)
	for i, n := range graph.Nodes {
		idx[n.ID] = i
	}

	// 构建邻接表，包含权重信息
	adj := make(map[string][]EdgeInfo)
	for _, e := range graph.Edges {
		weight := 1.0 // 默认权重
		if e.Weight != nil {
			if w, ok := e.Weight.(float64); ok {
				weight = w
			} else if w, ok := e.Weight.(int); ok {
				weight = float64(w)
			}
		}

		adj[e.From] = append(adj[e.From], EdgeInfo{To: e.To, Weight: weight})
		if graph.Type == "undirected" {
			adj[e.To] = append(adj[e.To], EdgeInfo{To: e.From, Weight: weight})
		}
	}

	// 确定起始节点
	startID := "node_0"
	if len(graph.Nodes) > 0 {
		startID = graph.Nodes[0].ID
	}

	// 初始化距离数组
	distances := make(map[string]float64)
	previous := make(map[string]string)
	visited := make(map[string]bool)

	for _, node := range graph.Nodes {
		distances[node.ID] = math.Inf(1)
	}
	distances[startID] = 0

	// 初始化优先队列
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// 将起始节点加入队列
	heap.Push(&pq, &Item{
		nodeID:   startID,
		distance: 0,
	})

	if sidx, ok := idx[startID]; ok {
		tracker.AddStep("设置起始节点距离为0", graph, []int{sidx})
		tracker.AddOperation(models.OpTypeUpdate, []int{sidx}, []interface{}{0}, "起始节点")
	}

	// Dijkstra主循环
	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*Item)
		currentNodeID := current.nodeID

		// 如果已访问过，跳过
		if visited[currentNodeID] {
			continue
		}

		visited[currentNodeID] = true

		if cidx, ok := idx[currentNodeID]; ok {
			tracker.SetPhase("处理节点")
			tracker.AddStep(fmt.Sprintf("处理节点 %s (距离: %.1f)",
				graph.Nodes[cidx].Label, distances[currentNodeID]), graph, []int{cidx})
			tracker.AddOperation(models.OpTypeAccess, []int{cidx}, nil, "访问节点")
		}

		// 遍历所有邻接节点
		for _, neighbor := range adj[currentNodeID] {
			neighborID := neighbor.To

			if visited[neighborID] {
				continue
			}

			// 计算新距离
			newDistance := distances[currentNodeID] + neighbor.Weight

			// 如果找到更短路径，更新距离
			if newDistance < distances[neighborID] {
				distances[neighborID] = newDistance
				previous[neighborID] = currentNodeID

				// 可视化边的松弛操作
				if cidx, ok1 := idx[currentNodeID]; ok1 {
					if nidx, ok2 := idx[neighborID]; ok2 {
						tracker.AddComparison(cidx, nidx, int(neighbor.Weight))
						tracker.AddStep(fmt.Sprintf("松弛边 %s->%s，新距离: %.1f",
							graph.Nodes[cidx].Label, graph.Nodes[nidx].Label, newDistance),
							graph, []int{nidx})
						tracker.AddOperation(models.OpTypeUpdate, []int{nidx},
							[]interface{}{newDistance}, "更新距离")
					}
				}

				// 将更新后的节点加入队列
				heap.Push(&pq, &Item{
					nodeID:   neighborID,
					distance: newDistance,
				})
			}
		}
	}

	// 构建路径结果
	paths := make(map[string]PathResult)
	for nodeID := range distances {
		path := []string{}
		current := nodeID

		// 重建路径
		for current != "" && current != startID {
			path = append([]string{current}, path...)
			current = previous[current]
		}

		if current == startID {
			path = append([]string{startID}, path...)
		}

		paths[nodeID] = PathResult{
			Distance: distances[nodeID],
			Path:     path,
		}
	}

	return map[string]interface{}{
		"distances": distances,
		"paths":     paths,
		"startNode": startID,
	}, nil
}

// EdgeInfo 边信息结构
type EdgeInfo struct {
	To     string
	Weight float64
}

// PathResult 路径结果
type PathResult struct {
	Distance float64  `json:"distance"`
	Path     []string `json:"path"`
}

// ValidateInput 验证图输入
func (d *Dijkstra) ValidateInput(data interface{}) error {
	if data == nil {
		return algorithms.ErrInvalidInput
	}
	switch g := data.(type) {
	case *models.GraphData:
		return d.validateGraph(g)
	case models.GraphData:
		return d.validateGraph(&g)
	default:
		return algorithms.ErrInvalidInput
	}
}

func (d *Dijkstra) validateGraph(g *models.GraphData) error {
	if len(g.Nodes) == 0 {
		return algorithms.ErrInvalidInput
	}
	// 检查是否有负权重边
	for _, edge := range g.Edges {
		if edge.Weight != nil {
			if w, ok := edge.Weight.(float64); ok && w < 0 {
				return fmt.Errorf("Dijkstra算法不支持负权重边")
			}
			if w, ok := edge.Weight.(int); ok && w < 0 {
				return fmt.Errorf("Dijkstra算法不支持负权重边")
			}
		}
	}
	return nil
}

// ProcessGraph 处理图（与Execute一致）
func (d *Dijkstra) ProcessGraph(graph *models.GraphData, tracker models.StepTracker) (interface{}, error) {
	return d.Execute(graph, tracker)
}

// GetGraphType 图类型
func (d *Dijkstra) GetGraphType() string { return "weighted" }

// GetComplexity 获取 Dijkstra 的时间与空间复杂度信息
func (d *Dijkstra) GetComplexity() algorithms.ComplexityInfo {
	return algorithms.ComplexityInfo{
		TimeComplexity: algorithms.ComplexityCase{
			Best:    "O((V+E)logV)",
			Average: "O((V+E)logV)",
			Worst:   "O((V+E)logV)",
		},
		SpaceComplexity: algorithms.ComplexityCase{
			Best:    "O(V)",
			Average: "O(V)",
			Worst:   "O(V)",
		},
	}
}