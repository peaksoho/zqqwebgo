package zqqwebgo

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
)

var (
	RouterRegister     map[string]string
	ControllerRegistor map[string]ControllerInterface
)

type ControllerInterface interface {
}

type Router struct {
	Segments   []string
	CtrlAction []string
}

func init() {
	ControllerRegistor = make(map[string]ControllerInterface)
	RouterRegister = make(map[string]string)
}

func (p *Router) Handler(w http.ResponseWriter, r *http.Request) bool {
	lenOfSegment := len(p.CtrlAction)

	//fmt.Println("CtrlAction:", p.CtrlAction)

	if 0 < lenOfSegment {
		ca1 := strings.Title(p.CtrlAction[0])
		ct := ControllerRegistor[ca1]
		cv := reflect.ValueOf(ct)
		if !cv.IsValid() {
			fmt.Fprintf(w, "No Controller ["+ca1+"]")
			return false
		}
		var method reflect.Value
		in := make([]reflect.Value, 2)
		in[0] = reflect.ValueOf(w)
		in[1] = reflect.ValueOf(r)
		method = cv.MethodByName("Init") //controller initialization
		method.Call(in)

		r.ParseForm() //解析参数
		method = cv.MethodByName("SetParams")
		method.Call(in)

		method = cv.MethodByName("SetSegment")
		method.Call([]reflect.Value{reflect.ValueOf(p.Segments)}) //将urlpath设为参数

		in = make([]reflect.Value, 0)

		defer func() { //输出部分
			method = cv.MethodByName("OutPutString") //如果设置了OutString就直接输出
			method.Call(in)
			method = cv.MethodByName("ParseTpl") //如果设置了TplNames就解析html模板
			method.Call(in)

			//static file server
			for prefix, staticDir := range StaticDir {
				if r.URL.Path == "/favicon.ico" {
					file := staticDir + r.URL.Path
					_, err := os.Stat(file)
					if err == nil {
						http.ServeFile(w, r, file)
					}
					return
				}
				if strings.HasPrefix(r.URL.Path, prefix) {
					file := staticDir + r.URL.Path[len(prefix):]
					_, err := os.Stat(file)
					if err == nil {
						http.ServeFile(w, r, file)
					}
					return
				}
			}
		}()

		ca2 := "Index"
		if 1 < lenOfSegment {
			ca2 = strings.Title(p.CtrlAction[1])
			method = cv.MethodByName(ca2)
		} else {
			method = cv.MethodByName("Index")
		}
		if !method.IsValid() {
			fmt.Fprintf(w, "No method ["+ca2+"] in the Controller ["+ca1+"]")
			return false
		} else {
			method.Call(in)
		}
		method = cv.MethodByName("Destructor")
		method.Call(in)
	}
	return true
}

func (p *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Segments = make([]string, 2)
	urlpath := strings.Trim(r.URL.Path, "/")
	p.Segments = strings.Split(urlpath, "/")

	//fmt.Println(urlpath)

	if 0 < len(RouterRegister) {
		p.CtrlAction = make([]string, 2)
		for router, ca := range RouterRegister {
			router = strings.Trim(router, "/")

			//fmt.Println("router:", router)
			//fmt.Println("urlpath:", strings.ToLower(urlpath))

			b, err := regexp.MatchString("^"+router+"$", strings.ToLower(urlpath))
			if err != nil {
				fmt.Println(err)
			}
			if b {
				ca = strings.Trim(ca, "/")
				p.CtrlAction = strings.Split(ca, "/")
				for i, v := range p.CtrlAction {
					p.CtrlAction[i] = strings.Title(v)
				}
				if 0 < len(p.CtrlAction) {
					res := p.Handler(w, r)
					if res {
						return
					}
				}
			}
		}
	}
	p.CtrlAction = p.Segments

	res := p.Handler(w, r)
	if res {
		return
	}
	http.NotFound(w, r)
	return
}
