package controllers

import (
	"github.com/astaxie/beego"
	// github.com/astaxie/beego/session
	// "github.com/astaxie/beego/orm"
	 "LoveHome/models"
	// "encoding/json"
)

type SessionController struct {
	beego.Controller
}

func (c *SessionController) MkRet(errcode string,resp map[string]interface{}){
    resp["errno"] = errcode
    resp["errmsg"] = models.RecodeText(errcode)
}

func (c *SessionController) RetData(resp map[string]interface{}){
    c.Data["json"] = resp
    c.ServeJSON()
}

func (this *SessionController) GetSessionData() {
    resp := make(map[string]interface{})
    // this.MkRet(models.RECODE_OK,resp)
    // 太奇葩了，加入这个errcode = 0 界面就没有按钮
    resp["errno"] = models.RECODE_DATAERR
    resp["errmsg"] = models.RecodeText(models.RECODE_DATAERR)
    defer this.RetData(resp)

    user := models.User{}
    name := this.GetSession("name")
    if name != nil{
        user.Name = name.(string)
        //this.MkRet(models.RECODE_OK,resp)
        resp["errno"] = models.RECODE_OK
        resp["errmsg"] = models.RecodeText(models.RECODE_OK)
        resp["data"] = user
    }
}

func  (this *SessionController) DeleteSession() {
    resp := make(map[string]interface{})
    defer this.RetData(resp)

    this.DelSession("name")
    resp["errno"] = models.RECODE_OK
    resp["errmsg"] = models.RecodeText(models.RECODE_OK)
}