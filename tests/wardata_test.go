package tests

import (
	. "bot/models"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func TestWarid(t *testing.T) {
	id, err := AddWarData("me", "enemy", 25)
	if err != nil {
		t.Fatalf("warid err:%s\n", err.Error())
	}
	t.Logf("warid is %d\n", id)

	wardata, err := GetWarData(24)
	if err != nil {
		t.Fatalf("GetWarData err:%s\n", err.Error())
	}
	t.Logf("GetWarData is %v\n", wardata)

	newbattlep := &Battle{}
	newbattlep.Init()
	newbattle := *newbattlep
	newcallp := &Caller{"me", 2, time.Now()}
	newbattle.Callers = append(newbattle.Callers, *newcallp)

	err = UpdateWarData(16, bson.M{"$set": bson.M{"battles.5": newbattle}})
	if err != nil {
		t.Fatalf("UpdateWarData err:%s\n", err.Error())
	}
	t.Logf("UpdateWarData is %v\n", "done")
}
