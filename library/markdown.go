package library

import (
    "github.com/gomarkdown/markdown"
    "github.com/gomarkdown/markdown/html"
    "github.com/gomarkdown/markdown/parser"
    gohtml "html/template"
)

/**
 * 将 markdown 文本渲染为可安全插入模板的 HTML
 * 使用 gomarkdown 解析,并开启常用扩展(表格/删除线/自动链接等)
 *
 * param: string md
 * return: template.HTML
 */
func RenderMarkdown(md string) gohtml.HTML {
    extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Tables | parser.FencedCode
    p := parser.NewWithExtensions(extensions)
    doc := p.Parse([]byte(md))

    opts := html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank}
    renderer := html.NewRenderer(opts)
    out := markdown.Render(doc, renderer)
    return gohtml.HTML(out)
}
