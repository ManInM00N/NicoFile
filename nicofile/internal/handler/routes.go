// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.5

package handler

import (
	"net/http"
	"time"

	file "main/nicofile/internal/handler/file"
	user "main/nicofile/internal/handler/user"
	"main/nicofile/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/from/:name",
				Handler: NicofileHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.UserExistMiddleware},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/file/checkchunk",
					Handler: file.CheckChunkHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/file/delete",
					Handler: file.FileDeleteHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/file/download",
					Handler: file.FileDownloadHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/file/list",
					Handler: file.FileListHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/file/mergechunk",
					Handler: file.MergeChunkHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/file/upload",
					Handler: file.FileUploadHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/file/uploadchunk",
					Handler: file.UploadChunkHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/user/login",
				Handler: user.UserLoginHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/user/register",
				Handler: user.UserRegisterHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
		rest.WithTimeout(3000*time.Millisecond),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.UserExistMiddleware},
			[]rest.Route{
				{
					Method:  http.MethodDelete,
					Path:    "/user/delete",
					Handler: user.DeleteUserHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/user/loadtoken",
					Handler: user.UserLoginTokenHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/user/newname",
					Handler: user.UserChangeNameHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/user/newpassword",
					Handler: user.UserChangePasswordHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/api/v1"),
		rest.WithTimeout(3000*time.Millisecond),
	)
}
