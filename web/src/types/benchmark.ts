// 性能测试相关类型定义

export interface BenchmarkTest {
  id: string;
  name: string;
  algorithmIds: string[];
  dataSizes: number[];
  dataType: string;
  testCount: number;
  parameters: Record<string, any>;
  status: string;
  createdAt: string;
  startedAt?: string;
  completedAt?: string;
  results: BenchmarkResult[];
  error?: string;
}

export interface BenchmarkResult {
  testId: string;
  algorithmId: string;
  algorithmName: string;
  dataSize: number;
  dataType: string;
  runIndex: number;
  executionTime: number;
  memoryUsage: number;
  operations: number;
  comparisons: number;
  swaps: number;
  success: boolean;
  error?: string;
  timestamp: string;
  metadata: ResultMetadata;
}

export interface ResultMetadata {
  cpuUsage: number;
  gcCount: number;
  gcTime: number;
  allocatedMem: number;
  systemMem: number;
  goRoutines: number;
}

export interface BenchmarkSummary {
  testId: string;
  algorithmId: string;
  dataSize: number;
  runCount: number;
  avgTime: number;
  minTime: number;
  maxTime: number;
  stdDev: number;
  avgMemory: number;
  avgOps: number;
  successRate: number;
  performance: PerformanceMetrics;
}

export interface PerformanceMetrics {
  timeComplexity: string;
  spaceComplexity: string;
  efficiency: number;
  stability: number;
  scalability: number;
}

export interface ComparisonResult {
  testIds: string[];
  algorithms: string[];
  dataSizes: number[];
  comparisons: AlgorithmComparison[];
  summary: ComparisonSummary;
  charts: Record<string, any>;
  createdAt: string;
}

export interface AlgorithmComparison {
  algorithmId: string;
  algorithmName: string;
  results: Record<number, BenchmarkSummary>;
  overallScore: number;
  strengths: string[];
  weaknesses: string[];
}

export interface ComparisonSummary {
  bestTime: string;
  bestMemory: string;
  bestStability: string;
  bestScalability: string;
  recommendation: string;
  reason: string;
}

// 性能测试请求类型
export interface BenchmarkRequest {
  algorithmIds: string[];
  dataSizes: number[];
  dataType: string;
  testCount?: number;
  parameters?: any;
}

// 性能测试响应类型
export interface BenchmarkResponse {
  success: boolean;
  testId: string;
  message: string;
}

export interface BenchmarkResultsResponse {
  success: boolean;
  data: BenchmarkTest;
}

export interface ComparisonResponse {
  success: boolean;
  data: ComparisonResult;
}

// 测试状态常量
export const TEST_STATUS = {
  PENDING: 'pending',
  RUNNING: 'running',
  COMPLETED: 'completed',
  FAILED: 'failed',
  CANCELLED: 'cancelled'
} as const;

export type TestStatus = typeof TEST_STATUS[keyof typeof TEST_STATUS];

// 数据类型常量
export const DATA_TYPES = {
  ARRAY: 'array',
  GRAPH: 'graph',
  TREE: 'tree',
  STRING: 'string',
  MATRIX: 'matrix'
} as const;

export type DataType = typeof DATA_TYPES[keyof typeof DATA_TYPES];

// 图表类型
export interface ChartData {
  labels: string[];
  datasets: ChartDataset[];
}

export interface ChartDataset {
  label: string;
  data: number[];
  backgroundColor?: string;
  borderColor?: string;
  borderWidth?: number;
}

// 性能对比配置
export interface BenchmarkConfig {
  maxExecutionTime: number;
  maxMemoryUsage: number;
  maxDataSize: number;
  timeoutDuration: number;
  warmupRuns: number;
  enableProfiling: boolean;
}
