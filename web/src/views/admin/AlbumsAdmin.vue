<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '@/api'

const albums = ref([])
const loading = ref(false)
const dialog = ref(false)
const saving = ref(false)
const form = reactive({ id: 0, albumName: '', isPublic: 1 })

async function load() {
  loading.value = true
  try {
    albums.value = (await api.adminAlbums()) || []
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

function openCreate() {
  Object.assign(form, { id: 0, albumName: '', isPublic: 1 })
  dialog.value = true
}

function openEdit(row) {
  Object.assign(form, { id: row.id, albumName: row.albumName, isPublic: row.isPublic })
  dialog.value = true
}

async function submit() {
  if (!form.albumName.trim()) { ElMessage.warning('请输入专辑名'); return }
  saving.value = true
  try {
    await api.adminAlbumSave({ id: form.id, albumName: form.albumName, isPublic: form.isPublic })
    ElMessage.success('已保存')
    dialog.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="bar">
    <h2 class="page-title">专辑</h2>
    <el-button type="primary" @click="openCreate">新建专辑</el-button>
  </div>

  <el-table :data="albums" v-loading="loading" border stripe>
    <el-table-column prop="id" label="ID" width="80" />
    <el-table-column prop="albumName" label="名称" min-width="180" />
    <el-table-column label="文章数" width="100" prop="articleTotal" />
    <el-table-column label="是否公开" width="100">
      <template #default="{ row }">
        <el-tag :type="row.isPublic === 1 ? 'success' : 'info'" size="small">
          {{ row.isPublic === 1 ? '公开' : '私有' }}
        </el-tag>
      </template>
    </el-table-column>
    <el-table-column label="操作" width="100">
      <template #default="{ row }">
        <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
      </template>
    </el-table-column>
  </el-table>

  <el-dialog v-model="dialog" :title="form.id ? '编辑专辑' : '新建专辑'" width="420px">
    <el-form label-width="80px">
      <el-form-item label="名称">
        <el-input v-model="form.albumName" />
      </el-form-item>
      <el-form-item label="公开">
        <el-switch v-model="form.isPublic" :active-value="1" :inactive-value="0" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialog = false">取消</el-button>
      <el-button type="primary" :loading="saving" @click="submit">保存</el-button>
    </template>
  </el-dialog>
</template>

<style scoped>
.bar { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.page-title { margin: 0; font-size: 20px; }
</style>
