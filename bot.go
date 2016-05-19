package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/luckily248/bot/handler"
	"github.com/luckily248/bot/models"
)

//main
func Run() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	http.HandleFunc("/bot", WarDataController)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
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
	go handle(rec)
	return
}
func handle(rec models.GMrecModel) {
	reptext, err := handler.HandlecocText(rec)
	fmt.Printf("reptextlen:%d\n", utf8.RuneCountInString(reptext))
	rep := &models.GMrepModel{}
	rep.InitbyGID(rec.Group_id)
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
