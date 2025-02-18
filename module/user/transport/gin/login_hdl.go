package ginuser

import (
	"net/http"
	"todololist/common"
	"todololist/component/tokenprovider"
	"todololist/module/user/biz"
	"todololist/module/user/model"
	"todololist/module/user/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Login(db *gorm.DB, tokenProvider tokenprovider.Provider) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginUserData model.UserLogin

		if err := c.ShouldBind(&loginUserData); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewSQLStore(db)
		md5 := common.NewMd5Hash()
		business := biz.NewLoginBusiness(store, tokenProvider, md5, 60*60*24*30)

		account, err := business.Login(c.Request.Context(), loginUserData)
		if err != nil {
			panic(err)
		}

		// Trả về token khi đăng nhập thành công
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(account))
	}
}
