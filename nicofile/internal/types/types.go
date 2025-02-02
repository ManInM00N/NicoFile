// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.5

package types

type Auth struct {
	AccessToken  string `json:"access_token"`
	AccessExpire int64  `json:"access_expire"`
}

type AuthResponse struct {
	Message  string `json:"message"`
	Token    string `json:"token"`
	Username string `json:"username"`
}

type CheckChunkRequest struct {
	FileName string   `json:"filename"`
	MD5      []string `json:"md5"`
	ChunkNum int      `json:"chunkNum"`
	FileMd5  string   `json:"fileMd5"`
	Ext      string   `json:"ext"`
}

type CheckChunkResponse struct {
	Error   bool   `json:"error,options=true|false"`
	Accept  int    `json:"accept,range=[0:]"`
	Message string `json:"message,optional"`
}

type CheckResponse struct {
	Error    bool   `json:"error"`
	Message  string `json:"message"`
	Username string `json:"username"`
}

type File struct {
	Id         uint   `json:"id"`
	PosterId   uint   `json:"posterId"`
	PosterName string `json:"posterName"`
	MD5        string `json:"md5"`
	Name       string `json:"name"`
	Ext        string `json:"ext"`
	Path       string `json:"path"`
	Size       int64  `json:"size"`
	Desc       string `json:"desc"`
	CreatedAt  string `json:"createdAt"`
}

type FileDeleteRequest struct {
	FileId int64 `form:"fileid"`
}

type FileDeleteResponse struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

type FileDownloadRequest struct {
	Url string `form:"url"`
}

type FileDownloadResponse struct {
}

type FileListRequest struct {
	Id   int64 `json:"id,optional"`
	Page int   `json:"page,range=[1:]"`
	Size int64 `json:"size,optional"`
}

type FileListResponse struct {
	List    []File `json:"list"`
	Num     int    `json:"num"`
	Error   bool   `json:"error"`
	Message string `json:"message"`
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

type LoginTokenRequest struct {
	Token string `json:"token"`
	Id    int64  `json:"id"`
}

type MergeChunkRequest struct {
	FileName    string `json:"filename"`
	MD5         string `json:"md5,optional"`
	Ext         string `json:"ext"`
	ChunkNum    int    `json:"chunkNum,range=[1:]"`
	Size        int64  `json:"size"`
	Description string `json:"description,optional"`
}

type MergeChunkResponse struct {
	Error   bool   `json:"error,options=true|false"`
	Message string `json:"message,optional"`
}

type NewNameRequest struct {
	Indentity string `json:"indentity,optional"`
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

type UploadChunkRequest struct {
	FileName   string `form:"filename"`
	MD5        string `form:"md5"`
	ChunkIndex int    `form:"chunkIndex"`
}

type UploadChunkResponse struct {
	Error   bool   `json:"error,options=true|false"`
	Message string `json:"message,optional"`
}
