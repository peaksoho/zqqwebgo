package zqqwebgo

/*
//#include <unistd.h>
//#include <stdio.h>
*/
//import "C"
import (
	"flag"
	"fmt"
	"github.com/peaksoho/zqqwebgo/session"
	"github.com/peaksoho/zqqwebgo/zqqjsongo"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
)

const VERSION = "0.0.1"

var (
	AppPath   string
	HttpAddr  string
	HttpPort  int
	StaticDir map[string]string
	ViewsPath string
	MaxMemory int64
	//related to session
	GlobalSessions       *session.Manager //GlobalSessions
	SessionOn            bool             // wheather auto start session,default is false
	SessionProvider      string           // default session provider  memory mysql redis
	SessionName          string           // sessionName cookie's name
	SessionGCMaxLifetime int64            // session's gc maxlifetime
	SessionSavePath      string           // session savepath if use mysql/redis/file this set to the connectinfo
)

func init() {
	AppPath, _ = os.Getwd()
	HttpAddr = ""
	HttpPort = 9090
	StaticDir = make(map[string]string)
	StaticDir["/static"] = "static"
	ViewsPath = "views/"
	MaxMemory = 1 << 26 //64MB
	SessionOn = true
	SessionProvider = "memory"
	SessionName = "zqqwebgosessionID"
	SessionGCMaxLifetime = 3600
	SessionSavePath = ""
}

func AddController(str string, c ControllerInterface) {
	ControllerRegistor[str] = c
}

func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func Run() {
	argPort := flag.String("p", "9090", "the port of web service!")
	argCfgPath := flag.String("c", "./config/server.json", "the config file path of web server! ")
	//dae := flag.Bool("d", false, "Whether or not to launch in the background (like a daemon)")
	flag.Parse()
	//if *dae {
	//	C.daemon(1, 1)
	//}
	//更新config值
	cfgFile := ""
	if *argCfgPath != "" {
		cfgFile = *argCfgPath
	} else {
		cfgFile = "./config/server.json"
	}

	isCfgPort := false
	if FileExist(cfgFile) {
		js, err := zqqjsongo.JsonFileDecode(cfgFile)
		if err != nil {
			log.Fatal(err)
		}
		if http_addr, _ := js.Get("HttpAddr").String(); HttpAddr != http_addr {
			HttpAddr = http_addr
		}
		if http_port, _ := js.Get("HttpPort").Int(); http_port > 0 && HttpPort != http_port {
			HttpPort = http_port
			isCfgPort = true
		}
		if static_dir, _ := js.Get("StaticDir").Map(); len(static_dir) > 0 {
			for k, v := range static_dir {
				switch vv := v.(type) {
				case string:
					StaticDir[k] = vv
				default:
				}
			}
		}
		if views_path, _ := js.Get("ViewsPath").String(); ViewsPath != views_path {
			ViewsPath = views_path
		}
		if max_memory, _ := js.Get("MaxMemory").Int64(); max_memory > 0 && MaxMemory != max_memory {
			MaxMemory = max_memory
		}
		if session_on, _ := js.Get("SessionOn").Bool(); SessionOn != session_on {
			SessionOn = session_on
		}
		if session_provider, _ := js.Get("SessionProvider").String(); SessionProvider != session_provider {
			SessionProvider = session_provider
		}
		if session_name, _ := js.Get("SessionName").String(); SessionName != session_name {
			SessionName = session_name
		}
		if session_gc_maxlifetime, _ := js.Get("SessionGCMaxLifetime").Int64(); SessionGCMaxLifetime != session_gc_maxlifetime {
			SessionGCMaxLifetime = session_gc_maxlifetime
		}
		if session_savepath, _ := js.Get("SessionSavePath").String(); SessionSavePath != session_savepath {
			SessionSavePath = session_savepath
		}
	}

	if isCfgPort == false && *argPort != "" {
		var err error
		HttpPort, err = strconv.Atoi(*argPort)
		if err != nil {
			log.Fatal("The port [" + *argPort + "] is wrong!")
		}
	}

	fmt.Println("The web server is running!\nThe Port is", HttpPort)
	if SessionOn {
		GlobalSessions, _ = session.NewManager(SessionProvider, SessionName, SessionGCMaxLifetime, SessionSavePath)
		go GlobalSessions.GC()
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	mux := &Router{}
	addr := fmt.Sprintf("%s:%d", HttpAddr, HttpPort)

	http.ListenAndServe(addr, mux)
}
