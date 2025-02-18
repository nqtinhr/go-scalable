package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"todololist/common"
)

const EntityName = "User"

type UserRole int

const (
	RoleUser UserRole = 1 << iota
	RoleAdmin
	RoleShipper
	RoleMod
)

// String() giúp UserRole chuyển đổi thành chuỗi
func (role UserRole) String() string {
	switch role {
	case RoleAdmin:
		return "admin"
	case RoleShipper:
		return "shipper"
	case RoleMod:
		return "mod"
	default:
		return "user"
	}
}

// Scan implements sql.Scanner để scan dữ liệu từ DB vào UserRole
func (role *UserRole) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("Failed to unmarshal JSONB value: %v", value))
	}

	var r UserRole
	roleValue := string(bytes)

	switch roleValue {
	case "user":
		r = RoleUser
	case "admin":
		r = RoleAdmin
	case "shipper":
		r = RoleShipper
	case "mod":
		r = RoleMod
	default:
		return errors.New("invalid role type")
	}

	*role = r
	return nil
}

// Value implements driver.Valuer để lưu UserRole vào DB
func (role *UserRole) Value() (driver.Value, error) {
	if role == nil {
		return nil, nil
	}
	return role.String(), nil
}

// MarshalJSON chuyển đổi UserRole thành JSON
func (role *UserRole) MarshalJSON() ([]byte, error) {
	return json.Marshal(role.String())
}

type User struct {
	common.SQLModel
	Email     string   `json:"email" gorm:"column:email;"`
	Password  string   `json:"-" gorm:"column:password;"`
	Salt      string   `json:"-" gorm:"column:salt;"`
	FirstName string   `json:"first_name" gorm:"column:first_name;"`
	LastName  string   `json:"last_name" gorm:"column:last_name;"`
	Phone     string   `json:"phone" gorm:"column:phone;"`
	Role      UserRole `json:"role" gorm:"column:role;"`
	Status    int      `json:"status" gorm:"column:status;"`
}

func (u *User) GetUserId() int {
	return u.Id
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRole() string {
	return u.Role.String()
}

func (User) TableName() string {
	return "users"
}

type UserCreate struct {
	common.SQLModel
	Email     string `json:"email" gorm:"column:email;"`
	Password  string `json:"-" gorm:"column:password;"`
	FirstName string `json:"first_name" gorm:"column:first_name;"`
	LastName  string `json:"last_name" gorm:"column:last_name;"`
	Role      string `json:"role" gorm:"column:role;"`
	Salt      string `json:"-" gorm:"column:salt;"`
}

func (UserCreate) TableName() string {
	return User{}.TableName()
}

type UserLogin struct {
	Email    string `json:"email" form:"email" gorm:"column:email;"`
	Password string `json:"password" form:"password" gorm:"column:password;"`
}

func (UserLogin) TableName() string {
	return User{}.TableName()
}

var (
	ErrEmailOrPasswordInvalid = common.NewCustomError(
		errors.New("email or password invalid"),
		"email or password invalid",
		"ErrUsernameOrPasswordInvalid",
	)

	ErrEmailExisted = common.NewCustomError(
		errors.New("email has already existed"),
		"email has already existed",
		"ErrEmailExisted",
	)
)
