package config

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xstart/app/services/config"
)

func List(c *gin.Context) {
	var auth config.ConfigList
	if err := c.ShouldBind(&auth); err != nil {
		c.JSON(200, gin.H{
			"code": 1, "message": err.Error(),
		})
		return
	}
	cmap, count := auth.GetList()

	c.JSON(200, gin.H{
		"message": "user",
		"code":    0,
		"data":    gin.H{"list": cmap, "total": count},
	})
}

func One(c *gin.Context) {
	var auth config.ConfigCreate
	id := c.Query("id")
	var err error
	auth.Id, err = strconv.Atoi(id)
	if id == "" || err != nil || auth.Id <= 0 {
		c.JSON(200, gin.H{
			"code": 1, "message": "request id is error",
		})
		return
	}
	one := auth.GetOne()
	if one["id"] == "" {
		c.JSON(200, gin.H{
			"code": 1, "message": "data not found",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "user",
		"code":    0,
		"data":    one,
	})
}

func Create(c *gin.Context) {
	var auth config.ConfigCreate
	if err := c.ShouldBind(&auth); err != nil {
		c.JSON(200, gin.H{
			"code": 1, "message": err.Error(),
		})
		return
	}
	auth.Id = 0
	repeat := auth.GetOne()
	if repeat["id"] != "" {
		c.JSON(200, gin.H{
			"code": 1, "message": "Duplicate data is not allowed",
		})
		return
	}

	cmap := auth.DoCreate()
	if cmap > 0 {
		c.JSON(200, gin.H{
			"message": "success", "code": 0,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "error", "code": 1,
	})
}

func Update(c *gin.Context) {
	var auth config.ConfigCreate
	if err := c.ShouldBind(&auth); err != nil {
		c.JSON(200, gin.H{
			"code": 1, "message": err.Error(),
		})
		return
	}
	name := auth.KeyStr
	auth.KeyStr = ""
	config := auth.GetOne()
	if config["id"] == "" {
		c.JSON(200, gin.H{
			"code": 1, "message": "data not found",
		})
		return
	}
	auth.KeyStr = name
	cmap := auth.DoUpdate()
	if cmap > 0 {
		c.JSON(200, gin.H{
			"message": "success", "code": 0,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "error", "code": 1,
	})
}

func Del(c *gin.Context) {
	var auth config.ConfigCreate
	auth.KeyStr = "del key str"
	auth.ValueStr = "del value str"
	if err := c.ShouldBind(&auth); err != nil {
		c.JSON(200, gin.H{
			"code": 1, "message": err.Error(),
		})
		return
	}
	//cannot be modified
	auth.KeyStr = ""
	config := auth.GetOne()
	if config["id"] == "" {
		c.JSON(200, gin.H{
			"code": 1, "message": "data not found",
		})
		return
	}
	cmap := auth.DoDel()
	if cmap > 0 {
		c.JSON(200, gin.H{
			"message": "success", "code": 0,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "error", "code": 1,
	})
}
