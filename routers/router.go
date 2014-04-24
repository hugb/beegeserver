package routers

import (
	"github.com/astaxie/beego"

	"github.com/hugb/beegeserver/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/ws", &controllers.WSController{})
}
