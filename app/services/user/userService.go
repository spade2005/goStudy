package user

import (
	"strings"
	"xstart/app/common/utils"
	"xstart/app/models"
)

type UserList struct {
	Start    int    `form:"start" binding:"gte=0"`
	Length   int    `form:"length" binding:"gte=0,lte=100"`
	UserName string `form:"username" binding:""`
	Phone    string `form:"phone" binding:""`
	Id       int    `form:"id" binding:""`
	Status   int    `form:"status" binding:"gte=0,lte=2"`
}

type UserCreate struct {
	Id       int    `form:"id" binding:""`
	UserName string `form:"username" binding:"required,min=3,max=100"`
	Password string `form:"userpass" binding:"required,max=20,min=4"`
	Phone    string `form:"phone" binding:""`
	Email    string `form:"email" binding:"required,email"`
	Nickname string `form:"nickname" binding:"required,min=1,max=50"`
	RoleId   int    `form:"role_id" binding:"required,gte=0"`
	Status   int    `form:"status" binding:"gte=0,lte=1"`
}

func (c UserList) GetList() ([]map[string]string, string) {
	model := models.UserModel{}
	model.Start = c.Start
	model.Length = c.Length
	model.Id = c.Id
	model.UserName = c.UserName
	model.Phone = c.Phone
	model.Status = c.Status
	list, count := model.QueryAll()
	return list, count
}

func (c UserCreate) GetOne() map[string]string {
	model := models.UserModel{}
	model.Start = 0
	model.Length = 1
	model.Id = c.Id
	model.UserName = c.UserName
	one := model.QueryOne()
	return one
}

func (c UserCreate) DoCreate() int64 {
	model := models.UserModel{}
	model.UserName = c.UserName
	model.Password = c.Password
	model.Phone = c.Phone
	model.Email = c.Email
	model.Nickname = c.Nickname
	model.RoleId = c.RoleId
	model.Status = c.Status
	lastId := model.Create()
	return lastId
}
func (c UserCreate) DoUpdate() int64 {
	model := models.UserModel{}
	model.Id = c.Id
	model.Password = c.Password
	model.Phone = c.Phone
	model.Email = c.Email
	model.Nickname = c.Nickname
	model.RoleId = c.RoleId
	model.Status = c.Status
	lastId := model.Update()
	return lastId
}
func (c UserCreate) DoDel() int64 {
	model := models.UserModel{}
	model.Id = c.Id
	lastId := model.Del()
	return lastId
}

func (c UserList) FormatList(cmap []map[string]string) []map[string]string {
	var tmp []string
	for _, v := range cmap {
		tmp = append(tmp, v["role_id"])
	}
	tmp = utils.UniqueArray(tmp)

	model := models.RoleModel{}
	model.Ids = strings.Join(tmp, ",")
	model.Length = 1000
	list, _ := model.QueryAll()
	var roleMap = make(map[string]string)
	for _, v := range list {
		roleMap["map_"+v["id"]] = v["name"]
	}

	var tmpRoot []map[string]string
	for _, fv := range cmap {
		fv["status_str"] = "启用"
		if fv["status"] == "1" {
			fv["status_str"] = "禁用"
		}
		fv["role_str"] = ""
		t, f := roleMap["map_"+fv["role_id"]]
		if f {
			fv["role_str"] = t
		}
		tmpRoot = append(tmpRoot, fv)
	}
	return tmpRoot
}

func (c UserCreate) FormatOne(cmap map[string]string) map[string]string {
	model := models.RoleModel{}
	model.Ids = cmap["role_id"]
	model.Length = 1
	role := model.QueryOne()

	cmap["status_str"] = "启用"
	if cmap["status"] == "1" {
		cmap["status_str"] = "禁用"
	}
	cmap["role_str"] = ""
	t, f := role["name"]
	if f {
		cmap["role_str"] = t
	}
	return cmap
}
