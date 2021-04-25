package common

// 临时用于传入渲染的数据
type Index struct {
	ThreadID         int32  `json:"thread_id"`
	ThreadCreatedAt  string `json:"thread_created_at"`
	PostCount        int32  `json:"post_count"`
	FirstPostContent string `json:"first_post_content"`
	TimePassed       string `json:"time_passed"`
	ThreadUpdatedAt  string `json:"thread_updated_at"`
}
