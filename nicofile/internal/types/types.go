// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.5

package types

type Article struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreatedAt  string `json:"createdat"`
	View       int64  `json:"view"`
	Like       int64  `json:"like"`
	AuthorId   int64  `json:"authorid"`
	AuthorName string `json:"authorname"`
	Cover      string `json:"cover"`
}

type ArticleCreateRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Cover   string `json:"cover"`
}

type ArticleDeleteRequest struct {
	Id int64 `path:"id"`
}

type ArticleDeleteResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type ArticleLikeRequest struct {
	Id int64 `json:"id"`
}

type ArticleListRequest struct {
	Page int   `json:"page,range=[1:]"`
	Size int64 `json:"size,optional"`
}

type ArticleListResponse struct {
	List     []Article `json:"list"`
	Num      int       `json:"num"`
	Error    bool      `json:"error"`
	Message  string    `json:"message"`
	AllPages int       `json:"allpages"`
	Page     int       `json:"page"`
}

type ArticleRequest struct {
	Id int64 `path:"id"`
}

type ArticleResponse struct {
	Error      bool   `json:"error"`
	Message    string `json:"message"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreatedAt  string `json:"createdat"`
	Cover      string `json:"cover"`
	View       string `json:"view"`
	Like       string `json:"like"`
	ArticleId  string `json:"articleid"`
	AuthorId   string `json:"authorid"`
	AuthorName string `json:"authorname"`
}

type ArticleSearchRequest struct {
	Keyword string `json:"keyword"`
	Page    int    `json:"page,range=[1:]"`
	Size    int64  `json:"size,optional"`
}

type ArticleUpdateRequest struct {
	Id      int64  `path:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Auth struct {
	AccessToken  string `json:"access_token"`
	AccessExpire int64  `json:"access_expire"`
}

type AuthResponse struct {
	Error    bool   `json:"error"`
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

type Comment struct {
	Id        int64     `json:"id"`
	IP        string    `json:"ip"`
	ArticleId int64     `json:"articleid"`
	Content   string    `json:"content"`
	CreatedAt string    `json:"createdat"`
	AuthorId  int64     `json:"authorid"`
	Author    string    `json:"author"`
	ParentId  int64     `json:"parentid"`
	Comments  []Comment `json:"comments"`
}

type CommentCreateRequest struct {
	ArticleId int64  `json:"articleid"`
	ParentId  int64  `json:"parentid,optional"`
	Content   string `json:"content"`
}

type CommentListRequest struct {
	ArticleId int64 `form:"articleid"`
	Page      int   `json:"page,range=[1:],optional"`
	Size      int64 `json:"size,optional"`
}

type CommentListResponse struct {
	List     []Comment `json:"list"`
	Num      int       `json:"num"`
	Error    bool      `json:"error"`
	Message  string    `json:"message"`
	AllPages int       `json:"allpages"`
	Page     int       `json:"page"`
}

type CommentResponse struct {
	Error    bool   `json:"error"`
	Message  string `json:"message"`
	Comment  string `json:"comment"`
	AuthorId string `json:"authorid"`
	Author   string `json:"author"`
}

type DeleteUserRequest struct {
	Userid int64 `form:"userid"`
}

type DeleteUserResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type DownloadIMGRequest struct {
	ID string `path:"id"` // 图片ID
}

type DownloadIMGResponse struct {
}

type File struct {
	Id             uint   `json:"id"`
	PosterId       uint   `json:"posterId"`
	PosterName     string `json:"posterName"`
	MD5            string `json:"md5"`
	Name           string `json:"name"`
	Ext            string `json:"ext"`
	Path           string `json:"path"`
	Size           int64  `json:"size"`
	Desc           string `json:"desc"`
	DonwloadCounts int64  `json:"downloadcounts"`
	CreatedAt      string `json:"createdAt"`
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
	List     []File `json:"list"`
	Num      int    `json:"num"`
	Error    bool   `json:"error"`
	Message  string `json:"message"`
	AllPages int    `json:"allpages"`
	Page     int    `json:"page"`
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
	NewName string `json:"newName"`
}

type NewNameResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type NewPasswordRequest struct {
	NewPassword string `json:"newPassword"`
}

type NewPasswordResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type Null struct {
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
	Ext        string `form:"ext"`
}

type UploadChunkResponse struct {
	Error   bool   `json:"error,options=true|false"`
	Message string `json:"message,optional"`
}

type UploadIMGRequest struct {
}

type UploadIMGResponse struct {
	URL string `json:"url"` // 图片访问URL
}
