package graph

import (
	"fmt"
	"gin/algorithms"
	"gin/models"
	"sort"
)

// Kruskal Kruskal最小生成树算法
type Kruskal struct {
	algorithms.BaseAlgorithm
}

// NewKruskal 创建Kruskal实例
func NewKruskal() *Kruskal {
	return &Kruskal{
		BaseAlgorithm: algorithms.BaseAlgorithm{
			ID:              "graph_kruskal",
			Name:            "Kruskal最小生成树算法",
			Category:        models.CategoryGraph,
			Description:     "基于边排序的最小生成树算法，使用并查集结构避免环路，适用于稀疏图。",
			TimeComplexity:  "O(ElogE)",
			SpaceComplexity: "O(V)",
			Parameters:      []models.Parameter{},
			Stable:          false,
			InPlace:         false,
			Adaptive:        false,
		},
	}
}

// UnionFind 并查集数据结构
type UnionFind struct {
	parent []int
	rank   []int
	size   int
}

// NewUnionFind 创建并查集
func NewUnionFind(size int) *UnionFind {
	uf := &UnionFind{
		parent: make([]int, size),
		rank:   make([]int, size),
		size:   size,
	}

	// 初始化，每个元素的父节点是自己
	for i := 0; i < size; i++ {
		uf.parent[i] = i
		uf.rank[i] = 0
	}

	return uf
}

// Find 查找元素的根节点（带路径压缩）
func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x]) // 路径压缩
	}
	return uf.parent[x]
}

// Union 合并两个集合（按秩合并）
func (uf *UnionFind) Union(x, y int) bool {
	rootX := uf.Find(x)
	rootY := uf.Find(y)

	if rootX == rootY {
		return false // 已经在同一个集合中
	}

	// 按秩合并
	if uf.rank[rootX] < uf.rank[rootY] {
		uf.parent[rootX] = rootY
	} else if uf.rank[rootX] > uf.rank[rootY] {
		uf.parent[rootY] = rootX
	} else {
		uf.parent[rootY] = rootX
		uf.rank[rootX]++
	}

	return true
}

// EdgeWithWeight 带权重的边
type EdgeWithWeight struct {
	From     string
	To       string
	FromIdx  int
	ToIdx    int
	Weight   float64
	Original *models.GraphEdge
}

// Execute 执行Kruskal算法
func (k *Kruskal) Execute(data interface{}, tracker models.StepTracker) (interface{}, error) {
	if err := k.ValidateInput(data); err != nil {
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
	tracker.AddStep("开始Kruskal最小生成树算法", graph, []int{})

	// 检查图类型
	if graph.Type == "directed" {
		return nil, algorithms.ErrInvalidInput // Kruskal只适用于无向图
	}

	// 构建节点索引
	idx := make(map[string]int)
	for i, n := range graph.Nodes {
		idx[n.ID] = i
	}

	// 提取所有边并排序
	edges := make([]EdgeWithWeight, 0, len(graph.Edges))
	for _, edge := range graph.Edges {
		weight := 1.0 // 默认权重
		if edge.Weight != nil {
			if w, ok := edge.Weight.(float64); ok {
				weight = w
			} else if w, ok := edge.Weight.(int); ok {
				weight = float64(w)
			}
		}

		fromIdx, fromExists := idx[edge.From]
		toIdx, toExists := idx[edge.To]

		if fromExists && toExists {
			edges = append(edges, EdgeWithWeight{
				From:     edge.From,
				To:       edge.To,
				FromIdx:  fromIdx,
				ToIdx:    toIdx,
				Weight:   weight,
				Original: &edge,
			})
		}
	}

	tracker.SetPhase("边排序")
	tracker.AddStep("按权重对所有边进行排序", graph, []int{})

	// 按权重排序边
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].Weight < edges[j].Weight
	})

	// 初始化并查集
	uf := NewUnionFind(len(graph.Nodes))
	mstEdges := make([]EdgeWithWeight, 0)
	totalWeight := 0.0

	tracker.SetPhase("构建最小生成树")

	// Kruskal主循环
	for _, edge := range edges {
		// 可视化当前考虑的边
		tracker.AddComparison(edge.FromIdx, edge.ToIdx, int(edge.Weight))
		tracker.AddStep("检查边 "+graph.Nodes[edge.FromIdx].Label+"->"+graph.Nodes[edge.ToIdx].Label+
			" (权重: "+formatWeight(edge.Weight)+")", graph, []int{edge.FromIdx, edge.ToIdx})

		// 检查是否会形成环
		if uf.Union(edge.FromIdx, edge.ToIdx) {
			// 没有形成环，添加到MST中
			mstEdges = append(mstEdges, edge)
			totalWeight += edge.Weight

			tracker.AddOperation(models.OpTypeInsert, []int{edge.FromIdx, edge.ToIdx},
				[]interface{}{edge.Weight}, "添加到MST")
			tracker.AddStep("边被添加到最小生成树，当前总权重: "+formatWeight(totalWeight),
				graph, []int{edge.FromIdx, edge.ToIdx})

			// 如果已经选择了V-1条边，MST构建完成
			if len(mstEdges) == len(graph.Nodes)-1 {
				break
			}
		} else {
			// 会形成环，跳过这条边
			tracker.AddOperation(models.OpTypeDelete, []int{edge.FromIdx, edge.ToIdx},
				nil, "跳过(形成环)")
		}
	}

	// 构建结果
	result := map[string]interface{}{
		"mstEdges":    convertMSTEdges(mstEdges),
		"totalWeight": totalWeight,
		"edgeCount":   len(mstEdges),
		"nodeCount":   len(graph.Nodes),
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

// MSTEdge MST边结果
type MSTEdge struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Weight float64 `json:"weight"`
	Label  string  `json:"label"`
}

// convertMSTEdges 转换MST边为结果格式
func convertMSTEdges(edges []EdgeWithWeight) []MSTEdge {
	result := make([]MSTEdge, len(edges))
	for i, edge := range edges {
		result[i] = MSTEdge{
			From:   edge.From,
			To:     edge.To,
			Weight: edge.Weight,
			Label:  edge.Original.Label,
		}
	}
	return result
}

// formatWeight 格式化权重显示
func formatWeight(weight float64) string {
	if weight == float64(int(weight)) {
		return fmt.Sprintf("%.0f", weight)
	}
	return fmt.Sprintf("%.1f", weight)
}

// ValidateInput 验证图输入
func (k *Kruskal) ValidateInput(data interface{}) error {
	if data == nil {
		return algorithms.ErrInvalidInput
	}
	switch g := data.(type) {
	case *models.GraphData:
		return k.validateGraph(g)
	case models.GraphData:
		return k.validateGraph(&g)
	default:
		return algorithms.ErrInvalidInput
	}
}

func (k *Kruskal) validateGraph(g *models.GraphData) error {
	if len(g.Nodes) == 0 {
		return algorithms.ErrInvalidInput
	}
	if len(g.Edges) == 0 {
		return algorithms.ErrInvalidInput
	}
	if g.Type == "directed" {
		return fmt.Errorf("Kruskal算法只适用于无向图")
	}
	return nil
}

// ProcessGraph 处理图（与Execute一致）
func (k *Kruskal) ProcessGraph(graph *models.GraphData, tracker models.StepTracker) (interface{}, error) {
	return k.Execute(graph, tracker)
}

// GetGraphType 图类型
func (k *Kruskal) GetGraphType() string { return "undirected_weighted" }

// GetComplexity 获取 Kruskal 的时间与空间复杂度信息
func (k *Kruskal) GetComplexity() algorithms.ComplexityInfo {
	return algorithms.ComplexityInfo{
		TimeComplexity: algorithms.ComplexityCase{
			Best:    "O(ElogE)",
			Average: "O(ElogE)",
			Worst:   "O(ElogE)",
		},
		SpaceComplexity: algorithms.ComplexityCase{
			Best:    "O(V)",
			Average: "O(V)",
			Worst:   "O(V)",
		},
	}
}