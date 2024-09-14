package api

import (
	"net/http"
	"strconv"
	"web-server/models"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// AddUser 添加用户
func AddUser(c *gin.Context) {
	var data models.User
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.Wrap(err, "invalid request data"),
		})
		return
	}

	if err := models.CheckUser(0, data.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.Wrap(err, "user check failed"),
		})
		return
	}

	err := models.CreateUser(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.Wrap(err, "user creation failed"),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "success",
		"data":    data,
	})
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.Wrap(err, "invalid id"),
		})
		return
	}

	if err := models.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.Wrap(err, "failed to delete user"),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User deleted successfully",
	})
}

// UpdateUser 修改用户
func UpdateUser(c *gin.Context) {
	var data models.User
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.Wrap(err, "invalid id"),
		})
		return
	}

	if err := models.CheckUser(uint(id), ""); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.Wrap(err, "user check failed"),
		})
		return
	}

	if err := models.UpdateUser(uint(id), &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.Wrap(err, "failed to update user"),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User updated successfully",
	})
}

// GetUser 查询用户
func GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.Wrap(err, "invalid id"),
		})
		return
	}
	data, err := models.GetUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.Wrap(err, "failed to get user"),
		})
		return
	}
	responseData := map[string]interface{}{
		"id":   data.ID,
		"name": data.Name,
		"role": data.Role,
	}
	c.JSON(http.StatusOK, responseData)
}

// GetUserPage 获取用户列表
func GetUserPage(c *gin.Context) {
	pageSizeStr := c.Query("pagesize")
	pageNumStr := c.Query("pagenum")
	username := c.Query("username")

	var pageSize, pageNum int
	var err error

	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 0 {
			pageSize = 10 // 最小10
		}
		if pageSize > 100 {
			pageSize = 100 // 最大100
		}
	} else {
		pageSize = 10 // 默认10
	}

	if pageNumStr != "" {
		pageNum, err = strconv.Atoi(pageNumStr)
		if err != nil || pageNum <= 0 {
			pageNum = 1 // 最小1
		}
	} else {
		pageNum = 1 // 默认1
	}

	data, total, err := models.GetUserPage(username, pageSize, pageNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.Wrap(err, "failed to get user page"),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    data,
		"total":   total,
		"message": "success",
	})
}
