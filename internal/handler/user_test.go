package handler

import (
	"bytes"
	"errors"
	"github.com/ChongYanOvO/little-blue-book/internal/domain"
	"github.com/ChongYanOvO/little-blue-book/internal/service"
	svcmock "github.com/ChongYanOvO/little-blue-book/internal/service/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name          string
		mock          func(ctl *gomock.Controller) service.UserService
		requestMethod string
		requestUrl    string
		requestBody   []byte
		wantCode      int
		wantBody      string
	}{
		{
			name: "注册成功",
			mock: func(ctl *gomock.Controller) service.UserService {
				us := svcmock.NewMockUserService(ctl)
				us.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "test@gmail.com",
					Password: "1qaz@WSX",
				}).Return(nil)
				return us
			},
			requestMethod: http.MethodPost,
			requestUrl:    "/users/signup",
			requestBody: []byte(
				`{
					"email": "test@gmail.com",
					"password": "1qaz@WSX",
					"confirmPassword": "1qaz@WSX"
				}`),
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "参数绑定失败",
			mock: func(ctl *gomock.Controller) service.UserService {
				us := svcmock.NewMockUserService(ctl)
				return us
			},
			requestMethod: http.MethodPost,
			requestUrl:    "/users/signup",
			requestBody: []byte(
				`{
					"email": "test@gmail.com",
					"password": "1qaz@WSX",
				}`),
			wantCode: http.StatusOK,
			wantBody: "参数绑定失败",
		},
		{
			name: "邮箱不正确",
			mock: func(ctl *gomock.Controller) service.UserService {
				us := svcmock.NewMockUserService(ctl)
				return us
			},
			requestMethod: http.MethodPost,
			requestUrl:    "/users/signup",
			requestBody: []byte(
				`{
					"email": "test",
					"password": "1qaz@WSX",
					"confirmPassword": "1qaz@WSX"
				}`),
			wantCode: http.StatusOK,
			wantBody: "邮箱不正确",
		},
		{
			name: "两次输入的密码不一致",
			mock: func(ctl *gomock.Controller) service.UserService {
				us := svcmock.NewMockUserService(ctl)
				return us
			},
			requestMethod: http.MethodPost,
			requestUrl:    "/users/signup",
			requestBody: []byte(
				`{
					"email": "test@gmail.com",
					"password": "1qaz@WSX",
					"confirmPassword": "1qaz@WSX2"
				}`),
			wantCode: http.StatusOK,
			wantBody: "两次输入的密码不一致",
		},
		{
			name: "密码格式错误",
			mock: func(ctl *gomock.Controller) service.UserService {
				us := svcmock.NewMockUserService(ctl)
				return us
			},
			requestMethod: http.MethodPost,
			requestUrl:    "/users/signup",
			requestBody: []byte(
				`{
					"email": "test@gmail.com",
					"password": "123456",
					"confirmPassword": "123456"
				}`),
			wantCode: http.StatusOK,
			wantBody: "密码必须包含数字、特殊字符，并且长度不能小于 8 位",
		},
		{
			name: "邮箱冲突",
			mock: func(ctl *gomock.Controller) service.UserService {
				us := svcmock.NewMockUserService(ctl)
				us.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "test@gmail.com",
					Password: "1qaz@WSX",
				}).Return(service.ErrUserDuplicateEmail)
				return us
			},
			requestMethod: http.MethodPost,
			requestUrl:    "/users/signup",
			requestBody: []byte(
				`{
					"email": "test@gmail.com",
					"password": "1qaz@WSX",
					"confirmPassword": "1qaz@WSX"
				}`),
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
		{
			name: "系统异常",
			mock: func(ctl *gomock.Controller) service.UserService {
				us := svcmock.NewMockUserService(ctl)
				us.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "test@gmail.com",
					Password: "1qaz@WSX",
				}).Return(errors.New("系统异常"))
				return us
			},
			requestMethod: http.MethodPost,
			requestUrl:    "/users/signup",
			requestBody: []byte(
				`{
					"email": "test@gmail.com",
					"password": "1qaz@WSX",
					"confirmPassword": "1qaz@WSX"
				}`),
			wantCode: http.StatusOK,
			wantBody: "系统异常",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			server := gin.Default()
			h := NewUserHandler(tc.mock(ctl), nil, nil)
			h.RegisterRoutes(server)

			request, err := http.NewRequest(tc.requestMethod, tc.requestUrl, bytes.NewBuffer(tc.requestBody))
			request.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)
			assert.Equal(t, http.StatusOK, tc.wantCode)
			assert.Equal(t, tc.wantBody, response.Body.String())
		})

	}
}

func TestUserHandler_Login(t *testing.T) {
	testCases := []struct {
		name          string
		mock          func(ctl *gomock.Controller) service.UserService
		requestMethod string
		requestUrl    string
		requestBody   []byte
		wantCode      int
		wantBody      string
	}{
		{
			name: "登录成功",
			mock: func(ctl *gomock.Controller) service.UserService {
				us := svcmock.NewMockUserService(ctl)
				us.EXPECT().Login(gomock.Any(), "test@gmail.com", "1qaz@WSX").
					Return(domain.User{
						Email:    "test@gmail.com",
						Password: "1qaz@WSX",
					}, nil)
				return us
			},
			requestMethod: http.MethodPost,
			requestUrl:    "/users/login",
			requestBody: []byte(
				`{
					"email": "test@gmail.com",
					"password": "1qaz@WSX"
				}`),
			wantCode: http.StatusOK,
			wantBody: "登录成功",
		},
		{
			name: "参数解析错误",
			mock: func(ctl *gomock.Controller) service.UserService {
				us := svcmock.NewMockUserService(ctl)
				return us
			},
			requestMethod: http.MethodPost,
			requestUrl:    "/users/login",
			requestBody: []byte(
				`{
					"email": "test@gmail.com",
				}`),
			wantCode: http.StatusOK,
			wantBody: "参数解析错误",
		},
		{
			name: "用户名或密码错误",
			mock: func(ctl *gomock.Controller) service.UserService {
				us := svcmock.NewMockUserService(ctl)
				us.EXPECT().Login(gomock.Any(), "test@gmail.com", "1qaz@WSX").
					Return(domain.User{}, service.ErrUserDuplicateEmail)
				return us
			},
			requestMethod: http.MethodPost,
			requestUrl:    "/users/login",
			requestBody: []byte(
				`{
					"email": "test@gmail.com",
					"password": "1qaz@WSX"
				}`),
			wantCode: http.StatusOK,
			wantBody: "用户名或密码错误",
		},
		{
			name: "用户不存在",
			mock: func(ctl *gomock.Controller) service.UserService {
				us := svcmock.NewMockUserService(ctl)
				us.EXPECT().Login(gomock.Any(), "test@gmail.com", "1qaz@WSX").
					Return(domain.User{}, service.ErrUserNotFound)
				return us
			},
			requestMethod: http.MethodPost,
			requestUrl:    "/users/login",
			requestBody: []byte(
				`{
					"email": "test@gmail.com",
					"password": "1qaz@WSX"
				}`),
			wantCode: http.StatusOK,
			wantBody: "系统异常",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			server := gin.Default()
			h := NewUserHandler(tc.mock(ctl), nil, nil)
			h.RegisterRoutes(server)

			request, err := http.NewRequest(tc.requestMethod, tc.requestUrl, bytes.NewBuffer(tc.requestBody))
			request.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)
			assert.Equal(t, http.StatusOK, tc.wantCode)
			assert.Equal(t, tc.wantBody, response.Body.String())
		})

	}
}
