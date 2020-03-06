package repository_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/post"
	"github.com/ezio1119/fishapp-post/post/repository"
	"github.com/ezio1119/fishapp-post/testutil"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (post.Repository, sqlmock.Sqlmock, func()) {
	db, mock, _ := testutil.NewDBMock(t)
	r := repository.NewPostRepository(db)
	return r, mock, func() {
		db.Close()
	}
}

func TestGetByID(t *testing.T) {
	r, mock, dbClose := setup(t)
	defer dbClose()

	tests := map[string]struct {
		arrange func(t *testing.T)
		assert  func(t *testing.T, p *models.Post, err error)
	}{
		"正常に取得できること": {
			arrange: func(t *testing.T) {
				p := &models.Post{ID: 1, Title: "title", Content: "content", UserID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}
				rows := sqlmock.NewRows([]string{"id", "title", "content", "user_id", "created_at", "updated_at"}).
					AddRow(p.ID, p.Title, p.Content, p.UserID, p.CreatedAt, p.UpdatedAt)
				query := regexp.QuoteMeta(`SELECT id, title, content, user_id, created_at, updated_at FROM posts WHERE id = ?`)
				mock.ExpectQuery(query).WithArgs(0).WillReturnRows(rows)
			},
			assert: func(t *testing.T, p *models.Post, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, p)
			},
		},
		"SQLエラー発生時、エラーをラップして返すこと": {
			arrange: func(t *testing.T) {
				query := regexp.QuoteMeta(`SELECT id, title, content, user_id, created_at, updated_at FROM posts WHERE id = ?`)
				mock.ExpectQuery(query).WithArgs(0).WillReturnError(fmt.Errorf("some sql error"))
			},
			assert: func(t *testing.T, p *models.Post, err error) {
				e := models.WrapOnPostRepoErr(fmt.Errorf("some sql error"))
				assert.Nil(t, p)
				assert.EqualError(t, err, e.Error())
			},
		},
		"レコードが見つからない場合、エラーを発生させラップして返すこと": {
			arrange: func(t *testing.T) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "user_id", "created_at", "updated_at"})
				query := regexp.QuoteMeta(`SELECT id, title, content, user_id, created_at, updated_at FROM posts WHERE id = ?`)
				mock.ExpectQuery(query).WithArgs(0).WillReturnRows(rows)
			},
			assert: func(t *testing.T, p *models.Post, err error) {
				e := models.WrapOnPostRepoErr(models.NewPostNotFoundErr(0))
				assert.Nil(t, p)
				assert.EqualError(t, err, e.Error())
			},
		},
	}
	for k, tt := range tests {
		t.Run(k, func(t *testing.T) {
			tt.arrange(t)
			p, err := r.GetByID(context.TODO(), 0)
			tt.assert(t, p, err)
		})
	}
}

func TestGetList(t *testing.T) {
	r, mock, dbClose := setup(t)
	defer dbClose()

	tests := map[string]struct {
		arrange func(t *testing.T) time.Time
		assert  func(t *testing.T, p []*models.Post, err error)
	}{
		"正常に取得できること": {
			arrange: func(t *testing.T) time.Time {
				now := time.Now()
				p1 := &models.Post{ID: 1, Title: "title", Content: "content", UserID: 1, CreatedAt: now, UpdatedAt: now}
				p2 := &models.Post{ID: 2, Title: "title 2", Content: "content 2", UserID: 2, CreatedAt: now, UpdatedAt: now}
				rows := sqlmock.NewRows([]string{"id", "title", "content", "user_id", "created_at", "updated_at"}).
					AddRow(p1.ID, p1.Title, p1.Content, p1.UserID, p1.CreatedAt, p1.UpdatedAt).
					AddRow(p2.ID, p2.Title, p2.Content, p2.UserID, p2.CreatedAt, p2.UpdatedAt)
				query := regexp.QuoteMeta(`SELECT id, title, content, user_id, created_at, updated_at FROM posts WHERE created_at > ? ORDER BY created_at DESC LIMIT ?`)
				mock.ExpectQuery(query).WithArgs(now, 0).WillReturnRows(rows)
				return now
			},
			assert: func(t *testing.T, p []*models.Post, err error) {
				assert.NoError(t, err)
				assert.Len(t, p, 2)
			},
		},
		"SQLエラーが発生した場合、エラーをラップして返すこと": {
			arrange: func(t *testing.T) time.Time {
				now := time.Now()
				query := regexp.QuoteMeta(`SELECT id, title, content, user_id, created_at, updated_at FROM posts WHERE created_at > ? ORDER BY created_at DESC LIMIT ?`)
				mock.ExpectQuery(query).WithArgs(now, 0).WillReturnError(fmt.Errorf("some sql error"))
				return now
			},
			assert: func(t *testing.T, p []*models.Post, err error) {
				e := models.WrapOnPostRepoErr(fmt.Errorf("some sql error"))
				assert.Nil(t, p)
				assert.EqualError(t, err, e.Error())
			},
		},
		"レコードが見つからない場合、エラーを発生させラップして返すこと": {
			arrange: func(t *testing.T) time.Time {
				now := time.Now()
				rows := sqlmock.NewRows([]string{"id", "title", "content", "user_id", "created_at", "updated_at"})
				query := regexp.QuoteMeta(`SELECT id, title, content, user_id, created_at, updated_at FROM posts WHERE created_at > ? ORDER BY created_at DESC LIMIT ?`)
				mock.ExpectQuery(query).WithArgs(now, 0).WillReturnRows(rows)
				return now
			},
			assert: func(t *testing.T, p []*models.Post, err error) {
				e := models.WrapOnPostRepoErr(models.NewPostsNotFoundErr())
				assert.Nil(t, p)
				assert.EqualError(t, err, e.Error())
			},
		},
	}
	for k, tt := range tests {
		t.Run(k, func(t *testing.T) {
			now := tt.arrange(t)
			p, err := r.GetList(context.TODO(), now, 0)
			tt.assert(t, p, err)
		})
	}
}

func TestCreate(t *testing.T) {
	r, mock, dbClose := setup(t)
	defer dbClose()

	tests := map[string]struct {
		arrange func(t *testing.T) *models.Post
		assert  func(t *testing.T, p *models.Post, err error)
	}{
		"正常に更新できること": {
			arrange: func(t *testing.T) *models.Post {
				now := time.Now()
				p := &models.Post{Title: "title", Content: "content", UserID: 1, CreatedAt: now, UpdatedAt: now}
				query := regexp.QuoteMeta(`INSERT posts SET title=?, content=?, user_id=?, created_at=?, updated_at=?`)
				mock.ExpectPrepare(query).ExpectExec().WithArgs(p.Title, p.Content, p.UserID, p.CreatedAt, p.UpdatedAt).WillReturnResult(sqlmock.NewResult(19, 1))
				return p
			},
			assert: func(t *testing.T, p *models.Post, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(19), p.ID)
			},
		},
		"SQLエラー発生時、エラーをラップして返すこと": {
			arrange: func(t *testing.T) *models.Post {
				query := regexp.QuoteMeta(`INSERT posts SET title=?, content=?, user_id=?, created_at=?, updated_at=?`)
				mock.ExpectPrepare(query).ExpectExec().WillReturnError(fmt.Errorf("some sql error"))
				return &models.Post{}
			},
			assert: func(t *testing.T, _ *models.Post, err error) {
				e := models.WrapOnPostRepoErr(fmt.Errorf("some sql error"))
				assert.EqualError(t, err, e.Error())
			},
		},
	}
	for k, tt := range tests {
		t.Run(k, func(t *testing.T) {
			p := tt.arrange(t)
			err := r.Create(context.TODO(), p)
			tt.assert(t, p, err)
		})
	}
}

func TestUpdate(t *testing.T) {
	r, mock, dbClose := setup(t)
	defer dbClose()

	tests := map[string]struct {
		arrange func(t *testing.T) *models.Post
		assert  func(t *testing.T, err error)
	}{
		"正常に更新できること": {
			arrange: func(t *testing.T) *models.Post {
				now := time.Now()
				p := &models.Post{ID: 1, Title: "title", Content: "content", CreatedAt: now}
				query := regexp.QuoteMeta(`UPDATE posts SET title=?, content=?, updated_at=? WHERE id = ?`)
				mock.ExpectPrepare(query).ExpectExec().WithArgs(p.Title, p.Content, p.UpdatedAt, p.ID).WillReturnResult(sqlmock.NewResult(0, 1))
				return p
			},
			assert: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		"SQLエラー発生時、エラーをラップして返すこと": {
			arrange: func(t *testing.T) *models.Post {
				query := regexp.QuoteMeta(`UPDATE posts SET title=?, content=?, updated_at=? WHERE id = ?`)
				mock.ExpectPrepare(query).ExpectExec().WillReturnError(fmt.Errorf("some sql error"))
				return &models.Post{}
			},
			assert: func(t *testing.T, err error) {
				e := models.WrapOnPostRepoErr(fmt.Errorf("some sql error"))
				assert.EqualError(t, err, e.Error())
			},
		},
		"更新時、影響する行が1行以外の場合、エラーを発生させラップして返すこと": {
			arrange: func(t *testing.T) *models.Post {
				query := regexp.QuoteMeta(`UPDATE posts SET title=?, content=?, updated_at=? WHERE id = ?`)
				mock.ExpectPrepare(query).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 2)) // 2行影響
				return &models.Post{}
			},
			assert: func(t *testing.T, err error) {
				e := models.WrapOnPostRepoErr(models.NewRowsAffectedErr(2))
				assert.EqualError(t, err, e.Error())
			},
		},
	}
	for k, tt := range tests {
		t.Run(k, func(t *testing.T) {
			p := tt.arrange(t)
			err := r.Update(context.TODO(), p)
			tt.assert(t, err)
		})
	}
}

func TestDelete(t *testing.T) {
	r, mock, dbClose := setup(t)
	defer dbClose()

	tests := map[string]struct {
		arrange func(t *testing.T) int64
		assert  func(t *testing.T, err error)
	}{
		"正常に更新できること": {
			arrange: func(t *testing.T) int64 {
				deleteID := int64(19)
				query := regexp.QuoteMeta(`DELETE FROM posts WHERE id = ?`)
				mock.ExpectPrepare(query).ExpectExec().WithArgs(deleteID).WillReturnResult(sqlmock.NewResult(0, 1))
				return deleteID
			},
			assert: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		"SQLエラー発生時、エラーをラップして返すこと": {
			arrange: func(t *testing.T) int64 {
				query := regexp.QuoteMeta(`DELETE FROM posts WHERE id = ?`)
				mock.ExpectPrepare(query).ExpectExec().WillReturnError(fmt.Errorf("some sql error"))
				return 0
			},
			assert: func(t *testing.T, err error) {
				e := models.WrapOnPostRepoErr(fmt.Errorf("some sql error"))
				assert.EqualError(t, err, e.Error())
			},
		},
		"更新時、影響する行が1行以外の場合、エラーを発生させラップして返すこと": {
			arrange: func(t *testing.T) int64 {
				query := regexp.QuoteMeta(`DELETE FROM posts WHERE id = ?`)
				mock.ExpectPrepare(query).ExpectExec().WillReturnResult(sqlmock.NewResult(0, 2)) // 2行影響
				return 0
			},
			assert: func(t *testing.T, err error) {
				e := models.WrapOnPostRepoErr(models.NewRowsAffectedErr(2))
				assert.EqualError(t, err, e.Error())
			},
		},
	}
	for k, tt := range tests {
		t.Run(k, func(t *testing.T) {
			id := tt.arrange(t)
			err := r.Delete(context.TODO(), id)
			tt.assert(t, err)
		})
	}
}
