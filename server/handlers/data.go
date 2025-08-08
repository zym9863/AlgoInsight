package handlers

import (
	"net/http"

	"gin/services"

	"github.com/gin-gonic/gin"
)

// GenerateDataRequest 生成数据请求结构
type GenerateDataRequest struct {
	DataType   string      `json:"dataType" binding:"required"` // array, graph, tree等
	Size       int         `json:"size" binding:"required"`     // 数据大小
	Pattern    string      `json:"pattern,omitempty"`           // random, sorted, reverse, nearly_sorted等
	Parameters interface{} `json:"parameters,omitempty"`        // 额外参数
}

// GenerateData 生成测试数据
func GenerateData(c *gin.Context) {
	var req GenerateDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"message": err.Error(),
		})
		return
	}

	// 验证数据大小
	if req.Size <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"message": "数据大小必须大于0",
		})
		return
	}

	if req.Size > 10000 { // 可以从配置中读取
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "数据规模过大",
			"message": "数据大小不能超过10000",
		})
		return
	}

	// 生成数据
	data, err := services.GenerateTestData(req.DataType, req.Size, req.Pattern, req.Parameters)
	if err != nil {
		if err == services.ErrUnsupportedDataType {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "不支持的数据类型",
				"message": "指定的数据类型不受支持",
			})
			return
		}

		if err == services.ErrInvalidPattern {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "无效的数据模式",
				"message": "指定的数据模式不受支持",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "数据生成失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"type":    req.DataType,
			"size":    req.Size,
			"pattern": req.Pattern,
			"content": data,
		},
	})
}

// GetDataPresets 获取预设数据集
func GetDataPresets(c *gin.Context) {
	dataType := c.Query("type") // 可选的数据类型过滤

	presets, err := services.GetDataPresets(dataType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取预设数据失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    presets,
		"count":   len(presets),
	})
}
