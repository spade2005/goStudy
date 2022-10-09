package auth

import (
	"strconv"
	"time"
	"xstart/app/common/utils"
	"xstart/app/models"
)

type AuthLogin struct {
	Username string `form:"username" binding:"required,max=20,gte=4"`
	Userpass string `form:"userpass" binding:"required,max=20,gte=4"`
	Ip       string
}

func (c AuthLogin) ValidUser() (map[string]string, string) {
	user := models.UserModel{}.GetUserByUserName(c.Username)
	if user["id"] == "" {
		return nil, "账号不存在"
	}
	if !utils.PasswordVerify(c.Userpass, user["password_hash"]) {
		return nil, "账号或者密码错误"
	}
	return user, ""
}

func (c AuthLogin) GetUserToken(user map[string]string) (map[string]string, string) {
	tokenModel := models.UserTokenModel{}
	token := tokenModel.SelectToken(user["id"])
	if token["id"] == "" {
		token = tokenModel.CreateToken(user["id"], c.Ip)
	}
	return token, ""
}
func (c AuthLogin) UpdateUserToken(token map[string]string, isUpdate bool) map[string]string {
	tokenModel := models.UserTokenModel{}
	if token["id"] == "" {
		token = tokenModel.CreateToken(token["mid"], c.Ip)
	}
	expireAt, _ := strconv.ParseInt(token["expire_at"], 10, 64)
	if expireAt <= time.Now().Unix()+1800 {
		token["ip"] = c.Ip
		token["isUpdate"] = ""
		if isUpdate {
			token["isUpdate"] = "1"
		}
		token = tokenModel.UpdateToken(token)
	}
	return token
}

func (c AuthLogin) ValidToken(tokenStr string) map[string]string {
	m := models.UserTokenModel{}
	token := m.SelectTokenByStr(tokenStr)
	if token["id"] == "" {
		return nil
	}
	return token
}

func (c AuthLogin) RemoveToken(token map[string]string) map[string]string {
	m := models.UserTokenModel{}
	m.RemoveToken(token)
	return token
}

func (c AuthLogin) CreateUser(user map[string]string) map[string]string {
	m := models.UserModel{}
	user["password_hash"], _ = utils.PasswordHash(user["password_hash"])
	id := m.CreateUser(user)
	if id == 0 {
		return nil
	}
	user["id"] = strconv.FormatInt(id, 10)
	return user
}

func (c AuthLogin) CheckToken(tokenStr string) map[string]string {
	if tokenStr == "" {
		return nil
	}
	m := models.UserTokenModel{}
	token := m.SelectTokenByStr(tokenStr)
	if token["id"] == "" {
		return nil
	}
	expireAt, _ := strconv.ParseInt(token["expire_at"], 10, 64)
	if expireAt < time.Now().Unix() {
		return nil
	}
	user := models.UserModel{}.GetUserByUserId(token["mid"])
	if user["id"] == "" {
		return nil
	}
	c.UpdateUserToken(token, false)
	user["token"] = token["token"]
	user["token_id"] = token["id"]
	return user
}

func (c AuthLogin) UpdatePassword(user map[string]string) int64 {
	m := models.UserModel{}
	user["password_hash"], _ = utils.PasswordHash(user["password_hash"])
	id := m.UpdatePassword(map[string]string{"password_hash": user["password_hash"]},
		map[string]string{"id": user["id"], "deleted": "0"})
	return id
}
