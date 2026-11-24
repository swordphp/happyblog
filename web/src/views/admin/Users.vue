<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '@/api'
import { fmtDate } from '@/utils/format'

const users = ref([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    users.value = (await api.adminUsers()) || []
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <h2 class="page-title">用户</h2>
  <el-table :data="users" v-loading="loading" border stripe>
    <el-table-column prop="nickName" label="昵称" min-width="140" />
    <el-table-column prop="accountEmail" label="邮箱" min-width="200" />
    <el-table-column label="注册时间" width="180">
      <template #default="{ row }">{{ fmtDate(row.createTime, true) }}</template>
    </el-table-column>
  </el-table>
</template>

<style scoped>
.page-title { margin: 0 0 16px; font-size: 20px; }
</style>
