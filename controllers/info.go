package controllers

import (
	"github.com/astaxie/beego"
	"zhenxiaoyuan/models"
)

type InfoController struct {
	beego.Controller
}

func (c *InfoController) Get()  {
	var hashInfo map[string]string
	hashInfo = models.TestGetHashTable()
	c.Data["title"] = hashInfo["title"]
	c.Data["content"] = hashInfo["content"]
	c.TplName = "info.tpl"
}