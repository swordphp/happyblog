<script setup>
import { reactive, ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { marked } from 'marked'
import { ElMessage } from 'element-plus'
import api from '@/api'

const props = defineProps({ id: { type: [String, Number], default: 0 } })
const router = useRouter()
const { t } = useI18n()

const isEdit = computed(() => Number(props.id) > 0)
const loading = ref(false)
const saving = ref(false)
const tagsList = ref([])
const albumsList = ref([])

const f = reactive({
  id: 0,
  title: '',
  uri: '',
  content: '',
  describe: '',
  keywords: '',
  headimage: '',
  tagIds: [],
  albumId: -1,
  pubStatus: 1,
  independPage: 0,
})

const preview = computed(() => marked.parse(f.content || ''))

async function load() {
  loading.value = true
  try {
    if (isEdit.value) {
      const data = await api.adminArticleEdit(props.id)
      const info = data.articleInfo || {}
      f.id = info.id || Number(props.id)
      f.title = info.title || ''
      f.uri = info.uri || ''
      f.content = info.content || ''
      f.describe = info.brief || ''
      f.keywords = info.keywords || ''
      f.headimage = info.headimage || ''
      f.pubStatus = info.pubStatus === 0 ? 0 : 1
      f.independPage = info.independPage === 1 ? 1 : 0
      f.albumId = data.relationAlbumId && data.relationAlbumId > 0 ? data.relationAlbumId : -1
      f.tagIds = (data.relationTagsIds ? String(data.relationTagsIds).split(',') : [])
        .map((x) => Number(x)).filter((x) => x > 0)
      tagsList.value = data.tagsList || []
      albumsList.value = data.albumsList || []
    } else {
      const [tags, albums] = await Promise.all([api.adminTags(), api.adminAlbums()])
      tagsList.value = tags || []
      albumsList.value = albums || []
    }
  } catch (e) {
    ElMessage.error(e.message || t('common.loadFailed'))
  } finally {
    loading.value = false
  }
}

async function save(pubStatus) {
  if (!f.title.trim()) { ElMessage.warning(t('article.needTitle')); return }
  saving.value = true
  try {
    await api.adminArticleSave({
      id: f.id || 0,
      title: f.title,
      uri: f.uri,
      content: f.content,
      describe: f.describe,
      keywords: f.keywords,
      headimage: f.headimage,
      albumId: f.albumId,
      tags: f.tagIds.join(','),
      pubStatus: pubStatus != null ? pubStatus : f.pubStatus,
      independPage: f.independPage,
    })
    ElMessage.success(t('common.saved'))
    router.push({ name: 'admin-articles' })
  } catch (e) {
    ElMessage.error(e.message || t('common.saveFailed'))
  } finally {
    saving.value = false
  }
}

async function onUpload(ev) {
  const file = ev.target.files && ev.target.files[0]
  if (!file) return
  try {
    const url = await api.upload(file, file.name)
    f.content += `\n![${file.name}](${url})\n`
    ElMessage.success(t('article.imageInserted'))
  } catch (e) {
    ElMessage.error(e.message || t('article.uploadFailed'))
  } finally {
    ev.target.value = ''
  }
}

onMounted(load)
</script>
<template>
  <div class="bar">
    <h2 class="page-title">{{ isEdit ? $t('article.editTitle') : $t('article.write') }}</h2>
    <div>
      <el-button @click="router.push({ name: 'admin-articles' })">{{ $t('common.back') }}</el-button>
      <el-button :loading="saving" @click="save(0)">{{ $t('article.saveDraft') }}</el-button>
      <el-button type="primary" :loading="saving" @click="save(1)">{{ $t('article.publish') }}</el-button>
    </div>
  </div>

  <el-form v-loading="loading" label-width="80px" class="editor">
    <el-form-item :label="$t('article.title')">
      <el-input v-model="f.title" :placeholder="$t('article.titlePh')" />
    </el-form-item>
    <el-form-item :label="$t('article.slug')">
      <el-input v-model="f.uri" :placeholder="$t('article.slugPh')" />
    </el-form-item>

    <el-form-item :label="$t('article.body')">
      <div class="md-toolbar">
        <label class="upload-btn">
          {{ $t('article.insertImage') }}
          <input type="file" accept="image/*" hidden @change="onUpload" />
        </label>
        <span class="muted">{{ $t('article.mdHint') }}</span>
      </div>
      <div class="md-split">
        <el-input v-model="f.content" type="textarea" :rows="18" resize="none" :placeholder="$t('article.mdPh')" />
        <div class="md-preview article-content" v-html="preview"></div>
      </div>
    </el-form-item>

    <el-form-item :label="$t('article.brief')">
      <el-input v-model="f.describe" type="textarea" :rows="2" :placeholder="$t('article.briefPh')" />
    </el-form-item>
    <el-form-item :label="$t('article.tags')">
      <el-select v-model="f.tagIds" multiple filterable :placeholder="$t('article.tagsPh')" style="width:100%">
        <el-option v-for="tag in tagsList" :key="tag.id" :label="tag.tagName" :value="tag.id" />
      </el-select>
    </el-form-item>
    <el-form-item :label="$t('article.album')">
      <el-select v-model="f.albumId" style="width:100%">
        <el-option :label="$t('article.noAlbum')" :value="-1" />
        <el-option v-for="a in albumsList" :key="a.id" :label="a.albumName" :value="a.id" />
      </el-select>
    </el-form-item>
    <el-form-item :label="$t('article.keywords')">
      <el-input v-model="f.keywords" :placeholder="$t('article.keywordsPh')" />
    </el-form-item>
    <el-form-item :label="$t('article.headimage')">
      <el-input v-model="f.headimage" :placeholder="$t('article.headimagePh')" />
    </el-form-item>
    <el-form-item :label="$t('article.independPage')">
      <el-switch v-model="f.independPage" :active-value="1" :inactive-value="0" />
      <span class="muted" style="margin-left:10px">{{ $t('article.independHint') }}</span>
    </el-form-item>
  </el-form>
</template>

<style scoped>
.bar { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.page-title { margin: 0; font-size: 20px; }
.editor { background: #fff; padding: 20px 24px; border-radius: 8px; }
.md-toolbar { display: flex; align-items: center; gap: 14px; margin-bottom: 8px; width: 100%; }
.upload-btn { cursor: pointer; color: #0d6efd; font-size: 14px; }
.md-split { display: flex; gap: 16px; width: 100%; }
.md-split > * { flex: 1; min-width: 0; }
.md-preview { border: 1px solid #e4e7ed; border-radius: 4px; padding: 10px 14px; overflow: auto; max-height: 420px; background: #fafafa; }
.muted { color: #999; font-size: 13px; }
</style>