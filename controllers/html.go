package controllers

import (
	"net/http"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type RootController struct {
	web.Controller
}

func (rc *RootController) Get() {
	logs.Debug("handle index request")

	// TODO: check session
	if true {
		// default find templete from views folder
		rc.TplName = "static/index.html"
		return
	}
	rc.Redirect("/login", http.StatusOK)
}

type IndexController struct {
	web.Controller
}

func (ic *IndexController) Get() {
	// default find templete from views folder
	ic.TplName = "static/index.html"
}
