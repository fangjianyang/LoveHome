package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"LoveHome/models"
	)

type OrderController struct {
	beego.Controller
}

func(this *OrderController)RetData(resp map[string]interface{}){
	this.Data["json"] = resp
	this.ServeJSON()
}


func(this *OrderController)GetOrders(){
    beego.Info("OrderController GetOrders")
    resp := make(map[string]interface{})
	defer this.RetData(resp)

	// 1 get user_id from session
	user_id := this.GetSession("user_id")

	//2 get param from url
	role := this.GetString("role")
	beego.Info("get param from user : ", role)
	if role == "custom"{
	    // get orders of this person
	    orders := []models.OrderHouse{}

	    ormHandel := orm.NewOrm()
	    qs := ormHandel.QueryTable("OrderHouse")

	    user := models.User{Id:user_id.(int)}
	    qs.Filter("user__id",user_id.(int)).All(&orders)


	    for _,order := range orders {
	        order.User = &user
	        ormHandel.LoadRelated(order,"User")
	    }

	    respData := make(map[string]interface{})
	    respData["order"] = orders
	    resp["data"] = respData
        resp["errno"] =0
        resp["errmsg"] = "success"
	    return
	} else if role == "landlord" {
	    // todo
	    return
	} else {
        resp["errno"] = models.RECODE_DBERR
        resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
        return
	}



}