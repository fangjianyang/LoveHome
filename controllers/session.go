package controllers

import (
	"github.com/astaxie/beego"
	// "github.com/astaxie/beego/session"
	 "github.com/astaxie/beego/orm"
	 "LoveHome/models"
	 "encoding/json"
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
    // 太奇葩了，if 这个errcode = 0 界面就没有按钮
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

func (this *SessionController) Login(){
    resp := make(map[string]interface{})
    resp["errno"] =0
    resp["errmsg"] = "success"
    defer this.RetData(resp)

    err1 := json.Unmarshal(this.Ctx.Input.RequestBody, &resp)
    if err1 != nil{
        resp["errno"] =12345
        resp["errmsg"] = "requestBody error"
        beego.Info("Unmarshal error")
        return
    }
    beego.Info(resp)
    // get user data
    ormHandel := orm.NewOrm()
    user := models.User{Name:resp["mobile"].(string)}
    userData := ormHandel.QueryTable("user");
    err := userData.Filter("mobile", resp["mobile"].(string)).One(&user)

    if err != nil{
        resp["errno"] =12345
        resp["errmsg"] = "requestBody error"
        return
    }

    // judge if the data is right
    beego.Info(user)
    if user.Password_hash != resp["password"] {
        resp["errno"] =12345
        resp["errmsg"] = "requestBody error"
        return
    }

    // add session
    if len(user.Name) > 0{
        this.SetSession("name",user.Name)
    }else {
        this.SetSession("name",resp["mobile"].(string))
    }
    this.SetSession("mobile",resp["mobile"].(string))
    this.SetSession("user_id",user.Id)

    // return result
    resp["errno"] =0
    resp["errmsg"] = "success"
}