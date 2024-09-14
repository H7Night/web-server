package models

import (
	"log"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20); not null" json:"name" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(500); not null" json:"password" validate:"required,min=6,max=120" label:"密码"`
	Role     int    `gorm:"type:int; default:2" json:"role" label:"角色"`
}

// CheckUser 检查
// 只传id进来，更新时查询id是否对应
// 只传name进来，新增时查询是否有同名
func CheckUser(id uint, name string) error {
	var user User

	if id != 0 {
		// 判断用户是否存在，unscoped 可查出已被删除用户和物理删除用户
		err = db.Unscoped().Where("id = ?", id).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("user not exist")
			}
			return errors.Wrap(err, "failed to find user by id")
		}
	} else if name != "" {
		err = db.Where("name = ?", name).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return errors.Wrap(err, "failed to find user by name")
		}
		if user.DeletedAt.Valid {
			return errors.New("user is deleted")
		}
	}

	// 判断时候有重名
	if name != "" && user.Name == name {
		return errors.New("username is already exists")
	}
	return nil
}

// CreateUser 新增
func CreateUser(data *User) error {
	if err := CheckUser(0, data.Name); err != nil {
		return errors.Wrap(err, "user check failed")
	}

	if err := db.Create(&data).Error; err != nil {
		return errors.Wrap(err, "faild to create user")
	}
	return nil
}

// DeleteUser 删除
func DeleteUser(id uint) error {
	if err := CheckUser(id, ""); err != nil {
		return errors.Wrap(err, "user check failed")
	}

	var user User
	if err := db.Select("id").Where("id = ?", id).Find(&user).Error; err != nil {
		return errors.Wrap(err, "failed to find user")
	}

	if err := db.Where("id = ?", id).Delete(&user).Error; err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	return nil
}

// UpdateUser 更新用户
func UpdateUser(id uint, data *User) error {
	if err := CheckUser(id, ""); err != nil {
		return errors.Wrap(err, "user check failed")
	}

	var user User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return errors.Wrap(err, "failed to find user")
	}

	// Update user fields if provided, otherwise use existing data
	if data.Name != "" {
		user.Name = data.Name
	}
	if data.Role != 0 {
		user.Role = data.Role
	}
	if data.Password != "" {
		user.Password = data.Password
	}

	if err := db.Save(&user).Error; err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	return nil
}

// GetUser 查询用户
func GetUser(id int) (User, error) {
	var user User
	if err := db.Limit(1).Where("id = ?", id).Find(&user).Error; err != nil {
		return user, errors.Wrap(err, "failed to get user")
	}
	return user, nil
}

// GetUserPage 获取用户列表
func GetUserPage(username string, pageSize int, pageNum int) ([]User, int64, error) {
	var users []User
	var total int64

	query := db.Select("id, name, role, created_at").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize)

	if username != "" {
		query = query.Where("name LIKE ?", "%"+username+"%")
	}

	if err := query.Find(&users).Count(&total).Error; err != nil {
		return users, 0, errors.Wrap(err, "failed to query users")
	}

	countQuery := db.Model(&User{})
	if username != "" {
		countQuery.Where("name LIKE ?", "%"+username+"%")
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return users, 0, errors.Wrap(err, "failed to count users")
	}

	return users, total, nil
}

// BeforeCreate 创建用户时密码加密；gorm的钩子
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
func CheckLogin(username, password string) (User, string, error) {
	var user User

	if err := db.Where("name = ?", username).First(&user).Error; err != nil {
		return user, "front", errors.Wrap(err, "user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return user, "front", errors.New("incorrect password")
	}

	if user.Role == 1 {
		return user, "back", nil
	}

	return user, "back", errors.New("invalid user role")
}
