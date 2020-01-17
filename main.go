package main

import (
	_ "LoveHome/routers"
	"github.com/astaxie/beego"
    "strings"
    "net/http"
    "github.com/astaxie/beego/context"
    _"LoveHome/models"
)


func ignoreStaticPath() {
	beego.InsertFilter("/", beego.BeforeRouter, TransparentStatic)
	beego.InsertFilter("/*", beego.BeforeRouter, TransparentStatic)
}

func TransparentStatic(ctx *context.Context) {
	orpath := ctx.Request.URL.Path
	beego.Debug("request url: ", orpath)
	//如果请求uri还有api字段,说明是指令应该取消静态资源路径重定向
	if strings.Index(orpath, "api") >= 0 {
		return
	}
	http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/html/"+ctx.Request.URL.Path)
}

func main() {
    ignoreStaticPath()
	beego.Run()
}

