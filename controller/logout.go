package controller

import (
	"GoProject/fudan_bbs/dal"
	"github.com/gin-gonic/gin"
	"net/http"
)

func logout(c *gin.Context) {
	// 读取 session
	s, err := session(c)
	if err != nil {
		msgErr(c, "读取会话错误:", err)
		return
	}

	// 删除 session
	if err := dal.DeleteSession(s.ID); err != nil {
		msgErr(c, "删除会话错误:", err)
		return
	}

	c.Redirect(http.StatusFound, "/login")
}
