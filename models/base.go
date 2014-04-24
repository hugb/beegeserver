package models

import (
	"github.com/astaxie/beego/logs"
)

var (
	log *logs.BeeLogger
)

func init() {
	log = logs.NewLogger(100)
	log.SetLogger("console", "")
}
