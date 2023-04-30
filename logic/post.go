package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	// 1. 生成post id
	p.ID = snowflake.GenID()
	// 2. 保存到数据库
	// 3. 返回
	err = mysql.CreatePost(p)
	if err != nil {
		return
	}
	return redis.CreatePost(p.ID, p.CommunityID)
}

func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	// 查询并组合我们接口想用的数据
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(pid) failed",
			zap.Int64("pid", pid),
			zap.Error(err))
		return
	}
	// 根据作者id查询作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.getUserById(post.AuthorID) failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}

	// 根据社区id查询社区详细信息
	community, err := mysql.GetCommunityDetail(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetail(post.CommunityID) failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}

	// 拼接数据
	data = &models.ApiPostDetail{
		AuthorName:      user.UserName,
		CommunityDetail: community,
		Post:            post,
	}
	return
}

func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList() failed", zap.Error(err))
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, size)
	for _, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.getUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}

		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetail(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}

		// 拼接数据
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			CommunityDetail: community,
			Post:            post,
		}
		data = append(data, postdetail)
	}

	return
}

// GetPostListNew 获取帖子列表， 支持按照时间或者投票进行排序， 也可支持按照社区进行查询
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	var ids []string
	if p.CommunityID != 0 {
		ids, err = redis.GetCommunityPostIDInOrder(p)
	} else {
		ids, err = redis.GetPostIDInOrder(p)
	}
	if err != nil {
		zap.L().Error("查询redis列表失败", zap.Error(err))
		return
	}
	// 2、根据帖子id列表从数据库查询帖子数据

	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		zap.L().Error("从mysql查询数据失败", zap.Error(err))
		return
	}

	// 提前查询好帖子的投票数据
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		zap.L().Error("查询帖子票数信息失败", zap.Error(err))
		return
	}

	// 将帖子的作者信息和社区信息处理
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.getUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}

		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetail(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}

		// 拼接数据
		postdetail := &models.ApiPostDetail{
			VoteNum:         voteData[idx],
			AuthorName:      user.UserName,
			CommunityDetail: community,
			Post:            post,
		}
		data = append(data, postdetail)
	}

	return
}
