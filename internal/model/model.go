package model

import "time"

type User struct {
	ID        int32
	EmailHash string
	Password  string
}

func (User) TableName() string {
	return "users"
}

type Thread struct {
	ID            int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
	UserID        int32
	UserCommented int32
}

func (Thread) TableName() string {
	return "threads"
}

type Session struct {
	ID        int32
	UserID    int32
	EmailHash string
}

func (Session) TableName() string {
	return "sessions"
}

type Post struct {
	ID        int32
	Content   string
	UserID    int32
	ThreadID  int32
	UserName  string
	CreatedAt time.Time
}

func (Post) TableName() string {
	return "posts"
}

type ThreadUser struct {
	ThreadID int32
	UserID   int32
	UserNum  int32
}

func (ThreadUser) TableName() string {
	return "thread_user"
}

type UserThread struct {
	UserID   int32
	ThreadID int32
}

func (UserThread) TableName() string {
	return "user_thread"
}
