package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"LoveHome/models"
	"encoding/json"
	"strconv"
)

type HousesController struct {
	beego.Controller
}

func(this *HousesController)RetData(resp map[string]interface{}){
	this.Data["json"] = resp
	this.ServeJSON()
}

func(this *HousesController)GetHouses(){
    resp := make(map[string]interface{})
    defer this.RetData(resp)

    // 1 get userid from session
    user_id := this.GetSession("user_id")

    houses := []models.House{}
    // get data from db
    ormHandel := orm.NewOrm()
    qs := ormHandel.QueryTable("house")
    num,err := qs.Filter("user__id",user_id.(int)).All(&houses)
    beego.Info("the data from the tble is : ",houses)
    if err != nil{
        resp["errno"] = models.RECODE_DBERR
        resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
        return
    }
    if num == 0 {
        resp["errno"] = models.RECODE_NODATA
        resp["errmsg"] = models.RecodeText(models.RECODE_NODATA)
        return
    }

    respData := make(map[string]interface{})
    if num != 0 {
        respData["houses"] = houses
    }
    resp["data"] = respData
    resp["errno"] = models.RECODE_OK
    resp["errmsg"] = models.RecodeText(models.RECODE_OK)
}

func(this *HousesController)PostHousesData(){
    resp := make(map[string]interface{})
    defer this.RetData(resp)

    // get data from client
    reqData := make(map[string]interface{})
    err1 := json.Unmarshal(this.Ctx.Input.RequestBody, &reqData)
    if err1 != nil{
        resp["errno"] =12345
        resp["errmsg"] = "requestBody error"
        beego.Info("Unmarshal error")
        return
    }
    // judge if this data is valid

    // store data to db
    house := models.House{}

    area_id,_ := strconv.Atoi(reqData["area_id"].(string))
    area :=models.Area{Id:area_id}
    house.Area = &area
    house.Title = reqData["title"].(string)
    price ,_ := strconv.Atoi(reqData["price"].(string))
    house.Price = price
    house.Address = reqData["address"].(string)
    roomCnt,_ := strconv.Atoi(reqData["room_count"].(string))
    house.Room_count = roomCnt
    acreage,_ := strconv.Atoi(reqData["acreage"].(string))
    house.Acreage = acreage
    house.Unit = reqData["unit"].(string)
    capacity,_ := strconv.Atoi(reqData["capacity"].(string))
    house.Capacity = capacity
    house.Beds = reqData["beds"].(string)

    deposit,_ := strconv.Atoi(reqData["deposit"].(string))
    house.Deposit = deposit
    minDay,_ := strconv.Atoi(reqData["min_days"].(string))
    house.Min_days = minDay
    maxDay,_ := strconv.Atoi(reqData["max_days"].(string))
    house.Max_days = maxDay
    faclities := []models.Facility{}
    for _,fid := range  reqData["facility"].([]interface{}){
        f_id,_:= strconv.Atoi(fid.(string))
        fac := models.Facility{Id:f_id}
        faclities = append(faclities,fac)
    }
    user_id := this.GetSession("user_id").(int)
    user := models.User{Id:user_id}
    house.User = &user
    ormHandel := orm.NewOrm()
    // err here
    houseId,dbErr := ormHandel.Insert(&house)
    if dbErr != nil {
        resp["errno"] = models.RECODE_DBERR
        resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
        beego.Error("Inert house data to db error!")
        return
    }

    house.Id = int(houseId)

    m2m := ormHandel.QueryM2M(&house,"Facilities")
    num,errM2M := m2m.Add(faclities)
    if errM2M !=nil{
        resp["errno"] = models.RECODE_DBERR
        resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
        beego.Error("Add faclities data to db error!")
        return
    }

    if num == 0 {
        resp["errno"] = models.RECODE_DBERR
        resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
         beego.Error("Add faclities data to db error!")
        return
    }
    beego.Info("Add faclities success num =  ",num)

    respData := make(map[string]interface{})
    respData["house_id"] = strconv.Itoa(house.Id)
    resp["errno"] = models.RECODE_OK
    resp["errmsg"] = models.RecodeText(models.RECODE_OK)
}

