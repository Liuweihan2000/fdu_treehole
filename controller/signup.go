package controller

import (
	"GoProject/fudan_bbs/internal/model"
	"GoProject/fudan_bbs/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// signUp() renders the sign up page
func signUp(c *gin.Context) {
	html(c, nil, "layout", "sign_up")
}

func signUpAction(c *gin.Context) {
	u := model.User{
		EmailHash: utils.MD5WithSalt(c.PostForm("email")),
		Password:  c.PostForm("password"),
	}
	// confirm password
	if u.Password != c.PostForm("confirmPassword") {
		msg(c, "两次密码不一致")
		return
	}
	// read the form and create a user
	if _, err := DaoInstance.CreateUser(u.EmailHash, u.Password); err != nil {
		msgErr(c, "创建用户错误", err)
		return
	}

	c.Redirect(http.StatusFound, "/login")
}
