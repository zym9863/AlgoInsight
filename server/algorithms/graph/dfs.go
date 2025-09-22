package graph

import (
	"gin/algorithms"
	"gin/models"
)

// DFS 深度优先搜索
type DFS struct {
	algorithms.BaseAlgorithm
}

// NewDFS 创建DFS实例
func NewDFS() *DFS {
	return &DFS{
		BaseAlgorithm: algorithms.BaseAlgorithm{
			ID:              "graph_dfs",
			Name:            "深度优先搜索 (DFS)",
			Category:        models.CategoryGraph,
			Description:     "沿着路径尽可能深入节点，遇到未访问的邻居继续递归。",
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

// Execute 执行DFS
func (d *DFS) Execute(data interface{}, tracker models.StepTracker) (interface{}, error) {
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
	tracker.AddStep("开始深度优先搜索", graph, []int{})

	// 构建索引和邻接
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
	if len(graph.Nodes) > 0 {
		startID = graph.Nodes[0].ID
	}

	visited := make(map[string]bool)
	order := make([]string, 0, len(graph.Nodes))

	var dfs func(u string)
	dfs = func(u string) {
		if visited[u] {
			return
		}
		visited[u] = true
		if ui, ok := idx[u]; ok {
			tracker.SetPhase("访问节点")
			tracker.AddStep("访问节点 "+graph.Nodes[ui].Label, graph, []int{ui})
			tracker.AddOperation(models.OpTypeAccess, []int{ui}, nil, "进入节点")
		}
		order = append(order, u)
		for _, v := range adj[u] {
			if !visited[v] {
				if ui, ok1 := idx[u]; ok1 {
					if vi, ok2 := idx[v]; ok2 {
						tracker.AddComparison(ui, vi, 0)
						tracker.AddStep("沿边深入", graph, []int{vi})
						tracker.AddOperation(models.OpTypeCall, []int{vi}, nil, "递归访问")
					}
				}
				dfs(v)
				if vi, ok := idx[v]; ok {
					tracker.SetPhase("回溯")
					tracker.AddStep("回溯到节点 "+graph.Nodes[vi].Label, graph, []int{vi})
				}
			}
		}
	}

	dfs(startID)

	return map[string]interface{}{
		"order": order,
	}, nil
}

// ValidateInput 验证图输入
func (d *DFS) ValidateInput(data interface{}) error {
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

func (d *DFS) validateGraph(g *models.GraphData) error {
	if len(g.Nodes) == 0 {
		return algorithms.ErrInvalidInput
	}
	return nil
}

// ProcessGraph 处理图（与Execute一致）
func (d *DFS) ProcessGraph(graph *models.GraphData, tracker models.StepTracker) (interface{}, error) {
	return d.Execute(graph, tracker)
}

// GetGraphType 图类型
func (d *DFS) GetGraphType() string { return "both" }

// GetComplexity 获取 DFS 的时间与空间复杂度信息
// 递归实现下栈深度最坏为 O(V)，时间复杂度在所有情况下均为 O(V+E)
func (d *DFS) GetComplexity() algorithms.ComplexityInfo {
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
