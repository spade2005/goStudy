package menu

import (
	"fmt"
	"xstart/app/common/utils"
	"xstart/app/models"
)

type MenuList struct {
	Start    int    `form:"start" binding:"gte=0"`
	Length   int    `form:"length" binding:"gte=0,lte=1000"`
	Name     string `form:"name" binding:""`
	Id       int    `form:"id" binding:""`
	MenuType int    `form:"menu_type" binding:"gte=0,lte=2"`
}

type MenuCreate struct {
	Id       int    `form:"id" binding:""`
	Name     string `form:"name" binding:"required,min=3,max=100"`
	ParentId int    `form:"parent_id" binding:"lte=65535"`
	MenuType int    `form:"menu_type" binding:"required,gte=1,lte=2"`
	SortBy   int    `form:"sort_by" binding:"required,gte=1,lte=65535"`
	Style    string `form:"style" binding:""`
	Router   string `form:"router" binding:"required,min=3,max=50"`
	Mark     string `form:"mark" binding:""`
	Status   int    `form:"status" binding:"gte=0,lte=1"`
}

func (c MenuList) GetList() ([]map[string]string, string) {
	if c.MenuType <= 1 {
		c.MenuType = 1
	}
	c.Start = 0
	c.Length = 1000 //Specific format

	model := models.MenuModel{}
	model.Start = c.Start
	model.Length = c.Length
	model.Id = c.Id
	model.Name = c.Name
	model.Type = c.MenuType
	list, count := model.QueryAll()
	return list, count
}

func (c MenuCreate) GetOne() map[string]string {
	model := models.MenuModel{}
	model.Start = 0
	model.Length = 1
	model.Id = c.Id
	model.Name = c.Name
	one := model.QueryOne()
	return one
}

func (c MenuCreate) DoCreate() int64 {
	model := models.MenuModel{}
	model.Name = c.Name
	model.ParentId = c.ParentId
	model.Type = c.MenuType
	model.SortBy = c.SortBy
	model.Style = c.Style
	model.Router = c.Router
	model.Mark = c.Mark
	model.Status = c.Status
	lastId := model.Create()
	return lastId
}
func (c MenuCreate) DoUpdate() int64 {
	model := models.MenuModel{}
	model.Id = c.Id
	model.Name = c.Name
	model.ParentId = c.ParentId
	model.Type = c.MenuType
	model.SortBy = c.SortBy
	model.Style = c.Style
	model.Router = c.Router
	model.Mark = c.Mark
	model.Status = c.Status
	lastId := model.Update()
	return lastId
}
func (c MenuCreate) DoDel() int64 {
	model := models.MenuModel{}
	model.Id = c.Id
	lastId := model.Del()
	return lastId
}

func (c MenuList) FormatMenu(menus []map[string]string) ([]map[string]string, map[string][]map[string]string) {
	if c.MenuType > 1 {
		fmt.Println("is auth menu")
		return menus, nil
	}
	var tmpRoot []map[string]string
	tmpTree := make(map[string][]map[string]string)
	for _, fv := range menus {
		fv["status_str"] = "启用"
		if fv["status"] == "1" {
			fv["status_str"] = "禁用"
		}
		fv["type_str"] = "权限"
		if fv["type"] == "1" {
			fv["type_str"] = "菜单"
		}
		fv["update_date"] = utils.DateFormatUnix(fv["update_at"])
		fv["create_date"] = utils.DateFormatUnix(fv["create_at"])
		if fv["parent_id"] == "0" {
			tmpRoot = append(tmpRoot, fv)
		} else {
			ind := "index_" + fv["parent_id"]
			tmpTree[ind] = append(tmpTree[ind], fv)
		}
	}
	return tmpRoot, tmpTree
}
