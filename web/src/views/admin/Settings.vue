<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'

const { t } = useI18n()
const groups = ref({})
const loading = ref(false)
const dialog = ref(false)
const saving = ref(false)
const form = reactive({ id: 0, ConfigGroup: '', ConfigName: '', ConfigValue: '', ConfigType: 'text', ConfigOrder: 0 })

// 网站基础信息(存在配置文件里,独立于下方的配置项表)
const siteForm = reactive({ sitename: '', sitedesc: '', siteurl: '', github: '', githubshow: 0 })
const siteSaving = ref(false)

async function loadSite() {
  try {
    const d = await api.adminSiteConfig()
    siteForm.sitename = d.sitename || ''
    siteForm.sitedesc = d.sitedesc || ''
    siteForm.siteurl = d.siteurl || ''
    siteForm.github = d.github || ''
    siteForm.githubshow = d.githubshow === '1' || d.githubshow === 1 ? 1 : 0
  } catch (e) { /* 忽略 */ }
}

async function saveSite() {
  if (!siteForm.sitename.trim()) { ElMessage.warning(t('setting.needSiteName')); return }
  siteSaving.value = true
  try {
    await api.adminSiteConfigSave({ ...siteForm })
    ElMessage.success(t('common.saved'))
  } catch (e) {
    ElMessage.error(e.message || t('common.saveFailed'))
  } finally {
    siteSaving.value = false
  }
}

async function load() {
  loading.value = true
  try {
    groups.value = (await api.adminSettings()) || {}
  } catch (e) {
    ElMessage.error(e.message || t('common.loadFailed'))
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
    ElMessage.success(t('common.saved'))
  } catch (e) {
    ElMessage.error(e.message || t('common.saveFailed'))
  }
}

async function removeRow(row) {
  try {
    await ElMessageBox.confirm(t('setting.confirmDel', { name: row.configName }), t('common.confirmTitle'), { type: 'warning' })
  } catch (e) { return }
  try {
    await api.adminSettingRemove(row.id)
    ElMessage.success(t('common.deleted'))
    load()
  } catch (e) {
    ElMessage.error(e.message || t('common.delFailed'))
  }
}

function openAdd() {
  Object.assign(form, { id: 0, ConfigGroup: '', ConfigName: '', ConfigValue: '', ConfigType: 'text', ConfigOrder: 0 })
  dialog.value = true
}

async function submitAdd() {
  if (!form.ConfigGroup.trim() || !form.ConfigName.trim()) {
    ElMessage.warning(t('setting.needGroupName')); return
  }
  saving.value = true
  try {
    await api.adminSettingAdd({ ...form })
    ElMessage.success(t('common.saved'))
    dialog.value = false
    load()
  } catch (e) {
    ElMessage.error(e.message || t('common.saveFailed'))
  } finally {
    saving.value = false
  }
}

onMounted(() => { loadSite(); load() })
</script>
<template>
  <div class="bar">
    <h2 class="page-title">{{ $t('setting.title') }}</h2>
    <el-button type="primary" @click="openAdd">{{ $t('setting.addItem') }}</el-button>
  </div>

  <el-card class="group-card" shadow="never">
    <template #header><b>{{ $t('setting.siteSection') }}</b></template>
    <el-form label-width="90px" style="max-width:620px">
      <el-form-item :label="$t('setting.siteName')">
        <el-input v-model="siteForm.sitename" :placeholder="$t('setting.siteNamePh')" />
      </el-form-item>
      <el-form-item :label="$t('setting.siteDesc')">
        <el-input v-model="siteForm.sitedesc" type="textarea" :rows="2" :placeholder="$t('setting.siteDescPh')" />
      </el-form-item>
      <el-form-item :label="$t('setting.siteUrl')">
        <el-input v-model="siteForm.siteurl" :placeholder="$t('setting.siteUrlPh')" />
      </el-form-item>
      <el-form-item :label="$t('setting.github')">
        <el-input v-model="siteForm.github" :placeholder="$t('setting.githubPh')" />
      </el-form-item>
      <el-form-item :label="$t('setting.githubShow')">
        <el-switch v-model="siteForm.githubshow" :active-value="1" :inactive-value="0" />
        <span class="muted" style="margin-left:10px">{{ $t('setting.githubShowHint') }}</span>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="siteSaving" @click="saveSite">{{ $t('common.save') }}</el-button>
      </el-form-item>
    </el-form>
  </el-card>

  <div v-loading="loading">
    <el-card v-for="(items, groupName) in groups" :key="groupName" class="group-card" shadow="never">
      <template #header><b>{{ groupName || $t('setting.noGroup') }}</b></template>
      <div v-for="row in items" :key="row.id" class="setting-row">
        <span class="name">{{ row.configName }}</span>
        <el-input v-model="row.configValue" class="val" />
        <el-button type="primary" @click="saveRow(row)">{{ $t('common.save') }}</el-button>
        <el-button type="danger" plain @click="removeRow(row)">{{ $t('common.del') }}</el-button>
      </div>
      <p v-if="!items.length" class="muted">{{ $t('setting.emptyGroup') }}</p>
    </el-card>
    <p v-if="!Object.keys(groups).length" class="muted">{{ $t('setting.empty') }}</p>
  </div>

  <el-dialog v-model="dialog" :title="$t('setting.addItem')" width="460px">
    <el-form label-width="70px">
      <el-form-item :label="$t('setting.group')"><el-input v-model="form.ConfigGroup" /></el-form-item>
      <el-form-item :label="$t('setting.cfgName')"><el-input v-model="form.ConfigName" /></el-form-item>
      <el-form-item :label="$t('setting.cfgValue')"><el-input v-model="form.ConfigValue" /></el-form-item>
      <el-form-item :label="$t('setting.cfgType')"><el-input v-model="form.ConfigType" placeholder="text" /></el-form-item>
      <el-form-item :label="$t('setting.cfgOrder')"><el-input v-model.number="form.ConfigOrder" type="number" /></el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialog = false">{{ $t('common.cancel') }}</el-button>
      <el-button type="primary" :loading="saving" @click="submitAdd">{{ $t('common.add') }}</el-button>
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