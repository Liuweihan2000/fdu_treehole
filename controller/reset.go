package controller

import (
	"GoProject/fudan_bbs/dal"
	"github.com/gin-gonic/gin"
	"net/http"
)

// reset 渲染密码重置页面
func reset(c *gin.Context) {
	html(c, nil, "layout", "reset")
}

// 密码重置页面 POST 方法的处理函数
func resetAction(c *gin.Context) {
	// 读取表单数据
	email := c.PostForm("email")
	password := c.PostForm("password")
	newPassword := c.PostForm("newPassword")

	// 根据 email 读取 user
	user, err := dal.QueryUserByEmail(email)
	if err != nil {
		msgErr(c, "通过邮箱读取用户错误: ", err)
		return
	}

	// 验证密码
	if !verify(password, user.Password) {
		msg(c, "密码错误")
		return
	}

	// 修改密码
	err = dal.ResetUserPassword(email, newPassword)
	if err != nil {
		msgErr(c, "修改密码错误", err)
		return
	}

	// 重定位到登录页面
	c.Redirect(http.StatusFound, "/login")
}
