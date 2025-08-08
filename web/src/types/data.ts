// 数据相关类型定义

export interface DataPreset {
  id: string;
  name: string;
  description: string;
  dataType: string;
  size: number;
  pattern: string;
  data: any;
  tags: string[];
  createdAt: string;
}

export interface DataGenerationRequest {
  dataType: string;
  size: number;
  pattern?: string;
  parameters?: Record<string, any>;
}

export interface DataGenerationResult {
  dataType: string;
  size: number;
  pattern: string;
  data: any;
  metadata: DataMetadata;
  generatedAt: string;
}

export interface DataMetadata {
  minValue?: any;
  maxValue?: any;
  range?: any;
  duplicates?: number;
  sorted?: boolean;
  reversed?: boolean;
  unique?: boolean;
  distribution?: string;
}

// 数组数据结构
export interface ArrayData {
  values: any[];
  type: string;
}

// 图数据结构
export interface GraphData {
  nodes: GraphNode[];
  edges: GraphEdge[];
  type: string;
}

export interface GraphNode {
  id: string;
  label: string;
  value: any;
  x: number;
  y: number;
}

export interface GraphEdge {
  from: string;
  to: string;
  weight: any;
  label: string;
}

// 树数据结构
export interface TreeData {
  root: TreeNode | null;
  type: string;
}

export interface TreeNode {
  id: string;
  value: any;
  children: TreeNode[];
  parent?: TreeNode;
  left?: TreeNode;
  right?: TreeNode;
  x: number;
  y: number;
  level: number;
}

// 矩阵数据结构
export interface MatrixData {
  values: any[][];
  rows: number;
  cols: number;
  type: string;
}

// 数据模式常量
export const DATA_PATTERNS = {
  RANDOM: 'random',
  SORTED: 'sorted',
  REVERSED: 'reversed',
  NEARLY_SORTED: 'nearly_sorted',
  FEW_UNIQUE: 'few_unique',
  MANY_DUPLICATES: 'many_duplicates',
  WORST_CASE: 'worst_case',
  BEST_CASE: 'best_case',
  AVERAGE_CASE: 'average_case'
} as const;

export type DataPattern = typeof DATA_PATTERNS[keyof typeof DATA_PATTERNS];

// 数据类型常量
export const DATA_TYPES = {
  ARRAY: 'array',
  GRAPH: 'graph',
  TREE: 'tree',
  STRING: 'string',
  MATRIX: 'matrix'
} as const;

export type DataType = typeof DATA_TYPES[keyof typeof DATA_TYPES];

// 数据模式显示名称
export const PATTERN_NAMES: Record<DataPattern, string> = {
  [DATA_PATTERNS.RANDOM]: '随机数据',
  [DATA_PATTERNS.SORTED]: '已排序',
  [DATA_PATTERNS.REVERSED]: '逆序',
  [DATA_PATTERNS.NEARLY_SORTED]: '近似排序',
  [DATA_PATTERNS.FEW_UNIQUE]: '少量唯一值',
  [DATA_PATTERNS.MANY_DUPLICATES]: '大量重复',
  [DATA_PATTERNS.WORST_CASE]: '最坏情况',
  [DATA_PATTERNS.BEST_CASE]: '最佳情况',
  [DATA_PATTERNS.AVERAGE_CASE]: '平均情况'
};

// 数据类型显示名称
export const DATA_TYPE_NAMES: Record<DataType, string> = {
  [DATA_TYPES.ARRAY]: '数组',
  [DATA_TYPES.GRAPH]: '图',
  [DATA_TYPES.TREE]: '树',
  [DATA_TYPES.STRING]: '字符串',
  [DATA_TYPES.MATRIX]: '矩阵'
};

// API响应类型
export interface DataGenerationResponse {
  success: boolean;
  data: {
    type: string;
    size: number;
    pattern: string;
    content: any;
  };
}

export interface DataPresetsResponse {
  success: boolean;
  data: DataPreset[];
  count: number;
}

// 数据输入配置
export interface DataInputConfig {
  allowCustomInput: boolean;
  allowFileUpload: boolean;
  allowPresets: boolean;
  allowGeneration: boolean;
  maxSize: number;
  supportedTypes: DataType[];
  supportedPatterns: DataPattern[];
}
