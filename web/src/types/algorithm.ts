// 算法相关类型定义

export interface Algorithm {
  id: string;
  name: string;
  category: string;
  description: string;
  timeComplexity: string;
  spaceComplexity: string;
  parameters: Parameter[];
  stable: boolean;
  inPlace: boolean;
  adaptive: boolean;
}

export interface Parameter {
  name: string;
  type: string;
  description: string;
  defaultValue: any;
  required: boolean;
  min?: number;
  max?: number;
  options?: string[];
}

export interface AlgorithmCategory {
  id: string;
  name: string;
  description: string;
  algorithms: Algorithm[];
}

// 算法类别常量
export const ALGORITHM_CATEGORIES = {
  SORTING: 'sorting',
  SEARCHING: 'searching',
  GRAPH: 'graph',
  TREE: 'tree',
  DYNAMIC_PROGRAMMING: 'dynamic_programming',
  GREEDY: 'greedy',
  BACKTRACKING: 'backtracking',
  DIVIDE_CONQUER: 'divide_conquer'
} as const;

export type AlgorithmCategoryType = typeof ALGORITHM_CATEGORIES[keyof typeof ALGORITHM_CATEGORIES];

// 算法类别显示名称
export const CATEGORY_NAMES: Record<AlgorithmCategoryType, string> = {
  [ALGORITHM_CATEGORIES.SORTING]: '排序算法',
  [ALGORITHM_CATEGORIES.SEARCHING]: '搜索算法',
  [ALGORITHM_CATEGORIES.GRAPH]: '图算法',
  [ALGORITHM_CATEGORIES.TREE]: '树算法',
  [ALGORITHM_CATEGORIES.DYNAMIC_PROGRAMMING]: '动态规划',
  [ALGORITHM_CATEGORIES.GREEDY]: '贪心算法',
  [ALGORITHM_CATEGORIES.BACKTRACKING]: '回溯算法',
  [ALGORITHM_CATEGORIES.DIVIDE_CONQUER]: '分治算法'
};

// API响应类型
export interface AlgorithmResponse {
  success: boolean;
  data: Algorithm[];
  count: number;
}

export interface AlgorithmInfoResponse {
  success: boolean;
  data: Algorithm;
}
