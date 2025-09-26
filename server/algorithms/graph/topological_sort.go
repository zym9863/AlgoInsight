package graph

import (
	"fmt"
	"gin/algorithms"
	"gin/models"
)

// TopologicalSort 拓扑排序算法
type TopologicalSort struct {
	algorithms.BaseAlgorithm
}

// NewTopologicalSort 创建拓扑排序实例
func NewTopologicalSort() *TopologicalSort {
	return &TopologicalSort{
		BaseAlgorithm: algorithms.BaseAlgorithm{
			ID:              "graph_topological_sort",
			Name:            "拓扑排序算法",
			Category:        models.CategoryGraph,
			Description:     "对有向无环图(DAG)进行拓扑排序，产生一个线性顺序，使得所有边都从前面的顶点指向后面的顶点。",
			TimeComplexity:  "O(V+E)",
			SpaceComplexity: "O(V)",
			Parameters: []models.Parameter{
				{
					Name:         "method",
					Type:         "string",
					Description:  "拓扑排序方法 (dfs: 基于DFS, kahn: Kahn算法)",
					DefaultValue: "kahn",
					Required:     false,
				},
			},
			Stable:   false,
			InPlace:  false,
			Adaptive: false,
		},
	}
}

// Execute 执行拓扑排序算法
func (t *TopologicalSort) Execute(data interface{}, tracker models.StepTracker) (interface{}, error) {
	if err := t.ValidateInput(data); err != nil {
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
	tracker.AddStep("开始拓扑排序算法", graph, []int{})

	// 检查图类型
	if graph.Type != "directed" {
		return nil, fmt.Errorf("拓扑排序只适用于有向图")
	}

	// 使用Kahn算法进行拓扑排序
	return t.kahnTopologicalSort(graph, tracker)
}

// kahnTopologicalSort Kahn算法实现拓扑排序
func (t *TopologicalSort) kahnTopologicalSort(graph *models.GraphData, tracker models.StepTracker) (interface{}, error) {
	// 构建节点索引和邻接表
	idx := make(map[string]int)
	for i, n := range graph.Nodes {
		idx[n.ID] = i
	}

	// 构建邻接表和入度数组
	adj := make(map[string][]string)
	inDegree := make(map[string]int)

	// 初始化入度为0
	for _, node := range graph.Nodes {
		inDegree[node.ID] = 0
	}

	// 构建邻接表并计算入度
	for _, edge := range graph.Edges {
		adj[edge.From] = append(adj[edge.From], edge.To)
		inDegree[edge.To]++
	}

	tracker.SetPhase("计算入度")
	tracker.AddStep("计算所有节点的入度", graph, []int{})

	// 找到所有入度为0的节点
	queue := make([]string, 0)
	for nodeID, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, nodeID)
			if nodeIdx, ok := idx[nodeID]; ok {
				tracker.AddStep("找到入度为0的节点: "+graph.Nodes[nodeIdx].Label, graph, []int{nodeIdx})
				tracker.AddOperation(models.OpTypeInsert, []int{nodeIdx}, []interface{}{nodeID}, "入度为0")
			}
		}
	}

	result := make([]string, 0, len(graph.Nodes))
	processedCount := 0

	tracker.SetPhase("拓扑排序")

	// Kahn算法主循环
	for len(queue) > 0 {
		// 取出队首节点
		current := queue[0]
		queue = queue[1:]

		// 添加到结果中
		result = append(result, current)
		processedCount++

		if currentIdx, ok := idx[current]; ok {
			tracker.AddStep("处理节点 "+graph.Nodes[currentIdx].Label+
				fmt.Sprintf(" (已处理: %d/%d)", processedCount, len(graph.Nodes)), graph, []int{currentIdx})
			tracker.AddOperation(models.OpTypeAccess, []int{currentIdx}, nil, "处理节点")
		}

		// 遍历当前节点的所有邻接节点
		for _, neighbor := range adj[current] {
			// 减少邻接节点的入度
			inDegree[neighbor]--

			if currentIdx, ok1 := idx[current]; ok1 {
				if neighborIdx, ok2 := idx[neighbor]; ok2 {
					tracker.AddComparison(currentIdx, neighborIdx, inDegree[neighbor])
					tracker.AddStep(fmt.Sprintf("减少节点 %s 的入度至 %d",
						graph.Nodes[neighborIdx].Label, inDegree[neighbor]), graph, []int{neighborIdx})
				}
			}

			// 如果邻接节点的入度变为0，加入队列
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
				if neighborIdx, ok := idx[neighbor]; ok {
					tracker.AddOperation(models.OpTypeInsert, []int{neighborIdx},
						[]interface{}{neighbor}, "入度变为0，入队")
				}
			}
		}
	}

	// 检查是否存在环
	hasCycle := len(result) != len(graph.Nodes)
	var cycleNodes []string

	if hasCycle {
		// 找出参与环的节点
		for nodeID, degree := range inDegree {
			if degree > 0 {
				cycleNodes = append(cycleNodes, nodeID)
			}
		}
	}

	// 构建结果
	resultMap := map[string]interface{}{
		"topologicalOrder": result,
		"isDAG":            !hasCycle,
		"processedCount":   len(result),
		"totalNodes":       len(graph.Nodes),
	}

	if hasCycle {
		resultMap["hasCycle"] = true
		resultMap["cycleNodes"] = cycleNodes
		resultMap["message"] = "图中存在环，无法进行拓扑排序"

		// 可视化环节点
		cycleIndices := make([]int, 0)
		for _, nodeID := range cycleNodes {
			if nodeIdx, ok := idx[nodeID]; ok {
				cycleIndices = append(cycleIndices, nodeIdx)
			}
		}
		if len(cycleIndices) > 0 {
			tracker.AddStep("检测到环，涉及节点: "+fmt.Sprintf("%v", cycleNodes), graph, cycleIndices)
		}
	} else {
		resultMap["hasCycle"] = false
		resultMap["message"] = "拓扑排序完成"
	}

	return resultMap, nil
}

// dfsTopologicalSort DFS算法实现拓扑排序（备用实现）
func (t *TopologicalSort) dfsTopologicalSort(graph *models.GraphData, tracker models.StepTracker) (interface{}, error) {
	// 构建节点索引和邻接表
	idx := make(map[string]int)
	for i, n := range graph.Nodes {
		idx[n.ID] = i
	}

	adj := make(map[string][]string)
	for _, edge := range graph.Edges {
		adj[edge.From] = append(adj[edge.From], edge.To)
	}

	// DFS状态：0-未访问，1-正在访问，2-已完成
	visited := make(map[string]int)
	result := make([]string, 0, len(graph.Nodes))
	hasCycle := false

	tracker.SetPhase("DFS遍历")

	// DFS递归函数
	var dfs func(nodeID string) bool
	dfs = func(nodeID string) bool {
		if visited[nodeID] == 1 {
			// 发现环
			return true
		}
		if visited[nodeID] == 2 {
			// 已经处理过
			return false
		}

		// 标记为正在访问
		visited[nodeID] = 1

		if nodeIdx, ok := idx[nodeID]; ok {
			tracker.AddStep("开始DFS访问节点 "+graph.Nodes[nodeIdx].Label, graph, []int{nodeIdx})
			tracker.AddOperation(models.OpTypeAccess, []int{nodeIdx}, nil, "DFS访问")
		}

		// 访问所有邻接节点
		for _, neighbor := range adj[nodeID] {
			if dfs(neighbor) {
				return true // 发现环
			}
		}

		// 标记为已完成
		visited[nodeID] = 2
		result = append([]string{nodeID}, result...) // 前插，保证拓扑顺序

		if nodeIdx, ok := idx[nodeID]; ok {
			tracker.AddStep("完成DFS访问节点 "+graph.Nodes[nodeIdx].Label, graph, []int{nodeIdx})
			tracker.AddOperation(models.OpTypeUpdate, []int{nodeIdx}, []interface{}{nodeID}, "DFS完成")
		}

		return false
	}

	// 对所有未访问的节点进行DFS
	for _, node := range graph.Nodes {
		if visited[node.ID] == 0 {
			if dfs(node.ID) {
				hasCycle = true
				break
			}
		}
	}

	// 构建结果
	resultMap := map[string]interface{}{
		"topologicalOrder": result,
		"isDAG":            !hasCycle,
		"processedCount":   len(result),
		"totalNodes":       len(graph.Nodes),
		"hasCycle":         hasCycle,
	}

	if hasCycle {
		resultMap["message"] = "图中存在环，无法进行拓扑排序"
	} else {
		resultMap["message"] = "拓扑排序完成"
	}

	return resultMap, nil
}

// ValidateInput 验证图输入
func (t *TopologicalSort) ValidateInput(data interface{}) error {
	if data == nil {
		return algorithms.ErrInvalidInput
	}
	switch g := data.(type) {
	case *models.GraphData:
		return t.validateGraph(g)
	case models.GraphData:
		return t.validateGraph(&g)
	default:
		return algorithms.ErrInvalidInput
	}
}

func (t *TopologicalSort) validateGraph(g *models.GraphData) error {
	if len(g.Nodes) == 0 {
		return algorithms.ErrInvalidInput
	}
	if g.Type != "directed" {
		return fmt.Errorf("拓扑排序只适用于有向图")
	}
	return nil
}

// ProcessGraph 处理图（与Execute一致）
func (t *TopologicalSort) ProcessGraph(graph *models.GraphData, tracker models.StepTracker) (interface{}, error) {
	return t.Execute(graph, tracker)
}

// GetGraphType 图类型
func (t *TopologicalSort) GetGraphType() string { return "directed" }

// GetComplexity 获取拓扑排序的时间与空间复杂度信息
func (t *TopologicalSort) GetComplexity() algorithms.ComplexityInfo {
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