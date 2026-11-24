<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'

const tags = ref([])
const loading = ref(false)
const newTag = ref('')
const adding = ref(false)

async function load() {
  loading.value = true
  try {
    tags.value = (await api.adminTags()) || []
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function add() {
  const name = newTag.value.trim()
  if (!name) { ElMessage.warning('请输入标签名'); return }
  adding.value = true
  try {
    await api.adminTagAdd(name)
    ElMessage.success('已添加')
    newTag.value = ''
    load()
  } catch (e) {
    ElMessage.error(e.message || '添加失败')
  } finally {
    adding.value = false
  }
}

async function remove(row) {
  try {
    await ElMessageBox.confirm(`删除标签「${row.tagName}」?`, '提示', { type: 'warning' })
  } catch (e) { return }
  try {
    await api.adminTagRemove(row.id)
    ElMessage.success('已删除')
    load()
  } catch (e) {
    ElMessage.error(e.message || '删除失败')
  }
}

onMounted(load)
</script>

<template>
  <h2 class="page-title">标签</h2>
  <div class="add-bar">
    <el-input v-model="newTag" placeholder="新标签名称" style="width:240px" @keyup.enter="add" />
    <el-button type="primary" :loading="adding" @click="add">添加</el-button>
  </div>

  <el-table :data="tags" v-loading="loading" border stripe style="max-width:560px">
    <el-table-column prop="id" label="ID" width="80" />
    <el-table-column prop="tagName" label="名称" />
    <el-table-column label="操作" width="100">
      <template #default="{ row }">
        <el-button link type="danger" @click="remove(row)">删除</el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<style scoped>
.page-title { margin: 0 0 16px; font-size: 20px; }
.add-bar { display: flex; gap: 10px; margin-bottom: 16px; }
</style>
