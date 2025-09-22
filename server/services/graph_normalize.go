package services

import (
    "fmt"
    "gin/models"
)

// normalizeGraphData 将任意输入尝试转换为 *models.GraphData
func normalizeGraphData(data interface{}) (*models.GraphData, error) {
    switch g := data.(type) {
    case *models.GraphData:
        return g, nil
    case models.GraphData:
        return &g, nil
    case map[string]interface{}:
        return mapToGraphData(g)
    default:
        return nil, fmt.Errorf("invalid graph data")
    }
}

func mapToGraphData(m map[string]interface{}) (*models.GraphData, error) {
    graph := &models.GraphData{Type: "directed"}
    // 解析节点
    if nodesVal, ok := m["nodes"]; ok {
        if nodesSlice, ok := nodesVal.([]interface{}); ok {
            nodes := make([]models.GraphNode, 0, len(nodesSlice))
            for _, nv := range nodesSlice {
                if nm, ok := nv.(map[string]interface{}); ok {
                    node := models.GraphNode{}
                    if id, ok := nm["id"].(string); ok { node.ID = id }
                    if label, ok := nm["label"].(string); ok { node.Label = label } else { node.Label = node.ID }
                    if val, ok := nm["value"]; ok { node.Value = val }
                    if x, ok := nm["x"].(float64); ok { node.X = x }
                    if y, ok := nm["y"].(float64); ok { node.Y = y }
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
                    if from, ok := em["from"].(string); ok { edge.From = from }
                    if to, ok := em["to"].(string); ok { edge.To = to }
                    if w, ok := em["weight"]; ok { edge.Weight = w }
                    if label, ok := em["label"].(string); ok { edge.Label = label }
                    edges = append(edges, edge)
                }
            }
            graph.Edges = edges
        }
    }
    if t, ok := m["type"].(string); ok { graph.Type = t }

    // 如果节点没有ID，则补上
    for i := range graph.Nodes {
        if graph.Nodes[i].ID == "" {
            graph.Nodes[i].ID = fmt.Sprintf("node_%d", i)
        }
        if graph.Nodes[i].Label == "" {
            graph.Nodes[i].Label = graph.Nodes[i].ID
        }
    }

    return graph, nil
}
