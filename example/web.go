package main

import (
	"zqqwebgo"
	"zqqwebgo/example/controllers"
)

func main() {
	zqqwebgo.AddController("Hello", &controllers.Hello{})     //必须,添加控制层后可以直接访问“/controllerName/actionName”，而不必注册路由
	zqqwebgo.RouterRegister["/hello/index"] = "/Hello/index/" //注册路由，可选，值为“/controllerName/actionName”
	zqqwebgo.RouterRegister["/(.*)"] = "/Hello/world/"        //注册路由，可选，值为“/controllerName/actionName”
	zqqwebgo.Run()
}
