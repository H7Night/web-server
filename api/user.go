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
	if validCode != errmsg.SUCCESS {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  validCode,
				"message": msg,
			},
		)
		c.Abort()
		return
	}
	code := models.CheckUser(nil, &data.Username)
	if code == errmsg.SUCCESS {
		models.CreateUser(&data)
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	code := models.DeleteUser(id)

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}
