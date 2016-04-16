package main

import (
	"bot/healthcheck"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
)

func main() {
	if beego.RunMode == "dev" {
		beego.DirectoryIndex = true
		beego.StaticDir["/swagger"] = "swagger"
	}
	toolbox.AddHealthCheck("database", &healthcheck.DatabaseCheck{})
	//collector.Run()
	beego.Run()
}
