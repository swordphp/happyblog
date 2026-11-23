package library

import "strings"

/**
 * 按字符(rune)安全截取字符串的前 n 个字符
 * 避免对多字节(如中文)内容按字节切割导致乱码或越界 panic
 *
 * param: string s
 * param: int    n
 * return: string
 */
func SubStr(s string, n int) string {
    if n <= 0 {
        return ""
    }
    r := []rune(s)
    if len(r) <= n {
        return s
    }
    return string(r[:n])
}

/**
 * 生成文章简介:去除 markdown/html 噪声后按字符截取
 *
 * param: string content
 * param: int    n
 * return: string
 */
func MakeBrief(content string, n int) string {
    cleaned := strings.NewReplacer(
        "\r", " ",
        "\n", " ",
        "#", "",
        "*", "",
        "`", "",
        ">", "",
    ).Replace(content)
    cleaned = strings.TrimSpace(cleaned)
    return SubStr(cleaned, n)
}
