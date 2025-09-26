package services

import (
	"fmt"
	"gin/models"
	"strings"
)

// normalizeGraphData 将任意输入尝试转换为 *models.GraphData
func normalizeGraphData(data interface{}) (*models.GraphData, error) {
	switch g := data.(type) {
	case *models.GraphData:
		return validateAndNormalizeGraph(g)
	case models.GraphData:
		return validateAndNormalizeGraph(&g)
	case map[string]interface{}:
		return mapToGraphData(g)
	default:
		return nil, fmt.Errorf("无效的图数据格式")
	}
}

// validateAndNormalizeGraph 验证和标准化图数据
func validateAndNormalizeGraph(graph *models.GraphData) (*models.GraphData, error) {
	if graph == nil {
		return nil, fmt.Errorf("图数据为空")
	}

	// 验证并标准化节点
	if err := validateAndNormalizeNodes(graph); err != nil {
		return nil, err
	}

	// 验证并标准化边
	if err := validateAndNormalizeEdges(graph); err != nil {
		return nil, err
	}

	// 验证图的连通性和一致性
	if err := validateGraphConsistency(graph); err != nil {
		return nil, err
	}

	return graph, nil
}

// validateAndNormalizeNodes 验证和标准化节点
func validateAndNormalizeNodes(graph *models.GraphData) error {
	if len(graph.Nodes) == 0 {
		return fmt.Errorf("图必须至少包含一个节点")
	}

	nodeIDs := make(map[string]bool)

	for i := range graph.Nodes {
		node := &graph.Nodes[i]

		// 确保节点ID不为空
		if strings.TrimSpace(node.ID) == "" {
			node.ID = fmt.Sprintf("node_%d", i)
		}

		// 检查节点ID是否重复
		if nodeIDs[node.ID] {
			return fmt.Errorf("节点ID重复: %s", node.ID)
		}
		nodeIDs[node.ID] = true

		// 确保标签不为空
		if strings.TrimSpace(node.Label) == "" {
			node.Label = node.ID
		}

		// 标准化节点值
		if node.Value == nil {
			node.Value = i
		}
	}

	return nil
}

// validateAndNormalizeEdges 验证和标准化边
func validateAndNormalizeEdges(graph *models.GraphData) error {
	// 创建节点ID映射以快速查找
	nodeExists := make(map[string]bool)
	for _, node := range graph.Nodes {
		nodeExists[node.ID] = true
	}

	// 用于检测重复边的集合
	edgeSet := make(map[string]bool)

	for i := range graph.Edges {
		edge := &graph.Edges[i]

		// 验证边的起点和终点
		if strings.TrimSpace(edge.From) == "" {
			return fmt.Errorf("边%d: 起点ID不能为空", i)
		}
		if strings.TrimSpace(edge.To) == "" {
			return fmt.Errorf("边%d: 终点ID不能为空", i)
		}

		// 检查节点是否存在
		if !nodeExists[edge.From] {
			return fmt.Errorf("边%d: 起点节点'%s'不存在", i, edge.From)
		}
		if !nodeExists[edge.To] {
			return fmt.Errorf("边%d: 终点节点'%s'不存在", i, edge.To)
		}

		// 生成边的唯一标识符
		edgeKey := edge.From + "->" + edge.To
		if graph.Type == "undirected" {
			// 对于无向图，确保边的方向统一
			if edge.From > edge.To {
				edgeKey = edge.To + "->" + edge.From
			}
		}

		// 检查是否有重复边
		if edgeSet[edgeKey] {
			if graph.Type == "undirected" {
				return fmt.Errorf("重复的无向边: %s <-> %s", edge.From, edge.To)
			} else {
				return fmt.Errorf("重复的有向边: %s -> %s", edge.From, edge.To)
			}
		}
		edgeSet[edgeKey] = true

		// 标准化权重
		if edge.Weight == nil {
			edge.Weight = 1 // 默认权重为1
		} else {
			// 尝试将权重转换为数值
			switch w := edge.Weight.(type) {
			case int:
				// 已经是整数，无需处理
			case float64:
				// 已经是浮点数，无需处理
			case string:
				// 尝试解析字符串为数值
				if w == "" {
					edge.Weight = 1
				}
				// 可以添加字符串到数值的转换逻辑
			default:
				edge.Weight = 1
			}
		}

		// 标准化标签
		if strings.TrimSpace(edge.Label) == "" {
			if edge.Weight != nil {
				edge.Label = fmt.Sprintf("%v", edge.Weight)
			} else {
				edge.Label = ""
			}
		}
	}

	return nil
}

// validateGraphConsistency 验证图的一致性
func validateGraphConsistency(graph *models.GraphData) error {
	// 检查图类型
	if graph.Type == "" {
		graph.Type = "directed" // 默认为有向图
	}

	validTypes := map[string]bool{
		"directed":   true,
		"undirected": true,
	}

	if !validTypes[graph.Type] {
		return fmt.Errorf("无效的图类型: %s, 支持的类型: directed, undirected", graph.Type)
	}

	// 对于只有一个节点的图，检查是否有自环
	if len(graph.Nodes) == 1 {
		for _, edge := range graph.Edges {
			if edge.From == edge.To {
				// 允许自环，但发出警告（在实际应用中可以记录日志）
				continue
			}
		}
	}

	// 检查是否存在孤立节点（可选的警告）
	connectedNodes := make(map[string]bool)
	for _, edge := range graph.Edges {
		connectedNodes[edge.From] = true
		connectedNodes[edge.To] = true
	}

	// 统计孤立节点数量（用于潜在的警告）
	isolatedCount := 0
	for _, node := range graph.Nodes {
		if !connectedNodes[node.ID] {
			isolatedCount++
		}
	}

	// 在这里可以根据需要处理孤立节点的情况
	// 例如，可以记录警告日志或返回信息

	return nil
}

func mapToGraphData(m map[string]interface{}) (*models.GraphData, error) {
	graph := &models.GraphData{Type: "directed"}

	// 解析节点
	if nodesVal, ok := m["nodes"]; ok {
		if nodesSlice, ok := nodesVal.([]interface{}); ok {
			nodes := make([]models.GraphNode, 0, len(nodesSlice))
			for i, nv := range nodesSlice {
				if nm, ok := nv.(map[string]interface{}); ok {
					node := models.GraphNode{}

					// 解析节点ID
					if id, ok := nm["id"].(string); ok {
						node.ID = strings.TrimSpace(id)
					}
					if node.ID == "" {
						node.ID = fmt.Sprintf("node_%d", i)
					}

					// 解析节点标签
					if label, ok := nm["label"].(string); ok {
						node.Label = strings.TrimSpace(label)
					}
					if node.Label == "" {
						node.Label = node.ID
					}

					// 解析节点值
					if val, ok := nm["value"]; ok {
						node.Value = val
					} else {
						node.Value = i
					}

					// 解析坐标（可选）
					if x, ok := nm["x"].(float64); ok {
						node.X = x
					}
					if y, ok := nm["y"].(float64); ok {
						node.Y = y
					}

					nodes = append(nodes, node)
				}
			}
			graph.Nodes = nodes
		}
	}

	// 解析边
	if edgesVal, ok := m["edges"]; ok {
		if edgesSlice, ok := edgesVal.([]interface{}); ok {
			edges := make([]models.GraphEdge, 0, len(edgesSlice))
			for _, ev := range edgesSlice {
				if em, ok := ev.(map[string]interface{}); ok {
					edge := models.GraphEdge{}

					// 解析起点和终点
					if from, ok := em["from"].(string); ok {
						edge.From = strings.TrimSpace(from)
					}
					if to, ok := em["to"].(string); ok {
						edge.To = strings.TrimSpace(to)
					}

					// 跳过无效的边
					if edge.From == "" || edge.To == "" {
						continue
					}

					// 解析权重
					if w, ok := em["weight"]; ok {
						edge.Weight = w
					} else {
						edge.Weight = 1 // 默认权重
					}

					// 解析标签
					if label, ok := em["label"].(string); ok {
						edge.Label = strings.TrimSpace(label)
					}

					edges = append(edges, edge)
				}
			}
			graph.Edges = edges
		}
	}

	// 解析图类型
	if t, ok := m["type"].(string); ok {
		graph.Type = strings.ToLower(strings.TrimSpace(t))
	}

	// 验证和标准化图数据
	return validateAndNormalizeGraph(graph)
}