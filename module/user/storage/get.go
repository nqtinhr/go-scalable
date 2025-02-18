package storage

import (
	"context"
	"errors"
	"todololist/common"
	"todololist/module/user/model"

	"gorm.io/gorm"
)

func (s *sqlStore) FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error) {
	db := s.db.Table(model.User{}.TableName())

	// Preload các thông tin bổ sung nếu có
	for _, info := range moreInfo {
		db = db.Preload(info)
	}

	var user model.User

	// Tìm user theo điều kiện
	if err := db.Where(conditions).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}

	return &user, nil
}
