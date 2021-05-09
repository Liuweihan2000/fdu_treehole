package controller

import (
	"GoProject/fudan_bbs/dal"
	"GoProject/fudan_bbs/utils"
	"github.com/gin-gonic/gin"
)

// verify 验证密码是否正确
func verify(password, hash string) bool {
	if utils.MD5(password) != hash {
		return false
	}
	return true
	//if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
	//	return false
	//}
	//return true
}

// logged 用于判断是否已经登录
func logged(c *gin.Context) bool {
	_, err := session(c)
	return err == nil
}

// session 根据 cookie 读取session
func session(c *gin.Context) (s dal.Session, err error) {
	// 读取 cookie，值为 session 的 uuid
	cookie, err := c.Request.Cookie("cookie")
	if err != nil {
		msgErr(c, "读取用户cookie错误:", err)
		return
	}

	// 根据 uuid 读取 session
	// fmt.Println(cookie.Value)
	s, err = dal.ReadSessionByEmailHash(cookie.Value)
	return
}
