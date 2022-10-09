package models

import (
	"strconv"
	"xstart/app/common"
	"xstart/app/common/utils"
)

type RoleModel struct {
	Id       int    `orm:"id,primary" json:"id"`         //
	Name     string `orm:"name"   json:"name"`           //
	Mark     string `orm:"mark"   json:"mark"`           //
	CreateAt int    `orm:"create_at"   json:"create_at"` //
	UpdateAt int    `orm:"update_at"   json:"update_at"` //
	Deleted  int    `orm:"deleted"   json:"deleted"`     //

	Start  int
	Length int
	Ids    string
}

func (cm RoleModel) TableName() string {
	return "com_role"
}

func (cm RoleModel) QueryAll() ([]map[string]string, string) {
	m := common.MySql{}.GetConn()
	field := []string{"id", "name", "mark", "create_at", "update_at"}
	cmap := make(map[string]string)
	cmap["deleted"] = "0"
	if cm.Id != 0 {
		cmap["id"] = strconv.FormatInt(int64(cm.Id), 10)
	}
	if cm.Name != "" {
		cmap["name like "] = "'%" + cm.Name + "%'"
	}
	if cm.Ids != "" {
		cmap["id IN "] = "(" + cm.Ids + ")"
	}
	countRes := m.Select(cm.TableName(), []string{"count(*) as num"}).Where(cmap).Limit(1).QueryRow()
	if countRes["num"] != "" && countRes["num"] != "0" {
		mc := m.Select(cm.TableName(), field).
			Where(cmap).
			Limit(cm.Length).Offset(cm.Start).
			OrderByString("id", "desc")

		res := mc.QueryAll()
		return res, countRes["num"]
	}
	return nil, countRes["num"]
}
func (cm RoleModel) QueryOne() map[string]string {
	m := common.MySql{}.GetConn()
	field := []string{"id", "name", "mark", "create_at", "update_at"}
	cmap := make(map[string]string)
	cmap["deleted"] = "0"
	if cm.Id != 0 {
		cmap["id"] = strconv.FormatInt(int64(cm.Id), 10)
	}
	if cm.Name != "" {
		cmap["name"] = cm.Name
	}

	mc := m.Select(cm.TableName(), field).Where(cmap).Limit(cm.Length).Offset(cm.Start)
	res := mc.QueryRow()
	return res
}

func (cm RoleModel) Create() int64 {
	m := common.MySql{}.GetConn()
	t := utils.DateNowUnix()
	table := make(map[string]string)
	table["name"] = cm.Name
	table["mark"] = cm.Mark
	table["create_at"] = t
	table["update_at"] = t
	table["deleted"] = "0"
	return m.Insert(cm.TableName(), table)
}

func (cm RoleModel) Update() int64 {
	m := common.MySql{}.GetConn()
	t := utils.DateNowUnix()
	table := make(map[string]string)
	table["name"] = cm.Name
	table["mark"] = cm.Mark
	table["update_at"] = t
	table["deleted"] = "0"
	return m.Where(map[string]string{"id": strconv.Itoa(cm.Id)}).Update(cm.TableName(), table)
}

func (cm RoleModel) Del() int64 {
	m := common.MySql{}.GetConn()
	table := make(map[string]string)
	table["deleted"] = "1"
	return m.Where(map[string]string{"id": strconv.Itoa(cm.Id)}).Update(cm.TableName(), table)
}
