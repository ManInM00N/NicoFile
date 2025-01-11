// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.5

package types

type AuthResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type File struct {
	Id                 int64  `json:"id"`
	Identity           string `json:"identity"`
	RepositoryIdentity string `json:"repository_identity"`
	Name               string `json:"name"`
	Ext                string `json:"ext"`
	Path               string `json:"path"`
	Size               int64  `json:"size"`
}

type FileCutUploadRequest struct {
	Hash string `json:"hash,optional"`
	Name string `json:"name,optional"`
	Ext  string `json:"ext,optional"`
	Size int64  `json:"size,optional"`
	Path string `json:"path,optional"`
}

type FileCutUploadResponse struct {
	Identity string `json:"identity"`
	Ext      string `json:"ext"`
	Name     string `json:"name"`
}

type FileDeleteRequest struct {
	Identity string `json:"identity"`
}

type FileDeleteResponse struct {
}

type FileDownloadRequest struct {
}

type FileDownloadResponse struct {
}

type FileListRequest struct {
	Id   int64 `json:"id,optional"`
	Page int   `json:"page,optional"`
	Size int   `json:"size,optional"`
}

type FileListResponse struct {
	List  []*File `json:"list"`
	Count int64   `json:"count"`
}

type FileMeta struct {
	Id           int64  `json:"id"`
	FileName     string `json:"filename"`
	FilePath     string `json:"filepath"`
	FileSize     int64  `json:"filesize"`
	UploadedSize int64  `json:"uploadedsize"`
	Status       string `json:"status"`
	CreatedAt    string `json:"createdat"`
}

type FileUploadRequest struct {
}

type FileUploadResponse struct {
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NewNameRequest struct {
	Indentity string `json:"indentity"`
	NewName   string `json:"newName"`
}

type NewNameResponse struct {
	Message string `json:"message"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Request struct {
	Name string `path:"name,options=you|me"`
}

type Response struct {
	Message string `json:"message"`
}
