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
	logger         *zap.Logger
	svc            service.ArticleService
	interactiveSvc service.InteractiveService
}

func NewArticleHandler(svc service.ArticleService, interactiveSvc service.InteractiveService, l *zap.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc:            svc,
		interactiveSvc: interactiveSvc,
		logger:         l,
	}
}

func (ah *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	ag := server.Group("/articles")
	ag.POST("/save", wrapper.WrapperBodyWitJwt[vo.CreateArticleRequest](ah.logger, ah.Save))
	ag.POST("/edit", wrapper.WrapperBodyWitJwt[vo.EditArticleRequest](ah.logger, ah.Edit))
	ag.POST("/publish", wrapper.WrapperBodyWitJwt[vo.PublishArticleRequest](ah.logger, ah.Publish))
	ag.POST("/list", wrapper.WrapperBodyWitJwt[vo.ListArticleRequest](ah.logger, ah.List))
	ag.POST("/like", wrapper.WrapperBodyWitJwt[vo.LikeArticleRequest](ah.logger, ah.Like))
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

	go func() {
		if er := ah.interactiveSvc.IncreaseReadCount(ctx, "article", articleId); er != nil {
			ah.logger.Error("浏览量增加失败", zap.Error(er))
		}
	}()

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

func (ah *ArticleHandler) Like(ctx *gin.Context, req vo.LikeArticleRequest, uc *jwt.UserClaims) (result.Result, error) {
	err := ah.interactiveSvc.IncreaseLikeCount(ctx, "article", req.Id, uc.Uid)
	if err != nil {
		ah.logger.Error("点赞失败", zap.Error(err))
		return result.FailWithMsg("点赞失败"), err
	}
	return result.SuccessWithMsg("点赞成功"), nil
}
