package handler

import (
	"bot/models"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var zerostar []string
var onestar []string
var twostar []string
var threestar []string

func init() {
	zerostar = []string{
		"Wow what a nub",
		"Do u even skitch bro?",
		"Nice try... for a handicap"}
	onestar = []string{
		"Good job.Nice scout",
		"Looting aint easy"}
	twostar = []string{
		"2 stars wins wars bro!",
		"50% ftw!"}
	threestar = []string{
		"Nice hit!!",
		"Well play"}
}

type CocbotHandler interface {
	getCommands() []string
	getHelp() string
	handle(text []string) (result string, err error)
}

var mainhandler *MainHandler

func HandlecocText(rec models.GMrecModel) (reptext string, err error) {
	reccom := strings.Split(rec.Text, " ")
	reccomfix := []string{}
	for _, s := range reccom {
		if s != "" {
			reccomfix = append(reccomfix, s)
		}
	}
	reptext = ""
	//fmt.Printf("reccom:%v\n", reccom)
	mainhandler = &MainHandler{}
	mainhandler.init(rec)
	for _, handler := range mainhandler.allcommands {
		for _, com := range handler.getCommands() {
			//fmt.Printf("com:%s\n", com)
			if strings.EqualFold(reccomfix[0], com) {
				reptext, err = handler.handle(reccomfix)
				return
			}
		}
	}
	err = errors.New("command false ,try ?help")
	return
}

type MainHandler struct {
	allcommands []CocbotHandler
	rec         models.GMrecModel
}

func (this *MainHandler) init(rec models.GMrecModel) {
	this.allcommands = make([]CocbotHandler, 0)
	this.allcommands = append(this.allcommands, &HelpHandler{})
	this.allcommands = append(this.allcommands, &NewwarHandler{})
	this.allcommands = append(this.allcommands, &EditwarHandler{})
	this.allcommands = append(this.allcommands, &DelwarHandler{})
	this.allcommands = append(this.allcommands, &ShowwarHandler{})
	this.allcommands = append(this.allcommands, &ScoutHandler{})
	this.allcommands = append(this.allcommands, &StarHandler{})
	this.allcommands = append(this.allcommands, &CallHandler{})
	this.allcommands = append(this.allcommands, &TimerHandler{})
	this.allcommands = append(this.allcommands, &OpenedwarHandler{})
	this.rec = rec
	return
}

func (this *MainHandler) getGroupName() string {
	mapforGroupName := map[string]string{"19624531": "luckbot"}
	for k, v := range mapforGroupName {
		//fmt.Printf("gid:%s\n", this.rec.Group_id)
		//fmt.Printf("gn:%s\n", v)
		//fmt.Printf("boolean:%t\n", this.rec.Group_id == k)
		if this.rec.Group_id == k {
			return v
		}
	}
	return ""
}

func printWarData(content *models.WarDataModel) string {
	return fmt.Sprintf("War #%d Created \n %s VS %s \n %d vs %d \n War starts %s", content.Id, content.TeamA, content.TeamB, len(content.Battles), len(content.Battles), content.Begintime.Format("3:04PM MST 1/2/2006"))
}

type HelpHandler struct {
}

func (this *HelpHandler) handle(text []string) (result string, err error) {
	resultslice := make([]string, 0)
	resultslice = append(resultslice, "commands list:")
	resultslice = append(resultslice, "--------------")
	for _, handler := range mainhandler.allcommands {
		resultslice = append(resultslice, handler.getHelp())
	}
	result = strings.Join(resultslice, "\n")
	return
}
func (this *HelpHandler) getCommands() []string {
	return []string{"!help", "!h"}
}
func (this *HelpHandler) getHelp() string {
	return "!help/!h \n for help"
}

type NewwarHandler struct {
}

func (this *NewwarHandler) handle(text []string) (result string, err error) {
	//fmt.Printf("len text:%d\n", len(text))
	if len(text) < 3 {
		err = errors.New("i need more info\n" + this.getHelp())
		return
	}
	cout, err := strconv.Atoi(text[1])
	if err != nil || cout < 0 {
		err = errors.New("arg2 must be number\n" + this.getHelp())
		return
	}
	groupname := mainhandler.getGroupName()
	if groupname == "" {
		err = errors.New("group not found groupid:" + mainhandler.rec.Group_id)
		return
	}
	enemyname := strings.Join(text[2:len(text)], " ")
	id, err := models.AddWarData(groupname, enemyname, cout)
	if err != nil {
		fmt.Println(err.Error())
		err = errors.New("server error")
	} else {
		content, err := models.GetWarData(id)
		if err == nil {
			fmt.Printf("now time:%v\n", time.Now())
			fmt.Printf("default time:%v\n", content.Begintime)
			result = printWarData(content)
		}
	}

	return

}
func (this *NewwarHandler) getCommands() []string {
	return []string{"!war"}
}
func (this *NewwarHandler) getHelp() string {
	return "!war [:number] [:enemyclanname] \n for new War \n usage: \n !war 25 enemy "
}

type ShowwarHandler struct {
}

func (this *ShowwarHandler) handle(text []string) (result string, err error) {
	groupname := mainhandler.getGroupName()
	if groupname == "" {
		err = errors.New("group not found groupid:" + mainhandler.rec.Group_id)
		return
	}
	content, err := models.GetWarDatabyclanname(groupname)
	if err != nil {
		fmt.Println(err.Error())
		err = errors.New("server error")
		return
	}
	if !content.IsEnable {
		result = "no war being"
		return
	}
	if len(text) == 1 {
		result = fmt.Sprintf("War #%d  \n %s VS %s \n %d vs %d \n War starts %s\n",
			content.Id,
			content.TeamA,
			content.TeamB,
			len(content.Battles),
			len(content.Battles),
			content.Begintime.Format("3:04PM MST 1/2/2006"))
		keys := []int{}
		for num, _ := range content.Battles {
			keys = append(keys, num)
		}
		battlesresult := ""
		sort.Ints(keys)
		for _, num := range keys {
			hightstar := -1
			hightstars := "ZZZ"
			lineresult := fmt.Sprintf("||%d.%s ", num+1, content.Battles[num].Scoutstate)
			for _, caller := range content.Battles[num].Callers {
				if caller.Starstate > -1 && caller.Starstate < 4 {
					if caller.Starstate > hightstar {
						hightstar = caller.Starstate
						hightstars = caller.GetStarstate()
					}
					lineresult = lineresult + fmt.Sprintf("|%s %s", caller.Callername, caller.GetStarstate())
				} else {
					if time.Now().After(caller.Calledtime.Add(6 * time.Hour)) {
						lineresult = lineresult + fmt.Sprintf("|%s expried", caller.Callername)
					} else {
						expried := caller.Calledtime.Add(6 * time.Hour).Sub(time.Now())
						lineresult = lineresult + fmt.Sprintf("|%s %dh%dm ", caller.Callername, int(expried.Hours()), int(math.Mod(expried.Minutes(), 60)))
					}
				}
			}
			battlesresult = battlesresult + hightstars + lineresult + "\n"
		}
		result = result + battlesresult
	} else {
		num, err := strconv.Atoi(text[1])
		if err != nil {
			err = errors.New("arg2 must be number\n" + this.getHelp())
		} else {
			if num > len(content.Battles) {
				err = errors.New(fmt.Sprintf("last one is %d\n", len(content.Battles)))
			} else {
				battle := content.Battles[num-1]
				result = fmt.Sprintf("%d.%s ", num, battle.Scoutstate)
				for _, caller := range battle.Callers {
					if time.Now().After(caller.Calledtime.Add(6 * time.Hour)) {
						result = result + fmt.Sprintf("%s expried")
					} else {
						expried := caller.Calledtime.Add(6 * time.Hour).Sub(time.Now())
						result = result + fmt.Sprintf("%s %sh%sm ", caller.Callername, expried.Hours(), expried.Minutes())
					}
				}
			}
		}
	}

	return
}
func (this *ShowwarHandler) getCommands() []string {
	return []string{"!show"}
}
func (this *ShowwarHandler) getHelp() string {
	return "!show [:number]\n for show current war condition ,show all or just what u want\n usage:\n !show\n !show 1"
}

type ScoutHandler struct {
}

func (this *ScoutHandler) handle(text []string) (result string, err error) {
	if len(text) < 2 {
		err = errors.New("i need more info\n" + this.getHelp())
		return
	}
	groupname := mainhandler.getGroupName()
	if groupname == "" {
		err = errors.New("group not found groupid:" + mainhandler.rec.Group_id)
		return
	}
	content, err := models.GetWarDatabyclanname(groupname)
	if err != nil {
		fmt.Println(err.Error())
		err = errors.New("server error")
		return
	}
	if !content.IsEnable {
		result = "no war being"
		return
	}
	num1, err := strconv.Atoi(text[1])
	if err != nil {
		err = errors.New("arg1 need a number \n" + this.getHelp())
		return
	}
	if content.Battles[num1].Scoutstate == "needscout" || content.Battles[num1].Scoutstate == "scouted" {
		result = fmt.Sprintf("#%d already %s", num1, content.Battles[num1].Scoutstate)
		return
	}

	err = models.UpdateWarData(content.Id, bson.M{"$set": bson.M{fmt.Sprintf("battles.%d.scoutstate", num1): "needscout"}})
	if err != nil {
		fmt.Println(err.Error())
		err = errors.New("server error")
		return
	}
	result = fmt.Sprintf("%s mark #%d need to scout", mainhandler.rec.Name, num1)
	return
}
func (this *ScoutHandler) getCommands() []string {
	return []string{"!scout"}
}
func (this *ScoutHandler) getHelp() string {
	return "!scout [:number]\n for someone need scout \n usage:\n !scout 1"
}

type CallHandler struct {
}

func (this *CallHandler) handle(text []string) (result string, err error) {
	if len(text) < 2 {
		err = errors.New("i need more info\n" + this.getHelp())
		return
	}
	groupname := mainhandler.getGroupName()
	if groupname == "" {
		err = errors.New("group not found groupid:" + mainhandler.rec.Group_id)
		return
	}
	content, err := models.GetWarDatabyclanname(groupname)
	if err != nil {
		fmt.Println(err.Error())
		err = errors.New("server error")
		return
	}
	if !content.IsEnable {
		result = "no war being"
		return
	}
	num1, err := strconv.Atoi(text[1])
	if err != nil {
		err = errors.New("arg1 need a number \n" + this.getHelp())
		return
	}
	newcallnum := -1
	newbattle := content.Battles[num1]
	for num, call := range newbattle.Callers {
		if call.Callername == mainhandler.rec.Name {
			newcallnum = num
		}
	}

	if newcallnum == -1 {
		newcallp := &models.Caller{mainhandler.rec.Name, -1, time.Now()}
		newbattle.Callers = append(newbattle.Callers, *newcallp)
	} else {
		newbattle.Callers[newcallnum].Calledtime = time.Now()
	}

	err = models.UpdateWarData(content.Id, bson.M{"$set": bson.M{fmt.Sprintf("battles.%d", num1): newbattle}})
	if err != nil {
		fmt.Println(err.Error())
		err = errors.New("server error")
		return
	}
	result = fmt.Sprintf("#%d called by %s", num1, mainhandler.rec.Name)
	return
}
func (this *CallHandler) getCommands() []string {
	return []string{"!call"}
}
func (this *CallHandler) getHelp() string {
	return "!call [:number]\n for call someone \n usage:\n !call 1"
}

type StarHandler struct {
}

func (this *StarHandler) handle(text []string) (result string, err error) {
	if len(text) < 3 {
		err = errors.New("i need more info\n" + this.getHelp())
		return
	}
	num1, err := strconv.Atoi(text[1])
	if err != nil || num1 < 0 {
		err = errors.New("arg1 need a number \n" + this.getHelp())
		return
	}
	num2, err := strconv.Atoi(text[2])
	if err != nil || num2 < 0 || num2 > 3 {
		err = errors.New("arg2 need a number \n" + this.getHelp())
		return
	}
	groupname := mainhandler.getGroupName()
	if groupname == "" {
		err = errors.New("group not found groupid:" + mainhandler.rec.Group_id)
		return
	}
	content, err := models.GetWarDatabyclanname(groupname)
	if err != nil {
		fmt.Println(err.Error())
		err = errors.New("server error")
		return
	}
	if !content.IsEnable {
		result = "no war being"
		return
	}

	newcallnum := -1
	newbattle := content.Battles[num1]

	for num, call := range newbattle.Callers {
		if call.Callername == mainhandler.rec.Name {
			newcallnum = num
		}
	}

	newbattle.Scouted()
	if newcallnum == -1 {
		newcallp := &models.Caller{mainhandler.rec.Name, num2, time.Now()}
		newbattle.Callers = append(newbattle.Callers, *newcallp)
	} else {
		newbattle.Callers[newcallnum].Calledtime = time.Now()
		newbattle.Callers[newcallnum].Starstate = num2
	}

	err = models.UpdateWarData(content.Id, bson.M{"$set": bson.M{fmt.Sprintf("battles.%d", num1): newbattle}})

	if err != nil {
		fmt.Println(err.Error())
		err = errors.New("server error")
		return
	}
	switch num2 {
	case 0:
		result = fmt.Sprintf("%s win %d star at #%d\n%s", mainhandler.rec.Name, num2, num1, zerostar[rand.Intn(len(zerostar))])
	case 1:
		result = fmt.Sprintf("%s win %d star at #%d\n%s", mainhandler.rec.Name, num2, num1, onestar[rand.Intn(len(onestar))])
	case 2:
		result = fmt.Sprintf("%s win %d stars at #%d\n%s", mainhandler.rec.Name, num2, num1, twostar[rand.Intn(len(twostar))])
	case 3:
		result = fmt.Sprintf("%s win %d stars at #%d\n%s", mainhandler.rec.Name, num2, num1, threestar[rand.Intn(len(threestar))])
	}

	return
}
func (this *StarHandler) getCommands() []string {
	return []string{"!star"}
}
func (this *StarHandler) getHelp() string {
	return "!star [:number] [:number]\n for finish attack someone \n usage:\n !star 1 0"
}

type DelwarHandler struct {
}

func (this *DelwarHandler) handle(text []string) (result string, err error) {

	if len(text) < 2 {
		err = errors.New("i need more info\n" + this.getHelp())
		return
	}
	num, err := strconv.Atoi(text[1])
	if err != nil || num < 0 {
		err = errors.New("arg2 must be number\n" + this.getHelp())
		return
	}
	err = models.DelWarDatabyWarid(num)
	result = fmt.Sprintf("War #%d deleted", num)
	return

}
func (this *DelwarHandler) getCommands() []string {
	return []string{"!del"}
}
func (this *DelwarHandler) getHelp() string {
	return "!del [:number] \n for del a War \n usage: \n !del 5 "
}

type EditwarHandler struct {
}

func (this *EditwarHandler) handle(text []string) (result string, err error) {
	groupname := mainhandler.getGroupName()
	if groupname == "" {
		err = errors.New("group not found groupid:" + mainhandler.rec.Group_id)
		return
	}
	if len(text) < 3 {
		err = errors.New("i need more info\n" + this.getHelp())
		return
	}
	num1, err := strconv.Atoi(text[1])
	if err != nil || num1 < 0 {
		err = errors.New("arg2 must be number\n" + this.getHelp())
		return
	}
	num2, err := strconv.Atoi(text[2])
	isTime := strings.HasSuffix(text[2], "am") || strings.HasSuffix(text[2], "pm")
	if (err != nil || num2 < 0) && !isTime {
		err = errors.New("arg3 must be number or time(endwith am/pm)\n" + this.getHelp())
		return
	}

	if err == nil {
		enemyname := strings.Join(text[3:len(text)], " ")
		battles := make([]models.Battle, num2)
		battlep := &models.Battle{}
		battlep.Init()
		for index := range battles {
			battles[index] = *battlep
		}
		err := models.UpdateWarData(num1, bson.M{"$set": bson.M{"teamb": enemyname, "battles": battles}})
		if err == nil {
			content, err := models.GetWarData(num1)
			if err != nil {
				fmt.Println(err.Error())
				err = errors.New("server error")
			} else if !content.IsEnable {
				result = "no war being"
			} else {
				result = fmt.Sprintf("War #%d Edited \n %s VS %s \n %d vs %d \n War starts %s", content.Id, content.TeamA, content.TeamB, len(content.Battles), len(content.Battles), content.Begintime.Format("3:04PM MST 1/2/2006"))
			}

		}
	} else if isTime {
		now := time.Now()
		var h int
		if strings.HasSuffix(text[2], "am") {
			fmt.Printf("trim time:%s\n", strings.Trim(text[2], "am"))
			h, err = strconv.Atoi(strings.Trim(text[2], "am"))
			if err != nil || h < 0 {
				err = errors.New("wrong time format")
			} else {
				mi := 0
				if h > 12 {
					mi = int(math.Mod(float64(h), 100))
					h = int(h / 100)
				}
				y, m, d := now.Date()
				if now.Hour()*60+now.Minute() > h*60+mi { //if befor now so will be tmw
					d++
				}
				newtime := time.Date(y, m, d, h, mi, 0, 0, time.Local)
				err := models.UpdateWarData(num1, bson.M{"$set": bson.M{"begintime": newtime}})
				if err == nil {
					content, err := models.GetWarData(num1)
					if err != nil {
						fmt.Println(err.Error())
						err = errors.New("server error")
					} else if !content.IsEnable {
						result = "no war being"
					} else {
						result = fmt.Sprintf("War #%d Edited \n %s VS %s \n %d vs %d \n War starts %s", content.Id, content.TeamA, content.TeamB, len(content.Battles), len(content.Battles), content.Begintime.Format("3:04PM MST 1/2/2006"))
					}
				}
			}
		} else {

			h, err = strconv.Atoi(strings.Trim(text[2], "pm"))
			fmt.Printf("trim time:%d\n", h)
			if err != nil || h < 0 {
				err = errors.New("wrong time format")
			} else {
				mi := 0
				if h > 12 {
					mi = int(math.Mod(float64(h), 100))
					h = int(h / 100)
				}
				h = h + 12
				y, m, d := now.Date()
				if now.Hour()*60+now.Minute() > h*60+mi { //if befor now so will be tmw
					d++
				}
				fmt.Printf("trim time mi:%d\n", mi)
				newtime := time.Date(y, m, d, h, mi, 0, 0, time.Local)
				err := models.UpdateWarData(num1, bson.M{"$set": bson.M{"begintime": newtime}})
				if err == nil {
					content, err := models.GetWarData(num1)
					if err != nil {
						fmt.Println(err.Error())
						err = errors.New("server error")
					} else if !content.IsEnable {
						result = "no war being"
					} else {
						result = fmt.Sprintf("War #%d Edited \n %s VS %s \n %d vs %d \n War starts %s", content.Id, content.TeamA, content.TeamB, len(content.Battles), len(content.Battles), content.Begintime.Format("3:04PM MST 1/2/2006"))
					}
				}
			}
		}

	} else {
		err = errors.New("wrong arg \n" + this.getHelp())
	}

	return

}
func (this *EditwarHandler) getCommands() []string {
	return []string{"!edit"}
}
func (this *EditwarHandler) getHelp() string {
	return "!edit [:number] [:number] [:clanname] /[:time](endwith am/pm)  \n for edit a War \n usage: \n !edit 5 25 enemy \n !edit 5 130am\n !edit 5 7pm"
}

type TimerHandler struct {
}

func (this *TimerHandler) handle(text []string) (result string, err error) {

	if len(text) < 2 {
		err = errors.New("i need more info\n" + this.getHelp())
		return
	}
	num, err := strconv.Atoi(text[1])
	if err != nil || num < 0 {
		err = errors.New("arg2 must be number\n" + this.getHelp())
		return
	}
	content, err := models.GetWarData(num)
	if err == nil {
		result = fmt.Sprintf("War #%d \n %s VS %s \n %d vs %d \n  %d hours till war start!", content.Id, content.TeamA, content.TeamB, len(content.Battles), len(content.Battles), int(content.Begintime.Sub(time.Now()).Hours()))
	}
	return

}
func (this *TimerHandler) getCommands() []string {
	return []string{"!timer"}
}
func (this *TimerHandler) getHelp() string {
	return "!timer [:number] \n for when a War begin \n usage: \n !timer 5 "
}

type OpenedwarHandler struct {
}

func (this *OpenedwarHandler) handle(text []string) (result string, err error) {
	groupname := mainhandler.getGroupName()
	if groupname == "" {
		err = errors.New("group not found groupid:" + mainhandler.rec.Group_id)
		return
	}
	content, err := models.GetWarDatabyclanname(groupname)
	if err != nil {
		fmt.Println(err.Error())
		err = errors.New("server error")
		return
	}
	if !content.IsEnable {
		result = "no war being"
		return
	}

	result = fmt.Sprintf("War #%d  \n %s VS %s \n %d vs %d \n War starts %s\n",
		content.Id,
		content.TeamA,
		content.TeamB,
		len(content.Battles),
		len(content.Battles),
		content.Begintime.Format("3:04PM MST 1/2/2006"))
	keys := []int{}
	for num, _ := range content.Battles {
		keys = append(keys, num)
	}
	battlesresult := ""
	sort.Ints(keys)
	for _, num := range keys {
		hightstar := -1
		hightstars := "ZZZ"
		called := false
		lineresult := fmt.Sprintf("||%d.%s ", num+1, content.Battles[num].Scoutstate)
		for _, caller := range content.Battles[num].Callers {
			if caller.Starstate > -1 && caller.Starstate < 4 {
				if caller.Starstate > hightstar {
					hightstar = caller.Starstate
					hightstars = caller.GetStarstate()
				}
				lineresult = lineresult + fmt.Sprintf("|%s %s", caller.Callername, caller.GetStarstate())
			} else {
				if time.Now().After(caller.Calledtime.Add(6 * time.Hour)) {
					continue
				} else {
					called = true
				}
			}
		}
		if hightstar != 3 && called == false {
			battlesresult = battlesresult + hightstars + lineresult + " open\n"
		}
	}
	result = result + battlesresult

	return
}
func (this *OpenedwarHandler) getCommands() []string {
	return []string{"!open"}
}
func (this *OpenedwarHandler) getHelp() string {
	return "!open \n for show the opening \n usage:\n !open"
}
