package services

import "errors"

// 定义服务层错误
var (
	// 算法相关错误
	ErrAlgorithmNotFound = errors.New("算法不存在")
	ErrInvalidAlgorithm  = errors.New("无效的算法")
	ErrAlgorithmTimeout  = errors.New("算法执行超时")
	
	// 输入数据相关错误
	ErrInvalidInput      = errors.New("输入数据无效")
	ErrDataSizeTooLarge  = errors.New("数据规模过大")
	ErrUnsupportedDataType = errors.New("不支持的数据类型")
	ErrInvalidPattern    = errors.New("无效的数据模式")
	
	// 可视化相关错误
	ErrSessionNotFound   = errors.New("可视化会话不存在")
	ErrStepNotFound      = errors.New("可视化步骤不存在")
	ErrSessionExpired    = errors.New("可视化会话已过期")
	
	// 性能测试相关错误
	ErrTestNotFound      = errors.New("性能测试不存在")
	ErrTestRunning       = errors.New("性能测试正在运行")
	ErrTestFailed        = errors.New("性能测试失败")
	ErrTooManyTests      = errors.New("并发测试数量过多")
	
	// 系统相关错误
	ErrInternalError     = errors.New("内部服务器错误")
	ErrResourceExhausted = errors.New("系统资源不足")
)
