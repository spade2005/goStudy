package config

import (
	"xstart/app/models"
)

type ConfigList struct {
	Start   int    `form:"start" binding:"gte=0"`
	Length  int    `form:"length" binding:"gte=0,lte=100"`
	Keyword string `form:"keyword" binding:""`
	Id      int    `form:"id" binding:""`
}

type ConfigCreate struct {
	Id       int    `form:"id" binding:""`
	KeyStr   string `form:"keystr" binding:"required,min=3,max=100"`
	ValueStr string `form:"valuestr" binding:"required,max=65535"`
}

func (c ConfigList) GetList() ([]map[string]string, string) {
	if c.Length <= 0 {
		c.Length = 10
	}
	configModel := models.ConfigModel{}
	configModel.Start = c.Start
	configModel.Length = c.Length
	configModel.Id = c.Id
	configModel.KeyStr = c.Keyword
	list, count := configModel.QueryAll()
	return list, count
}

func (c ConfigCreate) GetOne() map[string]string {
	configModel := models.ConfigModel{}
	configModel.Start = 0
	configModel.Length = 1
	configModel.Id = c.Id
	configModel.KeyStr = c.KeyStr
	one := configModel.QueryOne()
	return one
}

func (c ConfigCreate) DoCreate() int64 {
	configModel := models.ConfigModel{}
	configModel.KeyStr = c.KeyStr
	configModel.ValueStr = c.ValueStr
	lastId := configModel.Create()
	return lastId
}
func (c ConfigCreate) DoUpdate() int64 {
	configModel := models.ConfigModel{}
	configModel.Id = c.Id
	configModel.KeyStr = c.KeyStr
	configModel.ValueStr = c.ValueStr
	lastId := configModel.Update()
	return lastId
}
func (c ConfigCreate) DoDel() int64 {
	configModel := models.ConfigModel{}
	configModel.Id = c.Id
	lastId := configModel.Del()
	return lastId
}
