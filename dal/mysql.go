package dal

import (
	"GoProject/fudan_bbs/common"
	"GoProject/fudan_bbs/internal/model"
	"GoProject/fudan_bbs/utils"
	"context"
	"gorm.io/gorm"
	"time"
)

var db *gorm.DB

func MySQL() *gorm.DB {
	return db
}

// TODO: 这里需要改一下
func QueryUserByEmail(email string) (user model.User, err error) {
	err = MySQL().Where("email_hash = ?", email).First(&user).Error
	return
}

func QueryUserByID(ID int32) (user model.User, err error) {
	err = MySQL().Where("id = ?", ID).First(&user).Error
	return
}

func CreateUser(EmailHash string, Password string) (user model.User, err error) {
	user = model.User{
		EmailHash: EmailHash,
		Password:  utils.MD5(Password),
	}
	err = MySQL().Create(&user).Error
	return
}

func ResetUserPassword(Email string, newPassword string) (err error) {
	user := model.User{}
	// 调用这个函数之前已经处理过 email 不存在的错误了，所以不用处理下面这段代码可能的错误
	MySQL().Where("email = ?", Email).First(&user)
	user.Password = utils.Hash(newPassword)
	err = MySQL().Save(&user).Error
	return
}

func QueryUserByEmailHash(EmailHash string) (user model.User, err error) {
	err = MySQL().Where("email_hash = ?", EmailHash).First(&user).Error
	return
}

// sessions
func CreateSessionByID(ID int32) (session model.Session, err error) {
	session.UserID = ID
	err = MySQL().Create(&session).Error
	return
}

func ReadSessionByUUID(UUID string) (session model.Session, err error) {
	err = MySQL().Where("uuid = ?", UUID).First(&session).Error
	return
}

func ReadSessionByUserID(ID int32) (err error) {
	// 如果能找到对应记录，则返回nil
	session := model.Session{}
	err = MySQL().Where("user_id = ?", ID).First(&session).Error
	return
}

func ReadSessionByEmailHash(hash string) (session model.Session, err error) {
	err = MySQL().Where("email_hash = ?", hash).First(&session).Error
	return
}

func DeleteSession(ID int32) (err error) {
	session := model.Session{
		ID: ID,
	}
	err = MySQL().Delete(&session).Error
	return
}

func CreateSession(UserID int32, EmailHash string) (session model.Session, err error) {
	session.UserID = UserID
	session.EmailHash = EmailHash
	err = MySQL().Create(&session).Error
	return
}

// threads
func  ReadAllThreads() (threads []model.Thread, err error) {
	result := MySQL().Find(&threads)
	err = result.Error
	return
}

func CreateThread(thread *model.Thread) (err error) {
	err = MySQL().Create(&thread).Error
	return
}

func  ReadThreadByUUID(UUID string) (thread model.Thread, err error) {
	err = MySQL().Where("uuid = ?", UUID).First(&thread).Error
	return
}

func ReadThreadByID(ID int32) (thread model.Thread, err error) {
	err = MySQL().Where("id = ?", ID).First(&thread).Error
	return
}

func  UpdateThreadIndex(ID int32, UserCommented int32) (err error) {
	thread := model.Thread{}
	err = MySQL().Model(&thread).Where("id = ?", ID).Update("user_commented", UserCommented).Error
	return
}

func GetLastThreadID() (ID int32, err error) {
	thread := model.Thread{}
	err = MySQL().Last(&thread).Error
	ID = thread.ID
	return
}

func  ReadAllThreadsByTime() (threads []model.Thread, err error) {
	err = MySQL().Order("updated_at desc").Find(&threads).Error
	return
}

func GetBatchThreadsByTime() ([]*common.Index, error) {
	slice, err := Redis().LRange(context.Background(), "threads_refresh_every_30s", 0, 30).Result()
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
func LoadThreadsToRedis() error {
	ctx := context.Background()
	curLen, err := Redis().LLen(ctx, "threads_refresh_every_30s").Result()
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
	err = MySQL().Order("updated_at desc").Find(&threads).Limit(30).Error
	if err != nil {
		return err
	}
	indices := make([]string, 0)
	for i := len(threads) - 1; i >= 0; i-- {
		thread := threads[i]
		firstPost, err := ReadFirstPostByThreadID(int32(thread.ID))
		if err != nil {
			return err
		}
		index := &common.Index{
			ThreadID:         thread.ID,
			ThreadCreatedAt:  thread.CreatedAt.Format("2006-01-02 15:04:05"),
			PostCount:        thread.UserCommented,
			FirstPostContent: firstPost.Content,
			ThreadUpdatedAt:  thread.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		indexStr, _ := utils.Marshal(index)
		indices = append(indices, indexStr)
	}
	Redis().LPush(ctx, "threads_refresh_every_30s", indices)

	if needClear {
		err := Redis().LTrim(ctx, "threads_refresh_every_30s", 0, 30).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateThreadTimeByID(ID int32) (err error) {
	thread := model.Thread{}
	err = MySQL().Model(&thread).Where("id = ?", ID).Update("updated_at", time.Now()).Error
	return
}

func  ReadUserFollowedThreadID(ID int32) (pairs []model.UserThread, err error) {
	err = MySQL().Model(&model.UserThread{}).Select("thread_id").Where("user_id = ?", ID).Find(&pairs).Error
	return
}

func ReadUserFollowedThreads(ID int32) (threads []model.Thread, err error) {
	pairList, err := ReadUserFollowedThreadID(ID)
	if err != nil {
		return
	}
	list := make([]int32, len(pairList))
	for _, pair := range pairList {
		list = append(list, pair.ThreadID)
	}
	// utils.Debug(1)
	err = MySQL().Where("id IN (?)", list).Find(&threads).Error
	return
}

// posts
func CountByThreadID(ID int32) (count int64, err error) {
	err = MySQL().Model(&model.Post{}).Where("thread_id = ?", ID).Count(&count).Error
	return
}

func ReadPostsByThreadID(ID int32) (posts []model.Post, err error) {
	err = MySQL().Where("thread_id = ?", ID).Find(&posts).Error
	return
}

func CreatePost(post *model.Post) (err error) {
	err = MySQL().Create(&post).Error
	return
}

func ReadFirstPostByThreadID(ID int32) (post model.Post, err error) {
	err = MySQL().Where("thread_id = ?", ID).First(&post).Error
	return
}

// thread->user
func  CreateThreadUserPair(ThreadID, UserID, UserNum int32) (err error) {
	threadUser := model.ThreadUser{
		ThreadID: ThreadID,
		UserID:   UserID,
		UserNum:  UserNum,
	}
	err = MySQL().Create(&threadUser).Error
	return
}

func  QueryThreadUserPair(ThreadID, UserID int32) (index int32, err error) {
	tu := model.ThreadUser{}
	err = MySQL().Where("thread_id = ? AND user_id = ?", ThreadID, UserID).First(&tu).Error
	index = tu.UserNum
	return
}

// user->thread
func  CreateUserThreadPair(UserID, ThreadID int32) (err error) {
	pair := model.UserThread{
		UserID:   UserID,
		ThreadID: ThreadID,
	}
	err = MySQL().Create(&pair).Error
	return
}

func  QueryUserThreadPair(UserID, ThreadID int32) (err error) {
	pair := model.UserThread{}
	err = MySQL().Where("user_id = ? AND thread_id = ?", UserID, ThreadID).Find(&pair).Error
	return
}
