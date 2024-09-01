package test

import (
	"bytes"
	"encoding/json"
	"github.com/ChongYanOvO/little-blue-book/internal/handler"
	"github.com/ChongYanOvO/little-blue-book/internal/repository/dao"
	"github.com/ChongYanOvO/little-blue-book/wire"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ArticleTestSuite struct {
	suite.Suite
	Server *gin.Engine
	db     *gorm.DB
}

func (s *ArticleTestSuite) SetupSuite() {
	app, _ := wire.InitApp()
	s.Server = app.Server
	s.db, _ = wire.InitMysql()
	articleHandler, _ := wire.InitArticleHandler()
	articleHandler.RegisterRoutes(s.Server)
}

func (s *ArticleTestSuite) TearDownTest() {
	s.db.Exec("truncate table articles")
}

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (s *ArticleTestSuite) TestCreate() {
	t := s.T()
	testCases := []struct {
		name       string
		article    Article
		before     func(t *testing.T)
		after      func(t *testing.T)
		wantCode   int
		wantResult handler.Result[int64]
	}{
		{
			name: "新建成功",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
				var article dao.Article
				err := s.db.Where("title=?", "测试标题").Find(&article).Error
				assert.NoError(t, err)
				assert.True(t, article.CreateTime > 0)
				assert.True(t, article.UpdateTime > 0)
				article.CreateTime = 0
				article.UpdateTime = 0
				assert.Equal(t, dao.Article{
					Id:         1,
					Title:      "测试标题",
					Content:    "测试内容",
					CreateTime: 0,
					UpdateTime: 0,
				}, article)
			},
			article: Article{
				Title:   "测试标题",
				Content: "测试内容",
			},
			wantCode: http.StatusOK,
			wantResult: handler.Result[int64]{
				Code: http.StatusOK,
				Data: 1,
				Msg:  "success",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			req, err := json.Marshal(tc.article)
			assert.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, "/article/create", bytes.NewBuffer([]byte(req)))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			s.Server.ServeHTTP(response, request)
			assert.Equal(t, tc.wantCode, response.Code)
			if response.Code != 200 {
				return
			}
			var res handler.Result[int64]
			err = json.NewDecoder(response.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, res)
		})
	}
}

func (s *ArticleTestSuite) TestEdit() {
	t := s.T()
	testCases := []struct {
		name       string
		article    Article
		before     func(t *testing.T)
		after      func(t *testing.T)
		wantCode   int
		wantResult handler.Result[int64]
	}{
		{
			name: "修改成功",
			before: func(t *testing.T) {
				s.db.Create(&dao.Article{
					Id:         1,
					Title:      "旧的测试标题",
					Content:    "旧的测试内容",
					AuthorId:   123,
					CreateTime: 222,
					UpdateTime: 2222,
				})
			},
			after: func(t *testing.T) {
				var article dao.Article
				err := s.db.Where("title=?", "新的测试标题").Find(&article).Error
				assert.NoError(t, err)
				assert.True(t, article.CreateTime > 0)
				assert.True(t, article.UpdateTime > 0)
				article.CreateTime = 0
				article.UpdateTime = 0
				assert.Equal(t, dao.Article{
					Id:         1,
					Title:      "新的测试标题",
					Content:    "新的测试内容",
					AuthorId:   123,
					CreateTime: 0,
					UpdateTime: 0,
				}, article)
			},
			article: Article{
				Id:      1,
				Title:   "新的测试标题",
				Content: "新的测试内容",
			},
			wantCode: http.StatusOK,
			wantResult: handler.Result[int64]{
				Code: http.StatusOK,
				Data: 1,
				Msg:  "success",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			req, err := json.Marshal(tc.article)
			assert.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, "/article/edit", bytes.NewBuffer([]byte(req)))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			s.Server.ServeHTTP(response, request)
			assert.Equal(t, tc.wantCode, response.Code)
			if response.Code != 200 {
				return
			}
			var res handler.Result[int64]
			err = json.NewDecoder(response.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, res)
		})
	}
}

func TestArticle(t *testing.T) {
	suite.Run(t, new(ArticleTestSuite))
}
