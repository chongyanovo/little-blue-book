package handler

import (
	"github.com/ChongYanOvO/little-blue-book/internal/domain"
	"github.com/ChongYanOvO/little-blue-book/internal/handler/vo"
	"github.com/ChongYanOvO/little-blue-book/internal/service"
	"github.com/ChongYanOvO/little-blue-book/pkg/ginx/jwt"
	"github.com/ChongYanOvO/little-blue-book/pkg/ginx/result"
	"github.com/ChongYanOvO/little-blue-book/pkg/ginx/wrapper"
	"github.com/chongyanovo/zkit/slice"
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
	ag.POST("/save", wrapper.WrapperBodyWitJwt[vo.CreateArticleRequest](ah.logger, ah.Save))
	ag.POST("/edit", wrapper.WrapperBodyWitJwt[vo.EditArticleRequest](ah.logger, ah.Edit))
	ag.POST("/publish", wrapper.WrapperBodyWitJwt[vo.PublishArticleRequest](ah.logger, ah.Publish))
	ag.POST("/list", wrapper.WrapperBodyWitJwt[vo.ListArticleRequest](ah.logger, ah.List))
}

func (ah *ArticleHandler) Save(ctx *gin.Context, req vo.CreateArticleRequest, uc *jwt.UserClaims) (result.Result, error) {
	articleId, err := ah.svc.Save(ctx, &domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		ah.logger.Error("保存文章失败", zap.Error(err))
		return result.FailWithMsg("保存文章失败"), err
	}
	return result.SuccessWithData("保存文章成功", articleId), err
}

func (ah *ArticleHandler) Edit(ctx *gin.Context, req vo.EditArticleRequest, uc *jwt.UserClaims) (result.Result, error) {
	articleId, err := ah.svc.Save(ctx, &domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		return result.FailWithMsg("编辑文章失败"), err
	}
	return result.SuccessWithData("编辑文章成功", articleId), err
}

func (ah *ArticleHandler) Publish(ctx *gin.Context, req vo.PublishArticleRequest, uc *jwt.UserClaims) (result.Result, error) {
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
		return result.FailWithMsg("发布文章失败"), err
	}
	return result.SuccessWithData("发布文章成功", articleId), err
}

func (ah *ArticleHandler) List(ctx *gin.Context, req vo.ListArticleRequest, uc *jwt.UserClaims) (result.Result, error) {
	articles, err := ah.svc.List(ctx, req.Offset, req.Limit)
	if err != nil {
		ah.logger.Error("获取文章列表失败", zap.Error(err))
		return result.FailWithMsg("获取文章列表失败"), err
	}

	articleVos := slice.Map[domain.Article, vo.ArticleVo](articles, func(idx int, src domain.Article) vo.ArticleVo {
		return vo.ArticleVo{
			Id:         src.Id,
			Title:      src.Title,
			Abstract:   src.Abstract(),
			Content:    src.Content,
			AuthorId:   src.Author.Id,
			AuthorName: src.Author.Name,
			Status:     src.Status.ToUint8(),
		}
	})

	return result.SuccessWithData("获取文章列表成功", articleVos), nil
}
