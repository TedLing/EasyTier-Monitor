<template>
  <div class="flex items-center">
    <el-switch 
      v-model="autoRefresh" 
      inline-prompt 
      active-text="开" 
      inactive-text="关" 
    />
    <span class="ml-2 refresh-text">自动刷新（每10秒）</span>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue';

const emit = defineEmits(['refresh']);
const autoRefresh = ref(true);
let refreshTimer = null;

const REFRESH_INTERVAL = 10000; // 10秒刷新间隔

const startRefresh = () => {
  if (refreshTimer) return;
  
  // 立即刷新一次，然后再开始定时刷新
  emit('refresh');
  refreshTimer = window.setInterval(() => {
    emit('refresh');
  }, REFRESH_INTERVAL);
};

const stopRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer);
    refreshTimer = null;
  }
};

// 监听自动刷新状态变化
watch(autoRefresh, (newVal) => {
  if (newVal) {
    startRefresh();
  } else {
    stopRefresh();
  }
});

// 页面可见性变化处理
const handleVisibilityChange = () => {
  if (document.visibilityState === 'hidden') {
    stopRefresh();
  } else if (document.visibilityState === 'visible' && autoRefresh.value) {
    startRefresh();
  }
};

// 窗口焦点变化处理
const handleWindowBlur = () => {
  stopRefresh();
};

const handleWindowFocus = () => {
  if (autoRefresh.value) {
    startRefresh();
  }
};

onMounted(() => {
  // 添加事件监听
  document.addEventListener('visibilitychange', handleVisibilityChange);
  window.addEventListener('blur', handleWindowBlur);
  window.addEventListener('focus', handleWindowFocus);
});

onUnmounted(() => {
  // 清理资源和事件监听
  stopRefresh();
  document.removeEventListener('visibilitychange', handleVisibilityChange);
  window.removeEventListener('blur', handleWindowBlur);
  window.removeEventListener('focus', handleWindowFocus);
});
</script>

<style scoped>
.refresh-text {
  font-size: 18px;
  font-weight: bold;
}
</style>