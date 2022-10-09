package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"xstart/app/common"
	"xstart/app/controllers/auth"
	"xstart/app/controllers/config"
	"xstart/app/controllers/menu"
	"xstart/app/controllers/role"
	"xstart/app/controllers/user"
	authServe "xstart/app/services/auth"
)

func Run() {
	r := gin.Default()
	config := cors.DefaultConfig()
	allowOrigin, _ := common.ConfigObj.GetString("web", "allowOrigin")
	config.AllowOrigins = strings.Split(allowOrigin, ",")
	config.AllowHeaders = []string{"token"}
	r.Use(cors.New(config))
	r.Any("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	setRouter(r)

	host, _ := common.ConfigObj.GetString("web", "host")
	port, _ := common.ConfigObj.GetString("web", "port")
	r.Run(host + ":" + port) // 监听并在 0.0.0.0:8080 上启动服务
}

func setRouter(r *gin.Engine) {
	r.POST("/auth/login", auth.Login)
	r.Any("/auth/register", auth.Register)
	r.POST("/auth/logout", auth.Logout)

	t := r.Group("/admin", authAdmin())
	{
		t.GET("/user", user.Member)
		t.GET("/user/password", user.Password)

		t.GET("/user/list", user.List)
		t.POST("/user/create", user.Create)
		t.POST("/user/update", user.Update)
		t.POST("/user/del", user.Del)
		t.GET("/user/one", user.One)

		t.GET("/config/list", config.List)
		t.POST("/config/create", config.Create)
		t.POST("/config/update", config.Update)
		t.POST("/config/del", config.Del)
		t.GET("/config/one", config.One)

		t.GET("/menu/list", menu.List)
		t.POST("/menu/create", menu.Create)
		t.POST("/menu/update", menu.Update)
		t.POST("/menu/del", menu.Del)
		t.GET("/menu/one", menu.One)

		t.GET("/role/list", role.List)
		t.POST("/role/create", role.Create)
		t.POST("/role/update", role.Update)
		t.POST("/role/del", role.Del)
		t.GET("/role/one", role.One)
		t.GET("/role/menus", role.Menus)
		t.POST("/role/menus", role.SaveMenus)

	}
}

func authAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		//t := time.Now()

		var auth authServe.AuthLogin
		tokenStr := c.Request.Header.Get("token")
		// for test
		//if tokenStr == "" {
		//	tokenStr = c.DefaultQuery("token", "")
		//}
		user := auth.CheckToken(tokenStr)
		if user == nil || user["id"] == "" {
			c.JSON(200, gin.H{
				"code": 401, "message": "token check error",
			})
			c.Abort()
		}
		// 请求前
		c.Set("token_user", user)

		c.Next()

		// 请求后
		//latency := time.Since(t)
		//log.Print(latency)

		// 获取发送的 status
		//status := c.Writer.Status()
		//log.Println(status)
	}
}
