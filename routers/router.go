package routers

import (
	"zhenxiaoyuan/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.IndexController{})
    beego.Router("/api/article/:id", &controllers.ArticleController{})
}
