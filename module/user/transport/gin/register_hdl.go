package ginuser

import (
	"net/http"
	"todololist/common"
	"todololist/module/user/biz"
	"todololist/module/user/model"
	"todololist/module/user/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data model.UserCreate

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		store := storage.NewSQLStore(db)
		md5Hasher := common.NewMd5Hash()
		biz := biz.NewRegisterBusiness(store, md5Hasher)

		if err := biz.Register(c.Request.Context(), &data); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Trả về phản hồi thành công
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
	}
}
