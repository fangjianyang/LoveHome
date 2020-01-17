package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"LoveHome/models"
	"encoding/json"
)

type UserController struct {
	beego.Controller
}

func(this *UserController)RetData(resp map[string]interface{}){
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *UserController)Reg(){
	resp := make(map[string]interface{})
    resp["errno"] =0
    resp["errmsg"] = "success"
	defer this.RetData(resp)

    // get data from client
    err1 := json.Unmarshal(this.Ctx.Input.RequestBody, &resp)
    if err1 != nil{
    	resp["errno"] =12345
    	resp["errmsg"] = "requestBody error"
    	beego.Info("Unmarshal error")
        return
    }
    beego.Info(resp)

    // insert data into db  mobile:111 password:111 sms_code:111
    ormHandel := orm.NewOrm()
    user := models.User{}
    user.Password_hash = resp["password"].(string)
    user.Name = resp["mobile"].(string)
    user.Mobile = resp["mobile"].(string)

    _,err := ormHandel.Insert(&user)
    if err != nil{
        resp["errno"] = 40012345
        resp["errmsg"] = "reg err"
        beego.Info("Insert error")
        return
    }
     beego.Info("reg success")

     // set session
     this.SetSession("name",user.Name)
}