<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { fmtDate } from '@/utils/format'

const router = useRouter()
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
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function remove(row) {
  try {
    await ElMessageBox.confirm(`确认删除文章「${row.title}」?`, '提示', { type: 'warning' })
  } catch (e) { return }
  try {
    await api.adminArticleRemove(row.id)
    ElMessage.success('已删除')
    load(page.value)
  } catch (e) {
    ElMessage.error(e.message || '删除失败')
  }
}

onMounted(() => load(1))
</script>

<template>
  <div class="bar">
    <h2 class="page-title">文章</h2>
    <el-button type="primary" @click="router.push({ name: 'admin-article-new' })">写文章</el-button>
  </div>

  <el-table :data="list" v-loading="loading" border stripe>
    <el-table-column prop="id" label="ID" width="70" />
    <el-table-column prop="title" label="标题" min-width="220" />
    <el-table-column label="状态" width="90">
      <template #default="{ row }">
        <el-tag :type="row.pubStatus === 1 ? 'success' : 'info'" size="small">
          {{ row.pubStatus === 1 ? '已发布' : '草稿' }}
        </el-tag>
      </template>
    </el-table-column>
    <el-table-column label="创建时间" width="170">
      <template #default="{ row }">{{ fmtDate(row.createTime, true) }}</template>
    </el-table-column>
    <el-table-column label="操作" width="160" fixed="right">
      <template #default="{ row }">
        <el-button link type="primary" @click="router.push({ name: 'admin-article-edit', params: { id: row.id } })">编辑</el-button>
        <el-button link type="danger" @click="remove(row)">删除</el-button>
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
