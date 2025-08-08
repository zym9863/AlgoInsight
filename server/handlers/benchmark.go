package handlers

import (
	"net/http"

	"gin/services"

	"github.com/gin-gonic/gin"
)

// RunBenchmarkRequest 运行性能测试请求结构
type RunBenchmarkRequest struct {
	AlgorithmIDs []string    `json:"algorithmIds" binding:"required"`
	DataSizes    []int       `json:"dataSizes" binding:"required"`
	DataType     string      `json:"dataType" binding:"required"`
	TestCount    int         `json:"testCount,omitempty"`
	Parameters   interface{} `json:"parameters,omitempty"`
}

// RunBenchmark 运行性能测试
func RunBenchmark(c *gin.Context) {
	var req RunBenchmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"message": err.Error(),
		})
		return
	}

	// 验证请求参数
	if len(req.AlgorithmIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"message": "至少需要选择一个算法",
		})
		return
	}

	if len(req.DataSizes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"message": "至少需要指定一个数据规模",
		})
		return
	}

	// 设置默认测试次数
	if req.TestCount <= 0 {
		req.TestCount = 3
	}

	// 执行性能测试
	testID, err := services.RunBenchmarkTest(req.AlgorithmIDs, req.DataSizes, req.DataType, req.TestCount, req.Parameters)
	if err != nil {
		if err == services.ErrInvalidAlgorithm {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "算法无效",
				"message": "包含无效的算法ID",
			})
			return
		}

		if err == services.ErrDataSizeTooLarge {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "数据规模过大",
				"message": "数据规模超出系统限制",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "性能测试启动失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"testId":  testID,
		"message": "性能测试已启动",
	})
}

// GetBenchmarkResults 获取性能测试结果
func GetBenchmarkResults(c *gin.Context) {
	testID := c.Param("testId")

	if testID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"message": "测试ID不能为空",
		})
		return
	}

	results, err := services.GetBenchmarkResults(testID)
	if err != nil {
		if err == services.ErrTestNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "测试不存在",
				"message": "未找到指定的性能测试",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取测试结果失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    results,
	})
}

// CompareBenchmarksRequest 对比性能测试请求结构
type CompareBenchmarksRequest struct {
	TestIDs []string `json:"testIds" binding:"required"`
}

// CompareBenchmarks 对比多个性能测试结果
func CompareBenchmarks(c *gin.Context) {
	var req CompareBenchmarksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"message": err.Error(),
		})
		return
	}

	if len(req.TestIDs) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"message": "至少需要两个测试ID进行对比",
		})
		return
	}

	comparison, err := services.CompareBenchmarkResults(req.TestIDs)
	if err != nil {
		if err == services.ErrTestNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "测试不存在",
				"message": "包含不存在的测试ID",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "对比分析失败",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    comparison,
	})
}
