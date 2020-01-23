package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"LoveHome/models"
	"encoding/json"
	"path"
	"github.com/tedcy/fdfs_client"
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
    this.SetSession("user_id",user.Id)
    this.SetSession("mobile",user.Mobile)

}

func (this *UserController) PostAvatar(){
    resp := make(map[string]interface{})
    defer this.RetData(resp)
    fileData, hd, err := this.GetFile("avatar")
    if err != nil {
        resp["errno"] = 40012345
        resp["errmsg"] = "reg err"
        return
    }
    // get suffix of the file
    suffix := path.Ext(hd.Filename)
	fileBuffer := make([]byte, hd.Size)
	_, err = fileData.Read(fileBuffer)
	if err != nil {
		resp["errno"] = models.RECODE_REQERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
		return
	}

    // upload the file
    client, errClient := fdfs_client.NewClientWithConfig("conf/client.conf")
    defer client.Destory()
    if errClient != nil {
        beego.Error(errClient)
        resp["errno"] = models.RECODE_REQERR
        resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
        return
    }
    fileId, errUp := client.UploadByBuffer( fileBuffer, suffix[1:])
    if  errUp != nil {
        resp["errno"] = models.RECODE_REQERR
        resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
        return
    }
    //4.从session里拿到user_id
    user_id := this.GetSession("user_id")
    var user models.User
    //5.更新用户数据库中的内容
    ormHandel := orm.NewOrm()
    qs := ormHandel.QueryTable("user")
    qs.Filter("Id",user_id).One(&user)
    user.Avatar_url = fileId

    _,errUpDb := ormHandel.Update(&user)
    if errUpDb != nil{
        resp["errno"] = models.RECODE_REQERR
        resp["errmsg"] = models.RecodeText(models.RECODE_REQERR)
        return
    }

    urlMap:= make(map[string]string)
    urlMap["avatar_url"] = "http://192.168.4.143:8080/"+fileId
    resp["errno"] = models.RECODE_OK
    resp["errmsg"] = models.RecodeText(models.RECODE_OK)
    resp["data"] = urlMap
}

func (this *UserController)  GetUserData(){
     resp := make(map[string]interface{})
     defer this.RetData(resp)

     // 1 get userid from session
     user_id := this.GetSession("user_id")

     // 2 get data from db where user_id = "user_id"
     user := models.User{Id:user_id.(int)}
     ormHandel := orm.NewOrm()
     dbErr := ormHandel.Read(&user)
     if dbErr != nil{
        resp["errno"] = models.RECODE_DBERR
        resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
        return
     }

     resp["data"] = &user
     resp["errno"] = models.RECODE_OK
     resp["errmsg"] = models.RecodeText(models.RECODE_OK)
}

func (this *UserController)  UpdateName(){
     resp := make(map[string]interface{})
     defer this.RetData(resp)

     // gain the userid from session
     user_id := this.GetSession("user_id")
     beego.Info("Get user id from Session :",user_id)

     // get data from the client
     UserName := make(map[string]string)

     if err := json.Unmarshal(this.Ctx.Input.RequestBody, &UserName); err != nil {
        resp["errno"] = models.RECODE_PARAMERR
        resp["errmsg"] = models.RecodeText(models.RECODE_PARAMERR)
        beego.Error("Get user name failed !")
        return
     }

    // update the db where uerid = "xxx"
    ormHandel := orm.NewOrm()
    user := models.User{Id:user_id.(int)}
    dbErr := ormHandel.Read(&user)
    if dbErr != nil {
        resp["errno"] = models.RECODE_DBERR
        resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
        beego.Error("Read user info from mysql failed !")
        return
    }
    user.Name = UserName["name"]

    _, upErr :=ormHandel.Update(&user)
    if upErr != nil {
        resp["errno"] = models.RECODE_DBERR
        resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
        beego.Error("update user name to mysql failed !")
        return
    }

     // change the name of the session
     this.SetSession("name",UserName["name"])

     // package data to client
    resp["data"] = UserName
    resp["errno"] = models.RECODE_OK
    resp["errmsg"] = models.RecodeText(models.RECODE_OK)
}

func (this *UserController)  AuthUser(){
     resp := make(map[string]interface{})
     defer this.RetData(resp)

     // gain the userid from session
     user_id := this.GetSession("user_id")
     beego.Info("Get user id from Session :",user_id.(int))

    // update the db where uerid = "xxx"
    ormHandel := orm.NewOrm()
    user := models.User{Id:user_id.(int)}
    dbErr := ormHandel.Read(&user)
    if dbErr != nil {
        resp["errno"] = models.RECODE_DBERR
        resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
        beego.Error("Read user info from mysql failed !")
        return
    }

    // change session
    this.SetSession("user_id",user.Id)

    // package data to client
    resp["data"] = &user
    resp["errno"] = models.RECODE_OK
    resp["errmsg"] = models.RecodeText(models.RECODE_OK)
}

func (this *UserController)PostAuth(){
    resp := make(map[string]interface{})
    defer this.RetData(resp)

    // gain the userid from session
    user_id := this.GetSession("user_id")
    beego.Info("Get user id from Session :",user_id.(int))


    // get data from the client
    UserInfo := make(map[string]string)

    if err := json.Unmarshal(this.Ctx.Input.RequestBody, &UserInfo); err != nil {
        resp["errno"] = models.RECODE_PARAMERR
        resp["errmsg"] = models.RecodeText(models.RECODE_PARAMERR)
        beego.Error("Get user name failed !")
        return
    }

    // update the db where uerid = "xxx"
    ormHandel := orm.NewOrm()
    user := models.User{Id:user_id.(int)}
    dbErr := ormHandel.Read(&user)
    if dbErr != nil {
        resp["errno"] = models.RECODE_DBERR
        resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
        beego.Error("Read user info from mysql failed !")
        return
    }
    user.Real_name = UserInfo["real_name"]
    user.Id_card = UserInfo["id_card"]

    _, upErr :=ormHandel.Update(&user)
    if upErr != nil {
        resp["errno"] = models.RECODE_DBERR
        resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
        beego.Error("update user name to mysql failed !")
        return
    }

    // change session
    this.SetSession("user_id",user.Id)

    // package data to client
    resp["errno"] = models.RECODE_OK
    resp["errmsg"] = models.RecodeText(models.RECODE_OK)
}

