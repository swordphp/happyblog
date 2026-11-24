<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import api from '@/api'

const { t } = useI18n()
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
    ElMessage.error(e.message || t('common.loadFailed'))
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
  if (!form.albumName.trim()) { ElMessage.warning(t('album.needName')); return }
  saving.value = true
  try {
    await api.adminAlbumSave({ id: form.id, albumName: form.albumName, isPublic: form.isPublic })
    ElMessage.success(t('common.saved'))
    dialog.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message || t('common.saveFailed'))
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="bar">
    <h2 class="page-title">{{ $t('album.listTitle') }}</h2>
    <el-button type="primary" @click="openCreate">{{ $t('album.create') }}</el-button>
  </div>

  <el-table :data="albums" v-loading="loading" border stripe>
    <el-table-column prop="id" :label="$t('common.id')" width="80" />
    <el-table-column prop="albumName" :label="$t('common.name')" min-width="180" />
    <el-table-column :label="$t('album.count')" width="100" prop="articleTotal" />
    <el-table-column :label="$t('album.isPublic')" width="100">
      <template #default="{ row }">
        <el-tag :type="row.isPublic === 1 ? 'success' : 'info'" size="small">
          {{ row.isPublic === 1 ? $t('common.public') : $t('common.private') }}
        </el-tag>
      </template>
    </el-table-column>
    <el-table-column :label="$t('common.operate')" width="100">
      <template #default="{ row }">
        <el-button link type="primary" @click="openEdit(row)">{{ $t('common.edit') }}</el-button>
      </template>
    </el-table-column>
  </el-table>

  <el-dialog v-model="dialog" :title="form.id ? $t('album.editTitle') : $t('album.createTitle')" width="420px">
    <el-form label-width="80px">
      <el-form-item :label="$t('common.name')">
        <el-input v-model="form.albumName" />
      </el-form-item>
      <el-form-item :label="$t('common.public')">
        <el-switch v-model="form.isPublic" :active-value="1" :inactive-value="0" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialog = false">{{ $t('common.cancel') }}</el-button>
      <el-button type="primary" :loading="saving" @click="submit">{{ $t('common.save') }}</el-button>
    </template>
  </el-dialog>
</template>

<style scoped>
.bar { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.page-title { margin: 0; font-size: 20px; }
</style>
