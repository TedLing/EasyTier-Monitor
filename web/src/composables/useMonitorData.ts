import { ref, onMounted } from 'vue';
import { ElLoading } from 'element-plus';
import { get } from './useApi';
import { API_ENDPOINTS } from '../config/api';
import type { NodeInfo, Connector, Peer } from '../types';

// 获取当前时间的格式化函数
const formatCurrentTime = (): string => {
  const time = new Date();
  const year = time.getFullYear();
  const month = (time.getMonth() + 1).toString().padStart(2, '0');
  const date = time.getDate().toString().padStart(2, '0');
  const hour = time.getHours().toString().padStart(2, '0');
  const minute = time.getMinutes().toString().padStart(2, '0');
  const second = time.getSeconds().toString().padStart(2, '0');
  
  return `${year}-${month}-${date} ${hour}:${minute}:${second}`;
};

export const useMonitorData = () => {
  console.log('当前版本是新版本');
  // 状态管理
  const nodeInfo = ref<NodeInfo | null>(null);
  const connectors = ref<Connector[]>([]);
  const peers = ref<Peer[]>([]);
  const deviceCount = ref(0);
  const lastTime = ref<string | null>(null);
  const loadingInstance = ref<any>(null);
  const isInitialLoad = ref(true);

  // 获取节点信息
  const fetchNodeInfo = async () => {
    const data = await get<NodeInfo>(API_ENDPOINTS.NODE);
    nodeInfo.value = data;
  };

  // 获取服务器连接信息
  const fetchConnectors = async () => {
    const data = await get<Connector[]>(API_ENDPOINTS.CONNECTOR);
    connectors.value = data || [];
  };

  // 获取对等节点信息
  const fetchPeers = async () => {
    const data = await get<Peer[]>(API_ENDPOINTS.PEER);
    peers.value = data || [];
    deviceCount.value = (data || []).length;
    lastTime.value = formatCurrentTime();
  };

  // 获取所有数据
  const fetchAllData = async (showLoading = false) => {
    if (showLoading && isInitialLoad.value) {
      loadingInstance.value = ElLoading.service({
        lock: true,
        text: '加载中...',
        background: 'rgba(0, 0, 0, 0.7)',
      });
    }

    try {
      await Promise.all([fetchNodeInfo(), fetchConnectors(), fetchPeers()]);
    } finally {
      if (loadingInstance.value) {
        loadingInstance.value.close();
        loadingInstance.value = null;
      }
      isInitialLoad.value = false;
    }
  };

  // 组件挂载时自动获取数据
  onMounted(() => {
    fetchAllData(true);
  });

  return {
    // 状态
    nodeInfo,
    connectors,
    peers,
    deviceCount,
    lastTime,
    // 方法
    fetchAllData,
    fetchNodeInfo,
    fetchConnectors,
    fetchPeers,
  };
};
