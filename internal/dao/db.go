package dao

import (
	"GoProject/fudan_bbs/internal/model"
	"GoProject/fudan_bbs/internal/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

//func (d *Dao) QueryUserByID(ID int32) (UserInfo model.User) {
//
//	return
//}

// NewDB new a new gormDB
func NewDB() (ormDB *gorm.DB, err error) {
	ormDB, err = gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/bbs?charset=utf8mb4&parseTime=True&loc=Local")
	return
}

func (d dao) QueryUserByEmail(email string) (user model.User, err error) {
	err = d.mysql.Where("email = ?", email).First(&user).Error
	return
}

func (d dao) QueryUserByID(ID int32) (user model.User, err error) {
	err = d.mysql.Where("id = ?", ID).First(&user).Error
	return
}

func (d dao) QueryUserByName(name string) (user model.User, err error) {
	err = d.mysql.Where("username = ?", name).First(&user).Error
	return
}

func (d dao) CreateUser(EmailHash string, Password string) (user model.User, err error) {
	user = model.User{
		EmailHash: EmailHash,
		Password:  utils.MD5(Password),
	}
	err = d.mysql.Create(&user).Error
	return
}

func (d dao) ResetUserPassword(Email string, newPassword string) (err error) {
	user := model.User{}
	// 调用这个函数之前已经处理过 email 不存在的错误了，所以不用处理下面这段代码可能的错误
	d.mysql.Where("email = ?", Email).First(&user)
	user.Password = utils.Hash(newPassword)
	err = d.mysql.Save(&user).Error
	return
}

func (d dao) SetAdmin(ID int32) (err error) {
	user := model.User{}
	// 调用这个函数之前已经处理过 user 不存在的错误了，所以不用处理下面这段代码可能的错误
	d.mysql.Where("id = ?", ID).First(&user)
	err = d.mysql.Save(&user).Error
	return
}

func (d dao) QueryUserByEmailHash(EmailHash string) (user model.User, err error) {
	err = d.mysql.Where("email_hash = ?", EmailHash).First(&user).Error
	return
}

// sessions
func (d dao) CreateSessionByID(ID int32) (session model.Session, err error) {
	session.UserID = ID
	err = d.mysql.Create(&session).Error
	return
}

func (d dao) ReadSessionByUUID(UUID string) (session model.Session, err error) {
	err = d.mysql.Where("uuid = ?", UUID).First(&session).Error
	return
}

func (d dao) ReadSessionByUserID(ID int32) (err error) {
	// 如果能找到对应记录，则返回nil
	session := model.Session{}
	err = d.mysql.Where("user_id = ?", ID).First(&session).Error
	return
}

func (d dao) ReadSessionByEmailHash(hash string) (session model.Session, err error) {
	err = d.mysql.Where("email_hash = ?", hash).First(&session).Error
	return
}

func (d dao) DeleteSession(ID int32) (err error) {
	session := model.Session{
		ID: ID,
	}
	err = d.mysql.Delete(&session).Error
	return
}

func (d dao) CreateSession(UserID int32, EmailHash string) (session model.Session, err error) {
	session.UserID = UserID
	session.EmailHash = EmailHash
	err = d.mysql.Create(&session).Error
	return
}

// threads
func (d dao) ReadAllThreads() (threads []model.Thread, err error) {
	result := d.mysql.Find(&threads)
	err = result.Error
	return
}

func (d dao) CreateThread(thread model.Thread) (err error) {
	err = d.mysql.Create(&thread).Error
	return
}

func (d dao) ReadThreadByUUID(UUID string) (thread model.Thread, err error) {
	err = d.mysql.Where("uuid = ?", UUID).First(&thread).Error
	return
}

func (d dao) ReadThreadByID(ID int32) (thread model.Thread, err error) {
	err = d.mysql.Where("id = ?", ID).First(&thread).Error
	return
}

func (d dao) UpdateThreadIndex(ID int32, UserCommented int32) (err error) {
	thread := model.Thread{}
	err = d.mysql.Model(&thread).Where("id = ?", ID).Update("user_commented", UserCommented).Error
	return
}

func (d dao) GetLastThreadID() (ID int32, err error) {
	thread := model.Thread{}
	err = d.mysql.Last(&thread).Error
	ID = thread.ID
	return
}

func (d dao) ReadAllThreadsByTime() (threads []model.Thread, err error) {
	err = d.mysql.Order("updated_at desc").Find(&threads).Error
	return
}

func (d dao) UpdateThreadTimeByID(ID int32) (err error) {
	thread := model.Thread{}
	err = d.mysql.Model(&thread).Where("id = ?", ID).Update("updated_at", time.Now()).Error
	return
}

func (d dao) ReadUserFollowedThreadID(ID int32) (pairs []model.UserThread, err error) {
	err = d.mysql.Model(&model.UserThread{}).Select("thread_id").Where("user_id = ?", ID).Find(&pairs).Error
	return
}

func (d dao) ReadUserFollowedThreads(ID int32) (threads []model.Thread, err error) {
	pairList, err := d.ReadUserFollowedThreadID(ID)
	if err != nil {
		return
	}
	list := make([]int32, len(pairList))
	for _, pair := range pairList {
		list = append(list, pair.ThreadID)
	}
	// utils.Debug(1)
	err = d.mysql.Where("id IN (?)", list).Find(&threads).Error
	return
}

// posts
func (d dao) CountByThreadID(ID int32) (count int32, err error) {
	err = d.mysql.Model(&model.Post{}).Where("thread_id = ?", ID).Count(&count).Error
	return
}

func (d dao) ReadPostsByThreadID(ID int32) (posts []model.Post, err error) {
	err = d.mysql.Where("thread_id = ?", ID).Find(&posts).Error
	return
}

func (d dao) CreatePost(post model.Post) (err error) {
	err = d.mysql.Create(&post).Error
	return
}

func (d dao) ReadFirstPostByThreadID(ID int32) (post model.Post, err error) {
	err = d.mysql.Where("thread_id = ?", ID).First(&post).Error
	return
}

// thread->user
func (d dao) CreateThreadUserPair(ThreadID, UserID, UserNum int32) (err error) {
	threadUser := model.ThreadUser{
		ThreadID: ThreadID,
		UserID:   UserID,
		UserNum:  UserNum,
	}
	err = d.mysql.Create(&threadUser).Error
	return
}

func (d dao) QueryThreadUserPair(ThreadID, UserID int32) (index int32, err error) {
	tu := model.ThreadUser{}
	err = d.mysql.Where("thread_id = ? AND user_id = ?", ThreadID, UserID).First(&tu).Error
	index = tu.UserNum
	return
}

// user->thread
func (d dao) CreateUserThreadPair(UserID, ThreadID int32) (err error) {
	pair := model.UserThread{
		UserID:   UserID,
		ThreadID: ThreadID,
	}
	err = d.mysql.Create(&pair).Error
	return
}

func (d dao) QueryUserThreadPair(UserID, ThreadID int32) (err error) {
	pair := model.UserThread{}
	err = d.mysql.Where("user_id = ? AND thread_id = ?", UserID, ThreadID).Find(&pair).Error
	return
}
