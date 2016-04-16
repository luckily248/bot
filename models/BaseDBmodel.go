package models

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
)

type BaseDBmodel struct {
	session *mgo.Session
	db      *mgo.Database
	c       *mgo.Collection
}

func (this *BaseDBmodel) DBname() string {
	return "dataserver"
}

//成功初始化后必须调用  defer this.session.Close()
func (this *BaseDBmodel) init() (err error) {
	mgourl := beego.AppConfig.String("mgourl")
	newsession, err := mgo.Dial(mgourl)
	if err != nil {
		return
	}
	this.session = newsession
	this.session.SetMode(mgo.Monotonic, true)
	this.db = this.session.DB(this.DBname())
	return
}

func (this *BaseDBmodel) Check() (err error) {
	err = this.init()
	if err != nil {
		return
	}
	defer this.session.Close()
	return
}
