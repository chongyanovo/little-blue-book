package handler

import (
	"errors"
	"github.com/ChongYanOvO/little-blue-book/internal/domain"
	"github.com/ChongYanOvO/little-blue-book/internal/service"
	"github.com/ChongYanOvO/little-blue-book/pkg/ginx/jwt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var _ Handler = (*UserHandler)(nil)

type UserHandler struct {
	svc            service.UserService
	codeSvc        service.CodeService
	emailRegExp    *regexp.Regexp
	passwordRegExp *regexp.Regexp
	logger         *zap.Logger
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService, l *zap.Logger) *UserHandler {
	return &UserHandler{
		svc:            svc,
		codeSvc:        codeSvc,
		emailRegExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		logger:         l,
	}
}

func (uh *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", uh.SignUp)
	ug.POST("/login", uh.Login)
	ug.POST("/edit", uh.Edit)
	ug.GET("/profile", uh.Profile)
	ug.PUT("/login/code", uh.SendLoginSmsCode)
	ug.POST("/login/code", uh.LoginSms)
}

func (uh *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq

	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "参数绑定失败")
		return
	}

	isEmail, err := uh.emailRegExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "邮箱不正确")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	isPassword, err := uh.passwordRegExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含数字、特殊字符，并且长度不能小于 8 位")
		return
	}

	// 调用一些 svc 的方法
	err = uh.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
}

// Login 用户登录接口
func (uh *UserHandler) Login(ctx *gin.Context) {
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.String(http.StatusOK, "参数解析错误")
		return
	}
	user, err := uh.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserDuplicateEmail) {
			ctx.String(http.StatusOK, "用户名或密码错误")
		} else {
			ctx.String(http.StatusOK, "系统异常")
		}
	} else {
		if err := jwt.SetJwtToken(ctx, user.Id, user.Email); err != nil {
			uh.logger.Error("jwt设置错误")
			return
		}
		ctx.String(http.StatusOK, "登录成功")
	}

}

func (uh *UserHandler) Edit(ctx *gin.Context) {

}

// Profile 用户详情
func (uh *UserHandler) Profile(ctx *gin.Context) {
	uc, err := jwt.ExtractJwtClaims(ctx)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
	}
	user, err := uh.svc.Profile(ctx, uc.Uid)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
	}
	ctx.JSON(http.StatusOK, user)
}

// SendLoginSmsCode 登录验证码发送
func (uh *UserHandler) SendLoginSmsCode(ctx *gin.Context) {
	const biz = "login"
	type Request struct {
		Phone string `json:"phone"`
	}
	var req Request
	if err := ctx.Bind(&req); err != nil {
		return
	}
	err := uh.codeSvc.Send(ctx, biz, req.Phone)
	if err != nil {
		ctx.String(http.StatusOK, "登录验证码发送失败")
	} else {
		ctx.String(http.StatusOK, "登录验证码发送成功")
	}

}

// LoginSms 登录验证码校验
func (uh *UserHandler) LoginSms(ctx *gin.Context) {
	const biz = "login"
	type Request struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Request
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := uh.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if !ok {
		ctx.String(http.StatusOK, "验证码错误")
	} else if err != nil {
		ctx.String(http.StatusOK, "系统异常")
	} else {
		u, err := uh.svc.FindOrCreate(ctx, req.Phone)
		if err != nil {
			ctx.String(http.StatusOK, "系统异常")
			return
		}
		if err := jwt.SetJwtToken(ctx, u.Id, u.Email); err != nil {
			uh.logger.Error("jwt设置错误")
		}
		ctx.String(http.StatusOK, "登录成功")
	}

}
