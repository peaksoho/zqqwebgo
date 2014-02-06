package zqqwebgo

import (
	"encoding/json"
	"encoding/xml"
	//"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"zqqwebgo/session"
	"zqqwebgo/zqqjsongo"
)

type Controller struct {
	Segment    []string
	Params     url.Values
	ParamsGet  url.Values
	ParamsPost url.Values
	ChildName  string
	TplNames   string
	Layout     string
	OutString  string
	Ctx        *Context
	Data       map[interface{}]interface{}
	CruSession session.SessionStore
}

func (c *Controller) Init(w http.ResponseWriter, r *http.Request) {
	c.Data = make(map[interface{}]interface{})
	c.Layout = ""
	c.TplNames = ""
	c.OutString = ""
	c.Ctx = &Context{W: w, R: r}
}

func (c *Controller) SetSegment(segment []string) {
	c.Segment = segment
}

func (c *Controller) SetParams(w http.ResponseWriter, r *http.Request) {
	var err error
	c.Params = r.Form
	c.ParamsGet, err = url.ParseQuery(r.URL.RawQuery) //GET method
	if err != nil {
		c.OutErr(err)
		return
	}
	c.ParamsPost = r.PostForm
	/*c.ParamsPost = make(url.Values, 1) //POST method
	for k, v := range c.Params {
		_, ok := c.ParamsGet[k]
		if !ok {
			c.ParamsPost[k] = v
		}
	}*/
}

func (c *Controller) ParseTpl() { //解析模板文件
	if c.TplNames != "" {
		t, err := template.ParseFiles(ViewsPath + c.TplNames)
		c.OutErr(err)
		t.Execute(c.Ctx.W, c.Data)
	}
	return
}

func (c *Controller) OutPutString() { //直接返回响应
	if c.OutString != "" && c.TplNames == "" {
		io.WriteString(c.Ctx.W, c.OutString)
	}
	return
}

func (c *Controller) OutErr(err error) {
	if err != nil {
		http.Error(c.Ctx.W, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) Redirect(args ...interface{}) { //跳转
	var url string = "/"
	var code int = http.StatusFound
	if 0 < len(args) {
		for _, element := range args {
			switch value := element.(type) {
			case int:
				code = value
			case string:
				url = value
			}
		}
	}
	c.Ctx.Redirect(code, url)
}

func (c *Controller) Abort(code string) {
	panic(code)
}

func (c *Controller) ServeJson() {
	content, err := json.MarshalIndent(c.Data["json"], "", "  ")
	if err != nil {
		c.OutErr(err)
		return
	}
	c.Ctx.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
	c.Ctx.W.Header().Set("Content-Type", "application/json")
	c.Ctx.W.Write(content)
}

func (c *Controller) ServeXml() {
	content, err := xml.Marshal(c.Data["xml"])
	if err != nil {
		c.OutErr(err)
		return
	}
	c.Ctx.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
	c.Ctx.W.Header().Set("Content-Type", "application/xml")
	c.Ctx.W.Write(content)
}

func (c *Controller) Input() url.Values {
	ct := c.Ctx.R.Header.Get("Content-Type")
	if strings.Contains(ct, "multipart/form-data") {
		c.Ctx.R.ParseMultipartForm(MaxMemory) //64MB
	} else {
		c.Ctx.R.ParseForm()
	}
	return c.Ctx.R.Form
}

func (c *Controller) GetString(key string) string {
	return c.Input().Get(key)
}

func (c *Controller) GetInt(key string) (int64, error) {
	return strconv.ParseInt(c.Input().Get(key), 10, 64)
}

func (c *Controller) GetBool(key string) (bool, error) {
	return strconv.ParseBool(c.Input().Get(key))
}

func (c *Controller) GetFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.Ctx.R.FormFile(key)
}

func (c *Controller) SaveToFile(fromfile, tofile string) error {
	file, _, err := c.Ctx.R.FormFile(fromfile)
	if err != nil {
		return err
	}
	defer file.Close()
	f, err := os.OpenFile(tofile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, file)
	return nil
}

func (c *Controller) Destructor() {
	if c.CruSession != nil {
		c.CruSession.SessionRelease()
	}
}

func (c *Controller) StartSession() session.SessionStore {
	//if c.CruSession == nil {
	c.CruSession = GlobalSessions.SessionStart(c.Ctx.W, c.Ctx.R)
	//}
	return c.CruSession
}

func (c *Controller) SetSession(name interface{}, value interface{}) {
	if c.CruSession == nil {
		c.StartSession()
	}
	c.CruSession.Set(name, value)
}

func (c *Controller) GetSession(name interface{}) interface{} {
	if c.CruSession == nil {
		c.StartSession()
	}
	return c.CruSession.Get(name)
}

func (c *Controller) DelSession(name interface{}) {
	if c.CruSession == nil {
		c.StartSession()
	}
	c.CruSession.Delete(name)
}

func (c *Controller) LoadConfig(name string) *zqqjsongo.Json {
	filename := "./config/" + name + ".json"
	js, err := zqqjsongo.JsonFileDecode(filename)
	if err != nil {
		c.OutErr(err)
	}
	return js
}
