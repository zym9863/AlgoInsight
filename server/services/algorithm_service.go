package services

import (
	"gin/algorithms"
	"gin/algorithms/searching"
	"gin/algorithms/sorting"
	"gin/models"
)

// AlgorithmService 算法服务
type AlgorithmService struct {
	registry *algorithms.AlgorithmRegistry
}

// NewAlgorithmService 创建算法服务
func NewAlgorithmService() *AlgorithmService {
	service := &AlgorithmService{
		registry: algorithms.NewAlgorithmRegistry(),
	}
	
	// 注册算法
	service.registerAlgorithms()
	
	return service
}

// registerAlgorithms 注册所有算法
func (s *AlgorithmService) registerAlgorithms() {
	// 注册排序算法
	s.registry.Register(sorting.NewBubbleSort())
	s.registry.Register(sorting.NewQuickSort())
	
	// 注册搜索算法
	s.registry.Register(searching.NewBinarySearch())
	
	// 可以继续注册更多算法...
}

// GetAllAlgorithms 获取所有算法
func (s *AlgorithmService) GetAllAlgorithms() ([]*models.Algorithm, error) {
	algorithms := s.registry.GetAll()
	result := make([]*models.Algorithm, len(algorithms))
	
	for i, algorithm := range algorithms {
		result[i] = algorithm.GetInfo()
	}
	
	return result, nil
}

// GetAlgorithmsByCategory 按类别获取算法
func (s *AlgorithmService) GetAlgorithmsByCategory(category string) ([]*models.Algorithm, error) {
	algorithms := s.registry.GetByCategory(category)
	result := make([]*models.Algorithm, len(algorithms))
	
	for i, algorithm := range algorithms {
		result[i] = algorithm.GetInfo()
	}
	
	return result, nil
}

// GetAlgorithmInfo 获取算法详细信息
func (s *AlgorithmService) GetAlgorithmInfo(algorithmID string) (*models.Algorithm, error) {
	algorithm, exists := s.registry.Get(algorithmID)
	if !exists {
		return nil, ErrAlgorithmNotFound
	}
	
	return algorithm.GetInfo(), nil
}

// GetAlgorithm 获取算法实例（内部使用）
func (s *AlgorithmService) GetAlgorithm(algorithmID string) (algorithms.Algorithm, error) {
	algorithm, exists := s.registry.Get(algorithmID)
	if !exists {
		return nil, ErrAlgorithmNotFound
	}
	
	return algorithm, nil
}

// 全局算法服务实例
var algorithmService *AlgorithmService

// GetAllAlgorithms 获取所有算法（全局函数）
func GetAllAlgorithms() ([]*models.Algorithm, error) {
	if algorithmService == nil {
		algorithmService = NewAlgorithmService()
	}
	return algorithmService.GetAllAlgorithms()
}

// GetAlgorithmsByCategory 按类别获取算法（全局函数）
func GetAlgorithmsByCategory(category string) ([]*models.Algorithm, error) {
	if algorithmService == nil {
		algorithmService = NewAlgorithmService()
	}
	return algorithmService.GetAlgorithmsByCategory(category)
}

// GetAlgorithmInfo 获取算法信息（全局函数）
func GetAlgorithmInfo(algorithmID string) (*models.Algorithm, error) {
	if algorithmService == nil {
		algorithmService = NewAlgorithmService()
	}
	return algorithmService.GetAlgorithmInfo(algorithmID)
}
