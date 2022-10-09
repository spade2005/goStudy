package user

import (
	"github.com/gin-gonic/gin"
	"strconv"
	authServe "xstart/app/services/auth"
	"xstart/app/services/user"
	authUser "xstart/app/services/user"
)

func Member(c *gin.Context) {
	user, b := c.Get("token_user")
	if !b {
		c.JSON(200, gin.H{"code": 1, "message": "error"})
		c.Abort()
	}
	data := user.(map[string]string)

	var auth authUser.UserCreate
	auth.Id, _ = strconv.Atoi(data["id"])
	one := auth.GetOne()
	one = auth.FormatOne(one)

	c.JSON(200, gin.H{
		"message": "user", "data": one, "code": 0,
	})
}

func Password(c *gin.Context) {
	tokenUser, b := c.Get("token_user")
	if !b {
		c.JSON(200, gin.H{"code": 1, "message": "error"})
		return
	}
	data := tokenUser.(map[string]string)

	oldPass := c.DefaultPostForm("oldPass", "")
	newPass := c.DefaultPostForm("newPass", "")
	newPass2 := c.DefaultPostForm("newPassw", "")
	if oldPass == "" || newPass == "" || newPass2 != newPass {
		c.JSON(200, gin.H{"code": 1, "message": "post param error"})
		return
	}

	var auth authServe.AuthLogin
	auth.Username = data["username"]
	auth.Userpass = data["password_hash"]
	user, _ := auth.ValidUser()
	if user == nil || user["id"] == "" {
		c.JSON(200, gin.H{"code": 1, "message": "old pass error"})
		return
	}
	user["password_hash"] = newPass
	id := auth.UpdatePassword(user)
	if id == 0 {
		c.JSON(200, gin.H{"code": 1, "message": "update error"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "message": "success"})
}

/** crud **/

func List(c *gin.Context) {
	var auth user.UserList
	if err := c.ShouldBind(&auth); err != nil {
		c.JSON(200, gin.H{
			"code": 1, "message": err.Error(),
		})
		return
	}
	cmap, count := auth.GetList()

	cmap = auth.FormatList(cmap)

	c.JSON(200, gin.H{
		"message": "menu",
		"code":    0,
		"data":    gin.H{"list": cmap, "total": count},
	})
}

func One(c *gin.Context) {
	var auth user.UserCreate
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
		"message": "menu",
		"code":    0,
		"data":    one,
	})
}

func Create(c *gin.Context) {
	var auth user.UserCreate
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
	var auth user.UserCreate
	password := c.PostForm("password")
	auth.Password = "123456" //Ignore validation
	if err := c.ShouldBind(&auth); err != nil {
		c.JSON(200, gin.H{
			"code": 1, "message": err.Error(),
		})
		return
	}
	auth.UserName = ""
	config := auth.GetOne()
	if config["id"] == "" {
		c.JSON(200, gin.H{
			"code": 1, "message": "data not found",
		})
		return
	}
	if password != "" {
		if len(password) < 4 || len(password) > 20 {
			c.JSON(200, gin.H{
				"code": 1, "message": "Invalid password format",
			})
			return
		}
	}

	auth.Password = password
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
	var auth user.UserCreate
	id := c.PostForm("id")
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
