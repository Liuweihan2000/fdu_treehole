package controller

import (
	"GoProject/fudan_bbs/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"time"
)

// index 获取所有主题，并渲染 /index.html 页面
func index(c *gin.Context) {
	if !logged(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 填充数据
	// var data []Index

	indices, err := DaoInstance.GetBatchThreadsByTime()
	if err != nil {
		msgErr(c, "读取主题错误", err)
		return
	}
	for _, index := range indices {
		index.TimePassed = utils.GetTimeDiff(utils.GetHourDiffer(index.ThreadUpdatedAt, time.Now().Format("2006-01-02 15:04:05")))
	}
	//threads, err := DaoInstance.ReadAllThreadsByTime()
	//if err != nil {
	//	msgErr(c, "读取全部主题错误", err)
	//	return
	//}
	//for _, thread := range threads {
	//	var index Index
	//	index.ThreadCreatedAt = thread.CreatedAt.Format("2006-01-02 15:04:05")
	//	index.ThreadID = thread.ID
	//
	//	count := thread.UserCommented
	//	//count, err := DaoInstance.CountByThreadID(thread.ID)
	//	//if err != nil {
	//	//	msgErr(c, "通过主题ID读取帖子错误:", err)
	//	//	return
	//	//}
	//	index.PostCount = count
	//	//firstPost, _ := DaoInstance.ReadFirstPostByThreadID(index.ThreadID)
	//	//index.FirstPostContent = firstPost.Content
	//	timeDiff := utils.GetHourDiffer(thread.UpdatedAt.Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
	//	index.TimePassed = utils.GetTimeDiff(timeDiff)
	//
	//	data = append(data, index)
	//}

	html(c, indices, "layout", "index")
}

// err 根据传入的msg URL参数来渲染 /err 错误页面
func err(c *gin.Context) {
	msg := c.Query("msg")
	html(c, msg, "layout", "err")
}

// msgErr 组合传入的参数和错误信息重定向到 /err?msg= 错误页面
func msgErr(c *gin.Context, msg string, err error) {
	c.Redirect(http.StatusFound, "/err?msg="+msg+err.Error())
}

// msg 组合传入的参数和错误信息重定向到 /err?msg= 错误页面
func msg(c *gin.Context, msg string) {
	c.Redirect(http.StatusFound, "/err?msg="+msg)
}

// html 根据传入的数据 data 和应该渲染的文件 names 来渲染页面
func html(c *gin.Context, data interface{}, names ...string) {
	var files []string
	for _, f := range names {
		files = append(files, fmt.Sprintf("view/%s.html", f))
	}
	if err := template.Must(template.ParseFiles(files...)).ExecuteTemplate(c.Writer, "layout", data); err != nil {
		msgErr(c, "渲染模板错误", err)
	}
}
