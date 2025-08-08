// 可视化相关API服务

import { apiService, handleApiResponse } from './api';
import type {
  VisualizationRequest,
  VisualizationResponse,
  VisualizationResult,
  VisualizationStep
} from '../types/visualization';

// 执行算法可视化
export async function executeVisualization(request: VisualizationRequest): Promise<VisualizationResult> {
  const response = await apiService.post<VisualizationResponse>('/visualize/execute', request);
  handleApiResponse(response);
  return response.data;
}

// 获取可视化步骤
export async function getVisualizationStep(sessionId: string, stepId: string): Promise<VisualizationStep> {
  const response = await apiService.get<{ success: boolean; data: VisualizationStep }>(`/visualize/step/${sessionId}/${stepId}`);
  handleApiResponse(response);
  return response.data;
}

// 重置可视化
export async function resetVisualization(sessionId: string): Promise<void> {
  const response = await apiService.post<{ success: boolean; message: string }>('/visualize/reset', { sessionId });
  handleApiResponse(response);
}

// 可视化管理类
export class VisualizationManager {
  private currentResult: VisualizationResult | null = null;
  private currentStep = 0;
  private isPlaying = false;
  private playInterval: number | null = null;
  private speed = 1000; // 默认1秒一步

  // 执行可视化
  async execute(algorithmId: string, data: any, parameters?: any): Promise<VisualizationResult> {
    this.stop();
    this.reset();

    const request: VisualizationRequest = {
      algorithmId,
      data,
      parameters
    };

    this.currentResult = await executeVisualization(request);
    return this.currentResult;
  }

  // 播放动画
  play(): void {
    if (!this.currentResult || this.isPlaying) return;

    this.isPlaying = true;
    this.playInterval = window.setInterval(() => {
      if (this.currentStep >= this.currentResult!.totalSteps - 1) {
        this.stop();
        return;
      }
      this.nextStep();
    }, this.speed);
  }

  // 暂停动画
  pause(): void {
    this.isPlaying = false;
    if (this.playInterval) {
      clearInterval(this.playInterval);
      this.playInterval = null;
    }
  }

  // 停止动画
  stop(): void {
    this.pause();
  }

  // 下一步
  nextStep(): void {
    if (!this.currentResult) return;
    
    if (this.currentStep < this.currentResult.totalSteps - 1) {
      this.currentStep++;
    }
  }

  // 上一步
  prevStep(): void {
    if (this.currentStep > 0) {
      this.currentStep--;
    }
  }

  // 跳转到指定步骤
  goToStep(step: number): void {
    if (!this.currentResult) return;
    
    if (step >= 0 && step < this.currentResult.totalSteps) {
      this.currentStep = step;
    }
  }

  // 重置到第一步
  reset(): void {
    this.currentStep = 0;
    this.stop();
  }

  // 设置播放速度
  setSpeed(speed: number): void {
    this.speed = Math.max(100, Math.min(5000, speed)); // 限制在100ms-5s之间
    
    if (this.isPlaying) {
      this.pause();
      this.play();
    }
  }

  // 获取当前步骤
  getCurrentStep(): VisualizationStep | null {
    if (!this.currentResult || this.currentStep >= this.currentResult.steps.length) {
      return null;
    }
    return this.currentResult.steps[this.currentStep];
  }

  // 获取当前步骤索引
  getCurrentStepIndex(): number {
    return this.currentStep;
  }

  // 获取总步骤数
  getTotalSteps(): number {
    return this.currentResult?.totalSteps || 0;
  }

  // 获取可视化结果
  getResult(): VisualizationResult | null {
    return this.currentResult;
  }

  // 检查是否正在播放
  getIsPlaying(): boolean {
    return this.isPlaying;
  }

  // 检查是否可以播放
  canPlay(): boolean {
    return this.currentResult !== null && !this.isPlaying && this.currentStep < this.currentResult.totalSteps - 1;
  }

  // 检查是否可以暂停
  canPause(): boolean {
    return this.isPlaying;
  }

  // 检查是否可以前进
  canNext(): boolean {
    return this.currentResult !== null && this.currentStep < this.currentResult.totalSteps - 1;
  }

  // 检查是否可以后退
  canPrev(): boolean {
    return this.currentStep > 0;
  }

  // 获取进度百分比
  getProgress(): number {
    if (!this.currentResult || this.currentResult.totalSteps === 0) {
      return 0;
    }
    return (this.currentStep / (this.currentResult.totalSteps - 1)) * 100;
  }

  // 销毁管理器
  destroy(): void {
    this.stop();
    this.currentResult = null;
    this.currentStep = 0;
  }
}

// 创建全局可视化管理器实例
export const visualizationManager = new VisualizationManager();
