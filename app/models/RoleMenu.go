package models

import (
	"strconv"
	"xstart/app/common"
)

type RoleMenuModel struct {
	Id      int `orm:"id,primary" json:"id"`     //
	RoleId  int `orm:"role_id"   json:"role_id"` //
	MenuId  int `orm:"menu_id"   json:"menu_id"` //
	Type    int `orm:"type"   json:"type"`       //
	Deleted int `orm:"deleted"   json:"deleted"` //
}

func (cm RoleMenuModel) TableName() string {
	return "com_role_menu"
}

func (cm RoleMenuModel) QueryAll() []map[string]string {
	m := common.MySql{}.GetConn()
	field := []string{"id", "role_id", "menu_id"}
	cmap := make(map[string]string)
	cmap["deleted"] = "0"
	if cm.Id != 0 {
		cmap["id"] = strconv.FormatInt(int64(cm.Id), 10)
	}
	if cm.RoleId != 0 {
		cmap["role_id"] = strconv.FormatInt(int64(cm.RoleId), 10)
	}
	if cm.Type != 0 {
		cmap["type"] = strconv.FormatInt(int64(cm.Type), 10)
	}
	mc := m.Select(cm.TableName(), field).Where(cmap)
	res := mc.QueryAll()
	return res
}

func (cm RoleMenuModel) Del() int64 {
	m := common.MySql{}.GetConn()
	table := make(map[string]string)
	table["deleted"] = "1"
	return m.Where(map[string]string{"id": strconv.Itoa(cm.Id)}).Update(cm.TableName(), table)
}
func (cm RoleMenuModel) DelByMenu() int64 {
	m := common.MySql{}.GetConn()
	table := make(map[string]string)
	table["deleted"] = "1"
	return m.Where(map[string]string{"role_id": strconv.Itoa(cm.RoleId), "menu_id": strconv.Itoa(cm.MenuId)}).
		Update(cm.TableName(), table)
}
func (cm RoleMenuModel) CreateByMenu() int64 {
	m := common.MySql{}.GetConn()
	field := []string{"id", "deleted"}
	cmap := map[string]string{"role_id": strconv.Itoa(cm.RoleId), "menu_id": strconv.Itoa(cm.MenuId)}
	mc := m.Select(cm.TableName(), field).Where(cmap)
	one := mc.QueryRow()
	if one["id"] == "" {
		cmap["type"] = strconv.Itoa(cm.Type)
		cmap["deleted"] = "0"
		return m.Insert(cm.TableName(), cmap)
	}
	where := map[string]string{"id": one["id"]}
	table := map[string]string{"deleted": "0"}
	return m.Where(where).Update(cm.TableName(), table)
}
