package api

import (
	"net/http"
	"strconv"
	"web-server/models"
	"web-server/utils/errmsg"
	"web-server/utils/validator"

	"github.com/gin-gonic/gin"
)

// AddUser 添加用户
func AddUser(c *gin.Context) {
	var data models.User
	var msg string
	var validCode int
	_ = c.ShouldBindJSON(&data)

	msg, validCode = validator.Validate(&data)
	if validCode != errmsg.Success {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  validCode,
				"message": msg,
			},
		)
		c.Abort()
		return
	}
	code := models.CheckUser(0, data.Username)
	if code == errmsg.Success {
		models.CreateUser(&data)
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	code := models.DeleteUser(uint(id))

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// UpdateUser 修改用户
func UpdateUser(c *gin.Context) {
	var data models.User
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)

	code := models.CheckUser(uint(id), "")
	if code == errmsg.Success {
		code = models.UpdateUser(uint(id), &data)
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// GetUser 查询用户
func GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var maps = make(map[string]interface{})
	data, code := models.GetUser(id)

	maps["id"] = data.ID
	maps["username"] = data.Username
	maps["role"] = data.Role
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    maps,
			"total":   1,
			"message": errmsg.GetErrMsg(code),
		})
}

// GetUserPage 获取用户列表
func GetUserPage(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	username := c.Query("username")

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageNum <= 0:
		pageNum = 10
	case pageNum == 0:
		pageNum = 1
	}
	data, total, code := models.GetUserPage(username, pageSize, pageNum)

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"total":   total,
			"message": errmsg.GetErrMsg(code),
		},
	)
}
