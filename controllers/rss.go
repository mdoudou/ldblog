package controllers

import (
	"fmt"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/mdoudou/ldblog/helpers"
	"github.com/mdoudou/ldblog/models"
	"github.com/mdoudou/ldblog/system"
)

func RssGet(c *gin.Context) {
	now := helpers.GetCurrentTime()
	domain := system.GetConfiguration().Domain
	feed := &feeds.Feed{
		Title:       "LD'blog",
		Link:        &feeds.Link{Href: domain},
		Description: "LD'blog,talk about golang,java and so on.",
		Author:      &feeds.Author{Name: "LD'blog", Email: "5352396@qq.com"},
		Created:     now,
	}

	feed.Items = make([]*feeds.Item, 0)
	posts, err := models.ListPublishedPost("", 0, 0)
	if err != nil {
		seelog.Error(err)
		return
	}

	for _, post := range posts {
		item := &feeds.Item{
			Id:          fmt.Sprintf("%s/post/%d", domain, post.ID),
			Title:       post.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/post/%d", domain, post.ID)},
			Description: string(post.Excerpt()),
			Created:     now,
		}
		feed.Items = append(feed.Items, item)
	}
	rss, err := feed.ToRss()
	if err != nil {
		seelog.Error(err)
		return
	}
	c.Writer.WriteString(rss)
}
