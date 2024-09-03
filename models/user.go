package models

import (
	"errors"
	"web-server/utils/errmsg"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint   `gorm:"primary_key; auto_increment" json:"id"`
	Username string `gorm:"type:varchar(20); not null" json:"username" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(500); not null" json:"password" validate:"required,min=6,max=120" label:"密码"`
	Role     int    `gorm:"type:int; default:2" json:"role" validate:"required,gte=2" label:"角色"`
}

func CheckUser(id *uint, name *string) (code int) {
	var user User
	if id != nil {
		db.Where("id = ?", *id).First(&user)
	} else if name != nil {
		db.Where("username = ?", *name).First(&user)
	}
	if name != nil && user.Username == *name {
		return errmsg.ERROR_USERNAME_USED
	} else if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return errmsg.ERROR_USER_NOT_EXIST
	} else if db.Error != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

func CreateUser(data *User) (code int) {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

func DeleteUser(id int) (code int) {
	var user User
	err = db.Where("id = ?", id).Delete(&user).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
