package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// ParamSignUp 注册参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登陆参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamVoteData struct {
	// UserID 从当前请求中获取
	PostID    int64 `json:"post_id,string" binding:"required"`       // 帖子id
	Direction int8  `json:"direction,string" binding:"oneof=1 0 -1"` // 投票类型，赞成票(1) or 反对票(-1) or 取消投票(0). 此处不要设置required，否则传入0时会自动过滤， 然后被认为没有配置该字段
}

// ParamPostList 获取帖子列表参数
type ParamPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"` // 可以为空， 如果为空走普通查询， 如果不为空， 走社区查询
	Page        int64  `json:"page"  form:"page"`
	Size        int64  `json:"size"  form:"size"`
	Order       string `json:"order"  form:"order"`
}
