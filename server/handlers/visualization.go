package handlers

import (
	"net/http"

	"gin/services"

	"github.com/gin-gonic/gin"
)

// ExecuteVisualizationRequest 执行可视化请求结构
type ExecuteVisualizationRequest struct {
	AlgorithmID string      `json:"algorithmId" binding:"required"`
	Data        interface{} `json:"data" binding:"required"`
	Parameters  interface{} `json:"parameters,omitempty"`
}

// ExecuteVisualization 执行算法并返回可视化步骤
func ExecuteVisualization(c *gin.Context) {
	var req ExecuteVisualizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"message": err.Error(),
		})
		return
	}

	// 执行算法可视化
	result, err := services.ExecuteAlgorithmVisualization(req.AlgorithmID, req.Data, req.Parameters)
	if err != nil {
		if err == services.ErrAlgorithmNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "算法不存在",
				"message": "未找到指定的算法",
			})
			return
		}

		if err == services.ErrInvalidInput {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "输入数据无效",
				"message": "输入数据格式不正确或超出限制",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "算法执行失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetVisualizationStep 获取特定步骤的可视化数据
func GetVisualizationStep(c *gin.Context) {
	sessionID := c.Param("sessionId")
	stepID := c.Param("stepId")

	if sessionID == "" || stepID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"message": "会话ID和步骤ID不能为空",
		})
		return
	}

	step, err := services.GetVisualizationStep(sessionID, stepID)
	if err != nil {
		if err == services.ErrSessionNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "会话不存在",
				"message": "未找到指定的可视化会话",
			})
			return
		}

		if err == services.ErrStepNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "步骤不存在",
				"message": "未找到指定的可视化步骤",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取步骤数据失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    step,
	})
}

// ResetVisualizationRequest 重置可视化请求结构
type ResetVisualizationRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
}

// ResetVisualization 重置可视化状态
func ResetVisualization(c *gin.Context) {
	var req ResetVisualizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"message": err.Error(),
		})
		return
	}

	err := services.ResetVisualization(req.SessionID)
	if err != nil {
		if err == services.ErrSessionNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "会话不存在",
				"message": "未找到指定的可视化会话",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "重置可视化失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "可视化状态已重置",
	})
}
