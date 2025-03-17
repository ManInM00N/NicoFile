package article

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"main/nicofile/internal/logic/article"
	"main/nicofile/internal/svc"
	"main/nicofile/internal/types"
)

func ArticleUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ArticleUpdateRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := article.NewArticleUpdateLogic(r.Context(), svcCtx)
		resp, err := l.ArticleUpdate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
