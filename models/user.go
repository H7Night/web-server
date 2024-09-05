package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"web-server/utils/errmsg"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20); not null" json:"name" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(500); not null" json:"password" validate:"required,min=6,max=120" label:"密码"`
	Role     int    `gorm:"type:int; default:2" json:"role" default:"2" label:"角色"`
}

// CheckUser 检查
// 只传id进来，更新时查询id是否对应
// 只传name进来，新增时查询是否有同名
func CheckUser(id uint, name string) int {
	var user User
	var err error

	if id != 0 {
		// 判断用户是否存在，unscoped 可查出已被删除用户和物理删除用户
		err = db.Unscoped().Where("id = ?", id).First(&user).Error
	} else if name != "" {
		err = db.Where("name = ?", name).First(&user).Error
	}

	// 如果查不到这个用户，或用户的 delete_at 不为空，则证明用户已被删除
	if (errors.Is(err, gorm.ErrRecordNotFound) || (id != 0 && user.DeletedAt.Valid)) && name == "" {
		return errmsg.ErrorUserNotExist
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errmsg.Error
	}

	// 判断时候有重名
	if name != "" && user.Name == name {
		return errmsg.ErrorUsernameUsed
	}
	return errmsg.Success
}

// CreateUser 新增
func CreateUser(data *User) int {
	code := CheckUser(0, data.Name)
	if code != errmsg.Success {
		return errmsg.ErrorUsernameUsed
	}
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.Error
	}
	return errmsg.Success
}

// DeleteUser 删除
func DeleteUser(id uint) int {
	code := CheckUser(id, "")
	if code != errmsg.Success {
		return code
	}
	var user User
	db.Select("id").Where("id = ?", id).Find(&user)

	err = db.Where("id = ?", id).Delete(&user).Error
	if err != nil {
		return errmsg.Error
	}
	return errmsg.Success
}

// UpdateUser 更新用户
func UpdateUser(id uint, data *User) int {
	// 检查用户是否存在
	code := CheckUser(id, "")
	if code != errmsg.Success {
		return code
	}

	var user User
	db.Where("id = ?", id).First(&user)
	//没更新则用原数据
	if data.Name == "" {
		data.Name = user.Name
	}
	if data.Role == 0 {
		data.Role = user.Role
	}
	if data.Password == "" {
		data.Password = user.Password
	}

	// 更新用户数据
	user.Name = data.Name
	user.Role = data.Role
	user.Password = data.Password

	// 使用 Save 而不是 Updates 确保触发钩子
	err = db.Save(&user).Error
	if err != nil {
		return errmsg.Error
	}
	return errmsg.Success

}

// GetUser 查询用户
func GetUser(id int) (User, int) {
	var user User
	err := db.Limit(1).Where("id = ?", id).Find(&user).Error
	if err != nil {
		return user, errmsg.Error
	}
	return user, errmsg.Success
}

// GetUserPage 获取用户列表
func GetUserPage(username string, pageSize int, pageNum int) ([]User, int64, int) {
	var users []User
	var total int64

	query := db.Select("id,name,role,created_at").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize)

	if username != "" {
		query = query.Where("name like ?", "%"+username+"%")
	}

	if err := query.Find(&users).Count(&total).Error; err != nil {
		return users, 0, errmsg.Error
	}

	countQuery := db.Model(&User{})
	if username != "" {
		countQuery.Where("name like ?", "%"+username+"%")
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return users, 0, errmsg.Error
	}
	return users, total, errmsg.Success
}

// BeforeCreate 密码加密；gorm的钩子
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	u.Password = ScryptPw(u.Password)
	return nil
}

// BeforeUpdate 更新密码时加密钩子
func (u *User) BeforeUpdate(_ *gorm.DB) (err error) {
	u.Password = ScryptPw(u.Password)
	return nil
}

// ScryptPw 生成密码
func ScryptPw(password string) string {
	const cost = 10

	HashPw, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Fatal(err)
	}

	return string(HashPw)
}

// CheckLogin 登录校验
func CheckLogin(username string, password string, state string) (User, int) {
	var user User
	pwdErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	err := db.Where("name =?", username).First(&user).Error

	if pwdErr != nil && !errors.Is(pwdErr, gorm.ErrRecordNotFound) {
		return user, errmsg.ErrorPasswordWrong
	}
	if err != nil {
		return user, errmsg.ErrorUserNotExist
	}
	// 后台登录校验用户角色
	if state == "back" && user.Role != 1 {
		return user, errmsg.ErrorUserNoRight
	}
	return user, errmsg.Success
}
