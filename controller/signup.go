package controller

import (
	"GoProject/fudan_bbs/dal"
	"GoProject/fudan_bbs/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// signUp() renders the sign up page
func signUp(c *gin.Context) {
	html(c, nil, "layout", "sign_up")
}

func signUpAction(c *gin.Context) {
	u := dal.User{
		EmailHash: utils.MD5WithSalt(c.PostForm("email")),
		Password:  c.PostForm("password"),
	}
	// confirm password
	if u.Password != c.PostForm("confirmPassword") {
		msg(c, "两次密码不一致")
		return
	}
	// read the form and create a user
	if _, err := dal.CreateUser(u.EmailHash, u.Password); err != nil {
		msgErr(c, "创建用户错误", err)
		return
	}

	c.Redirect(http.StatusFound, "/login")
}
