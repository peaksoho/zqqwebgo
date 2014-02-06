package zqqwebgo

import (
	"fmt"
	"mime"
	"net/http"
	"strings"
	"time"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func (ctx *Context) WriteString(content string) {
	ctx.W.Write([]byte(content))
}

func (ctx *Context) Abort(status int, body string) {
	ctx.W.WriteHeader(status)
	ctx.W.Write([]byte(body))
}

func (ctx *Context) Redirect(status int, url_ string) {
	ctx.W.Header().Set("Location", url_)
	ctx.W.WriteHeader(status)
}

func (ctx *Context) NotModified() {
	ctx.W.WriteHeader(304)
}

func (ctx *Context) NotFound(message string) {
	ctx.W.WriteHeader(404)
	ctx.W.Write([]byte(message))
}

//Sets the content type by extension, as defined in the mime package.
//For example, ctx.ContentType("json") sets the content-type to "application/json"
func (ctx *Context) ContentType(ext string) {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	ctype := mime.TypeByExtension(ext)
	if ctype != "" {
		ctx.W.Header().Set("Content-Type", ctype)
	}
}

func (ctx *Context) SetHeader(hdr string, val string, unique bool) {
	if unique {
		ctx.W.Header().Set(hdr, val)
	} else {
		ctx.W.Header().Add(hdr, val)
	}
}

//Sets a cookie -- duration is the amount of time in seconds. 0 = forever
func (ctx *Context) SetCookie(name string, value string, age int64) {
	var utctime time.Time
	if age == 0 {
		// 2^31 - 1 seconds (roughly 2038)
		utctime = time.Unix(2147483647, 0)
	} else {
		utctime = time.Unix(time.Now().Unix()+age, 0)
	}
	cookie := fmt.Sprintf("%s=%s; Expires=%s; Path=/", name, value, webTime(utctime))
	ctx.SetHeader("Set-Cookie", cookie, true)
}
