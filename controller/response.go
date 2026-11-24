package controller

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// 统一的接口返回体与错误码。所有 JSON 接口都通过 respOK / respFail 返回,
// 保证 {errNo, errMsg, data} 结构一致,避免各处手写 gin.H 带来的拼写与字段不一致。
const (
    errCodeOK      = 0
    errCodeParam   = 10  // 参数错误
    errCodeNotFound = 404 // 资源不存在
    errCodeServer  = 500 // 服务端错误
)

// respOK 返回成功响应:{errNo:0, errMsg:"success", data:...}
func respOK(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, gin.H{"errNo": errCodeOK, "errMsg": "success", "data": data})
}

// respFail 返回业务错误响应:{errNo, errMsg, data:null}
func respFail(c *gin.Context, errNo int, errMsg string) {
    c.JSON(http.StatusOK, gin.H{"errNo": errNo, "errMsg": errMsg, "data": nil})
}
