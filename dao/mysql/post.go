package mysql

import (
	"bluebell/models"
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(
                 post_id,title,content,author_id,community_id)
    values (?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostByID 根据id查询单个帖子数据
func GetPostByID(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select 
    post_id,title,content,author_id,community_id,create_time 
	from post 
	where post_id=?`
	err = db.Get(post, sqlStr, pid)
	if err == sql.ErrNoRows {
		err = ErrorInvalidID
	}
	return
}

// GetPostList 查询帖子列表函数
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select  post_id,title,content,author_id,community_id,create_time 
	from post order by create_time desc
	limit ? , ?`
	posts = make([]*models.Post, 0, size)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}

// GetPostListByIDs 根据帖子id列表查询帖子数据
func GetPostListByIDs(ids []string) (posts []*models.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time
	from post
	where post_id in (?)
	order by FIND_IN_SET(post_id,?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return
	}
	query = db.Rebind(query)
	posts = make([]*models.Post, 0, len(ids))
	err = db.Select(&posts, query, args...)
	return
}
