package handler

import (
	"github.com/ChongYanOvO/little-blue-book/internal/domain"
	"github.com/ChongYanOvO/little-blue-book/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

var _ Handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	logger *zap.Logger
	svc    service.ArticleService
}

func NewArticleHandler(svc service.ArticleService, l *zap.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc:    svc,
		logger: l,
	}
}

func (ah *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	ag := server.Group("/article")
	ag.POST("/create", ah.Create)
	ag.POST("/edit", ah.Edit)
	ag.POST("/publish", ah.Publish)
}

func (ah *ArticleHandler) Create(ctx *gin.Context) {
	type CreateArticleRequest struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req CreateArticleRequest
	if err := ctx.Bind(&req); err != nil {
		ah.logger.Error("获取前端参数错误", zap.Error(err))
		return
	}
	uc, err := ExtractJwtClaims(ctx)
	if err != nil {
		ah.logger.Error("获取UserClaims错误", zap.Error(err))
		return
	}
	a, err := ah.svc.Save(ctx, domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, Result[int64]{
		Code: http.StatusOK,
		Msg:  "success",
		Data: a,
	})
}

func (ah *ArticleHandler) Edit(ctx *gin.Context) {
	type EditArticleRequest struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req EditArticleRequest
	if err := ctx.Bind(&req); err != nil {
		return
	}
	edit, err := ah.svc.Save(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: 0,
		},
	})
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, Result[int64]{
		Code: http.StatusOK,
		Msg:  "success",
		Data: edit,
	})
}

func (ah *ArticleHandler) Publish(ctx *gin.Context) {
	ah.svc.Publish(ctx)
}
