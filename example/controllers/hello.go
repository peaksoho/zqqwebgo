package controllers

import (
	"fmt"
	"zqqwebgo"
)

type Hello struct {
	zqqwebgo.Controller
}

func (this *Hello) Index() {
	//this.Params = this.Ctx.R.Form
	fmt.Println(this.Segment)
	fmt.Println(this.Params)
	fmt.Fprintf(this.Ctx.W, "Hello/Index")
}

func (this *Hello) World() {
	//this.Params = this.Ctx.R.Form
	fmt.Println("Segment:", this.Segment)
	fmt.Println("Params:", this.Params)
	fmt.Println("ParamsGet:", this.ParamsGet)
	fmt.Println("ParamsPost:", this.ParamsPost)

	this.TplNames = "world.html"
	this.Data["name"] = this.ParamsPost.Get("name")
	this.Data["age"] = this.ParamsPost.Get("age")
	this.Data["phone"] = this.ParamsPost.Get("phone")

	//fmt.Fprintf(this.Ctx.W, "Hello world!"+strings.Join(this.Segment, ","))
}
