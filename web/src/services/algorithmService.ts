// 算法相关API服务

import { apiService, handleApiResponse } from './api';
import type {
  Algorithm,
  AlgorithmResponse,
  AlgorithmInfoResponse,
  AlgorithmCategoryType
} from '../types/algorithm';

// 获取所有算法
export async function getAllAlgorithms(): Promise<Algorithm[]> {
  const response = await apiService.get<AlgorithmResponse>('/algorithms');
  handleApiResponse(response);
  return response.data;
}

// 按类别获取算法
export async function getAlgorithmsByCategory(category: AlgorithmCategoryType): Promise<Algorithm[]> {
  const response = await apiService.get<AlgorithmResponse>(`/algorithms/category/${category}`);
  handleApiResponse(response);
  return response.data;
}

// 获取算法详细信息
export async function getAlgorithmInfo(algorithmId: string): Promise<Algorithm> {
  const response = await apiService.get<AlgorithmInfoResponse>(`/algorithms/info/${algorithmId}`);
  handleApiResponse(response);
  return response.data;
}

// 算法缓存管理
class AlgorithmCache {
  private cache = new Map<string, Algorithm[]>();
  private algorithmCache = new Map<string, Algorithm>();
  private cacheExpiry = 5 * 60 * 1000; // 5分钟过期

  private isExpired(timestamp: number): boolean {
    return Date.now() - timestamp > this.cacheExpiry;
  }

  // 缓存所有算法
  cacheAllAlgorithms(algorithms: Algorithm[]): void {
    this.cache.set('all', algorithms);
    this.cache.set('all_timestamp', [Date.now() as any]);
    
    // 同时缓存单个算法
    algorithms.forEach(algorithm => {
      this.algorithmCache.set(algorithm.id, algorithm);
    });
  }

  // 获取缓存的所有算法
  getCachedAllAlgorithms(): Algorithm[] | null {
    const algorithms = this.cache.get('all');
    // 修复类型转换问题，确保 timestamp 是 number 类型
    const timestampArr = this.cache.get('all_timestamp') as number[] | undefined;
    const timestamp = timestampArr ? timestampArr[0] : undefined;
    
    if (algorithms && timestamp && !this.isExpired(timestamp)) {
      return algorithms;
    }
    
    return null;
  }

  // 缓存分类算法
  cacheCategoryAlgorithms(category: string, algorithms: Algorithm[]): void {
    this.cache.set(`category_${category}`, algorithms);
    this.cache.set(`category_${category}_timestamp`, [Date.now() as any]);
  }

  // 获取缓存的分类算法
  getCachedCategoryAlgorithms(category: string): Algorithm[] | null {
    const algorithms = this.cache.get(`category_${category}`);
    // 修复类型转换问题，确保 timestamp 是 number 类型
    const timestampArr = this.cache.get(`category_${category}_timestamp`) as number[] | undefined;
    const timestamp = timestampArr ? timestampArr[0] : undefined;
    
    if (algorithms && timestamp && !this.isExpired(timestamp)) {
      return algorithms;
    }
    
    return null;
  }

  // 缓存单个算法
  cacheAlgorithm(algorithm: Algorithm): void {
    this.algorithmCache.set(algorithm.id, algorithm);
  }

  // 获取缓存的算法
  getCachedAlgorithm(algorithmId: string): Algorithm | null {
    return this.algorithmCache.get(algorithmId) || null;
  }

  // 清除缓存
  clearCache(): void {
    this.cache.clear();
    this.algorithmCache.clear();
  }
}

// 创建缓存实例
const algorithmCache = new AlgorithmCache();

// 带缓存的算法服务
export class CachedAlgorithmService {
  // 获取所有算法（带缓存）
  async getAllAlgorithms(): Promise<Algorithm[]> {
    const cached = algorithmCache.getCachedAllAlgorithms();
    if (cached) {
      return cached;
    }

    const algorithms = await getAllAlgorithms();
    algorithmCache.cacheAllAlgorithms(algorithms);
    return algorithms;
  }

  // 按类别获取算法（带缓存）
  async getAlgorithmsByCategory(category: AlgorithmCategoryType): Promise<Algorithm[]> {
    const cached = algorithmCache.getCachedCategoryAlgorithms(category);
    if (cached) {
      return cached;
    }

    const algorithms = await getAlgorithmsByCategory(category);
    algorithmCache.cacheCategoryAlgorithms(category, algorithms);
    return algorithms;
  }

  // 获取算法信息（带缓存）
  async getAlgorithmInfo(algorithmId: string): Promise<Algorithm> {
    const cached = algorithmCache.getCachedAlgorithm(algorithmId);
    if (cached) {
      return cached;
    }

    const algorithm = await getAlgorithmInfo(algorithmId);
    algorithmCache.cacheAlgorithm(algorithm);
    return algorithm;
  }

  // 按类别分组算法
  async getAlgorithmsByCategories(): Promise<Record<string, Algorithm[]>> {
    const allAlgorithms = await this.getAllAlgorithms();
    const grouped: Record<string, Algorithm[]> = {};

    allAlgorithms.forEach(algorithm => {
      if (!grouped[algorithm.category]) {
        grouped[algorithm.category] = [];
      }
      grouped[algorithm.category].push(algorithm);
    });

    return grouped;
  }

  // 搜索算法
  async searchAlgorithms(query: string): Promise<Algorithm[]> {
    const allAlgorithms = await this.getAllAlgorithms();
    const lowerQuery = query.toLowerCase();

    return allAlgorithms.filter(algorithm =>
      algorithm.name.toLowerCase().includes(lowerQuery) ||
      algorithm.description.toLowerCase().includes(lowerQuery) ||
      algorithm.category.toLowerCase().includes(lowerQuery)
    );
  }

  // 清除缓存
  clearCache(): void {
    algorithmCache.clearCache();
  }
}

// 创建服务实例
export const cachedAlgorithmService = new CachedAlgorithmService();
