package dao

import (
	"GoProject/fudan_bbs/controller"
	"GoProject/fudan_bbs/internal/model"
	redis "github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

var Provider = wire.NewSet(NewDao, NewDB)

type Dao interface {

	// users
	QueryUserByEmail(email string) (model.User, error)
	QueryUserByID(ID int32) (model.User, error)
	QueryUserByName(name string) (model.User, error)
	CreateUser(EmailHash string, password string) (model.User, error)
	ResetUserPassword(Email string, newPassword string) error
	SetAdmin(ID int32) error
	QueryUserByEmailHash(EmailHash string) (model.User, error)

	// sessions
	CreateSessionByID(ID int32) (model.Session, error)
	ReadSessionByUUID(UUID string) (model.Session, error)
	ReadSessionByUserID(ID int32) error
	ReadSessionByEmailHash(Hash string) (model.Session, error)
	DeleteSession(ID int32) error
	CreateSession(UserID int32, EmailHash string) (model.Session, error)

	// threads
	ReadAllThreads() ([]model.Thread, error)
	CreateThread(thread model.Thread) error
	ReadThreadByUUID(UUID string) (model.Thread, error)
	ReadThreadByID(ID int32) (model.Thread, error)
	UpdateThreadIndex(ID int32, UserCommented int32) error
	GetLastThreadID() (int32, error)
	ReadAllThreadsByTime() ([]model.Thread, error)
	GetBatchThreadsByTime() ([]*controller.Index, error) // 直接返回渲染页面所需要的内容
	UpdateThreadTimeByID(ID int32) error
	ReadUserFollowedThreadID(ID int32) ([]model.UserThread, error)
	ReadUserFollowedThreads(Id int32) ([]model.Thread, error)

	// posts
	CountByThreadID(ID int32) (int32, error)
	ReadPostsByThreadID(ID int32) ([]model.Post, error)
	CreatePost(post model.Post) error
	ReadFirstPostByThreadID(ID int32) (model.Post, error)

	// thread->user
	CreateThreadUserPair(ThreadID, UserID, UserNum int32) error
	QueryThreadUserPair(ThreadID, UserID int32) (int32, error)

	// user->thread
	CreateUserThreadPair(UserID, ThreadID int32) error
	QueryUserThreadPair(UserID, ThreadID int32) error
}

// struct Dao implements interface dao
type dao struct {
	mysql *gorm.DB
	redis *redis.Client
}

// NewDao new a new Dao
func NewDao(ormDB *gorm.DB) (d Dao, err error) {
	d = &dao{
		mysql: ormDB,
	}
	return
}
