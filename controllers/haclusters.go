package controllers

import (
	"fmt"

	"github.com/beego/beego/v2/server/web"
)

type HAClustersController struct {
	web.Controller
}

func (hcc *HAClustersController) Get() {
	// TODO: handle get here
	fmt.Println("handle get request in HAClustersController.")
	result := struct {
		Data string
	}{
		Data: "hello",
	}
	hcc.Data["json"] = &result
	hcc.ServeJSON()
}
