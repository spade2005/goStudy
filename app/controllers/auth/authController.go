package auth

import (
	"github.com/gin-gonic/gin"
	authServe "xstart/app/services/auth"
)

func Login(c *gin.Context) {
	var auth authServe.AuthLogin
	if err := c.ShouldBind(&auth); err != nil {
		c.JSON(200, gin.H{
			"code": 1, "message": err.Error(),
		})
		return
	}
	user, message := auth.ValidUser()
	if user == nil || user["id"] == "" {
		c.JSON(200, gin.H{
			"code": 1, "message": message,
		})
		return
	}
	RequestIp := c.ClientIP()
	if RequestIp == "::1" {
		RequestIp = "127.0.0.1"
	}
	auth.Ip = RequestIp

	token, message := auth.GetUserToken(user)
	t := auth.UpdateUserToken(token, true)

	c.JSON(200, gin.H{
		"code": 0, "message": message,
		//"user": user,
		"token": t["token"],
	})
}

func Logout(c *gin.Context) {
	tokenStr := c.Request.Header.Get("token")
	if tokenStr == "" {
		c.JSON(200, gin.H{
			"code": 0, "message": "logout success",
		})
		return
	}
	var auth authServe.AuthLogin
	token := auth.ValidToken(tokenStr)
	if token == nil {
		c.JSON(200, gin.H{
			"code": 0, "message": "logout nil",
		})
		return
	}
	auth.RemoveToken(token)

	c.JSON(200, gin.H{
		"code": 0, "message": "logout success",
	})
}

func Register(c *gin.Context) {
	//test reg for first
	var auth authServe.AuthLogin
	auth.Username = "test4"
	tmp, _ := auth.ValidUser()
	if tmp["id"] != "" {
		c.JSON(200, gin.H{
			"message": "you are Registered", "code": 1,
		})
		return
	}
	data := map[string]string{
		"username": "test4", "password_hash": "test4", "nick_name": "test4",
		"role_id": "0", "status": "0", "phone": "", "email": "",
		"deleted": "0",
	}
	user := auth.CreateUser(data)
	c.JSON(200, gin.H{
		"message": "Register", "code": 0, "data": user,
	})
}
