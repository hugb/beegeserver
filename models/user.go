package models

import (
	"errors"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type User struct {
	Id       int64  `orm:"pk;auto"`
	Account  string `form:"account" valid:"Required" orm:"unique"`
	Password string `form:"password" valid:"Required"`
	Showname string `form:"-"`
	Remember string `form:"remember" orm:"-"`
}

func init() {
	orm.RegisterModelWithPrefix("beege_", new(User))
	WSSwitcher.RegisterHandler("login", login)
}

func (u *User) TableName() string {
	return "user"
}

func (this *User) Login() bool {
	this.Password = PwdHash(this.Password)
	err := orm.NewOrm().Read(this, "account", "password")
	if err != nil {
		return false
	} else {
		return true
	}
}

func checkUser(u *User) (err error) {
	valid := validation.Validation{}
	if b, _ := valid.Valid(&u); !b {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}
	return nil
}

func AddUser(u *User) (int64, error) {
	if err := checkUser(u); err != nil {
		return 0, err
	}
	u.Password = PwdHash(u.Password)
	return orm.NewOrm().Insert(u)
}

func login(conn *Connection, message *WSMessage) {
	defer conn.SendMessage(message)

	params, ok := message.Data.(map[string]interface{})
	if !ok {
		message.Code = 5
		return
	}
	user := &User{
		Account:  params["account"].(string),
		Password: params["password"].(string),
	}
	if !user.Login() {
		// Invalid account or password
		message.Code = 4
		return
	}

	message.Code = 0
	conn.auth = true
}
