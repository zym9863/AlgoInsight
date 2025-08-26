package graph

import (
	"gin/algorithms"
	"gin/models"
	"strconv"
)

// DFS 深度优先搜索算法
type DFS struct {
	algorithms.BaseAlgorithm
}

// NewDFS 创建深度优先搜索算法实例
func NewDFS() *DFS {
	return &DFS{
		BaseAlgorithm: algorithms.BaseAlgorithm{
			ID:              "dfs",
			Name:            "深度优先搜索",
			Category:        models.CategoryGraph,
			Description:     "深度优先搜索（DFS）是一种图遍历算法，从起始节点开始，尽可能深地搜索图的分支。使用栈（或递归）来实现。",
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

// Execute 执行深度优先搜索
func (dfs *DFS) Execute(data interface{}, tracker models.StepTracker) (interface{}, error) {
	// 验证输入
	if err := dfs.ValidateInput(data); err != nil {
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
	result, err := dfs.Search(graph, startNode, tracker)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Search 深度优先搜索实现
func (dfs *DFS) Search(graph *Graph, startNode int, tracker models.StepTracker) (map[string]interface{}, error) {
	if graph.Vertices == 0 {
		tracker.AddStep("图为空，搜索结束", nil, []int{})
		return map[string]interface{}{
			"visited": []int{},
			"order":   []int{},
		}, nil
	}

	tracker.SetPhase("初始化")
	tracker.AddStep("开始深度优先搜索，起始节点: "+strconv.Itoa(startNode), nil, []int{startNode})

	// 初始化数据结构
	visited := make(map[int]bool)
	stack := []int{startNode}
	visitOrder := []int{}
	
	tracker.AddStep("初始化栈，添加起始节点", nil, []int{startNode})
	tracker.AddNote("栈: [" + strconv.Itoa(startNode) + "]")

	// DFS主循环（迭代版本）
	for len(stack) > 0 {
		tracker.SetPhase("DFS遍历")
		
		// 从栈顶取出节点
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		
		// 显示栈状态
		stackStr := "["
		for i, node := range stack {
			if i > 0 {
				stackStr += ", "
			}
			stackStr += strconv.Itoa(node)
		}
		stackStr += "]"
		
		tracker.AddStep("从栈取出节点: "+strconv.Itoa(current), nil, []int{current})
		tracker.AddNote("当前栈: " + stackStr)

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
				
				// 将未访问的邻接节点加入栈（倒序添加以保持正确的DFS顺序）
				newNodes := []int{}
				for i := len(neighbors) - 1; i >= 0; i-- {
					neighbor := neighbors[i]
					if !visited[neighbor] && !dfs.inStack(stack, neighbor) {
						stack = append(stack, neighbor)
						newNodes = append(newNodes, neighbor)
					}
				}
				
				if len(newNodes) > 0 {
					tracker.AddStep("将未访问的邻接节点加入栈", nil, newNodes)
					
					// 更新栈显示
					newStackStr := "["
					for i, node := range stack {
						if i > 0 {
							newStackStr += ", "
						}
						newStackStr += strconv.Itoa(node)
					}
					newStackStr += "]"
					tracker.AddNote("更新后栈: " + newStackStr)
				}
			} else {
				tracker.AddNote("节点 " + strconv.Itoa(current) + " 没有邻接节点")
			}
		}
	}

	tracker.SetPhase("完成")
	tracker.AddStep("DFS搜索完成", nil, visitOrder)
	
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

// SearchRecursive 递归版本的DFS（可选实现）
func (dfs *DFS) SearchRecursive(graph *Graph, startNode int, tracker models.StepTracker) (map[string]interface{}, error) {
	visited := make(map[int]bool)
	visitOrder := []int{}
	
	tracker.SetPhase("初始化")
	tracker.AddStep("开始递归深度优先搜索，起始节点: "+strconv.Itoa(startNode), nil, []int{startNode})
	
	dfs.dfsRecursive(graph, startNode, visited, &visitOrder, tracker, 0)
	
	tracker.SetPhase("完成")
	tracker.AddStep("递归DFS搜索完成", nil, visitOrder)
	
	return map[string]interface{}{
		"visited":    visited,
		"order":      visitOrder,
		"startNode":  startNode,
		"totalNodes": len(visitOrder),
	}, nil
}

// dfsRecursive 递归DFS辅助函数
func (dfs *DFS) dfsRecursive(graph *Graph, node int, visited map[int]bool, visitOrder *[]int, tracker models.StepTracker, depth int) {
	// 标记当前节点为已访问
	visited[node] = true
	*visitOrder = append(*visitOrder, node)
	
	tracker.SetPhase("递归深度 " + strconv.Itoa(depth))
	tracker.AddStep("访问节点: "+strconv.Itoa(node), nil, []int{node})
	tracker.AddOperation(models.OpTypeAccess, []int{node}, []interface{}{node}, "递归访问节点")
	
	// 遍历所有邻接节点
	neighbors := graph.AdjList[node]
	if len(neighbors) > 0 {
		tracker.AddStep("检查节点 "+strconv.Itoa(node)+" 的邻接节点", nil, neighbors)
		
		for _, neighbor := range neighbors {
			if !visited[neighbor] {
				tracker.AddNote("递归访问未访问的邻接节点: " + strconv.Itoa(neighbor))
				dfs.dfsRecursive(graph, neighbor, visited, visitOrder, tracker, depth+1)
			}
		}
	}
}

// inStack 检查节点是否在栈中
func (dfs *DFS) inStack(stack []int, node int) bool {
	for _, n := range stack {
		if n == node {
			return true
		}
	}
	return false
}

// ValidateInput 验证输入数据
func (dfs *DFS) ValidateInput(data interface{}) error {
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
func (dfs *DFS) GetComplexity() algorithms.ComplexityInfo {
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