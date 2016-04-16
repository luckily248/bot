package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type WarIdsModel struct {
	BaseDBmodel
	name  string
	Warid int
}

func (this *WarIdsModel) Tablename() string {
	return "warids"
}

func (this *WarIdsModel) init() (err error) {
	err = this.BaseDBmodel.init()
	if err != nil {
		return
	}
	this.c = this.db.C(this.Tablename())
	return
}

func GetNextWarId() (id int, err error) {
	warids := &WarIdsModel{}
	err = warids.init()
	if err != nil {
		return
	}
	defer warids.session.Close()
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"warid": 1}},
		ReturnNew: true,
	}
	var newwarid *WarIdsModel
	_, err = warids.c.Find(bson.M{"name": "user"}).Apply(change, &newwarid)
	//fmt.Printf("warid:%d\n", newwarid.Warid)
	//fmt.Printf("info:%v\n", info)
	id = newwarid.Warid
	return
}
