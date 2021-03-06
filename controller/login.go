package controller

import (
	"GoProject/fudan_bbs/dal"
	"GoProject/fudan_bbs/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// login 渲染登录页面
func login(c *gin.Context) {
	html(c, nil, "layout", "login")
}

func loginAction(c *gin.Context) {
	// 根据 email 读取用户
	emailHash := utils.MD5WithSalt(c.PostForm("email"))
	user, err := dal.QueryUserByEmailHash(emailHash)
	// fmt.Println(emailHash)
	if err != nil {
		msgErr(c, "通过邮箱读取用户错误:", err)
		return
	}

	// 验证密码
	if !verify(c.PostForm("password"), user.Password) {
		msg(c, "密码不正确")
		return
	}

	// 建立 session
	// 如果 session 已经存在就不再添加了
	err = dal.ReadSessionByUserID(user.ID)
	if err != nil { // 数据库还没有对应的 session 记录
		fmt.Println(err)
		session, err := dal.CreateSession(user.ID, user.EmailHash)
		if err != nil {
			msgErr(c, "创建会话错误:", err)
			return
		}

		// 设置cookie
		c.SetCookie("cookie", session.EmailHash, 2147483647, "", "localhost", false, true)
	}

	c.Redirect(http.StatusFound, "/")
}
