package role

import (
	"strconv"
	"strings"
	"xstart/app/common/utils"
	"xstart/app/models"
)

type RoleList struct {
	Start  int    `form:"start" binding:"gte=0"`
	Length int    `form:"length" binding:"gte=0,lte=100"`
	Name   string `form:"name" binding:""`
	Id     int    `form:"id" binding:""`
}

type RoleCreate struct {
	Id   int    `form:"id" binding:""`
	Name string `form:"name" binding:"required,min=3,max=100"`
	Mark string `form:"mark" binding:""`
	Type int    `form:"type" binding:""`
}
type RoleMenuSave struct {
	Id      int    `form:"id" binding:""`
	MenuIds string `form:"menuIds" binding:"max=65535"`
	Type    int    `form:"type" binding:""`
}

func (c RoleList) GetList() ([]map[string]string, string) {
	model := models.RoleModel{}
	model.Start = c.Start
	model.Length = c.Length
	model.Id = c.Id
	model.Name = c.Name
	list, count := model.QueryAll()
	return list, count
}

func (c RoleCreate) GetOne() map[string]string {
	model := models.RoleModel{}
	model.Start = 0
	model.Length = 1
	model.Id = c.Id
	model.Name = c.Name
	one := model.QueryOne()
	return one
}

func (c RoleCreate) DoCreate() int64 {
	model := models.RoleModel{}
	model.Name = c.Name
	model.Mark = c.Mark
	lastId := model.Create()
	return lastId
}
func (c RoleCreate) DoUpdate() int64 {
	model := models.RoleModel{}
	model.Id = c.Id
	model.Name = c.Name
	model.Mark = c.Mark
	lastId := model.Update()
	return lastId
}
func (c RoleCreate) DoDel() int64 {
	model := models.RoleModel{}
	model.Id = c.Id
	lastId := model.Del()
	return lastId
}

func (c RoleCreate) GetMenus() []string {
	model := models.RoleMenuModel{}
	model.RoleId = c.Id
	model.Type = c.Type
	roleMenuList := model.QueryAll()
	ids := make([]string, 0)
	for _, v := range roleMenuList {
		ids = append(ids, v["menu_id"])
	}
	return ids
}

func (c RoleMenuSave) DoSaveMenus() []string {
	model := models.RoleMenuModel{}
	model.RoleId = c.Id
	model.Type = c.Type
	roleMenuList := model.QueryAll()
	ids := make([]string, 0)
	for _, v := range roleMenuList {
		ids = append(ids, v["menu_id"])
	}
	postIds := strings.Split(c.MenuIds, ",")

	needDel := utils.DiffArray(ids, postIds)
	needAdd := utils.DiffArray(postIds, ids)
	for _, v := range needDel {
		model.MenuId, _ = strconv.Atoi(v)
		if model.MenuId > 0 {
			model.DelByMenu()
		}
	}
	for _, v := range needAdd {
		model.MenuId, _ = strconv.Atoi(v)
		if model.MenuId > 0 {
			model.CreateByMenu()
		}
	}

	return ids
}
