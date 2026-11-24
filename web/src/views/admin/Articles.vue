<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { fmtDate } from '@/utils/format'

const router = useRouter()
const { t } = useI18n()
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const loading = ref(false)

async function load(p = 1) {
  loading.value = true
  try {
    const data = await api.adminArticles({ page: p })
    list.value = (data && data.list) || []
    total.value = (data && data.total) || 0
    page.value = p
  } catch (e) {
    ElMessage.error(e.message || t('common.loadFailed'))
  } finally {
    loading.value = false
  }
}

async function remove(row) {
  try {
    await ElMessageBox.confirm(t('article.confirmDel', { title: row.title }), t('common.confirmTitle'), { type: 'warning' })
  } catch (e) { return }
  try {
    await api.adminArticleRemove(row.id)
    ElMessage.success(t('common.deleted'))
    load(page.value)
  } catch (e) {
    ElMessage.error(e.message || t('common.delFailed'))
  }
}

onMounted(() => load(1))
</script>

<template>
  <div class="bar">
    <h2 class="page-title">{{ $t('article.listTitle') }}</h2>
    <el-button type="primary" @click="router.push({ name: 'admin-article-new' })">{{ $t('article.write') }}</el-button>
  </div>

  <el-table :data="list" v-loading="loading" border stripe>
    <el-table-column prop="id" :label="$t('common.id')" width="70" />
    <el-table-column prop="title" :label="$t('article.title')" min-width="220" />
    <el-table-column :label="$t('article.status')" width="90">
      <template #default="{ row }">
        <el-tag :type="row.pubStatus === 1 ? 'success' : 'info'" size="small">
          {{ row.pubStatus === 1 ? $t('article.published') : $t('article.draft') }}
        </el-tag>
      </template>
    </el-table-column>
    <el-table-column :label="$t('article.createTime')" width="170">
      <template #default="{ row }">{{ fmtDate(row.createTime, true) }}</template>
    </el-table-column>
    <el-table-column :label="$t('common.operate')" width="160" fixed="right">
      <template #default="{ row }">
        <el-button link type="primary" @click="router.push({ name: 'admin-article-edit', params: { id: row.id } })">{{ $t('common.edit') }}</el-button>
        <el-button link type="danger" @click="remove(row)">{{ $t('common.del') }}</el-button>
      </template>
    </el-table-column>
  </el-table>

  <el-pagination
    v-if="total > pageSize"
    class="pager"
    layout="prev, pager, next"
    :total="total"
    :page-size="pageSize"
    :current-page="page"
    @current-change="load"
  />
</template>

<style scoped>
.bar { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.page-title { margin: 0; font-size: 20px; }
.pager { margin-top: 16px; justify-content: flex-end; }
</style>
