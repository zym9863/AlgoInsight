package handlers

import (
	"net/http"

	"gin/services"

	"github.com/gin-gonic/gin"
)

// GetAlgorithms 获取所有可用算法列表
func GetAlgorithms(c *gin.Context) {
	algorithms, err := services.GetAllAlgorithms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取算法列表失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    algorithms,
		"count":   len(algorithms),
	})
}

// GetAlgorithmsByCategory 按类别获取算法
func GetAlgorithmsByCategory(c *gin.Context) {
	category := c.Param("category")

	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"message": "算法类别不能为空",
		})
		return
	}

	algorithms, err := services.GetAlgorithmsByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取算法列表失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"category": category,
		"data":     algorithms,
		"count":    len(algorithms),
	})
}

// GetAlgorithmInfo 获取算法详细信息
func GetAlgorithmInfo(c *gin.Context) {
	algorithmID := c.Param("id")

	if algorithmID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"message": "算法ID不能为空",
		})
		return
	}

	algorithm, err := services.GetAlgorithmInfo(algorithmID)
	if err != nil {
		if err == services.ErrAlgorithmNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "算法不存在",
				"message": "未找到指定的算法",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取算法信息失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    algorithm,
	})
}
