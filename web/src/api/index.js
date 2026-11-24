import axios from 'axios'

// 同源部署:SPA 与接口同一个 host,cookie 会话自动带上。
const http = axios.create({ baseURL: '/', withCredentials: true, timeout: 20000 })

// 把表单对象编码成 application/x-www-form-urlencoded(后端用 PostForm 读取)
function form(obj) {
  const p = new URLSearchParams()
  Object.entries(obj || {}).forEach(([k, v]) => {
    if (v !== undefined && v !== null) p.append(k, v)
  })
  return p
}

// 统一拆包:errNo === 0 返回 data,否则抛出带 errMsg 的错误。
async function call(promise) {
  const res = await promise
  const body = res.data
  if (body && typeof body === 'object' && 'errNo' in body) {
    if (body.errNo !== 0) {
      const err = new Error(body.errMsg || '请求失败')
      err.errNo = body.errNo
      err.body = body
      throw err
    }
    return body.data
  }
  return body
}

export const api = {
  // ---- 安装向导 ----
  installStatus: () => call(http.get('/api/install/status')),
  install: (payload) => call(http.post('/api/install', payload)),

  // ---- 读者端公开接口 ----
  siteInfo: () => call(http.get('/api/siteinfo')),
  pubArticleList: (page = 1) => call(http.get(`/api/article/list/${page}`)),
  pubArticleInfo: (id) => call(http.get(`/api/article/info/${id}`)),
  pubTags: () => call(http.get('/api/tags')),
  pubTagArticles: (id) => call(http.get(`/api/tag/${id}/articles`)),
  pubAlbums: (page = 1) => call(http.get(`/api/albums/${page}`)),
  pubAlbum: (id) => call(http.get(`/api/album/${id}`)),

  // ---- 管理端:会话 ----
  login: (email, password) => call(http.post('/admin/api/login', form({ email, password }))),
  logout: () => call(http.post('/admin/api/logout')),
  me: () => call(http.get('/admin/api/me')),

  // ---- 管理端:文章 ----
  adminArticles: (params) => call(http.get('/admin/api/articles', { params })),
  adminArticleEdit: (id) => call(http.get(`/admin/api/article/${id}`)),
  adminArticleSave: (data) => call(http.post('/admin/api/article/save', form(data))),
  adminArticleRemove: (id) => call(http.post('/admin/api/article/remove', form({ id }))),

  // ---- 管理端:标签 ----
  adminTags: () => call(http.get('/admin/api/tags')),
  adminTagAdd: (tagName) => call(http.post('/admin/api/tag/add', form({ tagName }))),
  adminTagRemove: (id) => call(http.post('/admin/api/tag/remove', form({ id }))),

  // ---- 管理端:专辑 ----
  adminAlbums: () => call(http.get('/admin/api/albums')),
  adminAlbum: (id) => call(http.get(`/admin/api/album/${id}`)),
  adminAlbumSave: (data) => call(http.post('/admin/api/album/save', form(data))),

  // ---- 管理端:用户 / 设置 ----
  adminUsers: () => call(http.get('/admin/api/users')),
  adminSiteConfig: () => call(http.get('/admin/api/siteconfig')),
  adminSiteConfigSave: (data) => call(http.post('/admin/api/siteconfig', form(data))),
  adminSettings: () => call(http.get('/admin/api/settings')),
  adminSetting: (id) => call(http.get('/admin/api/setting', { params: { id } })),
  adminSettingAdd: (data) => call(http.post('/admin/api/setting/add', form(data))),
  adminSettingSave: (data) => call(http.post('/admin/api/setting/save', form(data))),
  adminSettingRemove: (id) => call(http.post('/admin/api/setting/remove', form({ id }))),

  // ---- 上传(返回图片可访问 URL) ----
  async upload(file, filename) {
    const fd = new FormData()
    fd.append('file', file)
    fd.append('filename', filename || file.name)
    const data = await call(http.post('/admin/api/upload', fd))
    return data.viewUrl
  },
}

export default api
