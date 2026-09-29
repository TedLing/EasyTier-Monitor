<template>
  <el-collapse v-model="activeNames">
    <el-collapse-item title="当前节点信息" name="1">
      <el-form 
        :model="nodeInfo" 
        label-width="120px" 
        class="grid grid-cols-1 md:grid-cols-2 gap-4"
      >
        <el-form-item label="虚拟IP地址">
          <el-input v-model="nodeInfo.virtual_ip" readonly />
        </el-form-item>
        <el-form-item label="主机名">
          <el-input v-model="nodeInfo.hostname" readonly />
        </el-form-item>
        <el-form-item label="子网代理">
          <el-input v-model="nodeInfo.proxy_cidrs" readonly />
        </el-form-item>
        <el-form-item label="公网IPV4">
          <el-input v-model="nodeInfo.public_ipv4" readonly />
        </el-form-item>
        <el-form-item label="公网IPV6">
          <el-input v-model="nodeInfo.public_ipv6" readonly />
        </el-form-item>
        <!-- <el-form-item label="Nat类型">
          <el-input v-model="nodeInfo.udp_stun_type" readonly />
        </el-form-item> -->
        <el-form-item label="IPv4地址">
          <el-input v-model="nodeInfo.interface_ipv4" readonly />
        </el-form-item>
        <el-form-item label="IPv6地址">
          <el-input v-model="nodeInfo.interface_ipv6" readonly />
        </el-form-item>
        <el-form-item label="监听地址">
          <el-input v-model="nodeInfo.listeners" readonly />
        </el-form-item>
      </el-form>
    </el-collapse-item>
  </el-collapse>
</template>

<script setup>
import { ref } from 'vue';

// 定义组件属性
const props = defineProps({
  nodeInfo: {
    type: Object,
    default: () => ({})
  }
});

// 默认折叠面板
const activeNames = ref([]);
</script>