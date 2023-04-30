package mysql

import (
	"bluebell/models"
	"bluebell/settings"
	"testing"
)

func init() {
	// 注意单元测试此处一定要写测试用数据库， 不能使用生产库
	dbCfg := settings.MySQLConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "rootpass",
		DbName:       "bluebell",
		Port:         3306,
		MaxIdleConns: 10,
		MaxOpenConns: 20,
	}
	err := Init(&dbCfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := models.Post{
		ID:          90,
		AuthorID:    123,
		CommunityID: 2,
		Title:       "test",
		Content:     "just a test",
	}
	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("CreatePost insert record into mysql failed,err: %v\n", err)
	}
	t.Logf("CreatePost insert record into mysql success")
}
