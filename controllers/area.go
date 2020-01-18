package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"LoveHome/models"
	_"github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/cache"
	"time"
	"encoding/json"
	"fmt"
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


    cacheHandel, err := cache.NewCache("redis", `{"key":"lovehome","conn":":6379","dbNum":"0"}`)
    if err != nil{
        resp["errno"] =12110
        resp["errmsg"] = "query err"
        return
    }


    // 1 get data from cache,mush Ummarshal from json to the right fmt
    if areaData := cacheHandel.Get("area"); areaData != nil{
        var areas  []models.Area
        err1 := json.Unmarshal(areaData.([]byte), &areas)
        if err1 != nil{
            beego.Error("Ummarshal data from redis failed!")
            return
        }
        beego.Info("areaData : " , areas)
        resp["data"] = areas
        return
    }

    // 2 get data from db
    var area  []models.Area
    ormHandel := orm.NewOrm()

    //err := ormHandel.Read(&area)
    num,errDb := ormHandel.QueryTable("area").All(&area)
    if errDb != nil{
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
    beego.Info("the data from mysql is ",area)

    // 3 paking data as json return to client
    jsonStr,errFmt := json.Marshal(area)
    if errFmt != nil {
        beego.Info("encoding err")
        return
    }
    //beego.Info("JsonStr : ",jsonStr)
    fmt.Printf("JsonStr %s",jsonStr)
    errFmt = cacheHandel.Put("area",jsonStr,time.Second * 3600)
    if errFmt != nil{
        beego.Error("cache err")
        return
    }

}