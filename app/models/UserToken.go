package models

import (
	"strconv"
	"time"
	"xstart/app/common"
	"xstart/app/common/utils"
)

type UserTokenModel struct {
	Id       uint   `orm:"id,primary" json:"id"`         //
	Mid      int    `orm:"mid"   json:"mid"`             //
	Token    string `orm:"token"   json:"token"`         //
	ExpireAt int64  `orm:"expire_at"   json:"expire_at"` //
	Ip       string `orm:"ip"   json:"ip"`               //
	CreateAt uint   `orm:"create_at"  json:"create_at"`  //
	UpdateAt uint   `orm:"update_at"  json:"update_at"`  //
	Deleted  int    `orm:"deleted"  json:"deleted"`      //
}

func (cm UserTokenModel) TableName() string {
	return "com_user_token"
}

func (cm UserTokenModel) SelectToken(userId string) map[string]string {
	m := common.MySql{}.GetConn()
	field := []string{"id", "mid", "token", "expire_at", "ip"}
	res := m.Select(cm.TableName(), field).
		Where(map[string]string{"mid": userId, "deleted": "0"}).
		Limit(1).
		QueryRow()
	return res
}

func (cm UserTokenModel) CreateToken(userId string, ip string) map[string]string {
	m := common.MySql{}.GetConn()
	t := strconv.FormatInt(time.Now().Unix(), 10)
	data := map[string]string{
		"mid": userId, "token": "kv" + utils.RandStringBytesMaskImprSrcUnsafe(30),
		"expire_at": strconv.FormatInt(time.Now().Unix()+7200, 10),
		"ip":        ip, "create_at": t, "update_at": t,
	}
	id := m.Insert(cm.TableName(), data)
	data["id"] = strconv.FormatInt(id, 10)
	return data
}

func (cm UserTokenModel) UpdateToken(token map[string]string) map[string]string {
	m := common.MySql{}.GetConn()
	t := strconv.FormatInt(time.Now().Unix(), 10)
	data := map[string]string{
		"expire_at": strconv.FormatInt(time.Now().Unix()+7200, 10),
		"ip":        token["ip"], "update_at": t,
	}
	if token["isUpdate"] != "" {
		data["token"] = "kv" + utils.RandStringBytesMaskImprSrcUnsafe(30)
	}
	id := m.Where(map[string]string{"id": token["id"], "deleted": "0"}).
		Update(cm.TableName(), data)
	if id == 0 {
		//log update error
	}
	token["expire_at"] = data["expire_at"]
	token["update_at"] = data["update_at"]
	token["token"] = data["token"]
	return token
}

func (cm UserTokenModel) SelectTokenByStr(tokenStr string) map[string]string {
	m := common.MySql{}.GetConn()
	field := []string{"id", "mid", "token", "expire_at", "ip"}
	res := m.Select(cm.TableName(), field).
		Where(map[string]string{"token": tokenStr, "deleted": "0"}).
		Limit(1).
		QueryRow()
	return res
}

func (cm UserTokenModel) RemoveToken(token map[string]string) map[string]string {
	m := common.MySql{}.GetConn()
	t := strconv.FormatInt(time.Now().Unix(), 10)
	data := map[string]string{
		"expire_at": "0", "update_at": t,
	}
	id := m.Where(map[string]string{"id": token["id"], "deleted": "0"}).
		Update(cm.TableName(), data)
	if id == 0 {
		//log update error
	}
	return data
}
