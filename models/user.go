package models

import (
	"goIland/utils/errmsg"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int8
	Username string `gorm:"type:varchar(20);not null " json:"username" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(500);not null" json:"password" validate:"required,min=6,max=120" label:"密码"`
	Role     int    `gorm:"type:int;DEFAULT:2" json:"role" validate:"required,gte=2" label:"角色码"`
}

func CheckLogin(username string, password string) (User, int) {
	var user User
	var PasswordErr error

	db.Where("username = ?", username).First(&user)

	PasswordErr = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if user.ID == 0 {
		return user, errmsg.ERROR
	}
	if PasswordErr != nil {
		return user, errmsg.ERROR
	}
	if user.Role != 1 {
		return user, errmsg.ERROR
	}
	return user, errmsg.SUCCESS
}

func CheckLoginFront(username string, password string) (User, int) {
	var user User
	var PasswordErr error

	db.Where("username = ?", username).First(&user)

	PasswordErr = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if user.ID == 0 {
		return user, errmsg.ERROR
	}
	if PasswordErr != nil {
		return user, errmsg.ERROR
	}
	return user, errmsg.SUCCESS
}

func GetUsers(userName string, pageSize int, pageNum int) ([]User, int64) {
	var Users []User
	var total int64

	if userName != "" {
		db.Select("id, name, branch").Where(
			"name LIKE ?", userName+"%",
		).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&Users)
		db.Model(&Users).Where(
			"name LIKE ?", userName+"%",
		).Count(&total)
		return Users, total
	}
	db.Select("id,name").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&Users)
	db.Model(&Users).Count(&total)
	return Users, total
}
