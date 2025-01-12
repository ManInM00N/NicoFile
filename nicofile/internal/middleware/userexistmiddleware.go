package middleware

import (
	"encoding/json"
	"gorm.io/gorm"
	config2 "main/config"
	"main/model"
	"net/http"
)

type UserExistMiddleware struct {
}

func NewUserExistMiddleware() *UserExistMiddleware {
	return &UserExistMiddleware{}
}

func (m *UserExistMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		id, _ := r.Context().Value("UserId").(json.Number).Int64()

		err := config2.DB.First(&model.User{Model: gorm.Model{ID: uint(id)}}).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
		next(w, r)
	}
}
