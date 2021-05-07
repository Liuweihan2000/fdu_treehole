package dal

import (
	"GoProject/fudan_bbs/common"
	"GoProject/fudan_bbs/internal/model"
	"GoProject/fudan_bbs/internal/utils"
	"context"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

// NewDB new a new gormDB and a redis client
func NewDB() (d Dao, err error) {
	ormDB, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/bbs?charset=utf8mb4&parseTime=True&loc=Local")
	r := redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379",
			Password: "",
			DB: 0,
		},
	)
	d.mysql = ormDB
	d.redis = r
	return
}

func (d Dao) QueryUserByEmail(email string) (user model.User, err error) {
	err = d.mysql.Where("email = ?", email).First(&user).Error
	return
}

func (d Dao) QueryUserByID(ID int32) (user model.User, err error) {
	err = d.mysql.Where("id = ?", ID).First(&user).Error
	return
}

func (d Dao) QueryUserByName(name string) (user model.User, err error) {
	err = d.mysql.Where("username = ?", name).First(&user).Error
	return
}

func (d Dao) CreateUser(EmailHash string, Password string) (user model.User, err error) {
	user = model.User{
		EmailHash: EmailHash,
		Password:  utils.MD5(Password),
	}
	err = d.mysql.Create(&user).Error
	return
}

func (d Dao) ResetUserPassword(Email string, newPassword string) (err error) {
	user := model.User{}
	// 调用这个函数之前已经处理过 email 不存在的错误了，所以不用处理下面这段代码可能的错误
	d.mysql.Where("email = ?", Email).First(&user)
	user.Password = utils.Hash(newPassword)
	err = d.mysql.Save(&user).Error
	return
}

func (d Dao) SetAdmin(ID int32) (err error) {
	user := model.User{}
	// 调用这个函数之前已经处理过 user 不存在的错误了，所以不用处理下面这段代码可能的错误
	d.mysql.Where("id = ?", ID).First(&user)
	err = d.mysql.Save(&user).Error
	return
}

func (d Dao) QueryUserByEmailHash(EmailHash string) (user model.User, err error) {
	err = d.mysql.Where("email_hash = ?", EmailHash).First(&user).Error
	return
}

// sessions
func (d Dao) CreateSessionByID(ID int32) (session model.Session, err error) {
	session.UserID = ID
	err = d.mysql.Create(&session).Error
	return
}

func (d Dao) ReadSessionByUUID(UUID string) (session model.Session, err error) {
	err = d.mysql.Where("uuid = ?", UUID).First(&session).Error
	return
}

func (d Dao) ReadSessionByUserID(ID int32) (err error) {
	// 如果能找到对应记录，则返回nil
	session := model.Session{}
	err = d.mysql.Where("user_id = ?", ID).First(&session).Error
	return
}

func (d Dao) ReadSessionByEmailHash(hash string) (session model.Session, err error) {
	err = d.mysql.Where("email_hash = ?", hash).First(&session).Error
	return
}

func (d Dao) DeleteSession(ID int32) (err error) {
	session := model.Session{
		ID: ID,
	}
	err = d.mysql.Delete(&session).Error
	return
}

func (d Dao) CreateSession(UserID int32, EmailHash string) (session model.Session, err error) {
	session.UserID = UserID
	session.EmailHash = EmailHash
	err = d.mysql.Create(&session).Error
	return
}

// threads
func (d Dao) ReadAllThreads() (threads []model.Thread, err error) {
	result := d.mysql.Find(&threads)
	err = result.Error
	return
}

func (d Dao) CreateThread(thread *model.Thread) (err error) {
	err = d.mysql.Create(&thread).Error
	return
}

func (d Dao) ReadThreadByUUID(UUID string) (thread model.Thread, err error) {
	err = d.mysql.Where("uuid = ?", UUID).First(&thread).Error
	return
}

func (d Dao) ReadThreadByID(ID int32) (thread model.Thread, err error) {
	err = d.mysql.Where("id = ?", ID).First(&thread).Error
	return
}

func (d Dao) UpdateThreadIndex(ID int32, UserCommented int32) (err error) {
	thread := model.Thread{}
	err = d.mysql.Model(&thread).Where("id = ?", ID).Update("user_commented", UserCommented).Error
	return
}

func (d Dao) GetLastThreadID() (ID int32, err error) {
	thread := model.Thread{}
	err = d.mysql.Last(&thread).Error
	ID = thread.ID
	return
}

func (d Dao) ReadAllThreadsByTime() (threads []model.Thread, err error) {
	err = d.mysql.Order("updated_at desc").Find(&threads).Error
	return
}

func (d Dao) GetBatchThreadsByTime() ([]*common.Index, error) {
	slice, err := d.redis.LRange(context.Background(), "threads_refresh_every_30s", 0, 30).Result()
	if err != nil {
		return nil, err
	}
	res := make([]*common.Index, 0)
	for _, item := range slice {
		index := &common.Index{}
		if err := utils.UnMarshal(item, index); err != nil {
			return nil, err
		}
		res = append(res, index)
	}
	return res, nil
}

// 1 2 3 11 22 33
// 倒序插入redis
func (d Dao) LoadThreadsToRedis() error {
	ctx := context.Background()
	curLen, err := d.redis.LLen(ctx, "threads_refresh_every_30s").Result()
	if err != nil {
		return err
	}
	var needClear bool
	if curLen == 0 {
		needClear = false
	} else {
		needClear = true
	}

	threads := make([]*model.Thread, 0)
	err = d.mysql.Order("updated_at desc").Find(&threads).Limit(30).Error
	if err != nil {
		return err
	}
	indices := make([]string, 0)
	for i := len(threads) - 1; i >= 0; i-- {
		thread := threads[i]
		firstPost, err := d.ReadFirstPostByThreadID(int32(thread.ID))
		if err != nil {
			return err
		}
		index := &common.Index{
			ThreadID: thread.ID,
			ThreadCreatedAt: thread.CreatedAt.Format("2006-01-02 15:04:05"),
			PostCount: thread.UserCommented,
			FirstPostContent: firstPost.Content,
			ThreadUpdatedAt: thread.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		indexStr, _ := utils.Marshal(index)
		indices = append(indices, indexStr)
	}
	d.redis.LPush(ctx, "threads_refresh_every_30s", indices)

	if needClear {
		err := d.redis.LTrim(ctx, "threads_refresh_every_30s", 0, 30).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (d Dao) UpdateThreadTimeByID(ID int32) (err error) {
	thread := model.Thread{}
	err = d.mysql.Model(&thread).Where("id = ?", ID).Update("updated_at", time.Now()).Error
	return
}

func (d Dao) ReadUserFollowedThreadID(ID int32) (pairs []model.UserThread, err error) {
	err = d.mysql.Model(&model.UserThread{}).Select("thread_id").Where("user_id = ?", ID).Find(&pairs).Error
	return
}

func (d Dao) ReadUserFollowedThreads(ID int32) (threads []model.Thread, err error) {
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
func (d Dao) CountByThreadID(ID int32) (count int32, err error) {
	err = d.mysql.Model(&model.Post{}).Where("thread_id = ?", ID).Count(&count).Error
	return
}

func (d Dao) ReadPostsByThreadID(ID int32) (posts []model.Post, err error) {
	err = d.mysql.Where("thread_id = ?", ID).Find(&posts).Error
	return
}

func (d Dao) CreatePost(post *model.Post) (err error) {
	err = d.mysql.Create(&post).Error
	return
}

func (d Dao) ReadFirstPostByThreadID(ID int32) (post model.Post, err error) {
	err = d.mysql.Where("thread_id = ?", ID).First(&post).Error
	return
}

// thread->user
func (d Dao) CreateThreadUserPair(ThreadID, UserID, UserNum int32) (err error) {
	threadUser := model.ThreadUser{
		ThreadID: ThreadID,
		UserID:   UserID,
		UserNum:  UserNum,
	}
	err = d.mysql.Create(&threadUser).Error
	return
}

func (d Dao) QueryThreadUserPair(ThreadID, UserID int32) (index int32, err error) {
	tu := model.ThreadUser{}
	err = d.mysql.Where("thread_id = ? AND user_id = ?", ThreadID, UserID).First(&tu).Error
	index = tu.UserNum
	return
}

// user->thread
func (d Dao) CreateUserThreadPair(UserID, ThreadID int32) (err error) {
	pair := model.UserThread{
		UserID:   UserID,
		ThreadID: ThreadID,
	}
	err = d.mysql.Create(&pair).Error
	return
}

func (d Dao) QueryUserThreadPair(UserID, ThreadID int32) (err error) {
	pair := model.UserThread{}
	err = d.mysql.Where("user_id = ? AND thread_id = ?", UserID, ThreadID).Find(&pair).Error
	return
}