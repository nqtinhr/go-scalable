package biz

import (
	"context"
	"todololist/common"
	"todololist/component/tokenprovider"
	"todololist/module/user/model"
)

type LoginStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error)
}

type loginBusiness struct {
	storeUser     LoginStorage
	tokenProvider tokenprovider.Provider
	hasher        Hasher
	expiry        int
}

func NewLoginBusiness(storeUser LoginStorage, tokenProvider tokenprovider.Provider, hasher Hasher, expiry int) *loginBusiness {
	return &loginBusiness{
		storeUser:     storeUser,
		tokenProvider: tokenProvider,
		hasher:        hasher,
		expiry:        expiry,
	}
}

// 1. Find user, email
// 2. Hash pass from input and compare with pass in db
// 3. Provider: issue JWT token for client
// 3.1. Access token and refresh token
// 4. Return token(s)

func (business *loginBusiness) Login(ctx context.Context, data model.UserLogin) (tokenprovider.Token, error) {
	user, err := business.storeUser.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if err != nil {
		return nil, model.ErrEmailOrPasswordInvalid
	}

	// Hash mật khẩu nhập vào và so sánh với mật khẩu trong database
	passHashed := business.hasher.Hash(data.Password + user.Salt)
	if user.Password != passHashed {
		return nil, model.ErrEmailOrPasswordInvalid
	}

	// Tạo payload cho token
	payload := &common.TokenPayload{
		UId:   user.Id,
		URole: user.Role.String(),
	}

	// Tạo access token
	accessToken, err := business.tokenProvider.Generate(payload, business.expiry)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	// refreshToken, err := business.tokenProvider.Generate(payload, business.tkCfg.GetRtExp())
	// if err != nil {
	//     return nil, common.ErrInternal(err)
	// }

	// Trả về token
	// account := usermodel.NewAccount(accessToken, refreshToken)
	return accessToken, nil
}
