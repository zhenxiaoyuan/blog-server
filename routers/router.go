package routers

import (
	"zhenxiaoyuan/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.IndexController{})
    beego.Router("/api/ping", &controllers.PingController{})
    beego.Router("/api/info", &controllers.InfoController{})
}
