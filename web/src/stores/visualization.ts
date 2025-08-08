// 可视化相关状态管理

import { writable, derived, get } from 'svelte/store';
import type { 
  VisualizationResult, 
  VisualizationStep, 
  VisualizationState,
  VisualizationTheme 
} from '../types/visualization';
import { visualizationManager } from '../services/visualizationService';
import { DEFAULT_THEME } from '../types/visualization';

// 可视化结果状态
export const visualizationResult = writable<VisualizationResult | null>(null);
export const visualizationLoading = writable<boolean>(false);
export const visualizationError = writable<string | null>(null);

// 可视化控制状态
export const currentStep = writable<number>(0);
export const isPlaying = writable<boolean>(false);
export const isPaused = writable<boolean>(false);
export const playSpeed = writable<number>(1000); // 毫秒
export const autoPlay = writable<boolean>(false);

// 可视化主题
export const visualizationTheme = writable<VisualizationTheme>(DEFAULT_THEME);

// 派生状态：当前步骤数据
export const currentStepData = derived(
  [visualizationResult, currentStep],
  ([$result, $currentStep]) => {
    if (!$result || $currentStep >= $result.steps.length) {
      return null;
    }
    return $result.steps[$currentStep];
  }
);

// 派生状态：总步骤数
export const totalSteps = derived(
  visualizationResult,
  ($result) => $result?.totalSteps || 0
);

// 派生状态：进度百分比
export const progress = derived(
  [currentStep, totalSteps],
  ([$currentStep, $totalSteps]) => {
    if ($totalSteps === 0) return 0;
    return ($currentStep / ($totalSteps - 1)) * 100;
  }
);

// 派生状态：控制按钮状态
export const controlsState = derived(
  [visualizationResult, currentStep, totalSteps, isPlaying],
  ([$result, $currentStep, $totalSteps, $isPlaying]) => ({
    canPlay: $result !== null && !$isPlaying && $currentStep < $totalSteps - 1,
    canPause: $isPlaying,
    canNext: $result !== null && $currentStep < $totalSteps - 1,
    canPrev: $currentStep > 0,
    canReset: $result !== null && $currentStep > 0
  })
);

// 可视化操作函数
export const visualizationActions = {
  // 执行可视化
  async execute(algorithmId: string, data: any, parameters?: any) {
    visualizationLoading.set(true);
    visualizationError.set(null);

    try {
      const result = await visualizationManager.execute(algorithmId, data, parameters);
      
      visualizationResult.set(result);
      currentStep.set(0);
      isPlaying.set(false);
      isPaused.set(false);

      // 如果启用自动播放，开始播放
      if (get(autoPlay)) {
        this.play();
      }
    } catch (error) {
      console.error('Visualization execution failed:', error);
      visualizationError.set(error instanceof Error ? error.message : '可视化执行失败');
    } finally {
      visualizationLoading.set(false);
    }
  },

  // 播放动画
  play() {
    if (!visualizationManager.canPlay()) return;

    visualizationManager.play();
    isPlaying.set(true);
    isPaused.set(false);

    // 监听步骤变化
    const checkStep = () => {
      const newStep = visualizationManager.getCurrentStepIndex();
      currentStep.set(newStep);

      if (visualizationManager.getIsPlaying()) {
        requestAnimationFrame(checkStep);
      } else {
        isPlaying.set(false);
      }
    };
    requestAnimationFrame(checkStep);
  },

  // 暂停动画
  pause() {
    visualizationManager.pause();
    isPlaying.set(false);
    isPaused.set(true);
  },

  // 停止动画
  stop() {
    visualizationManager.stop();
    isPlaying.set(false);
    isPaused.set(false);
  },

  // 下一步
  nextStep() {
    visualizationManager.nextStep();
    currentStep.set(visualizationManager.getCurrentStepIndex());
  },

  // 上一步
  prevStep() {
    visualizationManager.prevStep();
    currentStep.set(visualizationManager.getCurrentStepIndex());
  },

  // 跳转到指定步骤
  goToStep(step: number) {
    this.stop();
    visualizationManager.goToStep(step);
    currentStep.set(visualizationManager.getCurrentStepIndex());
  },

  // 重置到第一步
  reset() {
    this.stop();
    visualizationManager.reset();
    currentStep.set(0);
    isPaused.set(false);
  },

  // 设置播放速度
  setSpeed(speed: number) {
    visualizationManager.setSpeed(speed);
    playSpeed.set(speed);
  },

  // 设置自动播放
  setAutoPlay(enabled: boolean) {
    autoPlay.set(enabled);
  },

  // 设置主题
  setTheme(theme: Partial<VisualizationTheme>) {
    visualizationTheme.update(current => ({ ...current, ...theme }));
  },

  // 重置主题
  resetTheme() {
    visualizationTheme.set(DEFAULT_THEME);
  },

  // 清除状态
  clear() {
    this.stop();
    visualizationResult.set(null);
    currentStep.set(0);
    isPlaying.set(false);
    isPaused.set(false);
    visualizationError.set(null);
    visualizationManager.destroy();
  }
};

// 可视化工具函数
export const visualizationUtils = {
  // 格式化执行时间
  formatExecutionTime(ms: number): string {
    if (ms < 1000) {
      return `${ms.toFixed(2)}ms`;
    } else if (ms < 60000) {
      return `${(ms / 1000).toFixed(2)}s`;
    } else {
      const minutes = Math.floor(ms / 60000);
      const seconds = ((ms % 60000) / 1000).toFixed(2);
      return `${minutes}m ${seconds}s`;
    }
  },

  // 格式化内存使用
  formatMemoryUsage(bytes: number): string {
    const units = ['B', 'KB', 'MB', 'GB'];
    let size = bytes;
    let unitIndex = 0;

    while (size >= 1024 && unitIndex < units.length - 1) {
      size /= 1024;
      unitIndex++;
    }

    return `${size.toFixed(2)} ${units[unitIndex]}`;
  },

  // 获取操作类型的显示名称
  getOperationTypeName(type: string): string {
    const names: Record<string, string> = {
      swap: '交换',
      move: '移动',
      insert: '插入',
      delete: '删除',
      compare: '比较',
      access: '访问',
      update: '更新',
      merge: '合并',
      split: '分割',
      partition: '分区'
    };
    return names[type] || type;
  },

  // 获取比较结果的显示文本
  getComparisonResultText(result: number): string {
    if (result < 0) return '小于';
    if (result > 0) return '大于';
    return '等于';
  },

  // 检查是否有高亮元素
  hasHighlights(step: VisualizationStep): boolean {
    return step.highlights && step.highlights.length > 0;
  },

  // 检查是否有比较操作
  hasComparisons(step: VisualizationStep): boolean {
    return step.comparisons && step.comparisons.length > 0;
  },

  // 检查是否有其他操作
  hasOperations(step: VisualizationStep): boolean {
    return step.operations && step.operations.length > 0;
  }
};
