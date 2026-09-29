<template>
  <div class="status-container">
    <!-- EasyTier Core 服务状态 -->
    <div class="et-status-card">
      <div class="et-service-section" :class="{ 'service-stopped': !status?.running }">
        <div class="et-service-header">
          <div class="et-service-title">EasyTier Core</div>
          <div class="et-status-badge" :class="status?.running ? 'et-status-running' : 'et-status-stopped'">
            <span class="et-status-dot"></span>
            <span>{{ status?.running ? '运行中' : '已停止' }}</span>
          </div>
        </div>

        <template v-if="status?.running">
          <div class="et-metrics-grid">
            <div class="et-metric-item">
              <div class="et-metric-icon">
                <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
                  <path d="M12 2l2.4 7.4H22l-6.2 4.5 2.4 7.4L12 16.8 5.8 21.3l2.4-7.4L2 9.4h7.6z"
                    stroke="var(--success)" stroke-width="1.8" stroke-linejoin="round" />
                </svg>
              </div>
              <div class="et-metric-info">
                <div class="et-metric-label">版本</div>
                <div class="et-metric-value">{{ status.version || '-' }}</div>
              </div>
            </div>

            <div class="et-metric-item">
              <div class="et-metric-icon">
                <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
                  <rect x="3" y="6" width="18" height="12" rx="2" stroke="var(--info)" stroke-width="2" />
                  <path d="M7 6V4M12 6V4M17 6V4M7 18v2M12 18v2M17 18v2" stroke="var(--info)"
                    stroke-width="2" stroke-linecap="round" />
                </svg>
              </div>
              <div class="et-metric-info">
                <div class="et-metric-label">主机名</div>
                <div class="et-metric-value" :title="status.hostname">{{ status.hostname || '-' }}</div>
              </div>
            </div>

            <div class="et-metric-item">
              <div class="et-metric-icon">
                <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
                  <circle cx="12" cy="12" r="9" stroke="var(--success)" stroke-width="2" />
                  <path d="M3 12h18M12 3c3 3 3 15 0 18M12 3c-3 3-3 15 0 18" stroke="var(--success)" stroke-width="1.5" />
                </svg>
              </div>
              <div class="et-metric-info">
                <div class="et-metric-label">虚拟网络 IPv4</div>
                <div class="et-metric-value">{{ status.virtual_ip || '-' }}</div>
              </div>
            </div>

            <div class="et-metric-item">
              <div class="et-metric-icon">
                <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
                  <circle cx="9" cy="12" r="3" stroke="var(--info)" stroke-width="2" />
                  <circle cx="17" cy="7" r="2.5" stroke="var(--success)" stroke-width="2" />
                  <circle cx="17" cy="17" r="2.5" stroke="var(--warning)" stroke-width="2" />
                  <path d="M11.5 10.8l3-2.6M11.5 13.2l3 2.6" stroke="var(--text-secondary)" stroke-width="1.5" />
                </svg>
              </div>
              <div class="et-metric-info">
                <div class="et-metric-label">在线节点</div>
                <div class="et-metric-value">{{ status.peer_count }}</div>
              </div>
            </div>

            <div class="et-metric-item">
              <div class="et-metric-icon">
                <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
                  <path d="M4 17l6-6 4 4 6-8" stroke="var(--info)" stroke-width="2"
                    stroke-linecap="round" stroke-linejoin="round" />
                  <path d="M15 7h5v5" stroke="var(--info)" stroke-width="2"
                    stroke-linecap="round" stroke-linejoin="round" />
                </svg>
              </div>
              <div class="et-metric-info">
                <div class="et-metric-label">已连服务器</div>
                <div class="et-metric-value">{{ status.connector }}</div>
              </div>
            </div>

            <div class="et-metric-item">
              <div class="et-metric-icon">
                <svg width="36" height="36" viewBox="0 0 24 24" fill="none">
                  <path d="M4 12a8 8 0 0116 0M7 12a5 5 0 0110 0" stroke="var(--warning)" stroke-width="2"
                    stroke-linecap="round" />
                  <circle cx="12" cy="12" r="1.8" fill="var(--warning)" />
                </svg>
              </div>
              <div class="et-metric-info">
                <div class="et-metric-label">代理网段</div>
                <div class="et-metric-value" :title="status.proxy_cidrs.join(', ')">
                  {{ status.proxy_cidrs.length ? status.proxy_cidrs.join(', ') : '-' }}
                </div>
              </div>
            </div>
          </div>

          <div class="et-listeners" v-if="status.listeners.length">
            <div class="et-listeners-label">监听地址</div>
            <div class="et-listener-tags">
              <span class="et-listener-tag" v-for="(l, i) in status.listeners" :key="i" :title="l">{{ l }}</span>
            </div>
          </div>
        </template>
      </div>
    </div>

    <!-- 实时流量 -->
    <div class="et-status-card" v-if="status?.running">
      <div class="et-service-header">
        <div class="et-service-title">实时流量（所有节点汇总）</div>
        <div class="et-speed-legend">
          <div class="legend-item">
            <div class="legend-box dl"></div>
            <span class="legend-label">下载</span>
            <span class="legend-speed dl">{{ rxSpeed }}/s</span>
          </div>
          <div class="legend-item">
            <div class="legend-box ul"></div>
            <span class="legend-label">上传</span>
            <span class="legend-speed ul">{{ txSpeed }}/s</span>
          </div>
        </div>
      </div>
      <canvas ref="chartCanvas" class="traffic-chart"></canvas>
      <div class="total-row">
        <div class="total-item">
          <div class="total-label">累计下载</div>
          <div class="total-value dl">{{ formatBytes(status.total_rx) }}</div>
        </div>
        <div class="total-item">
          <div class="total-label">累计上传</div>
          <div class="total-value ul">{{ formatBytes(status.total_tx) }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue';
import axios from 'axios';
import { API_ENDPOINTS } from '../config/api';
import type { Status, ApiResponse } from '../types';

const POLL_INTERVAL = 3000;

const status = ref<Status | null>(null);
const rxSpeed = ref('0 B');
const txSpeed = ref('0 B');
const chartCanvas = ref<HTMLCanvasElement | null>(null);

let timer: number | undefined;
let lastRx = 0;
let lastTx = 0;
let hasLast = false;
const labels: string[] = [];
const rxData: number[] = [];
const txData: number[] = [];

function formatBytes(bytes: number): string {
  if (!bytes || bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.min(Math.floor(Math.log(bytes) / Math.log(k)), sizes.length - 1);
  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`;
}

// 静默获取（轮询失败不弹错误提示）
async function fetchStatus(): Promise<Status | null> {
  try {
    const res = await axios.get<ApiResponse<Status>>(API_ENDPOINTS.STATUS, { timeout: 8000 });
    if (res.data.code === 0) return res.data.data;
    return null;
  } catch {
    return null;
  }
}

async function tick() {
  const data = await fetchStatus();
  if (!data) return;
  status.value = data;

  if (data.running) {
    let curRxSpeed = 0;
    let curTxSpeed = 0;
    if (hasLast) {
      // 计数器回绕或 core 重启时差值可能为负，归零处理
      if (data.total_rx >= lastRx) curRxSpeed = (data.total_rx - lastRx) / (POLL_INTERVAL / 1000);
      if (data.total_tx >= lastTx) curTxSpeed = (data.total_tx - lastTx) / (POLL_INTERVAL / 1000);
    }
    lastRx = data.total_rx;
    lastTx = data.total_tx;
    hasLast = true;

    rxSpeed.value = formatBytes(curRxSpeed);
    txSpeed.value = formatBytes(curTxSpeed);

    const now = new Date();
    const timeStr = `${now.getHours()}:${String(now.getMinutes()).padStart(2, '0')}:${String(now.getSeconds()).padStart(2, '0')}`;
    labels.push(timeStr);
    rxData.push(curRxSpeed);
    txData.push(curTxSpeed);
    if (labels.length > 20) {
      labels.shift();
      rxData.shift();
      txData.shift();
    }

    await nextTick();
    drawChart();
  } else {
    hasLast = false;
  }
}

function drawChart() {
  const canvas = chartCanvas.value;
  if (!canvas) return;
  const ctx = canvas.getContext('2d');
  if (!ctx) return;

  canvas.width = canvas.offsetWidth;
  canvas.height = 200;

  const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  if (rxData.length === 0) return;

  const paddingLeft = 70;
  const paddingRight = 20;
  const paddingTop = 20;
  const paddingBottom = 30;
  const width = canvas.width - paddingLeft - paddingRight;
  const height = canvas.height - paddingTop - paddingBottom;

  let maxValue = Math.max(...rxData, ...txData);
  if (maxValue === 0) maxValue = 100;
  maxValue *= 1.2;

  // 网格线
  ctx.strokeStyle = isDark ? '#3e3e3e' : '#e8e8e8';
  ctx.lineWidth = 1;
  ctx.beginPath();
  for (let i = 0; i <= 5; i++) {
    const y = paddingTop + (height / 5) * i;
    ctx.moveTo(paddingLeft, y);
    ctx.lineTo(paddingLeft + width, y);
  }
  ctx.stroke();

  // Y 轴刻度
  ctx.fillStyle = isDark ? '#b0b0b0' : '#666';
  ctx.font = '11px sans-serif';
  ctx.textAlign = 'right';
  ctx.textBaseline = 'middle';
  for (let i = 0; i <= 5; i++) {
    const y = paddingTop + (height / 5) * i;
    const value = maxValue * (1 - i / 5);
    ctx.fillText(`${formatBytes(value)}/s`, paddingLeft - 10, y);
  }

  function drawLine(dataArray: number[], color: string) {
    if (dataArray.length < 2) return;

    const gradient = ctx!.createLinearGradient(0, paddingTop, 0, paddingTop + height);
    gradient.addColorStop(0, `${color}40`);
    gradient.addColorStop(1, `${color}00`);

    ctx!.beginPath();
    ctx!.moveTo(paddingLeft, paddingTop + height);
    for (let i = 0; i < dataArray.length; i++) {
      const x = paddingLeft + (width / (dataArray.length - 1)) * i;
      const y = paddingTop + height - (dataArray[i] / maxValue) * height;
      ctx!.lineTo(x, y);
    }
    ctx!.lineTo(paddingLeft + width, paddingTop + height);
    ctx!.closePath();
    ctx!.fillStyle = gradient;
    ctx!.fill();

    ctx!.strokeStyle = color;
    ctx!.lineWidth = 2.5;
    ctx!.beginPath();
    for (let i = 0; i < dataArray.length; i++) {
      const x = paddingLeft + (width / (dataArray.length - 1)) * i;
      const y = paddingTop + height - (dataArray[i] / maxValue) * height;
      if (i === 0) ctx!.moveTo(x, y);
      else ctx!.lineTo(x, y);
    }
    ctx!.stroke();

    // 末端圆点
    const lastX = paddingLeft + width;
    const lastY = paddingTop + height - (dataArray[dataArray.length - 1] / maxValue) * height;
    ctx!.beginPath();
    ctx!.arc(lastX, lastY, 4, 0, 2 * Math.PI);
    ctx!.fillStyle = color;
    ctx!.fill();
  }

  drawLine(rxData, isDark ? '#66BB6A' : '#4CAF50');
  drawLine(txData, isDark ? '#42A5F5' : '#2196F3');

  // X 轴时间标签
  ctx.fillStyle = isDark ? '#b0b0b0' : '#666';
  ctx.font = '10px sans-serif';
  ctx.textAlign = 'center';
  const isMobile = canvas.width < 600;
  const maxLabels = isMobile ? 4 : 6;
  const step = Math.max(1, Math.floor(labels.length / maxLabels));
  for (let i = 0; i < labels.length; i += step) {
    const x = paddingLeft + (width / (labels.length - 1)) * i;
    ctx.fillText(labels[i], x, paddingTop + height + 20);
  }
}

onMounted(() => {
  tick();
  timer = window.setInterval(tick, POLL_INTERVAL);
});

onUnmounted(() => {
  if (timer) window.clearInterval(timer);
});
</script>

<style scoped>
.status-container {
  max-width: 1500px;
  margin: 0 auto;
  padding: 20px;
}

.et-status-card {
  background: var(--card-bg);
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid #e0e0e0;
  margin: 16px 0;
  color: #333;
}

:root {
  --card-bg: #ffffff;
}

@media (prefers-color-scheme: dark) {
  .et-status-card {
    background: #1e1e1e;
    border-color: #3e3e3e;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
    color: #e0e0e0;
  }
}

.et-service-section {
  padding: 0;
}

.et-service-section.service-stopped {
  padding: 0;
}

.et-service-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 12px;
}

.et-service-title {
  font-size: 15px;
  font-weight: 600;
}

.et-status-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.et-status-running {
  background: rgba(76, 175, 80, 0.1);
  color: #4caf50;
}

.et-status-stopped {
  background: rgba(244, 67, 54, 0.1);
  color: #f44336;
}

.et-status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  animation: et-pulse 2s infinite;
}

.et-status-running .et-status-dot {
  background: #4caf50;
}

.et-status-stopped .et-status-dot {
  background: #f44336;
}

@keyframes et-pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.et-metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
  margin-bottom: 12px;
}

.et-metric-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: transparent;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
}

@media (prefers-color-scheme: dark) {
  .et-metric-item {
    border-color: #3e3e3e;
  }
}

.et-metric-icon {
  flex-shrink: 0;
}

.et-metric-info {
  flex: 1;
  min-width: 0;
}

.et-metric-label {
  font-size: 11px;
  color: #666;
  margin-bottom: 2px;
}

@media (prefers-color-scheme: dark) {
  .et-metric-label {
    color: #b0b0b0;
  }
}

.et-metric-value {
  font-size: 13px;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.et-listeners {
  margin-top: 8px;
}

.et-listeners-label {
  font-size: 12px;
  color: #666;
  margin-bottom: 8px;
}

@media (prefers-color-scheme: dark) {
  .et-listeners-label {
    color: #b0b0b0;
  }
}

.et-listener-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.et-listener-tag {
  display: inline-block;
  padding: 4px 10px;
  font-size: 12px;
  background: rgba(33, 150, 243, 0.08);
  color: #2196f3;
  border-radius: 6px;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.et-speed-legend {
  display: flex;
  gap: 16px;
  font-size: 12px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.legend-box {
  width: 12px;
  height: 12px;
  border-radius: 2px;
}

.legend-box.dl {
  background: #4caf50;
}

.legend-box.ul {
  background: #2196f3;
}

.legend-label {
  color: #666;
}

.legend-speed {
  font-weight: 600;
}

.legend-speed.dl {
  color: #4caf50;
}

.legend-speed.ul {
  color: #2196f3;
}

.traffic-chart {
  width: 100%;
  height: 200px;
  display: block;
}

.total-row {
  display: flex;
  justify-content: space-around;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #e0e0e0;
}

@media (prefers-color-scheme: dark) {
  .total-row {
    border-color: #3e3e3e;
  }
}

.total-item {
  text-align: center;
}

.total-label {
  font-size: 11px;
  color: #666;
  margin-bottom: 4px;
}

@media (prefers-color-scheme: dark) {
  .total-label {
    color: #b0b0b0;
  }
}

.total-value {
  font-size: 14px;
  font-weight: 600;
}

.total-value.dl {
  color: #4caf50;
}

.total-value.ul {
  color: #2196f3;
}

@media (max-width: 768px) {
  .et-status-card {
    padding: 16px;
  }

  .et-metrics-grid {
    grid-template-columns: 1fr;
  }

  .et-speed-legend {
    width: 100%;
    justify-content: flex-start;
  }
}
</style>
