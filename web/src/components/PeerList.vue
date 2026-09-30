<template>
  <div>
    <div class="header-container">
      <div class="header-left">
        <h3>全部节点信息</h3>
        <span class="device-count">当前连接设备数[ <strong>{{ deviceCount }}</strong> ]个</span>
        <span class="last-date">最后刷新时间[ <strong>{{ lastTime }}</strong> ]</span>
      </div>
      <div class="auto-refresh-container">
        <auto-refresh @refresh="handleRefresh" />
      </div>
    </div>
    <el-table :data="peers" stripe border style="width: 100%">
      <el-table-column prop="id" label="节点ID" align="center" sortable width="150" />
      <el-table-column prop="ipv4" label="虚拟IPv4" align="center" sortable min-width="120">
        <template #default="scope">
          <div class="flex items-center justify-center">
            <span class="mr-2">{{ scope.row.ipv4 }}</span>
            <el-button 
              v-if="scope.row.ipv4" 
              link 
              type="primary" 
              @click="copyData(scope.row.ipv4)"
              size="small"
              style="font-size: 16px; line-height: 1;"
            >
              📋
            </el-button>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="hostname" label="主机名" align="center" sortable min-width="100" />
      <el-table-column prop="cost" label="穿透方式" align="center" sortable width="100" />
      <el-table-column prop="lat_ms" label="延迟(ms)" align="center" sortable width="100" />
      <el-table-column prop="loss_rate" label="丢包率(%)" align="center" sortable width="100" />
      <el-table-column 
        prop="rx_bytes" 
        label="接收(bytes)" 
        align="center" 
        sortable 
        min-width="120"
      />
      <el-table-column 
        prop="tx_bytes" 
        label="发送(bytes)" 
        align="center" 
        sortable 
        min-width="120"
      />
      <el-table-column prop="tunnel_proto" label="协议类型" align="center" sortable width="120" />
      <el-table-column prop="nat_type" label="Nat类型" align="center" sortable width="140">
        <template #default="scope">
          {{ getNatTypeDisplay(scope.row.nat_type) }}
        </template>
      </el-table-column>
      <el-table-column prop="version" label="内核版本" align="center" sortable width="120" />
    </el-table>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { ElMessage } from 'element-plus';
import AutoRefresh from './AutoRefresh.vue';

// NAT类型映射
const natTypes = {
  'Unknown': '未知',
  'OpenInternet': '开放互联网',
  'NoPAT': '无PAT转换',
  'FullCone': '全锥形',
  'Restricted': '限制性',
  'PortRestricted': '端口限制性',
  'Symmetric': '对称型',
  'SymUdpFirewall': '对称型UDP防火墙',
  'SymmetricEasyInc': '对称型（递增端口）',
  'SymmetricEasyDec': '对称型（递减端口）'
};

// 获取NAT类型显示名称
const getNatTypeDisplay = (type) => {
  return `${natTypes[type] || type} (${type})`;
};

// 定义组件属性
const props = defineProps({
  peers: {
    type: Array,
    default: () => []
  },
  deviceCount: {
    type: Number,
    default: 0
  },
  lastTime: {
    type: String,
    default: null
  }
});



const emit = defineEmits(['refresh']);

// 处理刷新事件
const handleRefresh = () => {
  emit('refresh');
};

// 复制IPv4地址到剪贴板
const copyData = async (ipv4) => {
  try {
    await navigator.clipboard.writeText(ipv4);
    ElMessage.success('已复制IPv4地址');
  } catch (error) {
    // 降级方案
    const oInput = document.createElement('input');
    oInput.value = ipv4;
    oInput.style.position = 'fixed';
    oInput.style.opacity = '0';
    document.body.appendChild(oInput);
    oInput.select();
    
    try {
      document.execCommand('Copy');
      ElMessage.success('已复制IPv4地址');
    } catch (fallbackError) {
      ElMessage.error('复制失败，请手动复制');
    } finally {
      document.body.removeChild(oInput);
    }
  }
};
</script>

<style scoped>
.header-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-left h3 {
  margin: 0 24px 0 0;
}

.header-container .auto-refresh-container {
  margin-left: 20px;
}

.device-count {
  font-size: 18px;
  margin-right: 24px;
}

.last-date {
  font-size: 18px;
}
</style>