<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'

const { t } = useI18n()
const tags = ref([])
const loading = ref(false)
const newTag = ref('')
const adding = ref(false)

async function load() {
  loading.value = true
  try {
    tags.value = (await api.adminTags()) || []
  } catch (e) {
    ElMessage.error(e.message || t('common.loadFailed'))
  } finally {
    loading.value = false
  }
}

async function add() {
  const name = newTag.value.trim()
  if (!name) { ElMessage.warning(t('tag.needName')); return }
  adding.value = true
  try {
    await api.adminTagAdd(name)
    ElMessage.success(t('common.saved'))
    newTag.value = ''
    load()
  } catch (e) {
    ElMessage.error(e.message || t('common.saveFailed'))
  } finally {
    adding.value = false
  }
}

async function remove(row) {
  try {
    await ElMessageBox.confirm(t('tag.confirmDel', { name: row.tagName }), t('common.confirmTitle'), { type: 'warning' })
  } catch (e) { return }
  try {
    await api.adminTagRemove(row.id)
    ElMessage.success(t('common.deleted'))
    load()
  } catch (e) {
    ElMessage.error(e.message || t('common.delFailed'))
  }
}

onMounted(load)
</script>

<template>
  <h2 class="page-title">{{ $t('tag.title') }}</h2>
  <div class="add-bar">
    <el-input v-model="newTag" :placeholder="$t('tag.newPh')" style="width:240px" @keyup.enter="add" />
    <el-button type="primary" :loading="adding" @click="add">{{ $t('common.add') }}</el-button>
  </div>

  <el-table :data="tags" v-loading="loading" border stripe style="max-width:560px">
    <el-table-column prop="id" :label="$t('common.id')" width="80" />
    <el-table-column prop="tagName" :label="$t('common.name')" />
    <el-table-column :label="$t('common.operate')" width="100">
      <template #default="{ row }">
        <el-button link type="danger" @click="remove(row)">{{ $t('common.del') }}</el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<style scoped>
.page-title { margin: 0 0 16px; font-size: 20px; }
.add-bar { display: flex; gap: 10px; margin-bottom: 16px; }
</style>
