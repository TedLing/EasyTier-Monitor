import axios from 'axios';
import { ElMessage } from 'element-plus';
import type { AxiosRequestConfig } from 'axios';
import type { ApiResponse } from '../types';

// 创建axios实例
const apiClient = axios.create({
  timeout: 10000,
});

// 统一错误处理
const handleError = (error: any): never => {
  let errorMessage = '请求失败，请检查网络或接口';
  
  if (error.response) {
    // 服务器返回错误响应
    errorMessage = `请求失败: ${error.response.data?.msg || error.response.statusText}`;
  } else if (error.request) {
    // 请求发送但未收到响应
    errorMessage = '未收到服务器响应，请检查网络连接';
  } else {
    // 请求配置错误
    errorMessage = `请求配置错误: ${error.message}`;
  }
  
  ElMessage.error(errorMessage);
  throw error;
};

// 通用请求方法
export const request = async <T = any>(
  config: AxiosRequestConfig
): Promise<T> => {
  try {
    const response = await apiClient(config);
    const apiResponse = response.data as ApiResponse<T>;
    
    if (apiResponse.code === 0) {
      return apiResponse.data;
    } else {
      ElMessage.error(apiResponse.msg);
      throw new Error(apiResponse.msg);
    }
  } catch (error) {
    return handleError(error);
  }
};

// GET请求封装
export const get = async <T = any>(
  url: string,
  params?: any
): Promise<T> => {
  return request<T>({
    method: 'GET',
    url,
    params,
  });
};

// POST请求封装
export const post = async <T = any>(
  url: string,
  data?: any
): Promise<T> => {
  return request<T>({
    method: 'POST',
    url,
    data,
  });
};
