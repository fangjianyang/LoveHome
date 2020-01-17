package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"LoveHome/models"
	//"encoding/json"
)

type AreaController struct {
	beego.Controller
}

func (c *AreaController) RetData(resp map[string]interface{}){
        c.Data["json"] = resp
        c.ServeJSON()
}

func (c *AreaController) GetArea() {
    beego.Info("AreaController GetArea")

    resp := make(map[string]interface{})
    resp["errno"] =0
    resp["errmsg"] = "success"
    defer c.RetData(resp)

    // 1 get data from session

    // 2 get data from db
    var area  []models.Area
    ormHandel := orm.NewOrm()

    //err := ormHandel.Read(&area)
    num,err := ormHandel.QueryTable("area").All(&area)
    if err != nil{
        beego.Info("query data from db error!")

        resp["errno"] =12110
        resp["errmsg"] = "query err"
        return
    }
    if num == 0{
        beego.Info("query no data from db !")

        resp["errno"] =12111
        resp["errmsg"] = "query no data"
        return
    }



    resp["data"] = area
    //beego.Info("the data from mysql is ",area)

    // 3 paking data as json return to client

}