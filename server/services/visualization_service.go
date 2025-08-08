package services

import (
	"gin/models"
	"sync"
	"time"
	"crypto/rand"
	"encoding/hex"
)

// VisualizationService 可视化服务
type VisualizationService struct {
	sessions map[string]*models.VisualizationSession
	mutex    sync.RWMutex
	algorithmService *AlgorithmService
}

// NewVisualizationService 创建可视化服务
func NewVisualizationService() *VisualizationService {
	return &VisualizationService{
		sessions: make(map[string]*models.VisualizationSession),
		algorithmService: NewAlgorithmService(),
	}
}

// ExecuteAlgorithmVisualization 执行算法可视化
func (s *VisualizationService) ExecuteAlgorithmVisualization(algorithmID string, data interface{}, parameters interface{}) (*models.VisualizationResult, error) {
	// 获取算法实例
	algorithm, err := s.algorithmService.GetAlgorithm(algorithmID)
	if err != nil {
		return nil, err
	}

	// 验证输入数据
	if err := algorithm.ValidateInput(data); err != nil {
		return nil, ErrInvalidInput
	}

	// 创建会话
	sessionID := s.generateSessionID()
	session := &models.VisualizationSession{
		ID:          sessionID,
		AlgorithmID: algorithmID,
		InputData:   data,
		Parameters:  parameters,
		Steps:       make([]models.VisualizationStep, 0),
		Status:      models.StatusRunning,
		CreatedAt:   time.Now(),
	}

	// 保存会话
	s.mutex.Lock()
	s.sessions[sessionID] = session
	s.mutex.Unlock()

	// 创建步骤追踪器
	tracker := models.NewStepTracker()

	// 执行算法
	startTime := time.Now()
	outputData, err := algorithm.Execute(data, tracker)
	executionTime := time.Since(startTime)

	// 更新会话状态
	s.mutex.Lock()
	if err != nil {
		session.Status = models.StatusError
		session.Error = err.Error()
	} else {
		session.Status = models.StatusCompleted
		completedAt := time.Now()
		session.CompletedAt = &completedAt
	}
	session.Steps = tracker.GetSteps()
	s.mutex.Unlock()

	if err != nil {
		return nil, err
	}

	// 构建结果
	result := &models.VisualizationResult{
		SessionID:     sessionID,
		AlgorithmID:   algorithmID,
		InputData:     data,
		OutputData:    outputData,
		Steps:         tracker.GetSteps(),
		TotalSteps:    len(tracker.GetSteps()),
		ExecutionTime: executionTime,
		Statistics:    tracker.GetStats(),
	}

	return result, nil
}

// GetVisualizationStep 获取特定步骤
func (s *VisualizationService) GetVisualizationStep(sessionID, stepID string) (*models.VisualizationStep, error) {
	s.mutex.RLock()
	session, exists := s.sessions[sessionID]
	s.mutex.RUnlock()

	if !exists {
		return nil, ErrSessionNotFound
	}

	// 解析步骤ID
	stepIndex := 0
	for i, step := range session.Steps {
		if step.StepID == stepIndex {
			stepIndex = i
			break
		}
	}

	if stepIndex >= len(session.Steps) {
		return nil, ErrStepNotFound
	}

	return &session.Steps[stepIndex], nil
}

// ResetVisualization 重置可视化状态
func (s *VisualizationService) ResetVisualization(sessionID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return ErrSessionNotFound
	}

	// 重置会话状态
	session.Steps = make([]models.VisualizationStep, 0)
	session.Status = models.StatusRunning
	session.Error = ""
	session.CompletedAt = nil

	return nil
}

// generateSessionID 生成会话ID
func (s *VisualizationService) generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// CleanupExpiredSessions 清理过期会话
func (s *VisualizationService) CleanupExpiredSessions() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	expiration := 24 * time.Hour // 24小时过期

	for sessionID, session := range s.sessions {
		if now.Sub(session.CreatedAt) > expiration {
			delete(s.sessions, sessionID)
		}
	}
}

// 全局可视化服务实例
var visualizationService *VisualizationService

// ExecuteAlgorithmVisualization 执行算法可视化（全局函数）
func ExecuteAlgorithmVisualization(algorithmID string, data interface{}, parameters interface{}) (*models.VisualizationResult, error) {
	if visualizationService == nil {
		visualizationService = NewVisualizationService()
	}
	return visualizationService.ExecuteAlgorithmVisualization(algorithmID, data, parameters)
}

// GetVisualizationStep 获取可视化步骤（全局函数）
func GetVisualizationStep(sessionID, stepID string) (*models.VisualizationStep, error) {
	if visualizationService == nil {
		visualizationService = NewVisualizationService()
	}
	return visualizationService.GetVisualizationStep(sessionID, stepID)
}

// ResetVisualization 重置可视化（全局函数）
func ResetVisualization(sessionID string) error {
	if visualizationService == nil {
		visualizationService = NewVisualizationService()
	}
	return visualizationService.ResetVisualization(sessionID)
}
