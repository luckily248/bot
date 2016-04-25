package main

import (
	"bot/handler"
	"bot/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func main() {
	id, err := models.AddWarData("my", "ee", 25)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("addwar:" + strconv.Itoa(id))

	content, err := models.GetWarDatabyclanname("my")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("getbyname:" + strconv.Itoa(content.Id) + "," + content.TeamA + "," + content.TeamB + "," + strconv.Itoa(content.BattleLen) + "," + content.Timestamp.String())

	content, err = models.GetWarData(content.Id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("getbyid:" + content.TeamA + "," + content.TeamB + "," + strconv.Itoa(content.BattleLen) + "," + content.Timestamp.String())

	content.TeamB = "123123"
	err = models.UpdateWarData(content)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("updatedwar")

	err = models.UpdateBattleCountbyId(content.Id, 30)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("updatecounted")

	err = models.UpdateBattle(content.Id, 15, "needscout")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("updatedbattle")

	battles, err := models.GetAllBattlebyId(content.Id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for num, battle := range battles {
		fmt.Println("battles:" + strconv.Itoa(num) + "," + battle.Scoutstate)
	}

	caller := &models.Caller{}
	caller.WarId = content.Id
	caller.BattleNo = 18
	caller.Callername = "luck"
	caller.Starstate = -1
	caller.Calledtime = time.Now()
	err = models.AddCaller(caller)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("addcaller")

	caller.Starstate = 2
	err = models.UpdateCaller(caller)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("updatecaller")

	acallers, err := models.GetAllCallerbyId(content.Id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for num, callers := range acallers {
		for num1, caller1 := range callers {
			fmt.Println("callers:" + strconv.Itoa(num) + "," + strconv.Itoa(num1) + "," + caller1.Callername)
		}
	}

	//err = models.DelWarDatabyWarid(content.Id)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//fmt.Println("delwar")

	http.HandleFunc("/bot", WarDataController)
	http.ListenAndServe(":8888", nil)
}
func WarDataController(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Printf("otherMethod:%s\n", r.Method)
		return
	}
	var rec models.GMrecModel
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	}
	if err := json.Unmarshal(body, &rec); err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	}
	fmt.Printf("rec:%v\n", rec)
	if rec.Text == "" {
		fmt.Printf("is empty\n")
		return
	}
	if !strings.HasPrefix(rec.Text, "!") {
		return
	}
	reptext, err := handler.HandlecocText(rec)
	fmt.Printf("reptextlen:%d\n", utf8.RuneCountInString(reptext))
	rep := &models.GMrepModel{}
	rep.Init()
	if err != nil {
		rep.SetText(err.Error())
		fmt.Printf("err:%s\n", err.Error())
	} else {
		rep.SetText(reptext)
		fmt.Printf("ob:%v\n", rep)
	}
	buff, err := json.Marshal(rep)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		return
	}
	fmt.Println(string(buff))
	go httpPost(buff)
	return
}
func httpPost(rep []byte) {
	resp, err := http.Post("https://api.groupme.com/v3/bots/post",
		"application/x-www-form-urlencoded",
		bytes.NewReader(rep))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}
