package config

import (
	"os"
)

// Config 应用配置结构
type Config struct {
	// 服务器配置
	Port        string `json:"port"`
	Environment string `json:"environment"`
	
	// 数据库配置（如果需要）
	DatabaseURL string `json:"database_url"`
	
	// 算法执行配置
	MaxExecutionTime int `json:"max_execution_time"` // 秒
	MaxDataSize      int `json:"max_data_size"`      // 最大数据集大小
	
	// 性能测试配置
	BenchmarkTimeout int `json:"benchmark_timeout"` // 秒
	MaxConcurrentTests int `json:"max_concurrent_tests"`
}

// Load 加载配置
func Load() *Config {
	config := &Config{
		Port:               getEnv("PORT", "8080"),
		Environment:        getEnv("ENVIRONMENT", "development"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		MaxExecutionTime:   getEnvInt("MAX_EXECUTION_TIME", 30),
		MaxDataSize:        getEnvInt("MAX_DATA_SIZE", 10000),
		BenchmarkTimeout:   getEnvInt("BENCHMARK_TIMEOUT", 60),
		MaxConcurrentTests: getEnvInt("MAX_CONCURRENT_TESTS", 5),
	}
	
	return config
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 获取整数类型的环境变量
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		// 这里可以添加字符串到整数的转换逻辑
		// 为了简化，暂时返回默认值
		return defaultValue
	}
	return defaultValue
}
