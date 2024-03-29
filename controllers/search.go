package controllers

import (
	"net/http"
	"strconv"

	"math"

	"github.com/gin-gonic/gin"
	"github.com/mdoudou/ldblog/models"
	"github.com/mdoudou/ldblog/system"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func SearchGet(c *gin.Context) {
	var (
		pageIndex int
		pageSize  = system.GetConfiguration().PageSize
		total     int
		page      string
		err       error
		posts     []*models.Post
		policy    *bluemonday.Policy
	)
	page = c.Query("page")
	pageIndex, _ = strconv.Atoi(page)
	if pageIndex <= 0 {
		pageIndex = 1
	}
	searchname := c.DefaultQuery("name", "")
	posts, err = models.SearchListPublishedPost("", pageIndex, pageSize, searchname)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	total, err = models.CountPostByTag("")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	policy = bluemonday.StrictPolicy()
	for _, post := range posts {
		post.Tags, _ = models.ListTagByPostId(strconv.FormatUint(uint64(post.ID), 10))
		post.Body = policy.Sanitize(string(blackfriday.MarkdownCommon([]byte(post.Body))))
	}
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"posts":           posts,
		"tags":            models.MustListTag(),
		"archives":        models.MustListPostArchives(),
		"links":           models.MustListLinks(),
		"user":            user,
		"pageIndex":       pageIndex,
		"totalPage":       int(math.Ceil(float64(total) / float64(pageSize))),
		"path":            c.Request.URL.Path,
		"maxReadPosts":    models.MustListMaxReadPost(),
		"maxCommentPosts": models.MustListMaxCommentPost(),
	})
}
