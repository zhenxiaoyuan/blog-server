package controllers

import ( 
	"github.com/astaxie/beego"
	"zhenxiaoyuan/models"
)

type PingController struct {
	beego.Controller
}

func (c *PingController) Get() {
	c.Data["ping"] = models.ExampleNewClient()
	c.TplName = "ping.tpl"
}