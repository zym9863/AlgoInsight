package graph

import (
	"container/heap"
	"fmt"
	"gin/algorithms"
	"gin/models"
	"math"
)

// Prim Prim最小生成树算法
type Prim struct {
	algorithms.BaseAlgorithm
}

// NewPrim 创建Prim实例
func NewPrim() *Prim {
	return &Prim{
		BaseAlgorithm: algorithms.BaseAlgorithm{
			ID:              "graph_prim",
			Name:            "Prim最小生成树算法",
			Category:        models.CategoryGraph,
			Description:     "基于顶点的最小生成树算法，从任意顶点开始逐步扩展MST，适用于稠密图。",
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

// PrimEdge Prim算法中的边
type PrimEdge struct {
	from   string
	to     string
	weight float64
	index  int
}

// PrimPriorityQueue Prim算法优先队列
type PrimPriorityQueue []*PrimEdge

func (pq PrimPriorityQueue) Len() int { return len(pq) }

func (pq PrimPriorityQueue) Less(i, j int) bool {
	return pq[i].weight < pq[j].weight
}

func (pq PrimPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PrimPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	edge := x.(*PrimEdge)
	edge.index = n
	*pq = append(*pq, edge)
}

func (pq *PrimPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	edge := old[n-1]
	old[n-1] = nil
	edge.index = -1
	*pq = old[0 : n-1]
	return edge
}

// Execute 执行Prim算法
func (p *Prim) Execute(data interface{}, tracker models.StepTracker) (interface{}, error) {
	if err := p.ValidateInput(data); err != nil {
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
	tracker.AddStep("开始Prim最小生成树算法", graph, []int{})

	// 检查图类型
	if graph.Type == "directed" {
		return nil, fmt.Errorf("Prim算法只适用于无向图")
	}

	// 构建节点索引和邻接表
	idx := make(map[string]int)
	for i, n := range graph.Nodes {
		idx[n.ID] = i
	}

	// 构建邻接表，包含权重信息
	adj := make(map[string][]PrimEdgeInfo)
	for _, edge := range graph.Edges {
		weight := 1.0 // 默认权重
		if edge.Weight != nil {
			if w, ok := edge.Weight.(float64); ok {
				weight = w
			} else if w, ok := edge.Weight.(int); ok {
				weight = float64(w)
			}
		}

		adj[edge.From] = append(adj[edge.From], PrimEdgeInfo{
			To:     edge.To,
			Weight: weight,
			Label:  edge.Label,
		})
		adj[edge.To] = append(adj[edge.To], PrimEdgeInfo{
			To:     edge.From,
			Weight: weight,
			Label:  edge.Label,
		})
	}

	// 确定起始节点
	startID := "node_0"
	if len(graph.Nodes) > 0 {
		startID = graph.Nodes[0].ID
	}

	// 初始化
	inMST := make(map[string]bool)
	minWeight := make(map[string]float64)
	parent := make(map[string]string)
	mstEdges := make([]PrimMSTEdge, 0)
	totalWeight := 0.0

	// 初始化所有节点的最小权重为无穷大
	for _, node := range graph.Nodes {
		minWeight[node.ID] = math.Inf(1)
	}
	minWeight[startID] = 0

	// 初始化优先队列
	pq := make(PrimPriorityQueue, 0)
	heap.Init(&pq)

	// 将起始节点的所有邻接边加入队列
	if sidx, ok := idx[startID]; ok {
		tracker.AddStep("选择起始节点 "+graph.Nodes[sidx].Label, graph, []int{sidx})
		tracker.AddOperation(models.OpTypeUpdate, []int{sidx}, []interface{}{startID}, "起始节点")
	}

	inMST[startID] = true

	// 将起始节点的所有邻接边加入优先队列
	for _, neighbor := range adj[startID] {
		if !inMST[neighbor.To] {
			heap.Push(&pq, &PrimEdge{
				from:   startID,
				to:     neighbor.To,
				weight: neighbor.Weight,
			})
		}
	}

	tracker.SetPhase("构建最小生成树")

	// Prim主循环
	for pq.Len() > 0 && len(mstEdges) < len(graph.Nodes)-1 {
		// 取出权重最小的边
		minEdge := heap.Pop(&pq).(*PrimEdge)

		// 如果目标节点已经在MST中，跳过这条边
		if inMST[minEdge.to] {
			continue
		}

		// 可视化当前选择的边
		if fromIdx, ok1 := idx[minEdge.from]; ok1 {
			if toIdx, ok2 := idx[minEdge.to]; ok2 {
				tracker.AddComparison(fromIdx, toIdx, int(minEdge.weight))
				tracker.AddStep(fmt.Sprintf("选择最小权重边 %s->%s (权重: %s)",
					graph.Nodes[fromIdx].Label, graph.Nodes[toIdx].Label,
					formatWeight(minEdge.weight)), graph, []int{fromIdx, toIdx})
			}
		}

		// 将目标节点加入MST
		inMST[minEdge.to] = true
		parent[minEdge.to] = minEdge.from
		totalWeight += minEdge.weight

		// 添加边到MST结果中
		mstEdges = append(mstEdges, PrimMSTEdge{
			From:   minEdge.from,
			To:     minEdge.to,
			Weight: minEdge.weight,
		})

		if toIdx, ok := idx[minEdge.to]; ok {
			tracker.AddOperation(models.OpTypeInsert, []int{toIdx},
				[]interface{}{minEdge.to}, "加入MST")
			tracker.AddStep(fmt.Sprintf("节点 %s 加入MST，当前总权重: %s",
				graph.Nodes[toIdx].Label, formatWeight(totalWeight)),
				graph, []int{toIdx})
		}

		// 将新加入节点的所有未访问邻接边加入优先队列
		for _, neighbor := range adj[minEdge.to] {
			if !inMST[neighbor.To] {
				heap.Push(&pq, &PrimEdge{
					from:   minEdge.to,
					to:     neighbor.To,
					weight: neighbor.Weight,
				})
			}
		}
	}

	// 构建结果
	result := map[string]interface{}{
		"mstEdges":    mstEdges,
		"totalWeight": totalWeight,
		"edgeCount":   len(mstEdges),
		"nodeCount":   len(graph.Nodes),
		"startNode":   startID,
	}

	// 检查是否构成连通图
	if len(mstEdges) < len(graph.Nodes)-1 {
		result["isConnected"] = false
		result["message"] = "图不连通，无法构成最小生成树"
	} else {
		result["isConnected"] = true
	}

	return result, nil
}

// PrimEdgeInfo Prim边信息
type PrimEdgeInfo struct {
	To     string
	Weight float64
	Label  string
}

// PrimMSTEdge Prim MST边结果
type PrimMSTEdge struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Weight float64 `json:"weight"`
}

// ValidateInput 验证图输入
func (p *Prim) ValidateInput(data interface{}) error {
	if data == nil {
		return algorithms.ErrInvalidInput
	}
	switch g := data.(type) {
	case *models.GraphData:
		return p.validateGraph(g)
	case models.GraphData:
		return p.validateGraph(&g)
	default:
		return algorithms.ErrInvalidInput
	}
}

func (p *Prim) validateGraph(g *models.GraphData) error {
	if len(g.Nodes) == 0 {
		return algorithms.ErrInvalidInput
	}
	if len(g.Edges) == 0 {
		return algorithms.ErrInvalidInput
	}
	if g.Type == "directed" {
		return fmt.Errorf("Prim算法只适用于无向图")
	}
	return nil
}

// ProcessGraph 处理图（与Execute一致）
func (p *Prim) ProcessGraph(graph *models.GraphData, tracker models.StepTracker) (interface{}, error) {
	return p.Execute(graph, tracker)
}

// GetGraphType 图类型
func (p *Prim) GetGraphType() string { return "undirected_weighted" }

// GetComplexity 获取 Prim 的时间与空间复杂度信息
func (p *Prim) GetComplexity() algorithms.ComplexityInfo {
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