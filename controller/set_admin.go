package controller

//
//import (
//	"github.com/gin-gonic/gin"
//	"net/http"
//)
//
//// 渲染设置管理员页面，首先需要判断是否登录 && 进行操作的是否是管理员
//func setAdmin(c *gin.Context) {
//	// 判断是否登录
//	if !logged(c) {
//		msg(c, "请先登录再进行操作！")
//		return
//	}
//
//	// 读取 session
//	s, err := session(c)
//	if err != nil {
//		msgErr(c, "读取session错误: ", err)
//		return
//	}
//
//	// 根据 session 里的 user_id 读取 user
//	user, err := DaoInstance.QueryUserByID(s.UserID)
//	if err != nil {
//		msgErr(c, "通过ID读取用户错误:", err)
//		return
//	}
//
//	html(c, nil, "layout", "set_admin")
//}
//
//func setAdminAction(c *gin.Context) {
//	username := c.PostForm("username")
//	confirmUsername := c.PostForm("confirmUsername")
//
//	// 用户名确认
//	if username != confirmUsername {
//		msg(c, "两次用户名不一致")
//		return
//	}
//
//	// 根据 username 读取 user
//	user, err := DaoInstance.QueryUserByName(username)
//	if err != nil {
//		msg(c, "该用户不存在")
//		return
//	}
//
//	// 更新管理员权限
//	err = DaoInstance.SetAdmin(user.ID)
//	if err != nil {
//		msgErr(c, "更新管理员错误:", err)
//		return
//	}
//
//	c.Redirect(http.StatusFound, "/")
//}
