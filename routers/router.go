package routers

import (
	"LoveHome/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/api/v1.0/areas", &controllers.AreaController{},"get:GetArea")
    beego.Router("/api/v1.0/houses/index",&controllers.HouseIndexController{},"get:GetHouseIndex")
    beego.Router("/api/v1.0/session",&controllers.SessionController{},"get:GetSessionData;delete:DeleteSession")
    beego.Router("/api/v1.0/users",&controllers.UserController{},"post:Reg")
    beego.Router("/api/v1.0/sessions",&controllers.SessionController{},"post:Login")
    beego.Router("/api/v1.0/user/avatar",&controllers.UserController{},"post:PostAvatar")
    beego.Router("/api/v1.0/user",&controllers.UserController{},"get:GetUserData")
    beego.Router("/api/v1.0/user/name",&controllers.UserController{},"put:UpdateName")
    beego.Router("/api/v1.0/user/auth",&controllers.UserController{},"get:AuthUser;post:PostAuth")
    beego.Router("/api/v1.0/user/houses",&controllers.HousesController{},"get:GetHouses")
    beego.Router("/api/v1.0/houses",&controllers.HousesController{},"post:PostHousesData")
    //beego.Router("/api/v1.0/houses/images",&controllers.HousesController{},"post:PostHousesImages")

    // house detail info /api/v1.0/houses/4
    beego.Router("/api/v1.0/houses/?:id",&controllers.HousesController{},"get:GetHouseDetail")

    beego.Router("/api/v1.0/user/orders",&controllers.OrderController{},"get:GetOrders")
}
