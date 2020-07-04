package models

import (
    "fmt"
    "github.com/jinzhu/gorm"
    lib "happyblog/library"
    "log"
    "net/url"
    "strconv"
)


var ConnInstance gorm.DB //被初始化的数据库连接


/**
 * 解析配置文件,并且初始化数据库连接
 */
func init(){
    tmpConf ,_:= lib.ReadWebConfig()
    webConf := *tmpConf
    dbHost := webConf["database.hostname"]
    dbPort,_ := strconv.ParseInt(webConf["database.port"],10,32)
    dbUserName := webConf["database.username"]
    dbPassword := webConf["database.password"]
    loc := webConf["timezone"]
    loc = url.QueryEscape(loc)
    conMsg := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",dbUserName, dbPassword, "tcp", dbHost, dbPort, "happyblog?parseTime=true&loc="+loc)
    db,err := gorm.Open("mysql",conMsg)
    if err != nil {
        log.Fatal(fmt.Sprintf("DB CONNECT ERR ,CONN STR : %s,err:%v",conMsg,err))
    }
    ConnInstance = *db
    maxIdle,_ := strconv.Atoi(webConf["database.maxidel"])
    maxOpen,_ := strconv.Atoi(webConf["database.maxopen"])
    ConnInstance.DB().SetMaxIdleConns(maxIdle)
    ConnInstance.DB().SetMaxOpenConns(maxOpen)
    db.LogMode(true)
}