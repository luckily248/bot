package models

import (
	"errors"
	"time"
)

type WarDataModel struct {
	BasePQDBmodel
	Id        int    `bson:"_id" form:"-" `
	TeamA     string `form:"TeamA"`
	TeamB     string `form:"TeamB"`
	BattleLen int
	IsEnable  bool
	Timestamp time.Time
	Begintime time.Time
}
type Battle struct {
	WarId      int
	BattleNo   int
	Scoutstate string //noscout needscout scouted
}
type Caller struct {
	WarId      int
	BattleNo   int
	Callername string
	Starstate  int
	Calledtime time.Time
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

func (this *Caller) Init() {
	this.Callername = ""
	this.Starstate = -1
	this.Calledtime = time.Now()
	return
}
func (this *Caller) GetStarstate() string {
	switch this.Starstate {
	case -1:
		return "\U0001F4A4\U0001F4A4\U0001F4A4"
	case 0:
		return "\U00002734\U00002734\U00002734"
	case 1:
		return "\U00002B50\U00002734\U00002734"
	case 2:
		return "\U00002B50\U00002B50\U00002734"
	case 3:
		return "\U00002B50\U00002B50\U00002B50"

	}
	return "\U0001F4A4\U0001F4A4\U0001F4A4"

}

func (this *WarDataModel) Tablename() string {
	return "wardatamodel"
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

	rows := wardata.DB.QueryRow("INSERT INTO wardatamodel(TeamA,TeamB,BattleLen,IsEnable,Timestamp,Begintime) VALUES($1,$2,$3,$4,$5,$6) RETURNING id", teama, teamb, cout, true, time.Now(), time.Now().Add(23*time.Hour))
	err = rows.Scan(&id)
	if err != nil {
		return
	}
	battle := &Battle{}
	battle.Init()
	for i := 1; i < cout+1; i++ {
		r, err := wardata.DB.Query("INSERT INTO Battle(WarId,BattleNo,Scoutstate) VALUES($1,$2,$3)", id, i, battle.Scoutstate)
		if err != nil {
			break
		}
		err = r.Close()
		if err != nil {
			break
		}
	}
	return
}
func AddCaller(caller *Caller) (err error) {
	wardata := &WarDataModel{}
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.DB.Close()

	stmt1, err := wardata.DB.Prepare("INSERT INTO Caller(WarId,BattleNo,Callername,Starstate,Calledtime) VALUES($1,$2,$3,$4,$5)")
	if err != nil {
		return
	}
	_, err = stmt1.Exec(caller.WarId, caller.BattleNo, caller.Callername, caller.Starstate, caller.Calledtime)
	if err != nil {
		return
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
	content = &WarDataModel{}
	err = rows.Scan(&content.Id, &content.TeamA, &content.TeamB, &content.IsEnable, &content.Timestamp, &content.Begintime, &content.BattleLen)
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
	rows := wardata.DB.QueryRow("SELECT * FROM WarDataModel WHERE TeamA=$1 ORDER BY Timestamp DESC", clanname)
	content = &WarDataModel{}
	err = rows.Scan(&content.Id, &content.TeamA, &content.TeamB, &content.IsEnable, &content.Timestamp, &content.Begintime, &content.BattleLen)
	if err != nil {
		return
	}
	return
}
func GetAllBattlebyId(warid int) (battles []*Battle, err error) {
	wardata := &WarDataModel{}
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.DB.Close()

	rows, err := wardata.DB.Query("SELECT * FROM Battle WHERE WarId=$1 ORDER By BattleNO ASC", warid)
	if err != nil {
		return
	}
	for rows.Next() {
		battle := &Battle{}
		err = rows.Scan(&battle.WarId, &battle.Scoutstate, &battle.BattleNo)
		if err != nil {
			break
		}
		battles = append(battles, battle)
	}
	return
}
func GetAllCallerbyId(warid int) (callers map[int][]*Caller, err error) {
	wardata := &WarDataModel{}
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.DB.Close()

	rows, err := wardata.DB.Query("SELECT * FROM Caller WHERE WarId=$1 ORDER BY BattleNO ASC", warid)
	if err != nil {
		return
	}
	callers = map[int][]*Caller{}
	for rows.Next() {
		caller := &Caller{}
		err = rows.Scan(&caller.WarId, &caller.Callername, &caller.Starstate, &caller.Calledtime, &caller.BattleNo)
		if err != nil {
			break
		}
		_, exists := callers[caller.BattleNo]
		if !exists {
			callers[caller.BattleNo] = make([]*Caller, 0)
		}
		callers[caller.BattleNo] = append(callers[caller.BattleNo], caller)
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
	_, err = stmt.Exec(warid)
	if err != nil {
		return
	}
	stmt, err = wardata.DB.Prepare("delete from Battle where WarId=$1")
	if err != nil {
		return
	}
	_, err = stmt.Exec(warid)
	if err != nil {
		return
	}
	stmt, err = wardata.DB.Prepare("delete from Caller where WarId=$1")
	if err != nil {
		return
	}
	_, err = stmt.Exec(warid)
	if err != nil {
		return
	}
	return
}
func DelCallbyNo(warid int, position int) (err error) {
	wardata := &WarDataModel{}
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.DB.Close()
	stmt, err := wardata.DB.Prepare("delete from Caller where WarId=$1 AND BattleNo=$2")
	if err != nil {
		return
	}
	result, err := stmt.Exec(warid, position)
	if err != nil {
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return
	} else if affected == 0 {
		err = errors.New("no call there")
		return
	}
	return
}
func DelCallbyid(warid int, position int, callername string) (err error) {
	wardata := &WarDataModel{}
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.DB.Close()
	stmt, err := wardata.DB.Prepare("delete from Caller where WarId=$1 AND BattleNo=$2 AND Callername=$3")
	if err != nil {
		return
	}
	result, err := stmt.Exec(warid, position, callername)
	if err != nil {
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return
	} else if affected == 0 {
		err = errors.New("you dont have any call on it")
		return
	}
	return
}

func UpdateWarData(wardata *WarDataModel) (err error) {
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.DB.Close()
	stmt, err := wardata.DB.Prepare("update WarDataModel set TeamA=$1,TeamB=$2,IsEnable=$3,Timestamp=$4,Begintime=$5 where id=$6")
	if err != nil {
		return
	}
	_, err = stmt.Exec(wardata.TeamA, wardata.TeamB, wardata.IsEnable, wardata.Timestamp, wardata.Begintime, wardata.Id)
	if err != nil {
		return
	}
	return
}
func UpdateBattleCountbyId(warid int, cout int) (err error) {
	wardata := &WarDataModel{}
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.DB.Close()

	stmt, err := wardata.DB.Prepare("update WarDataModel set BattleLen=$1 where id=$2")
	if err != nil {
		return
	}
	_, err = stmt.Exec(cout, warid)
	if err != nil {
		return
	}

	stmt, err = wardata.DB.Prepare("delete from Battle where WarId=$1")
	if err != nil {
		return
	}
	_, err = stmt.Exec(warid)
	if err != nil {
		return
	}
	battle := &Battle{}
	battle.Init()
	for i := 1; i < cout+1; i++ {
		r, err := wardata.DB.Query("INSERT INTO Battle(WarId,BattleNo,Scoutstate) VALUES($1,$2,$3)", warid, i, battle.Scoutstate)
		if err != nil {
			break
		}
		err = r.Close()
		if err != nil {
			break
		}
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
	stmt, err := wardata.DB.Prepare("update Battle set Scoutstate=$1 where WarId=$2 AND BattleNo=$3")
	if err != nil {
		return
	}
	_, err = stmt.Exec(scoutstate, warid, battleno)
	if err != nil {
		return
	}
	return
}
func UpdateCaller(caller *Caller) (err error) {
	wardata := &WarDataModel{}
	err = wardata.init()
	if err != nil {
		return
	}
	defer wardata.DB.Close()
	stmt, err := wardata.DB.Prepare("update Caller set Starstate=$1,Calledtime=$2 where WarId=$3 AND BattleNo=$4 AND Callername=$5")
	if err != nil {
		return
	}
	_, err = stmt.Exec(caller.Starstate, caller.Calledtime, caller.WarId, caller.BattleNo, caller.Callername)
	if err != nil {
		return
	}
	return
}
