package role

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xstart/app/services/role"
)

func List(c *gin.Context) {
	var auth role.RoleList
	if err := c.ShouldBind(&auth); err != nil {
		c.JSON(200, gin.H{
			"code": 1, "message": err.Error(),
		})
		return
	}
	cmap, count := auth.GetList()

	c.JSON(200, gin.H{
		"message": "success",
		"code":    0,
		"data":    gin.H{"list": cmap, "total": count},
	})
}

func One(c *gin.Context) {
	var auth role.RoleCreate
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
	var auth role.RoleCreate
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
	var auth role.RoleCreate
	if err := c.ShouldBind(&auth); err != nil {
		c.JSON(200, gin.H{
			"code": 1, "message": err.Error(),
		})
		return
	}
	name := auth.Name
	auth.Name = ""
	config := auth.GetOne()
	if config["id"] == "" {
		c.JSON(200, gin.H{
			"code": 1, "message": "data not found",
		})
		return
	}
	auth.Name = name
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
	var auth role.RoleCreate
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

/**获取角色菜单 */
func Menus(c *gin.Context) {
	var auth role.RoleCreate
	id := c.Query("id")
	var err error
	auth.Id, err = strconv.Atoi(id)
	auth.Type, _ = strconv.Atoi(c.DefaultQuery("type", "1"))
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
	auth.Type = 1
	if auth.Type == 2 {
		auth.Type = 2
	}
	cmap := auth.GetMenus()
	c.JSON(200, gin.H{
		"message": "success",
		"code":    0,
		"data":    cmap,
	})
}

func SaveMenus(c *gin.Context) {
	var auth role.RoleMenuSave
	if err := c.ShouldBind(&auth); err != nil {
		c.JSON(200, gin.H{
			"code": 1, "message": err.Error(),
		})
		return
	}

	var authRole role.RoleCreate
	authRole.Id = auth.Id
	one := authRole.GetOne()
	if one["id"] == "" {
		c.JSON(200, gin.H{
			"code": 1, "message": "data not found",
		})
		return
	}
	auth.Type = 1
	if auth.Type == 2 {
		auth.Type = 2
	}
	cmap := auth.DoSaveMenus()
	c.JSON(200, gin.H{
		"message": "success",
		"code":    0,
		"data":    cmap,
	})
}
