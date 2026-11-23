package controller

import (
    "github.com/gin-gonic/gin"
    "github.com/gorilla/feeds"
    . "happyblog/models"
    "net/http"
    "strconv"
    "strings"
    "time"
)

// RssController 生成站点的 RSS 订阅
type RssController struct{
    controller
}

var CRss RssController

/**
 * 输出 RSS 2.0 订阅内容,取最新的已发布文章
 *
 * param: *gin.Context c
 */
func (ctrl RssController) Rss(c *gin.Context) {
    siteName := webConf["sitename"]
    siteUrl := strings.TrimRight(webConf["siteurl"],"/")
    siteDesc := webConf["sitedesc"]
    if siteDesc == "" {
        siteDesc = siteName
    }

    feed := &feeds.Feed{
        Title:       siteName,
        Link:        &feeds.Link{Href: siteUrl},
        Description: siteDesc,
        Created:     time.Now(),
    }

    modelArticle := new(Article)
    articles,_ := modelArticle.GetArticleListView(1) //最新一页已发布文章
    for _, a := range articles {
        articleLink := siteUrl + "/article/" + strconv.Itoa(a.Id)
        author := ""
        if a.AuthorInfo.NickName != "" {
            author = a.AuthorInfo.NickName
        }
        feed.Items = append(feed.Items, &feeds.Item{
            Title:       a.Title,
            Link:        &feeds.Link{Href: articleLink},
            Description: a.Brief,
            Author:      &feeds.Author{Name: author},
            Created:     a.CreateTime,
            Updated:     a.UpdateTime,
            Id:          articleLink,
        })
    }

    rss, err := feed.ToRss()
    if err != nil {
        c.String(http.StatusInternalServerError, "generate rss error")
        return
    }
    c.Header("Content-Type", "application/rss+xml; charset=utf-8")
    c.String(http.StatusOK, rss)
}
