package models

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

func DBInit() {
	orm.Debug = true

	orm.RegisterDataBase("default", "sqlite3", "data.db", 30)
	orm.RunSyncdb("default", true, true)

	o := orm.NewOrm()
	o.Using("default")

	user := &User{
		Account:  "admin@admin.com",
		Password: "admin",
		Showname: "管理员",
	}

	h := md5.New()
	h.Write([]byte(user.Password))
	user.Password = hex.EncodeToString(h.Sum(nil))

	id, err := o.Insert(user)
	if err != nil {
		log.Info("%s.", err.Error())
	} else {
		log.Info("User added,ID is:%d.", id)
	}
}
