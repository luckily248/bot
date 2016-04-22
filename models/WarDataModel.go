package models

import (
	"time"
)

type WarDataModel struct {
	BasePQDBmodel
	Id        int    `bson:"_id" form:"-" `
	TeamA     string `form:"TeamA"`
	TeamB     string `form:"TeamB"`
	IsEnable  bool
	Timestamp time.Time
	Begintime time.Time
}
type Battle struct {
	WarId      int
	BattleNo   int
	Scoutstate string //noscout needscout scouted
}

func (this *Battle) Init() {
	this.Scoutstate = "noscout"
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
	WarId      int
	BattleNo   int
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
	err = this.BasePQDBmodel.init()
	if err != nil {
		return
	}
	return
}

func AddWarData(teama string, teamb string, cout int) (id int, err error) {
	wardata := &WarDataModel{}
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.DB.Close()

	rows := wardata.DB.QueryRow(`INSERT INTO WarDataModel(TeamA,TeamB,IsEnable,Timestamp,Begintime) VALUES($1,$2,$3,$4,$5) RETURNING id`, teama, teamb, true, time.Now(), time.Now().Add(23*time.Hour))
	err = rows.Scan(&id)
	if err != nil {
		return
	}
	battle := &Battle{}
	battle.Init()
	for i := 1; i < 7; i++ {
		stmt1, err := wardata.DB.Prepare("INSERT INTO Battle(WarId,BattleNo,Scoutstate) VALUES($1,$2,$3)")
		if err != nil {
			break
		}
		res1, err := stmt1.Exec(id, battle.Scoutstate)
		if err != nil {
			break
		}
	}
	return
}

func GetWarData(warid int) (content *WarDataModel, err error) {
	wardata := &WarDataModel{}
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.DB.Close()
	rows := wardata.DB.QueryRow("SELECT * FROM WarDataModel WHERE ID=$1", warid)
	err = rows.Scan(wardata.Id, wardata.TeamA, wardata.TeamB, wardata.IsEnable, wardata.Timestamp, wardata.Begintime)
	if err != nil {
		return
	}
	return
}
func GetWarDatabyclanname(clanname string) (content *WarDataModel, err error) {
	wardata := &WarDataModel{}
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.DB.Close()
	rows := wardata.DB.QueryRow("SELECT * FROM WarDataModel LIMIT 1 ORDER Timestamp DESC")
	err = rows.Scan(wardata.Id, wardata.TeamA, wardata.TeamB, wardata.IsEnable, wardata.Timestamp, wardata.Begintime)
	if err != nil {
		return
	}
	return
}
func DelWarDatabyWarid(warid int) (err error) {
	wardata := &WarDataModel{}
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.DB.Close()
	stmt, err := wardata.DB.Prepare("delete from WarDataModel where ID=$1")
	if err != nil {
		return
	}
	res, err := stmt.Exec(warid)
	if err != nil {
		return
	}
	stmt, err = wardata.DB.Prepare("delete from Battle where WarId=$1")
	if err != nil {
		return
	}
	res, err = stmt.Exec(warid)
	if err != nil {
		return
	}
	stmt, err = wardata.DB.Prepare("delete from Caller where WarId=$1")
	if err != nil {
		return
	}
	res, err = stmt.Exec(warid)
	if err != nil {
		return
	}
	return
}

func UpdateWarData(warid int, wardata WarDataModel) (err error) {
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.DB.Close()
	stmt, err := wardata.DB.Prepare("update WarDataModel set Scoutstate=$1 where WarId=$2")
	if err != nil {
		return
	}
	res, err := stmt.Exec(scoutstate, warid, battleno)
	if err != nil {
		return
	}
	return
}
func UpdateBattle(warid int, battleno int, scoutstate string) (err error) {
	wardata := &WarDataModel{}
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.DB.Close()
	stmt, err := wardata.DB.Prepare("update Battle set Scoutstate=$1 where WarId=$2 BattleNo=$3")
	if err != nil {
		return
	}
	res, err := stmt.Exec(scoutstate, warid, battleno)
	if err != nil {
		return
	}
	return
}
