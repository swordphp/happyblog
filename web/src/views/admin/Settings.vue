<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'

const groups = ref({})   // { groupName: [setting, ...] }
const loading = ref(false)
const dialog = ref(false)
const saving = ref(false)
const form = reactive({ id: 0, ConfigGroup: '', ConfigName: '', ConfigValue: '', ConfigType: 'text', ConfigOrder: 0 })

async function load() {
  loading.value = true
  try {
    groups.value = (await api.adminSettings()) || {}
  } catch (e) {
    ElMessage.error(e.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function saveRow(row) {
  try {
    await api.adminSettingSave({
      id: row.id,
      ConfigGroup: row.group,
      ConfigName: row.configName,
      ConfigValue: row.configValue,
      ConfigType: row.type,
      ConfigOrder: row.order,
    })
    ElMessage.success('已保存')
  } catch (e) {
    ElMessage.error(e.message || '保存失败')
  }
}

async function removeRow(row) {
  try {
    await ElMessageBox.confirm(`删除配置「${row.configName}」?`, '提示', { type: 'warning' })
  } catch (e) { return }
  try {
    await api.adminSettingRemove(row.id)
    ElMessage.success('已删除')
    load()
  } catch (e) {
    ElMessage.error(e.message || '删除失败')
  }
}

function openAdd() {
  Object.assign(form, { id: 0, ConfigGroup: '', ConfigName: '', ConfigValue: '', ConfigType: 'text', ConfigOrder: 0 })
  dialog.value = true
}

async function submitAdd() {
  if (!form.ConfigGroup.trim() || !form.ConfigName.trim()) {
    ElMessage.warning('分组和名称不能为空'); return
  }
  saving.value = true
  try {
    await api.adminSettingAdd({ ...form })
    ElMessage.success('已添加')
    dialog.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message || '添加失败')
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>
<template>
  <div class="bar">
    <h2 class="page-title">设置</h2>
    <el-button type="primary" @click="openAdd">新增配置项</el-button>
  </div>

  <div v-loading="loading">
    <el-card v-for="(items, groupName) in groups" :key="groupName" class="group-card" shadow="never">
      <template #header><b>{{ groupName || '未分组' }}</b></template>
      <div v-for="row in items" :key="row.id" class="setting-row">
        <span class="name">{{ row.configName }}</span>
        <el-input v-model="row.configValue" class="val" />
        <el-button type="primary" @click="saveRow(row)">保存</el-button>
        <el-button type="danger" plain @click="removeRow(row)">删除</el-button>
      </div>
      <p v-if="!items.length" class="muted">该分组暂无配置项。</p>
    </el-card>
    <p v-if="!Object.keys(groups).length" class="muted">还没有任何配置项。</p>
  </div>

  <el-dialog v-model="dialog" title="新增配置项" width="460px">
    <el-form label-width="70px">
      <el-form-item label="分组"><el-input v-model="form.ConfigGroup" /></el-form-item>
      <el-form-item label="名称"><el-input v-model="form.ConfigName" /></el-form-item>
      <el-form-item label="值"><el-input v-model="form.ConfigValue" /></el-form-item>
      <el-form-item label="类型"><el-input v-model="form.ConfigType" placeholder="text" /></el-form-item>
      <el-form-item label="排序"><el-input v-model.number="form.ConfigOrder" type="number" /></el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialog = false">取消</el-button>
      <el-button type="primary" :loading="saving" @click="submitAdd">添加</el-button>
    </template>
  </el-dialog>
</template>

<style scoped>
.bar { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.page-title { margin: 0; font-size: 20px; }
.group-card { margin-bottom: 16px; }
.setting-row { display: flex; align-items: center; gap: 10px; margin-bottom: 10px; }
.setting-row .name { width: 160px; color: #555; flex-shrink: 0; }
.setting-row .val { flex: 1; }
.muted { color: #999; font-size: 13px; }
</style>

