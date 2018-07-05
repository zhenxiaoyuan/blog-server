package controllers

import (
	"github.com/astaxie/beego"
	"zhenxiaoyuan/models"
)

type ArticleController struct {
	beego.Controller
}

func (c *ArticleController) Get()  {
	if c.Ctx.Input.Param(":id") == "all" {
		c.Data["json"] = models.GetAllArticles()
	} else {
		c.Data["json"] = models.GetOneArticle(c.Ctx.Input.Param(":id"))
	}
	c.ServeJSON()
}