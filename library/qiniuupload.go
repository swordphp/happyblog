package library

import (
	"bytes"
	"context"
	"errors"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"time"
)

type QiNiuUpload struct{
	bucket string
	ak string
	sk string
	upToken string
	tokenMakeTime time.Time
	zone *storage.Region
	https bool
	useCdnDomain bool
	pathTo string
}

var QiNiuUploader = new(QiNiuUpload)

func init(){
	tmpConf ,_ := ReadWebConfig()
	WebConf := *tmpConf
	QiNiuUploader.ak = WebConf["qiniuoss.ak"]
	QiNiuUploader.sk = WebConf["qiniuoss.sk"]
	QiNiuUploader.bucket = WebConf["qiniuoss.bucket"]
	zoneInfo,_ := getZone(WebConf["qiniuoss.zone"])
	QiNiuUploader.zone = zoneInfo
	if WebConf["qiniuoss.https"] == "1"{
		QiNiuUploader.https = true
	} else {
		QiNiuUploader.https = false
	}

	if WebConf["qiniuoss.usecdndomain"] == "1" {
		QiNiuUploader.useCdnDomain = true
	} else {
		QiNiuUploader.useCdnDomain = false
	}
	QiNiuUploader.pathTo = WebConf["qiniuoss.pathto"]
}


/**
 * 用户获取上传用的票据
 *
 * return: string
 */
func (uploadObj  QiNiuUpload)  getToken () (upToken string){
	duration := (time.Now().Unix() - uploadObj.tokenMakeTime.Unix())/1000
	if uploadObj.upToken != "" && duration < 7200 {
		upToken = uploadObj.upToken
	} else {
		bucket := uploadObj.bucket
		putPolicy := storage.PutPolicy{
			Scope: bucket,
		}
		putPolicy.Expires = 7200//上传票据两小时有效
		mac := qbox.NewMac(uploadObj.ak, uploadObj.sk)
		upToken = putPolicy.UploadToken(mac)
		uploadObj.upToken = upToken
		uploadObj.tokenMakeTime = time.Now()
	}
	return
}

/**
 * 用于解析配置文件里面的地域
 *
 * param: string zone
 * return: *storage.Region
 * return: error
 */
func getZone(zone string) (region *storage.Region,err error){
	switch zone {
		case "ZoneHuanan":
			region = &storage.ZoneHuanan
		case "ZoneHuadong":
			region = &storage.ZoneHuadong
		case "ZoneHuabei":
			region = &storage.ZoneHuabei
		case "ZoneBeimei":
			region = &storage.ZoneBeimei
		case "ZoneXinjiapo":
			region = &storage.ZoneXinjiapo
		default:
			region = &storage.ZoneHuadong
			err = errors.New("zone string error")
	}
	return
}

/**
 * 将文件上传到七牛云存储 *
 * param: string filename
 * param: []byte data
 * param: int64  dataLen
 * return: string
 * return: error
 */
func (uploadObj QiNiuUpload) UpLoadFile (filename string,data []byte,dataLen int64) (distName string,err error){
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = uploadObj.zone
	// 是否使用https域名
	cfg.UseHTTPS = uploadObj.https
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	ret := storage.PutRet{}
	formUploader := storage.NewFormUploader(&cfg)
	upToken := uploadObj.getToken()
	filename = uploadObj.pathTo + filename
	distName = ""
	putExtra := storage.PutExtra{}
	err = formUploader.Put(context.Background(),
		&ret, upToken, filename, bytes.NewReader(data), dataLen, &putExtra)
	distName = ret.Key
	return
}
