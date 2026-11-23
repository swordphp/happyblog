package library

import (
    "os"
    "path"
    "path/filepath"
)

// LocalUpload 把文件保存到服务器本地磁盘,通过静态路由对外访问
type LocalUpload struct{
    saveDir string // 本地保存目录,例如 ./assets/uploads
    urlPre  string // 对外访问 URL 前缀,例如 /assets/uploads/
}

var LocalUploader = &LocalUpload{saveDir: "./assets/uploads", urlPre: "/assets/uploads/"}

/**
 * 将文件保存到本地磁盘
 *
 * param: string filename 目标文件名(含扩展名)
 * param: []byte data     文件内容
 * return: string         保存后的文件名(相对,不含目录)
 * return: error
 */
func (u LocalUpload) UpLoadFile(filename string, data []byte) (distName string, err error) {
    if err = os.MkdirAll(u.saveDir, 0755); err != nil {
        return "", err
    }
    //只取文件名部分,避免路径穿越
    safeName := filepath.Base(filename)
    fullPath := filepath.Join(u.saveDir, safeName)
    if err = os.WriteFile(fullPath, data, 0644); err != nil {
        return "", err
    }
    return safeName, nil
}

/**
 * 返回本地文件的对外访问 URL
 *
 * param: string distName
 * return: string
 */
func (u LocalUpload) ViewUrl(distName string) string {
    return u.urlPre + path.Base(distName)
}
