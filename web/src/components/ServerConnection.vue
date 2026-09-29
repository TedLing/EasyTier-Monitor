<template>
  <el-collapse v-model="activeNames">
    <el-collapse-item title="服务器连接情况" name="1">
      <el-table :data="connectors" stripe border style="width: 100%">
        <el-table-column 
          prop="url" 
          label="服务器地址" 
          align="center" 
          min-width="200"
        />
        <el-table-column 
          prop="status" 
          label="连接情况" 
          align="center" 
          sortable 
          width="150"
        >
          <template #default="scope">
            <el-tag 
              :type="scope.row.status === '连接成功' ? 'success' : 'danger'"
            >
              {{ scope.row.status === '连接成功' ? '已连接' : '未连接' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-collapse-item>
  </el-collapse>
</template>

<script setup>
import { ref } from 'vue';

// 定义组件属性
const props = defineProps({
  connectors: {
    type: Array,
    default: () => []
  }
});

// 默认折叠面板
const activeNames = ref([]);
</script>