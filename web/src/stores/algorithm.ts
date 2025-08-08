// 算法相关状态管理

import { writable, derived, get } from 'svelte/store';
import type { Algorithm, AlgorithmCategoryType } from '../types/algorithm';
import { cachedAlgorithmService } from '../services/algorithmService';

// 算法列表状态
export const algorithms = writable<Algorithm[]>([]);
export const algorithmsLoading = writable<boolean>(false);
export const algorithmsError = writable<string | null>(null);

// 当前选中的算法
export const selectedAlgorithm = writable<Algorithm | null>(null);

// 算法类别状态
export const selectedCategory = writable<AlgorithmCategoryType | null>(null);

// 算法搜索状态
export const searchQuery = writable<string>('');
export const searchResults = writable<Algorithm[]>([]);

// 派生状态：按类别分组的算法
export const algorithmsByCategory = derived(
  algorithms,
  ($algorithms) => {
    const grouped: Record<string, Algorithm[]> = {};
    $algorithms.forEach(algorithm => {
      if (!grouped[algorithm.category]) {
        grouped[algorithm.category] = [];
      }
      grouped[algorithm.category].push(algorithm);
    });
    return grouped;
  }
);

// 派生状态：过滤后的算法
export const filteredAlgorithms = derived(
  [algorithms, selectedCategory, searchQuery],
  ([$algorithms, $selectedCategory, $searchQuery]) => {
    let filtered = $algorithms;

    // 按类别过滤
    if ($selectedCategory) {
      filtered = filtered.filter(algorithm => algorithm.category === $selectedCategory);
    }

    // 按搜索查询过滤
    if ($searchQuery.trim()) {
      const query = $searchQuery.toLowerCase();
      filtered = filtered.filter(algorithm =>
        algorithm.name.toLowerCase().includes(query) ||
        algorithm.description.toLowerCase().includes(query) ||
        algorithm.category.toLowerCase().includes(query)
      );
    }

    return filtered;
  }
);

// 算法操作函数
export const algorithmActions = {
  // 加载所有算法
  async loadAlgorithms() {
    algorithmsLoading.set(true);
    algorithmsError.set(null);

    try {
      const data = await cachedAlgorithmService.getAllAlgorithms();
      algorithms.set(data);
    } catch (error) {
      console.error('Failed to load algorithms:', error);
      algorithmsError.set(error instanceof Error ? error.message : '加载算法失败');
    } finally {
      algorithmsLoading.set(false);
    }
  },

  // 按类别加载算法
  async loadAlgorithmsByCategory(category: AlgorithmCategoryType) {
    algorithmsLoading.set(true);
    algorithmsError.set(null);

    try {
      const data = await cachedAlgorithmService.getAlgorithmsByCategory(category);
      algorithms.set(data);
      selectedCategory.set(category);
    } catch (error) {
      console.error('Failed to load algorithms by category:', error);
      algorithmsError.set(error instanceof Error ? error.message : '加载算法失败');
    } finally {
      algorithmsLoading.set(false);
    }
  },

  // 选择算法
  selectAlgorithm(algorithm: Algorithm | null) {
    selectedAlgorithm.set(algorithm);
  },

  // 设置搜索查询
  setSearchQuery(query: string) {
    searchQuery.set(query);
  },

  // 搜索算法
  async searchAlgorithms(query: string) {
    if (!query.trim()) {
      searchResults.set([]);
      return;
    }

    try {
      const results = await cachedAlgorithmService.searchAlgorithms(query);
      searchResults.set(results);
    } catch (error) {
      console.error('Failed to search algorithms:', error);
      searchResults.set([]);
    }
  },

  // 清除搜索
  clearSearch() {
    searchQuery.set('');
    searchResults.set([]);
  },

  // 重置状态
  reset() {
    algorithms.set([]);
    selectedAlgorithm.set(null);
    selectedCategory.set(null);
    searchQuery.set('');
    searchResults.set([]);
    algorithmsLoading.set(false);
    algorithmsError.set(null);
  },

  // 清除缓存并重新加载
  async refresh() {
    cachedAlgorithmService.clearCache();
    await this.loadAlgorithms();
  }
};

// 算法工具函数
export const algorithmUtils = {
  // 根据ID查找算法
  findAlgorithmById(id: string): Algorithm | null {
    const $algorithms = get(algorithms);
    return $algorithms.find(algorithm => algorithm.id === id) || null;
  },

  // 获取算法的复杂度信息
  getComplexityInfo(algorithm: Algorithm) {
    return {
      time: algorithm.timeComplexity,
      space: algorithm.spaceComplexity,
      stable: algorithm.stable,
      inPlace: algorithm.inPlace,
      adaptive: algorithm.adaptive
    };
  },

  // 检查算法是否有参数
  hasParameters(algorithm: Algorithm): boolean {
    return algorithm.parameters && algorithm.parameters.length > 0;
  },

  // 获取必需参数
  getRequiredParameters(algorithm: Algorithm) {
    return algorithm.parameters?.filter(param => param.required) || [];
  },

  // 获取可选参数
  getOptionalParameters(algorithm: Algorithm) {
    return algorithm.parameters?.filter(param => !param.required) || [];
  }
};
