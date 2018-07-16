package controllers

import (
	// "encoding/json"
	// "fmt"
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

func (c *ArticleController) Post()  {
	// inputs := c.Input()
	// var ob models.Object
	// json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
	// fmt.Println(c.Ctx.Input.RequestBody)
	if c.Ctx.Input.Param(":id") == "add" {
		c.Data["json"] = models.AddArticle(c.Ctx.Input.RequestBody)
	} else if c.Ctx.Input.Param(":id") == "update" {
		c.Data["json"] = models.UpdateArticle(c.Ctx.Input.RequestBody)
	} else if c.Ctx.Input.Param(":id") == "delete" {
		c.Data["json"] = models.DeleteArticle(c.Ctx.Input.RequestBody)
	}
	c.ServeJSON()

}