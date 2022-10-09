package models

import (
	"strconv"
	"xstart/app/common"
	"xstart/app/common/utils"
)

type UserModel struct {
	Id       int    `orm:"id,primary" json:"id"`         //
	UserName string `orm:"username"   json:"username"`   //
	Password string `orm:"password"   json:"password"`   //
	Phone    string `orm:"phone"   json:"phone"`         //
	Email    string `orm:"email"   json:"email"`         //
	Nickname string `orm:"nick_name"   json:"nick_name"` //
	RoleId   int    `orm:"role_id"   json:"role_id"`     //
	Status   int    `orm:"status"   json:"status"`       //
	CreateAt int    `orm:"create_at"  json:"create_at"`  //
	UpdateAt int    `orm:"update_at"  json:"update_at"`  //
	Deleted  int    `orm:"deleted"  json:"deleted"`      //

	Start  int
	Length int
}

func (cm UserModel) TableName() string {
	return "com_user"
}

func (cm UserModel) GetUserByUserName(username string) map[string]string {
	m := common.MySql{}.GetConn()
	field := []string{"id", "username", "password_hash", "nick_name", "role_id", "status"}
	res := m.Select(cm.TableName(), field).
		Where(map[string]string{"username": username, "deleted": "0"}).
		Limit(1).
		QueryRow()
	return res
}
func (cm UserModel) GetUserByUserId(id string) map[string]string {
	m := common.MySql{}.GetConn()
	field := []string{"id", "username", "password_hash", "nick_name", "role_id", "status"}
	res := m.Select(cm.TableName(), field).
		Where(map[string]string{"id": id, "deleted": "0"}).
		Limit(1).
		QueryRow()
	return res
}

func (cm UserModel) CreateUser(user map[string]string) int64 {
	m := common.MySql{}.GetConn()
	t := utils.DateNowUnix()
	user["create_at"] = t
	user["update_at"] = t
	return m.Insert(cm.TableName(), user)
}

func (cm UserModel) UpdatePassword(user map[string]string, where map[string]string) int64 {
	m := common.MySql{}.GetConn()
	user["update_at"] = utils.DateNowUnix()
	return m.Where(where).Update(cm.TableName(), user)
}

/** curd ***/

func (cm UserModel) QueryAll() ([]map[string]string, string) {
	m := common.MySql{}.GetConn()
	field := []string{"id", "username", "phone", "email", "nick_name", "role_id", "status", "create_at", "update_at"}
	cmap := make(map[string]string)
	cmap["deleted"] = "0"
	if cm.Id != 0 {
		cmap["id"] = strconv.FormatInt(int64(cm.Id), 10)
	}
	if cm.Status != 0 {
		cmap["status"] = strconv.FormatInt(int64(cm.Status), 10)
	}
	if cm.Phone != "" {
		cmap["phone"] = cm.Phone
	}
	if cm.UserName != "" {
		cmap["username like "] = "'%" + cm.UserName + "%'"
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
func (cm UserModel) QueryOne() map[string]string {
	m := common.MySql{}.GetConn()
	field := []string{"id", "username", "phone", "email", "nick_name", "role_id", "status", "create_at", "update_at"}
	cmap := make(map[string]string)
	cmap["deleted"] = "0"
	if cm.Id != 0 {
		cmap["id"] = strconv.FormatInt(int64(cm.Id), 10)
	}
	if cm.UserName != "" {
		cmap["username"] = cm.UserName
	}

	mc := m.Select(cm.TableName(), field).Where(cmap).Limit(cm.Length).Offset(cm.Start)
	res := mc.QueryRow()
	return res
}

func (cm UserModel) Create() int64 {
	m := common.MySql{}.GetConn()
	t := utils.DateNowUnix()
	passwordHash, _ := utils.PasswordHash(cm.Password)
	table := map[string]string{
		"username": cm.UserName, "password_hash": passwordHash, "nick_name": cm.Nickname,
		"phone": cm.Phone, "email": cm.Email,
		"role_id": strconv.FormatInt(int64(cm.RoleId), 10), "status": strconv.FormatInt(int64(cm.Status), 10),
		"create_at": t, "update_at": t, "deleted": "0",
	}
	return m.Insert(cm.TableName(), table)
}

func (cm UserModel) Update() int64 {
	m := common.MySql{}.GetConn()
	t := utils.DateNowUnix()
	table := map[string]string{
		"nick_name": cm.Nickname, "phone": cm.Phone, "email": cm.Email,
		"role_id": strconv.FormatInt(int64(cm.RoleId), 10), "status": strconv.FormatInt(int64(cm.Status), 10),
		"update_at": t,
	}
	if cm.Password != "" {
		passwordHash, _ := utils.PasswordHash(cm.Password)
		table["password_hash"] = passwordHash
	}
	return m.Where(map[string]string{"id": strconv.Itoa(cm.Id)}).Update(cm.TableName(), table)
}

func (cm UserModel) Del() int64 {
	m := common.MySql{}.GetConn()
	table := make(map[string]string)
	table["deleted"] = "1"
	return m.Where(map[string]string{"id": strconv.Itoa(cm.Id)}).Update(cm.TableName(), table)
}
