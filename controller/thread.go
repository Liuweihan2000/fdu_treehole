package controller

import (
	"GoProject/fudan_bbs/common"
	"GoProject/fudan_bbs/dal"
	"GoProject/fudan_bbs/internal/model"
	"GoProject/fudan_bbs/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// createThread 渲染新建主题页面
func createThread(c *gin.Context) {
	if !logged(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	html(c, nil, "layout", "create_thread")
}

// createThreadAction 根据表单来创建主题
func createThreadAction(c *gin.Context) {
	// 读取 session
	s, err := session(c)
	if err != nil {
		msgErr(c, "读取会话错误:", err)
		return
	}

	// 读取 session 对应的用户
	user, err := dal.QueryUserByID(s.UserID)
	if err != nil {
		msgErr(c, "通过ID读取用户错误:", err)
		return
	}

	// TODO: 开启一个事务，将下面的三个操作放在同一个事务中作为原子操作
	// 为该用户建立主题
	now := time.Now()
	t := model.Thread{
		CreatedAt:     now,
		UpdatedAt:     now,
		UserID:        user.ID,
		UserCommented: 1,
	}
	if err = dal.CreateThread(&t); err != nil {
		msgErr(c, "创建主题错误:", err)
		return
	}
	threadID := t.ID
	// 创建主题的同时创建第一条回复
	post := model.Post{
		Content:   c.PostForm("content"),
		UserID:    user.ID,
		ThreadID:  threadID,
		UserName:  "洞主",
		CreatedAt: time.Now(),
	}
	if err = dal.CreatePost(&post); err != nil {
		msgErr(c, "创建主题错误:", err)
		return
	}

	_ = dal.CreateThreadUserPair(threadID, user.ID, 0)
	c.Redirect(http.StatusFound, "/")
}

// 临时用于传入Thread数据
type ReadThread struct {
	ThreadID        int32
	ThreadCreatedAt string
	PostCount       int32
	Posts           []Post
	Followed        bool
}

type Post struct {
	ID        int32
	Content   string
	UserID    int32
	ThreadID  int32
	UserName  string
	CreatedAt string
}

// readThread 渲染 /threads/read?thread_id= 页面
func readThread(c *gin.Context) {
	if !logged(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 填充数据
	var data ReadThread
	ID, _ := strconv.Atoi(c.Query("thread_id"))
	thread, err := dal.ReadThreadByID(int32(ID))
	if err != nil {
		msgErr(c, "通过ID读取帖子错误:", err)
		return
	}
	data.ThreadCreatedAt = thread.CreatedAt.Format("2006-01-02 15:04:05")
	data.ThreadID = thread.ID

	count, err := dal.CountByThreadID(thread.ID)
	if err != nil {
		count = 0
	}
	data.PostCount = int32(count)
	s, _ := session(c)
	err = dal.QueryUserThreadPair(s.UserID, int32(ID))
	if err == nil { // 找到记录，说明用户已经收藏了这个帖子
		data.Followed = true
	} else {
		data.Followed = false
	}

	// 如果上面的 count == 0 的话说明这个话题下还没有帖子
	if count != 0 {
		posts, err := dal.ReadPostsByThreadID(thread.ID)
		if err != nil {
			msgErr(c, "通过主题ID读取帖子错误:", err)
			return
		}
		for _, post := range posts {
			formatedPost := Post{
				ID:        post.ID,
				Content:   post.Content,
				UserName:  post.UserName,
				UserID:    post.UserID,
				ThreadID:  post.ThreadID,
				CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
			}
			data.Posts = append(data.Posts, formatedPost)
		}
	}

	html(c, data, "layout", "read_thread")
}

func searchThread(c *gin.Context) {
	if !logged(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	// 填充数据
	var data common.Index

	ID, err := strconv.Atoi(c.Query("thread_id"))
	if err != nil {
		msg(c, "抱歉，目前仅支持id搜索，请输入整数id")
		return
	}

	utils.Debug(ID)

	t, err := dal.ReadThreadByID(int32(ID))
	if err != nil {
		msgErr(c, "读取主题错误", err)
		return
	}

	data.ThreadCreatedAt = t.CreatedAt.Format("2006-01-02 15:04:05")
	data.ThreadID = t.ID

	count, err := dal.CountByThreadID(t.ID)
	if err != nil {
		msgErr(c, "通过主题ID读取帖子错误:", err)
		return
	}
	data.PostCount = int32(count)
	firstPost, _ :=dal.ReadFirstPostByThreadID(data.ThreadID)
	data.FirstPostContent = firstPost.Content
	timeDiff := utils.GetHourDiffer(t.UpdatedAt.Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
	data.TimePassed = utils.GetTimeDiff(timeDiff)

	html(c, data, "layout", "search_thread")
}

func followThreadAction(c *gin.Context) {
	if !logged(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	threadID, _ := strconv.Atoi(c.PostForm("thread_id"))

	s, err := session(c)
	if err != nil {
		msgErr(c, "读取会话错误:", err)
		return
	}

	userID := s.UserID

	if err = dal.CreateUserThreadPair(userID, int32(threadID)); err != nil {
		msgErr(c, "收藏出错:", err)
		return
	}

	location := fmt.Sprintf("/threads/read?thread_id=%d", threadID)

	c.Redirect(http.StatusFound, location)
}

func readFollowThread(c *gin.Context) {
	if !logged(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var data []common.Index
	s, _ := session(c)

	threads, err := dal.ReadUserFollowedThreads(s.UserID)
	if err != nil {
		msgErr(c, "读取全部主题错误", err)
		return
	}
	for _, t := range threads {
		var index common.Index
		index.ThreadCreatedAt = t.CreatedAt.Format("2006-01-02 15:04:05")
		index.ThreadID = t.ID

		count, err := dal.CountByThreadID(t.ID)
		if err != nil {
			msgErr(c, "通过主题ID读取帖子错误:", err)
			return
		}
		index.PostCount = int32(count)
		firstPost, _ := dal.ReadFirstPostByThreadID(index.ThreadID)
		index.FirstPostContent = firstPost.Content
		timeDiff := utils.GetHourDiffer(t.UpdatedAt.Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
		index.TimePassed = utils.GetTimeDiff(timeDiff)

		data = append(data, index)
	}

	html(c, data, "layout", "index")

}
