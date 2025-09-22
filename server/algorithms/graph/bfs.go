package graph

import (
	"container/list"
	"gin/algorithms"
	"gin/models"
)

// BFS 广度优先搜索
type BFS struct {
	algorithms.BaseAlgorithm
}

// NewBFS 创建BFS实例
func NewBFS() *BFS {
	return &BFS{
		BaseAlgorithm: algorithms.BaseAlgorithm{
			ID:              "graph_bfs",
			Name:            "广度优先搜索 (BFS)",
			Category:        models.CategoryGraph,
			Description:     "从起始节点开始逐层遍历图的所有可达节点。",
			TimeComplexity:  "O(V+E)",
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

// Execute 执行BFS
func (b *BFS) Execute(data interface{}, tracker models.StepTracker) (interface{}, error) {
	if err := b.ValidateInput(data); err != nil {
		return nil, err
	}

	graph, ok := data.(*models.GraphData)
	if !ok {
		if g2, ok2 := data.(models.GraphData); ok2 { // 支持非指针传入
			graph = &g2
		} else {
			return nil, algorithms.ErrInvalidInput
		}
	}

	tracker.SetPhase("初始化")
	tracker.AddStep("开始广度优先搜索", graph, []int{})

	// 构建邻接表和索引
	idx := make(map[string]int)
	for i, n := range graph.Nodes {
		idx[n.ID] = i
	}
	adj := make(map[string][]string)
	for _, e := range graph.Edges {
		adj[e.From] = append(adj[e.From], e.To)
		if graph.Type == "undirected" {
			adj[e.To] = append(adj[e.To], e.From)
		}
	}

	startID := "node_0"
	// 如果存在第一个节点，默认从第一个节点开始
	if len(graph.Nodes) > 0 {
		startID = graph.Nodes[0].ID
	}

	visited := make(map[string]bool)
	order := make([]string, 0, len(graph.Nodes))

	q := list.New()
	q.PushBack(startID)
	visited[startID] = true

	if sidx, ok := idx[startID]; ok {
		tracker.AddStep("入队起始节点", graph, []int{sidx})
		tracker.AddOperation(models.OpTypeUpdate, []int{sidx}, []interface{}{startID}, "起始节点入队")
	}

	for q.Len() > 0 {
		front := q.Front()
		v := front.Value.(string)
		q.Remove(front)

		if vi, ok := idx[v]; ok {
			tracker.SetPhase("访问节点")
			tracker.AddStep("访问节点 "+graph.Nodes[vi].Label, graph, []int{vi})
			tracker.AddOperation(models.OpTypeAccess, []int{vi}, nil, "出队访问")
		}

		order = append(order, v)

		for _, to := range adj[v] {
			if !visited[to] {
				visited[to] = true
				q.PushBack(to)
				// 可视化边的遍历
				if vi, ok1 := idx[v]; ok1 {
					if tj, ok2 := idx[to]; ok2 {
						// 用比较表示边的探索
						tracker.AddComparison(vi, tj, 0)
						tracker.AddStep("发现新节点并入队", graph, []int{tj})
						tracker.AddOperation(models.OpTypeInsert, []int{tj}, []interface{}{to}, "节点入队")
					}
				}
			}
		}
	}

	return map[string]interface{}{
		"order": order,
	}, nil
}

// ValidateInput 验证图输入
func (b *BFS) ValidateInput(data interface{}) error {
	if data == nil {
		return algorithms.ErrInvalidInput
	}
	switch g := data.(type) {
	case *models.GraphData:
		return b.validateGraph(g)
	case models.GraphData:
		return b.validateGraph(&g)
	default:
		return algorithms.ErrInvalidInput
	}
}

func (b *BFS) validateGraph(g *models.GraphData) error {
	if len(g.Nodes) == 0 {
		return algorithms.ErrInvalidInput
	}
	return nil
}

// ProcessGraph 处理图（与Execute一致）
func (b *BFS) ProcessGraph(graph *models.GraphData, tracker models.StepTracker) (interface{}, error) {
	return b.Execute(graph, tracker)
}

// GetGraphType 图类型
func (b *BFS) GetGraphType() string { return "both" }

// GetComplexity 获取 BFS 的时间与空间复杂度信息
// Best/Average/Worst 时间复杂度均为 O(V+E)，空间复杂度为 O(V)
func (b *BFS) GetComplexity() algorithms.ComplexityInfo {
	return algorithms.ComplexityInfo{
		TimeComplexity: algorithms.ComplexityCase{
			Best:    "O(V+E)",
			Average: "O(V+E)",
			Worst:   "O(V+E)",
		},
		SpaceComplexity: algorithms.ComplexityCase{
			Best:    "O(V)",
			Average: "O(V)",
			Worst:   "O(V)",
		},
	}
}
