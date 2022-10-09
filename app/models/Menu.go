package models

import (
	"strconv"
	"xstart/app/common"
	"xstart/app/common/utils"
)

type MenuModel struct {
	Id       int    `orm:"id,primary" json:"id"`         //
	Name     string `orm:"name"   json:"name"`           //
	ParentId int    `orm:"parent_id"   json:"parent_id"` //
	Type     int    `orm:"type"   json:"type"`           //
	SortBy   int    `orm:"sort_by"   json:"sort_by"`     //
	Style    string `orm:"style"   json:"style"`         //
	Router   string `orm:"router"   json:"router"`       //
	Mark     string `orm:"mark"   json:"mark"`           //
	Status   int    `orm:"status"   json:"status"`       //
	CreateAt int    `orm:"create_at"   json:"create_at"` //
	UpdateAt int    `orm:"update_at"   json:"update_at"` //
	Deleted  int    `orm:"deleted"   json:"deleted"`     //

	Start  int
	Length int
}

func (cm MenuModel) TableName() string {
	return "com_menu"
}

func (cm MenuModel) QueryAll() ([]map[string]string, string) {
	m := common.MySql{}.GetConn()
	field := []string{"id", "name", "parent_id", "type", "sort_by", "style", "router", "create_at", "update_at", "status"}
	cmap := make(map[string]string)
	cmap["deleted"] = "0"
	if cm.Id != 0 {
		cmap["id"] = strconv.FormatInt(int64(cm.Id), 10)
	}
	if cm.Name != "" {
		cmap["name like "] = "'%" + cm.Name + "%'"
	}
	if cm.Type != 0 {
		cmap["type"] = strconv.FormatInt(int64(cm.Type), 10)
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
func (cm MenuModel) QueryOne() map[string]string {
	m := common.MySql{}.GetConn()
	field := []string{"id", "name", "parent_id", "type", "sort_by", "style", "router", "create_at", "update_at", "status", "mark"}
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

func (cm MenuModel) Create() int64 {
	m := common.MySql{}.GetConn()
	t := utils.DateNowUnix()
	table := make(map[string]string)
	table["name"] = cm.Name
	table["parent_id"] = strconv.Itoa(cm.ParentId)
	table["type"] = strconv.Itoa(cm.Type)
	table["sort_by"] = strconv.Itoa(cm.SortBy)
	table["style"] = cm.Style
	table["router"] = cm.Router
	table["mark"] = cm.Mark
	table["status"] = strconv.Itoa(cm.Status)
	table["create_at"] = t
	table["update_at"] = t
	table["deleted"] = "0"
	return m.Insert(cm.TableName(), table)
}

func (cm MenuModel) Update() int64 {
	m := common.MySql{}.GetConn()
	t := utils.DateNowUnix()
	table := make(map[string]string)
	table["name"] = cm.Name
	table["parent_id"] = strconv.Itoa(cm.ParentId)
	table["type"] = strconv.Itoa(cm.Type)
	table["sort_by"] = strconv.Itoa(cm.SortBy)
	table["style"] = cm.Style
	table["router"] = cm.Router
	table["mark"] = cm.Mark
	table["status"] = strconv.Itoa(cm.Status)
	table["update_at"] = t
	return m.Where(map[string]string{"id": strconv.Itoa(cm.Id)}).Update(cm.TableName(), table)
}

func (cm MenuModel) Del() int64 {
	m := common.MySql{}.GetConn()
	table := make(map[string]string)
	table["deleted"] = "1"
	return m.Where(map[string]string{"id": strconv.Itoa(cm.Id)}).Update(cm.TableName(), table)
}
