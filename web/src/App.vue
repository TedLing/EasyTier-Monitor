<template>
  <div class="container">
    <peer-list :peers="peers" :device-count="deviceCount" :last-time="lastTime" v-if="peers.length" @refresh="() => fetchAllData(false)" />
    <server-connection :connectors="connectors" v-if="connectors.length" />
    <node-info :node-info="nodeInfo" v-if="nodeInfo" />
  </div>
</template>

<script setup>
import { useMonitorData } from './composables/useMonitorData';
import NodeInfo from './components/NodeInfo.vue';
import ServerConnection from './components/ServerConnection.vue';
import PeerList from './components/PeerList.vue';

// 使用数据管理composable
const { 
  nodeInfo, 
  connectors, 
  peers, 
  deviceCount, 
  lastTime, 
  fetchAllData 
} = useMonitorData();
</script>

<style scoped>
.container {
  max-width: 1500px;
  margin: 0 auto;
  padding: 20px;
}
</style>