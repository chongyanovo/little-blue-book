package accesslog

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/atomic"
	"io"
	"time"
)

// AccessLog 请求日志
type AccessLog struct {
	//http 请求类型
	Method string
	//url 整个请求的url
	Url string
	//请求体
	RequestBody string
	//响应体
	ResponseBody string
	//处理时间
	Duration string
	//状态码
	Status int
}

// MiddleWareBuilder 日志打印中间件
type MiddleWareBuilder struct {
	allowRequestBody  *atomic.Bool
	allowResponseBody *atomic.Bool
	maxLength         *atomic.Int64
	loggerFunc        func(ctx context.Context, log *AccessLog)
}

// NewBuilder 创建日志打印中间件构造器
func NewBuilder(loggerFunc func(ctx context.Context, log *AccessLog)) *MiddleWareBuilder {
	return &MiddleWareBuilder{
		allowRequestBody:  atomic.NewBool(false),
		allowResponseBody: atomic.NewBool(false),
		maxLength:         atomic.NewInt64(1024),
		loggerFunc:        loggerFunc,
	}
}

// AllowRequestBody 是否打印请求体
func (b *MiddleWareBuilder) AllowRequestBody(flag bool) *MiddleWareBuilder {
	b.allowRequestBody.Store(flag)
	return b
}

// AllowResponseBody 是否打印响应体
func (b *MiddleWareBuilder) AllowResponseBody(flag bool) *MiddleWareBuilder {
	b.allowResponseBody.Store(flag)
	return b
}

// MaxLength 打印的最大长度
func (b *MiddleWareBuilder) MaxLength(length int) *MiddleWareBuilder {
	b.maxLength.Store(int64(length))
	return b
}

// Build 构建中间件
func (b *MiddleWareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			//请求处理开始时间
			start = time.Now()
			//url
			url = ctx.Request.URL.String()
			//url 长度
			urlLen = int64(len(url))
			//运行打印的最大长度
			maxLength = b.maxLength.Load()
			//是否打印请求体
			allowRequestBody = b.allowRequestBody.Load()
			//是否打印响应体
			allowResponseBody = b.allowResponseBody.Load()
		)
		if urlLen >= maxLength {
			url = url[:maxLength]
		}

		log := &AccessLog{
			Method: ctx.Request.Method,
			Url:    url,
		}
		if allowRequestBody && ctx.Request.Body != nil {
			body, _ := ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
			if int64(len(body)) >= maxLength {
				body = body[:maxLength]
			}
			//注意资源的消耗
			log.RequestBody = string(body)
		}

		if allowResponseBody {
			ctx.Writer = responseWriter{
				ResponseWriter: ctx.Writer,
				accessLog:      log,
				maxLength:      maxLength,
			}
		}

		defer func() {
			log.Duration = time.Since(start).String()
			//日志打印
			b.loggerFunc(ctx, log)
		}()
		ctx.Next()
	}
}

// responseWriter 装饰gin.ResponseWriter获取响应体
type responseWriter struct {
	gin.ResponseWriter
	accessLog *AccessLog
	maxLength int64
}

// WriteHeader 获取状态码
func (r responseWriter) WriteHeader(statusCode int) {
	r.accessLog.Status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Write 获取响应体
func (r responseWriter) Write(data []byte) (int, error) {
	curLen := int64(len(data))
	if curLen >= r.maxLength {
		data = data[:r.maxLength]
	}
	r.accessLog.ResponseBody = string(data)
	return r.ResponseWriter.Write(data)
}
