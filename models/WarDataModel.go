package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type WarDataModel struct {
	BaseDBmodel
	Id        int    `bson:"_id" form:"-" `
	TeamA     string `form:"TeamA"`
	TeamB     string `form:"TeamB"`
	Battles   []Battle
	IsEnable  bool
	Timestamp time.Time
	Begintime time.Time
}
type Battle struct {
	Scoutstate string //noscout needscout scouted
	Callers    []Caller
}

func (this *Battle) Init() {
	this.Scoutstate = "noscout"
	this.Callers = make([]Caller, 0)
	return
}
func (this *Battle) Needscout() {
	this.Scoutstate = "needscout"
	return
}
func (this *Battle) Scouted() {
	this.Scoutstate = "scouted"
	return
}

type Caller struct {
	Callername string
	Starstate  int
	Calledtime time.Time
}

func (this *Caller) Init() {
	this.Callername = ""
	this.Starstate = -1
	this.Calledtime = time.Now()
	return
}
func (this *Caller) GetStarstate() string {
	switch this.Starstate {
	case -1:
		return "ZZZ"
	case 0:
		return "XXX"
	case 1:
		return "OXX"
	case 2:
		return "OOX"
	case 3:
		return "OOO"

	}
	return "ZZZ"

}

func (this *WarDataModel) Tablename() string {
	return "wardata"
}

func (this *WarDataModel) init() (err error) {
	err = this.BaseDBmodel.init()
	if err != nil {
		return
	}
	this.c = this.db.C(this.Tablename())
	this.TeamA = "teamA"
	this.TeamB = "teamB"
	this.IsEnable = true
	this.Timestamp = time.Now()
	return
}

func AddWarData(teama string, teamb string, cout int) (id int, err error) {
	wardata := &WarDataModel{}
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.session.Close()
	id, err = GetNextWarId()
	if err != nil {
		return
	}
	wardata.Id = id
	wardata.TeamA = teama
	wardata.TeamB = teamb
	wardata.Begintime = time.Now().Add(23 * time.Hour)
	wardata.Battles = make([]Battle, cout)
	battlep := &Battle{}
	battlep.Init()
	for index := range wardata.Battles {
		wardata.Battles[index] = *battlep
	}
	err = wardata.c.Insert(wardata)
	return
}

func GetWarData(warid int) (content *WarDataModel, err error) {
	wardata := &WarDataModel{}
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.session.Close()
	//err = wardata.c.Find(bson.M{"teama": wardata.TeamA}).Sort("-timestamp").One(&content)
	err = wardata.c.FindId(warid).One(&content)
	return
}
func GetWarDatabyclanname(clanname string) (content *WarDataModel, err error) {
	wardata := &WarDataModel{}
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.session.Close()
	err = wardata.c.Find(bson.M{"teama": clanname}).Sort("-timestamp").One(&content)
	return
}
func DelWarDatabyWarid(warid int) (err error) {
	wardata := &WarDataModel{}
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.session.Close()
	err = wardata.c.RemoveId(warid)
	return
}

func UpdateWarData(warid int, updata interface{}) (err error) {
	wardata := &WarDataModel{}
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.session.Close()
	err = wardata.c.UpdateId(warid, updata)
	return
}
