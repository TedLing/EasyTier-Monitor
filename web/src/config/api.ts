// API配置文件

// API基础URL（优先使用环境变量，如果没有则使用相对路径）
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api';

// API端点
export const API_ENDPOINTS = {
  NODE: `${API_BASE_URL}/node`,
  CONNECTOR: `${API_BASE_URL}/connector`,
  PEER: `${API_BASE_URL}/peer`,
};
