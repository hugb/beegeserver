package main

import (
	"github.com/astaxie/beego"

	"github.com/hugb/beegeserver/models"

	_ "github.com/hugb/beegeserver/routers"
)

func main() {
	models.DBInit()

	beego.Run()
}
