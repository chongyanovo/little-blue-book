package handler

import (
	"github.com/ChongYanOvO/little-blue-book/internal/domain"
	"github.com/ChongYanOvO/little-blue-book/internal/service"
	"github.com/ChongYanOvO/little-blue-book/pkg/ginx"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	ag := server.Group("/articles")
	ag.POST("/save", ginx.WrapperBodyWitJwt[CreateArticleRequest](ah.logger, ah.Save))
	ag.POST("/edit", ginx.WrapperBodyWitJwt[EditArticleRequest](ah.logger, ah.Edit))
	ag.POST("/publish", ginx.WrapperBodyWitJwt[PublishArticleRequest](ah.logger, ah.Publish))
}

type CreateArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (ah *ArticleHandler) Save(ctx *gin.Context, req CreateArticleRequest, uc *ginx.UserClaims) (ginx.Result, error) {
	articleId, err := ah.svc.Save(ctx, &domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		ah.logger.Error("保存文章失败", zap.Error(err))
		return ginx.FailWithMsg("保存文章失败"), err
	}
	return ginx.SuccessWithData("保存文章成功", articleId), err
}

type EditArticleRequest struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (ah *ArticleHandler) Edit(ctx *gin.Context, req EditArticleRequest, uc *ginx.UserClaims) (ginx.Result, error) {
	articleId, err := ah.svc.Save(ctx, &domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		return ginx.FailWithMsg("编辑文章失败"), err
	}
	return ginx.SuccessWithData("编辑文章成功", articleId), err
}

type PublishArticleRequest struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (ah *ArticleHandler) Publish(ctx *gin.Context, req PublishArticleRequest, uc *ginx.UserClaims) (ginx.Result, error) {
	articleId, err := ah.svc.Publish(ctx, &domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		ah.logger.Error("发布文章失败", zap.Error(err))
		return ginx.FailWithMsg("发布文章失败"), err
	}
	return ginx.SuccessWithData("发布文章成功", articleId), err
}
