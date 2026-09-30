<template>
  <div class="connection-container">
    <div class="et-conninfo-container">
      <div class="et-conninfo-header">
        <div class="et-conninfo-title">连接信息</div>
        <div class="et-conninfo-actions">
          <auto-refresh @refresh="refreshConnInfo" />
          <button class="et-refresh-btn" :class="{ loading: isRefreshing }" :disabled="isRefreshing"
            @click="refreshConnInfo">
            <svg class="et-refresh-icon" width="18" height="18" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2">
              <path d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2" />
            </svg>
            <span>{{ isRefreshing ? '刷新中...' : '刷新' }}</span>
          </button>
        </div>
      </div>

      <div class="et-tabs-container">
        <div class="et-tabs">
          <button v-for="t in tabs" :key="t.id" class="et-tab"
            :class="{ active: currentTab === t.id }" @click="switchTab(t.id)">
            {{ t.label }}
          </button>
        </div>
      </div>

      <div v-for="t in tabs" :key="t.id" class="et-tab-content"
        :class="{ active: currentTab === t.id }" v-html="renderTab(t.id)"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import axios from 'axios';
import AutoRefresh from '../components/AutoRefresh.vue';
import { API_ENDPOINTS } from '../config/api';
import type { ConnectionInfo, ApiResponse } from '../types';

const tabs = [
  { id: 'node', label: '节点信息', key: 'node' },
  { id: 'peer', label: '节点列表', key: 'peer' },
  { id: 'route', label: '路由信息', key: 'route' },
  { id: 'connector', label: '连接器', key: 'connector' },
  { id: 'stun', label: 'STUN', key: 'stun' },
  { id: 'peercenter', label: '节点中心', key: 'peer_center' },
  { id: 'vpnportal', label: 'VPN 门户', key: 'vpn_portal' },
  { id: 'proxy', label: '代理', key: 'proxy' },
  { id: 'acl', label: 'ACL', key: 'acl' },
  { id: 'mappedlistener', label: '映射监听', key: 'mapped_listener' },
  { id: 'stats', label: '统计信息', key: 'stats' },
  { id: 'cmdline', label: '启动命令', key: 'cmdline' },
] as const;

type TabId = typeof tabs[number]['id'];
type InfoKey = typeof tabs[number]['key'];

const data = ref<ConnectionInfo | null>(null);
const currentTab = ref<TabId>('node');
const isRefreshing = ref(false);

const idToKey: Record<TabId, InfoKey> = tabs.reduce((acc, t) => {
  acc[t.id] = t.key;
  return acc;
}, {} as Record<TabId, InfoKey>);

// 中文翻译表（移植官方 translations）
const translations: Record<string, string> = {
  'ipv4': 'IPv4 地址',
  'hostname': '主机名',
  'cost': '路由开销',
  'lat(ms)': '延迟(ms)',
  'loss': '丢包率',
  'rx': '接收流量',
  'tx': '发送流量',
  'tunnel': '隧道协议',
  'NAT': 'NAT 类型',
  'version': '版本',
  'Virtual IP': '虚拟 IP',
  'Hostname': '主机名',
  'Proxy CIDRs': '代理网段',
  'Peer ID': '节点 ID',
  'Public IPv4': '公网 IPv4',
  'Public IPv6': '公网 IPv6',
  'Listener': '监听地址',
  'proxy_cidrs': '代理网段',
  'next_hop_ipv4': '下一跳 IPv4',
  'next_hop_hostname': '下一跳主机名',
  'next_hop_lat': '下一跳延迟',
  'path_len': '路径长度',
  'path_latency': '路径延迟',
  'next_hop_ipv4_lat_first': '首选下一跳 IPv4 延迟',
  'next_hop_hostname_lat_first': '首选下一跳主机名延迟',
  'path_len_lat_first': '首选路径长度',
  'path_latency_lat_first': '首选路径延迟',
  'UDP Stun Type': 'UDP NAT 类型',
  'Interface IPv4': '网卡 IPv4',
  'Interface IPv6': '网卡 IPv6',
  'src': '源',
  'dst': '目标',
  'start_time': '开始时间',
  'state': '状态',
  'transport_type': '传输类型',
  'last_update_time': '最后更新时间',
  'node_id': '节点 ID',
  'direct_peers': '直连节点',
  'Metric Name': '指标名称',
  'Value': '数值',
  'Labels': '标签',
  'status': '状态',
  'client_config_start': '客户端配置起始',
  'client_config_end': '客户端配置结束',
  'Local': '本地',
  'PortRestricted': '端口受限',
  'NoPat': '无 PAT',
  'DIRECT': '直连',
  'unknown': '未知',
  'Connected': '已连接',
  'Connecting': '连接中',
  'Disconnected': '已断开',
  'Symmetric': '对称型',
  'Restricted': '受限',
  'FullCone': '全锥型',
  'AddressRestricted': '地址受限',
  'p2p': '点对点',
  'PublicServer_': '服务器',
  'StunInfo': 'STUN 信息',
  'udp_nat_type': 'UDP NAT 类型',
  'tcp_nat_type': 'TCP NAT 类型',
  'Error': '错误',
  'Invalid service name': '无效服务名',
  'proto name': '协议名',
  'response': '响应',
  'ListConnectorResponse': '连接器列表响应',
  'connectors': '连接器',
  'Connector': '连接器',
  'url': 'URL',
  'Url': 'URL',
  'stun info': 'STUN 信息',
  'public_ip': '公网 IP',
  'min_port': '最小端口',
  'max_port': '最大端口',
  'portal_name': '门户名称',
  'wireguard': 'WireGuard',
  'connected_clients': '已连接客户端',
  'AclStats': 'ACL 统计',
  'Global': '全局',
  'CacheMaxSize': '缓存最大容量',
  'CacheSize': '缓存大小',
  'ConnTrack': '连接跟踪',
  'Rules': '规则',
  'ListMappedListenerResponse': '映射监听列表响应',
  'ERROR: Wireguard VPN Portal Not Started': '错误：WireGuard VPN 门户未启动',
  'no running instances found': '未找到运行中的实例',
  'compression_bytes_rx_before': '压缩前接收字节',
  'compression_bytes_rx_after': '压缩后接收字节',
  'compression_bytes_tx_before': '压缩前发送字节',
  'compression_bytes_tx_after': '压缩后发送字节',
  'peer_rpc_client_rx': '节点 RPC 客户端接收',
  'peer_rpc_client_tx': '节点 RPC 客户端发送',
  'peer_rpc_server_rx': '节点 RPC 服务端接收',
  'peer_rpc_server_tx': '节点 RPC 服务端发送',
  'peer_rpc_duration_ms': '节点 RPC 耗时',
  'traffic_bytes_rx': '流量字节接收',
  'traffic_bytes_tx': '流量字节发送',
  'traffic_bytes_self_rx': '自身流量字节接收',
  'traffic_bytes_self_tx': '自身流量字节发送',
  'traffic_packets_rx': '流量报文接收',
  'traffic_packets_tx': '流量报文发送',
  'traffic_packets_self_rx': '自身流量报文接收',
  'traffic_packets_self_tx': '自身流量报文发送',
  'traffic_packets_forwarded': '转发流量报文',
  'traffic_bytes_forwarded': '转发流量字节',
  'network_name': '网络名称',
  'dst_peer_id': '目标节点 ID',
  'src_peer_id': '源节点 ID',
  'method_name': '方法名',
  'service_name': '服务名',
  'success': '成功',
  'error': '错误',
  'failure': '失败',
  'timeout': '超时',
  'STUN': 'STUN',
  'ACL': 'ACL',
  'OspfRouteRpc': 'OSPF 路由 RPC',
  'PeerCenterRpc': '节点中心 RPC',
  'Caused by': '原因',
  'failed to get peer manager client': '获取节点管理器客户端失败',
  'failed to get peer center client': '获取节点中心客户端失败',
  'failed to connect to server': '连接服务器失败',
  'Connection refused': '连接被拒绝',
  'scheme': '协议',
  'cannot_be_a_base': '不能作为基础地址',
  'username': '用户名',
  'password': '密码',
  'host': '主机',
  'port': '端口',
  'path': '路径',
  'query': '查询参数',
  'fragment': '片段',
  'Some': 'Some',
  'None': 'None',
  'Domain': 'Domain',
  'Failed to connect to RPC server. Please check if the RPC port is set to 15888 in the configuration.': '无法连接到 RPC 服务器，请检查配置中的 RPC 端口是否为 15888。',
  'Configuration File Content': '配置文件内容',
  'sync_route_info': '同步路由信息',
  'get_global_peer_map': '获取全局节点映射',
  'report_peers': '上报节点',
  'traffic_control_bytes_rx': '流量控制字节接收',
  'traffic_control_bytes_tx': '流量控制字节发送',
  'traffic_control_bytes_rx_by_instance': '各实例流量控制字节接收',
  'traffic_control_bytes_tx_by_instance': '各实例流量控制字节发送',
  'traffic_control_packets_rx': '流量控制报文接收',
  'traffic_control_packets_tx': '流量控制报文发送',
  'traffic_control_packets_rx_by_instance': '各实例流量控制报文接收',
  'traffic_control_packets_tx_by_instance': '各实例流量控制报文发送',
  'from_instance_id': '来源实例 ID',
  'to_instance_id': '目标实例 ID',
};

function tr(text: string): string {
  return translations[text] ?? text;
}

function escapeHtml(text: string): string {
  const div = document.createElement('div');
  div.appendChild(document.createTextNode(text));
  return div.innerHTML;
}

function renderCopyablePre(content: string): string {
  const id = 'et_copy_' + Math.random().toString(36).slice(2);
  return '<div class="et-command-panel">' +
    '<div class="et-command-toolbar">' +
    `<button type="button" class="et-copy-btn" data-copy-id="${id}">复制</button>` +
    '</div>' +
    `<div id="${id}" class="et-pre">${escapeHtml(content)}</div>` +
    '</div>';
}

// 事件委托：复制按钮
function handleCopyClick(e: MouseEvent) {
  const btn = (e.target as HTMLElement).closest('.et-copy-btn') as HTMLButtonElement | null;
  if (!btn) return;
  const id = btn.getAttribute('data-copy-id');
  if (!id) return;
  const el = document.getElementById(id);
  if (!el) return;
  const text = el.textContent || '';
  const markCopied = () => {
    const old = btn.textContent;
    btn.textContent = '已复制';
    window.setTimeout(() => { btn.textContent = old; }, 1500);
  };
  if (navigator.clipboard && window.isSecureContext) {
    navigator.clipboard.writeText(text).then(markCopied);
    return;
  }
  const ta = document.createElement('textarea');
  ta.value = text;
  ta.setAttribute('readonly', 'readonly');
  ta.style.position = 'fixed';
  ta.style.left = '-9999px';
  document.body.appendChild(ta);
  ta.select();
  try { document.execCommand('copy'); markCopied(); } finally {
    document.body.removeChild(ta);
  }
}

function parseTable(content: string, tabName: string): string {
  if (!content || content.trim() === '') {
    return '<div class="et-empty">暂无数据</div>';
  }

  if (content.indexOf('错误：') === 0) {
    return `<div class="et-error">${escapeHtml(content)}</div>`;
  }

  if (content.indexOf('ERROR: Wireguard VPN Portal Not Started') !== -1) {
    return `<div class="et-error">${escapeHtml(tr('ERROR: Wireguard VPN Portal Not Started'))}</div>`;
  }

  if (content.indexOf('no running instances found') !== -1) {
    return `<div class="et-error">${escapeHtml(tr('no running instances found'))}</div>`;
  }

  if (content.indexOf('Error: Invalid service name: "PeerCenterRpc"') !== -1) {
    return '<div class="et-error">节点中心服务不可用，该功能需要网络中存在多个节点。</div>';
  }

  if (tabName === 'cmdline') {
    let html = renderCopyablePre(content);
    if (data.value && data.value.config_file) {
      html +=
        '<div class="et-pre" style="margin-top:16px;border-top:2px solid var(--border-color);padding-top:16px;">' +
        `<div style="font-weight:600;margin-bottom:8px;color:var(--tab-active);">${tr('Configuration File Content')}:</div>` +
        `<pre style="margin:0;white-space:pre-wrap;overflow-wrap:anywhere;word-break:break-word;">${escapeHtml(data.value.config_file)}</pre>` +
        '</div>';
    }
    return html;
  }

  if (content.indexOf('failed to get peer manager client') !== -1 ||
    content.indexOf('failed to get peer center client') !== -1 ||
    content.indexOf('Connection refused') !== -1) {
    const translatedContent = content
      .replace(/failed to get peer manager client/g, tr('failed to get peer manager client'))
      .replace(/failed to get peer center client/g, tr('failed to get peer center client'))
      .replace(/failed to connect to server:/g, tr('failed to connect to server') + ':')
      .replace(/Connection refused/g, tr('Connection refused'))
      .replace(/\bscheme:/g, tr('scheme') + ':')
      .replace(/cannot_be_a_base:/g, tr('cannot_be_a_base') + ':')
      .replace(/\busername:/g, tr('username') + ':')
      .replace(/\bpassword:/g, tr('password') + ':')
      .replace(/\bhost:/g, tr('host') + ':')
      .replace(/\bport:/g, tr('port') + ':')
      .replace(/\bpath:/g, tr('path') + ':')
      .replace(/\bquery:/g, tr('query') + ':')
      .replace(/\bfragment:/g, tr('fragment') + ':')
      .replace(/Some\(/g, tr('Some') + '(')
      .replace(/None/g, tr('None'))
      .replace(/Domain\(/g, tr('Domain') + '(');

    return '<div class="et-error">' +
      `<div style="font-weight:600;margin-bottom:8px;">${tr('Failed to connect to RPC server. Please check if the RPC port is set to 15888 in the configuration.')}</div>` +
      `<div style="font-size:12px;opacity:0.9;white-space:pre-wrap;font-family:monospace;line-height:1.6;">${escapeHtml(translatedContent)}</div>` +
      '</div>';
  }

  const lines = content.split(/\r?\n/).filter(l => l.trim() !== '');

  if (lines.length === 0 || lines[0].indexOf('|') === -1) {
    let translatedContent = content;
    for (const key in translations) {
      const escaped = key.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
      const regex = new RegExp('\\b' + escaped + '\\b', 'gi');
      translatedContent = translatedContent.replace(regex, translations[key]);
    }
    translatedContent = translatedContent.replace(/(\w+)=/g, (_m, key) =>
      (translations[key] ? tr(key) : key) + '=');
    translatedContent = translatedContent.replace(/\bListener\s+\d+\b/g, match =>
      match.replace('Listener', tr('Listener')));
    translatedContent = translatedContent.replace(/PublicServer_/g,
      '<span style="background:#e3f2fd;color:#1976d2;padding:2px 4px;border-radius:4px;font-size:11px;font-weight:500;">' + tr('PublicServer_') + '</span>');
    translatedContent = translatedContent.replace(/###############\s*([^#]+)\s*###############/g,
      (_match, p1) => '############### ' + tr(p1.trim()) + ' ###############');

    return '<div class="et-pre">' +
      (translatedContent.indexOf('<span') !== -1 ? translatedContent : escapeHtml(translatedContent)) +
      '</div>';
  }

  const isNodeTable = tabName === 'node';
  const isCmdlineTab = tabName === 'cmdline';

  const extractCells = (line: string): string[] => {
    const cells: string[] = [];
    const matches = line.match(/\|([^|]*)/g);
    if (matches) {
      for (const m of matches) cells.push(m.replace(/^\|/, '').trim());
      if (cells.length > 0 && cells[cells.length - 1] === '') cells.pop();
    }
    return cells;
  };

  const firstCells = extractCells(lines[0]);
  const isTwoColumn = firstCells.length === 2;

  const tableClass = isNodeTable ? 'et-table node-table'
    : isTwoColumn ? 'et-table two-column' : 'et-table';
  let html = `<div class="et-table-container"><table class="${tableClass}">`;
  let rowIndex = 0;

  for (const line of lines) {
    if (/^\|[\s-|]+\|$/.test(line)) continue;
    const cells = extractCells(line);
    if (cells.length === 0) continue;

    html += '<tr>';
    for (let cellContent of cells) {
      const tag = (rowIndex === 0 && !isNodeTable && !isTwoColumn) ? 'th' : 'td';

      if (!isCmdlineTab) {
        if (/^Listener\s+\d+$/.test(cellContent)) {
          cellContent = cellContent.replace('Listener', tr('Listener'));
        } else if (cellContent.indexOf('PublicServer_') !== -1) {
          cellContent = cellContent.replace(/PublicServer_/g,
            '<span style="background:#e3f2fd;color:#1976d2;padding:2px 4px;border-radius:4px;font-size:11px;font-weight:500;">' + tr('PublicServer_') + '</span>');
        } else {
          cellContent = cellContent.replace(/(\w+)=([^,\s]+)/g, (_match, key, value) => {
            const tk = translations[key] ? tr(key) : key;
            const tv = translations[value] ? tr(value) : value;
            return tk + '=' + tv;
          });
          cellContent = tr(cellContent);
        }
      }

      const display = cellContent.trim() === '' ? '-' : cellContent;
      if (display.indexOf('<span') !== -1) {
        html += `<${tag}>${display}</${tag}>`;
      } else {
        html += `<${tag}>${escapeHtml(display)}</${tag}>`;
      }
    }
    html += '</tr>';
    rowIndex++;
  }

  html += '</table></div>';
  return html;
}

function renderTab(tabId: TabId): string {
  if (!data.value) {
    return '<div class="et-loading"><div class="et-loading-spinner"></div>正在采集数据...</div>';
  }
  const key = idToKey[tabId];
  return parseTable(data.value[key] ?? '', tabId);
}

function switchTab(tabName: TabId) {
  if (isRefreshing.value) return;
  currentTab.value = tabName;
}

async function fetchConnection(): Promise<ConnectionInfo | null> {
  try {
    const res = await axios.get<ApiResponse<ConnectionInfo>>(API_ENDPOINTS.CONNECTION, { timeout: 15000 });
    if (res.data.code === 0) return res.data.data;
    return null;
  } catch {
    return null;
  }
}

async function refreshConnInfo() {
  if (isRefreshing.value) return;
  isRefreshing.value = true;
  try {
    const d = await fetchConnection();
    if (d) data.value = d;
  } finally {
    isRefreshing.value = false;
  }
}

onMounted(() => {
  document.addEventListener('click', handleCopyClick);
  // 首次加载；后续定时刷新由 AutoRefresh 开关控制
  refreshConnInfo();
});

onUnmounted(() => {
  document.removeEventListener('click', handleCopyClick);
});
</script>

<style scoped>
.connection-container {
  max-width: 1500px;
  margin: 0 auto;
  padding: 20px;
}
</style>

<style>
.et-conninfo-container {
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid #ecf0f1;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.et-conninfo-header {
  padding: 12px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #ecf0f1;
}

.et-conninfo-title {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
}

.et-conninfo-actions {
  display: flex;
  align-items: center;
  gap: 20px;
}

.et-conninfo-actions .refresh-text {
  font-size: 14px;
  font-weight: 500;
  color: #2c3e50;
}

.et-refresh-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: #3498db;
  color: white;
  border: 1px solid #3498db;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.2s;
}

.et-refresh-btn:hover { opacity: 0.9; }
.et-refresh-btn:active { transform: scale(0.95); }
.et-refresh-btn:disabled { opacity: 0.6; cursor: not-allowed; transform: none; }

.et-refresh-icon { transition: transform 0.6s cubic-bezier(0.4, 0, 0.2, 1); }
.et-refresh-btn.loading .et-refresh-icon { animation: et-spin 1s linear infinite; }

@keyframes et-spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.et-tabs-container {
  border-bottom: 1px solid #ecf0f1;
  overflow-x: auto;
}

.et-tabs-container::-webkit-scrollbar { height: 6px; }
.et-tabs-container::-webkit-scrollbar-track { background: #ecf0f1; border-radius: 3px; }
.et-tabs-container::-webkit-scrollbar-thumb { background: #7f8c8d; border-radius: 3px; }

.et-tabs {
  display: flex;
  min-width: max-content;
}

.et-tab {
  padding: 16px 24px;
  cursor: pointer;
  border: none;
  background: none;
  color: #7f8c8d;
  font-size: 14px;
  font-weight: 600;
  white-space: nowrap;
  border-bottom: 3px solid transparent;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.et-tab:hover {
  color: #3498db;
  background: #f8f9fa;
}

.et-tab.active {
  color: #3498db;
  border-bottom-color: #3498db;
  background: #f8f9fa;
}

.et-tab-content {
  display: none;
  padding: 24px;
  min-height: 200px;
  animation: et-fade-in 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.et-tab-content.active { display: block; }

@keyframes et-fade-in {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.et-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: #7f8c8d;
}

.et-loading-spinner {
  width: 24px;
  height: 24px;
  border: 3px solid #ecf0f1;
  border-top: 3px solid #3498db;
  border-radius: 50%;
  animation: et-spin 1s linear infinite;
  margin-right: 12px;
}

.et-table-container {
  overflow-x: auto;
  border-radius: 12px;
  border: 1px solid #ecf0f1;
  background: #ffffff;
}

.et-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
  border: 2px solid #ecf0f1;
}

.et-table th {
  background: #f0f0f0;
  color: #2c3e50;
  padding: 10px 12px;
  text-align: center;
  font-weight: 600;
  border: 2px solid #ccc;
  white-space: nowrap;
}

.et-table td {
  padding: 8px 12px;
  text-align: center;
  border: 2px solid #ddd;
  color: #2c3e50;
  white-space: nowrap;
  background: #ffffff;
}

.et-table tr:nth-child(2n) td { background: #ffffff; }
.et-table tr:nth-child(2n+1) td { background: #f5f5f5; }
.et-table tr:hover td { background: rgba(52, 152, 219, 0.1); }

.et-table.node-table td:first-child {
  background: #f0f0f0 !important;
  font-weight: 600;
  text-align: left;
  border-right: 3px solid #999;
}

.et-table.node-table td:not(:first-child) { text-align: left; }

.et-table.two-column td:first-child {
  background: #e8e8e8 !important;
  font-weight: 600;
  text-align: left;
  border-right: 3px solid #999 !important;
}

.et-table.two-column tr:nth-child(2n) td:last-child { background: #ffffff !important; }
.et-table.two-column tr:nth-child(2n+1) td:last-child { background: #f5f5f5 !important; }

.et-empty {
  text-align: center;
  color: #7f8c8d;
  padding: 60px 20px;
  font-size: 16px;
  font-style: italic;
}

.et-error {
  color: #e74c3c;
  padding: 20px;
  text-align: center;
  background: rgba(231, 76, 60, 0.1);
  border-radius: 12px;
  border: 1px solid rgba(231, 76, 60, 0.2);
}

.et-pre {
  background: #f8f9fa;
  border: 1px solid #ecf0f1;
  padding: 16px;
  border-radius: 12px;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
  word-break: break-word;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  color: #2c3e50;
  overflow-x: auto;
  line-height: 1.5;
}

.et-command-panel { position: relative; }
.et-command-toolbar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 8px;
}

.et-copy-btn {
  display: inline-flex;
  align-items: center;
  padding: 5px 10px;
  background: #3498db;
  color: #fff;
  border: 1px solid #3498db;
  border-radius: 6px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 500;
  line-height: 1.4;
}

.et-copy-btn:hover { opacity: 0.9; }

/* 暗色模式 */
@media (prefers-color-scheme: dark) {
  .et-conninfo-container {
    background: #2c3e50;
    border-color: #34495e;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  }

  .et-conninfo-header { border-bottom-color: #34495e; }
  .et-conninfo-title { color: #ecf0f1; }
  .et-conninfo-actions .refresh-text { color: #ecf0f1; }
  .et-tabs-container { border-bottom-color: #34495e; }
  .et-tabs-container::-webkit-scrollbar-track { background: #34495e; }

  .et-tab { color: #bdc3c7; }
  .et-tab:hover, .et-tab.active { background: #34495e; }

  .et-table-container { background: #2c3e50; border-color: #34495e; }
  .et-table { border-color: #34495e; }

  .et-table th {
    background: #2a2a2a;
    color: #ecf0f1;
    border-color: #444;
  }

  .et-table td {
    background: transparent;
    border-color: #333;
    color: #ecf0f1;
  }

  .et-table tr:nth-child(2n) td { background: transparent; }
  .et-table tr:nth-child(2n+1) td { background: rgba(255, 255, 255, 0.05); }

  .et-table.node-table td:first-child,
  .et-table.two-column td:first-child {
    background: #2a2a2a !important;
    border-right-color: #444 !important;
  }

  .et-table.two-column tr:nth-child(2n) td:last-child { background: rgba(255, 255, 255, 0.05) !important; }
  .et-table.two-column tr:nth-child(2n+1) td:last-child { background: transparent !important; }

  .et-loading { color: #bdc3c7; }
  .et-loading-spinner { border-color: #34495e; }

  .et-pre {
    background: #34495e;
    border-color: #34495e;
    color: #ecf0f1;
  }
}

@media (max-width: 768px) {
  .et-conninfo-header {
    flex-direction: column;
    gap: 16px;
  }

  .et-tab { padding: 12px 16px; font-size: 13px; }
  .et-tab-content { padding: 16px; }

  .et-table th, .et-table td {
    padding: 8px 12px;
    font-size: 12px;
  }
}
</style>
