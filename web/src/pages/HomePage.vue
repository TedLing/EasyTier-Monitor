<template>
  <div class="container">
    <peer-list :peers="peers" :device-count="deviceCount" :last-time="lastTime" v-if="peers.length" @refresh="() => fetchAllData(false)" />
  </div>
</template>

<script setup lang="ts">
import { useMonitorData } from '../composables/useMonitorData';
import PeerList from '../components/PeerList.vue';

const {
  peers, deviceCount, lastTime, fetchAllData
} = useMonitorData();
</script>

<style scoped>
.container {
  max-width: 1500px;
  margin: 0 auto;
  padding: 20px;
}
</style>
