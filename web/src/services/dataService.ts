// 数据相关API服务
import { apiService, handleApiResponse } from './api';
import type { DataGenerationRequest, DataGenerationResponse, DataPreset, DataPresetsResponse } from '../types/data';

export async function generateData(req: DataGenerationRequest): Promise<any> {
  const response = await apiService.post<DataGenerationResponse>('/data/generate', req);
  handleApiResponse(response);
  return response.data.content;
}

export async function getDataPresets(type?: string): Promise<DataPreset[]> {
  const endpoint = type ? `/data/presets?type=${encodeURIComponent(type)}` : '/data/presets';
  const response = await apiService.get<DataPresetsResponse>(endpoint);
  handleApiResponse(response);
  return response.data;
}
