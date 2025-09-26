// 可视化相关类型定义

export interface VisualizationStep {
  stepId: number;
  description: string;
  data: any;
  highlights: Array<number | string>;
  comparisons: Comparison[];
  operations: Operation[];
  metadata: StepMetadata;
}

export interface Comparison {
  index1: number;
  index2: number;
  result: number; // -1: <, 0: =, 1: >
  type: string;
}

export interface Operation {
  type: string;
  indices: number[];
  values: any[];
  description: string;
}

export interface StepMetadata {
  executionTime: number;
  memoryUsage: number;
  complexity: string;
  phase: string;
  notes: string[];
}

export interface VisualizationResult {
  sessionId: string;
  algorithmId: string;
  inputData: any;
  outputData: any;
  steps: VisualizationStep[];
  totalSteps: number;
  executionTime: number;
  memoryUsage: number;
  statistics: ExecutionStats;
}

export interface ExecutionStats {
  comparisons: number;
  swaps: number;
  moves: number;
  accesses: number;
}

// 可视化请求类型
export interface VisualizationRequest {
  algorithmId: string;
  data: any;
  parameters?: any;
}

// 可视化响应类型
export interface VisualizationResponse {
  success: boolean;
  data: VisualizationResult;
}

// 可视化控制状态
export interface VisualizationState {
  isPlaying: boolean;
  isPaused: boolean;
  currentStep: number;
  totalSteps: number;
  speed: number; // 播放速度 (ms)
  autoPlay: boolean;
}

// 操作类型常量
export const OPERATION_TYPES = {
  SWAP: 'swap',
  MOVE: 'move',
  INSERT: 'insert',
  DELETE: 'delete',
  COMPARE: 'compare',
  ACCESS: 'access',
  UPDATE: 'update',
  MERGE: 'merge',
  SPLIT: 'split',
  PARTITION: 'partition'
} as const;

export type OperationType = typeof OPERATION_TYPES[keyof typeof OPERATION_TYPES];

// 可视化主题
export interface VisualizationTheme {
  primaryColor: string;
  secondaryColor: string;
  highlightColor: string;
  compareColor: string;
  swapColor: string;
  backgroundColor: string;
  textColor: string;
  borderColor: string;
}

// 默认主题
export const DEFAULT_THEME: VisualizationTheme = {
  primaryColor: '#3b82f6',
  secondaryColor: '#64748b',
  highlightColor: '#f59e0b',
  compareColor: '#ef4444',
  swapColor: '#10b981',
  backgroundColor: '#ffffff',
  textColor: '#1f2937',
  borderColor: '#e5e7eb'
};
