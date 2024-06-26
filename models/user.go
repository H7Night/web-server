package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"web-server/utils/errmsg"
)

type User struct {
	gorm.Model
	ID       int8   `gorm:"primary_key;auto_increment" json:"id"`
	Username string `gorm:"type:varchar(20);not null" json:"username" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(500);not null" json:"password" validate:"required,min=6,max=120" label:"密码"`
	Role     int    `gorm:"type:int;DEFAULT:2" json:"role" validate:"required,gte=2" label:"角色码"`
}

func CheckLogin(username string, password string) (User, int) {
	var user User
	var PasswordErr error

	db.Where("username = ?", username).First(&user)

	PasswordErr = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if user.ID == 0 {
		return user, errmsg.ERROR_USER_NOT_EXIST
	}
	if PasswordErr != nil {
		return user, errmsg.ERROR_PASSWORD_WRONG
	}
	if user.Role != 1 {
		return user, errmsg.ERROR_USER_NO_RIGHT
	}
	return user, errmsg.SUCCESS
}

func CheckLoginFront(username string, password string) (User, int) {
	var user User
	var PasswordErr error

	db.Where("username = ?", username).First(&user)

	PasswordErr = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if user.ID == 0 {
		return user, errmsg.ERROR_USER_NOT_EXIST
	}
	if PasswordErr != nil {
		return user, errmsg.ERROR_PASSWORD_WRONG
	}
	return user, errmsg.SUCCESS
}

func GetUsers(userName string, pageSize int, pageNum int) ([]User, int64) {
	var Users []User
	var total int64

	if userName != "" {
		db.Select("id, username, role").Where(
			"username LIKE ?", userName+"%",
		).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&Users)
		db.Model(&Users).Where(
			"username LIKE ?", userName+"%",
		).Count(&total)
		return Users, total
	}
	db.Select("id, username, role").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&Users)
	db.Model(&Users).Count(&total)
	return Users, total
}
