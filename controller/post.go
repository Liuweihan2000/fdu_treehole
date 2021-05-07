package controller

import (
	"GoProject/fudan_bbs/internal/model"
	"GoProject/fudan_bbs/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// createPostAction 用来提交表单
func createPostAction(c *gin.Context) {
	if !logged(c) {
		msg(c, "未登录无法新建帖子")
		return
	}

	ID, _ := strconv.Atoi(c.PostForm("thread_id"))

	s, err := session(c)
	if err != nil {
		msgErr(c, "读取会话错误:", err)
		return
	}

	thread, err := DaoInstance.ReadThreadByID(int32(ID))
	if err != nil {
		msgErr(c, "通过id读取主题错误:", err)
		return
	}

	// user由session决定，而不是thread
	user, err := DaoInstance.QueryUserByID(s.UserID)
	if err != nil {
		msgErr(c, "通过user_id读取user错误:", err)
		return
	}

	var NickName string
	userNum := thread.UserCommented

	index, err := DaoInstance.QueryThreadUserPair(thread.ID, user.ID)
	if err != nil { // 新来的评论，添加thread_user pair
		NickName = utils.NameList[userNum]
		// 更新 thread
		err = DaoInstance.UpdateThreadIndex(int32(ID), userNum+1)
		if err != nil {
			msgErr(c, "更新 thread 出错: ", err)
		}
		// 添加 thread-user pair
		err = DaoInstance.CreateThreadUserPair(thread.ID, user.ID, userNum)
		if err != nil {
			msgErr(c, "添加thread-user pair出错: ", err)
		}
	} else { // 该用户已经在这个 thread 下面发表过评论
		NickName = utils.NameList[index]
	}

	post := model.Post{
		Content:   c.PostForm("content"),
		UserID:    user.ID,
		ThreadID:  thread.ID,
		UserName:  NickName,
		CreatedAt: time.Now(),
	}

	// 创建回复
	err = DaoInstance.CreatePost(&post)
	if err != nil {
		msgErr(c, "回复帖子出错:", err)
		return
	}

	err = DaoInstance.UpdateThreadTimeByID(thread.ID)
	if err != nil {
		msgErr(c, "更新帖子时间出错:", err)
		return
	}

	c.Redirect(http.StatusFound, "/threads/read?thread_id="+strconv.Itoa(int(thread.ID)))
}
