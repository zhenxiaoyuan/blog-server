package routers

import (
	"zhenxiaoyuan/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.IndexController{})
}
