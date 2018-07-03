package controllers

import ( 
	"github.com/astaxie/beego"
	// "encoding/json"
	// "zhenxiaoyuan/models"
)

type PingController struct {
	beego.Controller
}

type pingPong struct {
	ping string
}
type Server struct {
	ServerName string
	ServerIP   string
}
func (c *PingController) Get() {
	// var myServer Server
	// myServer.ServerIP = "127.0.0.1"
	// myServer.ServerName = "ServerName"
	// b, err := json.Marshal(myServer)
	// if err != nil {
	// 	return
	// }
	var myPing pingPong
	myPing.ping = "pong"

	c.Data["json"] = &myPing
	// pp := pingPong {
	// 	ping: "pong",
	// }
	// // pp.ping = models.ExampleNewClient()
	// c.Data["json"] = &pp
	// // c.Data["hello"] = "world"
	// // c.Data["json"] = map[string]string{"ping": "pong"}
	// // c.Data[""]
	// // c.Data[""]
	// c.TplName = "ping.tpl"
	c.ServeJSON()
	// c.Data["ping"] = models.ExampleNewClient()
	
}