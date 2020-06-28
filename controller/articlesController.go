package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	lib "happyblog/library"
	. "happyblog/models"
	"log"
	"net/http"
	"strconv"
)

type ArticleController struct{
	controller
}

var CArticle ArticleController


func init(){
	webConf ,_ := lib.ReadWebConfig()
	CArticle.webConf = *webConf
	menu,_ := lib.ReadLanguageConfig(CArticle.webConf["language"])
	CArticle.menu = *menu
}
/**
 * 文章列表页面
 *
 * param: *gin.Context c
 */
func (ctrl ArticleController) Articles(c *gin.Context){
	obString := ctrl.makeOrderStr(c.Query("obstring"),c.Query("obstatus"))
	nowPage,_ := strconv.Atoi(c.DefaultQuery("page","1"))
	ModelArticle := new(Article)
	articleList ,err := ModelArticle.GetArticlesList(obString,nowPage)
	if err != nil {
		log.Fatal(fmt.Sprintf("ARTICLE LIST GET ERROR :err:%v",err))
	}

	totalRows,_ := ModelArticle.GetArticlesTotal()
	c.HTML(http.StatusOK, "index", gin.H{
		"title": ctrl.webConf["sitename"],
		"userInfo": ctrl.GetCacheUinfo(c),
		"menu":ctrl.menu,
		"articles":1,
		"articleslist":articleList,
		"obstatus" : c.DefaultQuery("obstatus","desc"),
		"ob" : c.DefaultQuery("obstring","ctime"),
		"curPage":nowPage,
		"totalRows":totalRows,
	})
}

/**
 * 用于组织ORDER BY 字段
 *
 * param: string ob
 * param: string obstatus
 * return: string
 */
func (ctrl ArticleController) makeOrderStr(ob string,obStatus string) (obString string) {
	if ob == "ctime" {
		if obStatus == "asc" {
			return "createTime ASC"
		} else {
			return "createTime DESC"
		}
	}
	if ob == "mtime" {
		if obStatus == "asc" {
			return "updateTime ASC"
		} else {
			return "updateTime DESC"
		}
	}
	if ob == "pub" {
		if obStatus == "asc" {
			return "pubStatus ASC"
		} else {
			return "pubStatus DESC"
		}
	}
	return "createTime DESC"
}
/**
 * 编辑文章的方法
 *
 * param: *gin.Context c
 */
func (ctrl ArticleController) ArticleEdit(c *gin.Context) {
	articleId,_ := strconv.Atoi(c.Param("id"))
	//如果获取到articleId  则去读取articleInfo
	ModelArticle := new(Article)
	articleInfo,err := ModelArticle.GetArticleInfo(articleId)
	if err != nil {
		log.Println("get article info err")
	}
	ModelRelationArticleAlbum := new(RelationArticleAlbums)
	relationAlbumId := ModelRelationArticleAlbum.GetBelongAlbumByArticleId(articleId)
	if err != nil {
		log.Println("get albumInfo err")
	}
	ModelAlbums := new(Album)
	albumsList,_ := ModelAlbums.GetAlbumList(-1)//获取所有的专辑信息
	//如果获取到了articleid ,则是编辑旧文章,读取旧文章内容.
	c.HTML(http.StatusOK, "index", gin.H{
		"title": ctrl.webConf["sitename"],
		"userInfo": ctrl.GetCacheUinfo(c),
		"menu":ctrl.menu,
		"articleEdit":1,
		"articleInfo":articleInfo,
		"albumsList":albumsList,//专辑列表
		"relationAlbumId" : relationAlbumId,//关联专辑ID
	})
}