package controllers

import (
	"bot/handler"
	"bot/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/astaxie/beego"
)

// get wardata
type WarDataController struct {
	beego.Controller
}

func (c *WarDataController) URLMapping() {
	c.Mapping("NewWar", c.NewWar)
	c.Mapping("GetWar", c.GetWar)
	c.Mapping("Bot", c.Bot)
}

// @Title Post
// @Description post new WarDataModel by TeamA TeamB
// @Param	TeamA		form 	string	true		"new war TeamA"
// @Param	TeamB		form 	string	true		"new war TeamB"
// @Success 200 {string} string
// @Failure 404  not found
// @router /newwar [post]
func (c *WarDataController) NewWar() {
	wardata := models.WarDataModel{}
	if err := c.ParseForm(&wardata); err != nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = err.Error()
		fmt.Printf("err:%s\n", err.Error())
	}
	fmt.Printf("teama:%s\n", wardata.TeamA)
	fmt.Printf("teamb:%s\n", wardata.TeamB)
	if wardata.TeamA == "" || wardata.TeamB == "" {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = string("someteam empty")
	} else {
		id, err := models.AddWarData(wardata.TeamA, wardata.TeamB, 25)
		if err != nil {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = err.Error()
			fmt.Printf("err:%s\n", err.Error())
		} else {
			c.Data["json"] = id
			fmt.Printf("id:%s\n", id)
		}
	}
	c.ServeJson()
}

// @Title Get
// @Description get WarDataModel by clanname
// @Param	clanname		path 	string	true		"The key for clan"
// @Success 200 {object} models.WarDataModel
// @Failure 404  not found
// @Failure 403 :clanname is empty
// @router /:clanname [get]
func (c *WarDataController) GetWar() {

	clanname := c.Ctx.Input.Params[":clanname"]
	fmt.Printf("clanname:%s\n", clanname)
	if clanname == "" {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = string("is empty")
	} else {
		ob, err := models.GetWarDatabyclanname(clanname)
		if err != nil {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = err.Error()
			fmt.Printf("err:%s\n", err.Error())
		} else {
			c.Data["json"] = ob
			fmt.Printf("ob:%v\n", ob)
		}
	}
	c.ServeJson()
}

// @Title Post
// @Description get rep from callerbot
// @Param   body        body    models.GMrecModel   true        "The object content"
// @Success 200 {object} models.GMrepModel
// @Failure 403  body is empty
// @router /bot [post]
func (c *WarDataController) Bot() {
	body := c.Ctx.Input.RequestBody
	fmt.Printf("body:%s\n", body)
	var rec models.GMrecModel
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &rec); err != nil {
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
	httpPost(buff)
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
