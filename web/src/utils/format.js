// 把后端返回的时间字符串(RFC3339)格式化为 YYYY-MM-DD 或带时分。
export function fmtDate(s, withTime = false) {
  if (!s) return ''
  const d = new Date(s)
  if (isNaN(d.getTime())) return String(s).slice(0, 10)
  const p = (n) => String(n).padStart(2, '0')
  const base = `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())}`
  if (!withTime) return base
  return `${base} ${p(d.getHours())}:${p(d.getMinutes())}`
}
