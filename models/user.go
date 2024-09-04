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

// CheckUser 检查
// 只传id进来，更新时查询id是否对应
// 只传name进来，新增时查询是否有同名
func CheckUser(id *uint, name *string) int {
	var user User
	var err error

	if id != nil {
		// 判断用户是否存在，unscoped 可查出已被删除用户和物理删除用户
		err = db.Unscoped().Where("id = ?", *id).First(&user).Error
	} else if name != nil {
		err = db.Where("username = ?", *name).First(&user).Error
	}

	// 如果查不到这个用户，或用户的 delete_at 不为空，则证明用户已被删除
	if errors.Is(err, gorm.ErrRecordNotFound) || (id != nil && user.DeletedAt.Valid) {
		return errmsg.ERROR_USER_NOT_EXIST
	} else if err != nil {
		return errmsg.ERROR
	}

	// 判断时候有重名
	if name != nil && user.Username == *name {
		return errmsg.ERROR_USERNAME_USED
	}
	return errmsg.SUCCESS
}

// CreateUser 新增
func CreateUser(data *User) int {
	code := CheckUser(nil, &data.Username)
	if code != errmsg.SUCCESS {
		return errmsg.ERROR_USERNAME_USED
	}
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// DeleteUser 删除
func DeleteUser(id uint) int {
	code := CheckUser(&id, nil)
	if code != errmsg.SUCCESS {
		return code
	}
	var user User
	db.Select("id").Where("id = ?", id).Find(&user)

	err = db.Where("id = ?", id).Delete(&user).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// UpdateUser 更新用户
func UpdateUser(id uint, data *User) int {

	code := CheckUser(&id, nil)
	if code == errmsg.SUCCESS {
		updateData := map[string]interface{}{
			"username": data.Username,
			"role":     data.Role,
		}
		err = db.Model(&User{}).Where("id = ?", id).Updates(updateData).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errmsg.ERROR_USER_NOT_EXIST
			}
			return errmsg.ERROR
		}
		return errmsg.SUCCESS
	} else {
		return code
	}
}

// GetUser 查询用户
func GetUser(id int) (User, int) {
	var user User
	err := db.Limit(1).Where("id = ?", id).Find(&user).Error
	if err != nil {
		return user, errmsg.ERROR
	}
	return user, errmsg.SUCCESS
}

func GetUserPage(username string, pageSize int, pageNum int) ([]User, int64, int) {
	var users []User
	var total int64

	query := db.Select("id,username,role,created_at").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize)

	if username != "" {
		query = query.Where("username like ?", "%"+username+"%")
	}

	if err := query.Find(&users).Count(&total).Error; err != nil {
		return users, 0, errmsg.ERROR
	}

	countQuery := db.Model(&User{})
	if username != "" {
		countQuery.Where("username like ?", "%"+username+"%")
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return users, 0, errmsg.ERROR
	}
	return users, total, errmsg.SUCCESS
}
