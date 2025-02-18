package model

import (
	"errors"
	"strings"
	"todololist/common"
)

var (
	ErrTitleCannotEmpty = errors.New("title cannot be empty")
	ErrItemIsDeleted    = errors.New("item is deleted")
)

const (
	EntityName = "Item"
)

type TodoItem struct {
	common.SQLModel            // embed struct
	UserId          int        `json:"user_id" gorm:"column:user_id;"`
	Title           string     `json:"title" gorm:"column:title;"`
	Description     string     `json:"description" gorm:"column:description;"`
	Status          ItemStatus `json:"status" gorm:"column:status;"`
}

func (TodoItem) TableName() string { return "todo_items" }

type TodoItemCreation struct {
	Id          int    `json:"-" gorm:"column:id;"`
	UserId      int    `json:"-" gorm:"column:user_id;"`
	Title       string `json:"title" gorm:"column:title;"`
	Description string `json:"description" gorm:"column:description;"`
	Status      string `json:"status" gorm:"column:status;"`
}

func (i *TodoItemCreation) Validate() error {
	i.Title = strings.TrimSpace(i.Title)

	if i.Title == "" {
		return ErrTitleCannotEmpty
	}

	return nil
}

func (TodoItemCreation) TableName() string { return TodoItem{}.TableName() }

type TodoItemUpdate struct {
	Title       *string `json:"title" gorm:"column:title;"`
	Description *string `json:"description" gorm:"column:description;"`
	Status      *string `json:"status" gorm:"column:status;"`
}

func (TodoItemUpdate) TableName() string { return TodoItem{}.TableName() }
