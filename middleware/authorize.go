package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"todololist/common"
	"todololist/component/tokenprovider"
	"todololist/module/user/model"

	"github.com/gin-gonic/gin"
)

type AuthenStore interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error)
}

func ErrWrongAuthHeader(err error) *common.AppError {
	return common.NewCustomError(err,
		fmt.Sprintf("wrong authentication header"),
		fmt.Sprint("ErrWrongAuthHeader"),
	)
}

// extractTokenFromHeaderString tách token từ header Authorization
func extractTokenFromHeaderString(s string) (string, error) {
	parts := strings.Split(s, " ")
	// "Authorization": "Bearer {token}"
	if len(parts) < 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
		return "", ErrWrongAuthHeader(nil)
	}
	return parts[1], nil
}

// RequiredAuth
// 1. Get token from header
// 2. Validate token and parse to payload
// 3. From the token payload, we use user_id to find from DB
func RequiredAuth(authStore AuthenStore, tokenProvider tokenprovider.Provider) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Lấy token từ header
		token, err := extractTokenFromHeaderString(c.GetHeader("Authorization"))
		if err != nil {
			panic(err)
		}

		// Xác thực token
		payload, err := tokenProvider.Validate(token)
		if err != nil {
			panic(err)
		}

		// Tìm người dùng dựa trên payload
		user, err := authStore.FindUser(c.Request.Context(), map[string]interface{}{"id": payload.UserId})
		if err != nil {
			panic(err)
		}

		if user.Status == 0 {
			panic(common.ErrNoPermission(errors.New("user has been deleted or banned")))
		}

		c.Set(common.CurrentUser, user)
		c.Next()
	}
}
