package controllers

import (
	"github.com/astaxie/beego"
)

type HouseIndexController struct {
	beego.Controller
}

func(this*HouseIndexController)RetData(resp map[string]interface{}){
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this*HouseIndexController)GetHouseIndex(){
	resp := make(map[string]interface{})

	resp["errno"] =12345
	resp["errmsg"] = "success"
	this.RetData(resp)
}
