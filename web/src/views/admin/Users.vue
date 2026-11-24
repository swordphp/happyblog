<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import api from '@/api'
import { fmtDate } from '@/utils/format'

const { t } = useI18n()
const users = ref([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    users.value = (await api.adminUsers()) || []
  } catch (e) {
    ElMessage.error(e.message || t('common.loadFailed'))
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <h2 class="page-title">{{ $t('user.title') }}</h2>
  <el-table :data="users" v-loading="loading" border stripe>
    <el-table-column prop="nickName" :label="$t('user.nickname')" min-width="140" />
    <el-table-column prop="accountEmail" :label="$t('user.email')" min-width="200" />
    <el-table-column :label="$t('user.registerTime')" width="180">
      <template #default="{ row }">{{ fmtDate(row.createTime, true) }}</template>
    </el-table-column>
  </el-table>
</template>

<style scoped>
.page-title { margin: 0 0 16px; font-size: 20px; }
</style>
