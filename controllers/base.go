package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/beego/i18n"
)

const (
	APP_VERSION = "0.0.1.140306"
)

var (
	log *logs.BeeLogger
)

func init() {
	log = logs.NewLogger(100)
	log.SetLogger("console", "")
}

type baseController struct {
	beego.Controller
	i18n.Locale
}

func (this *baseController) Prepare() {
	log.Debug("=========================== request ==============================")
}

func (this *baseController) Render() error {
	return nil
}
