package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"web-server/models"
	"web-server/utils/errmsg"
	"web-server/utils/validator"
)

// AddUser 添加用户
func AddUser(c *gin.Context) {
	var user models.User
	var msg string
	var validCode int
	_ = c.ShouldBindJSON(&user)

	msg, validCode = validator.Validate(&user)
	if validCode != errmsg.SUCCESS {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  validCode,
			"message": msg,
		})
		c.Abort()
		return
	}
	code := models.CheckUser(user.Username)
	if code != errmsg.SUCCESS {
		models.CreateUser(&user)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := models.DeleteUser(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

func GetUsers(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	username := c.Query("username")

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}
	if pageNum == 0 {
		pageNum = 1
	}
	data, total := models.GetUsers(username, pageSize, pageNum)
	code := errmsg.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}
