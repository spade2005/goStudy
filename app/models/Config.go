package models

import (
	"strconv"
	"xstart/app/common"
)

type ConfigModel struct {
	Id       int    `orm:"id,primary" json:"id"`     //
	KeyStr   string `orm:"key_str"   json:"name"`    //
	ValueStr string `orm:"value_str"   json:"value"` //
	Deleted  int    `orm:"deleted"   json:"deleted"` //
	Start    int
	Length   int
}

func (cm ConfigModel) TableName() string {
	return "com_config"
}

func (cm ConfigModel) QueryAll() ([]map[string]string, string) {
	m := common.MySql{}.GetConn()
	field := []string{"id", "key_str", "value_str"}
	cmap := make(map[string]string)
	cmap["deleted"] = "0"
	if cm.Id != 0 {
		cmap["id"] = strconv.FormatInt(int64(cm.Id), 10)
	}
	if cm.KeyStr != "" {
		cmap["key_str like "] = "'%" + cm.KeyStr + "%'"
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
func (cm ConfigModel) QueryOne() map[string]string {
	m := common.MySql{}.GetConn()
	field := []string{"id", "key_str", "value_str"}
	cmap := make(map[string]string)
	cmap["deleted"] = "0"
	if cm.Id != 0 {
		cmap["id"] = strconv.FormatInt(int64(cm.Id), 10)
	}
	if cm.KeyStr != "" {
		cmap["key_str"] = cm.KeyStr
	}

	mc := m.Select(cm.TableName(), field).Where(cmap).Limit(cm.Length).Offset(cm.Start)
	res := mc.QueryRow()
	return res
}

func (cm ConfigModel) Create() int64 {
	m := common.MySql{}.GetConn()
	table := make(map[string]string)
	table["key_str"] = cm.KeyStr
	table["value_str"] = cm.ValueStr
	table["deleted"] = "0"
	return m.Insert(cm.TableName(), table)
}

func (cm ConfigModel) Update() int64 {
	m := common.MySql{}.GetConn()
	table := make(map[string]string)
	table["value_str"] = cm.ValueStr
	return m.Where(map[string]string{"id": strconv.Itoa(cm.Id)}).Update(cm.TableName(), table)
}

func (cm ConfigModel) Del() int64 {
	m := common.MySql{}.GetConn()
	table := make(map[string]string)
	table["deleted"] = "1"
	return m.Where(map[string]string{"id": strconv.Itoa(cm.Id)}).Update(cm.TableName(), table)
}
