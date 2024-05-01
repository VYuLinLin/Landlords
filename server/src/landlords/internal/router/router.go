package router

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/googollee/go-socket.io"
	"landlords/internal/api"
	"landlords/internal/conf"
	"log"
	"net/http"
)

// responseMsg 接口出参统一格式
type responseMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
type NullStruct struct{}
type RequestMsg = map[string]interface{}
type HandleFunc func(http.ResponseWriter, *http.Request)

func parseResponse(r *http.Request) (res interface{}, err error) {
	if r.Method == "GET" {
		res = r.URL.Query()
	}
	if r.Method == "POST" {
		// 第一种解析方式
		//b, err := io.ReadAll(r.Body)
		//if err == nil {
		//	json.Unmarshal(b, &res)
		//}
		// 第二种解析方式
		err = json.NewDecoder(r.Body).Decode(&res)
	}

	return res, err
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := &responseMsg{0, "success", NullStruct{}}
	fmt.Println("URL", r.URL)

	response, err1 := parseResponse(r)
	if err1 != nil {
		data.Code = 1
		data.Msg = err1.Error()
		log.Println(err1)
	}
	logs.Info("http 请求数据体[%s]: %s", r.URL.Path, response)
	var res any
	var err2 error
	var loginApi = &api.LoginApi{}
	var roomApi = &api.RoomApi{}
	switch r.URL.Path {
	case "/register":
		res, err2 = loginApi.RegisterHandler(response)
	case "/login":
		res, err2 = loginApi.LoginHandler(response)
	case "/logout":
		_, err2 = loginApi.LogoutHandler(response)
	case "/joinRoom":
		res, err2 = roomApi.JoinRoom(response)
	case "/exitRoom":
		err2 = roomApi.ExitRoom(response)
	case "/getTable":
		res, err2 = roomApi.GetTable(response)
	case "/ws/":
		err2 = api.StartServeWs(w, r)
		logs.Error("/ws", err2)
		return
	}

	if err2 != nil {
		data.Code = 1
		data.Msg = err2.Error()
	} else {
		data.Data = res
	}
	strParams, err2 := json.Marshal(data)
	if err2 != nil {
		data.Code = 1
		data.Msg = err2.Error()
		data.Data = nil
		strParams, _ = json.Marshal(data)
	}
	logs.Info("http 返回数据体[%s]: %s", r.URL.Path, strParams)
	w.Write(strParams)
}
func init() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("connect server success.")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "notify", func(s socketio.Conn, msg responseMsg) {
		//var (
		//	status   = 1
		//	msgIndex = msg["msgIndex"]
		//	msgType  = msg["msgType"]
		//	data     = make(Msgs)
		//)
		//
		//fmt.Println("notify:", msg)
		//log.Println(
		//	"notify:",
		//	msg,
		//)
		//if len(msgType.(string)) == 0 {
		//	fmt.Println(msgType)
		//	status = 0
		//	data["data"] = "请求错误"
		//} else {
		//	var p = msg["data"]
		//	var err error
		//	switch msgType {
		//	case "register":
		//		data["data"], err = api.Register(p)
		//	case "login":
		//		data["data"], err = api.LoginHandler(p)
		//	case "JoinRoom":
		//		data["data"], err = api.JoinRoom(p)
		//	case "startgame":
		//		data["data"], err = api.JoinRoom(p)
		//	}
		//	if err != nil {
		//		status = 0
		//		data["data"] = err.Error()
		//	}
		//}
		//data["status"] = status
		//data["msgType"] = msgType
		//data["msgIndex"] = msgIndex
		//
		//fmt.Println("emit notify: ", data)
		//s.Emit("notify", data)
	})

	// server.OnEvent("/chat", "msg", func(s router.Conn, msg string) string {
	// 	s.SetContext(msg)
	// 	return "recv " + msg
	// })

	// server.OnEvent("/", "bye", func(s router.Conn) string {
	// 	last := s.Context().(string)
	// 	s.Emit("bye", last)
	// 	s.Close()
	// 	return last
	// })

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed: ", reason)
	})

	go server.Serve()
	defer server.Close()

	http.HandleFunc("/login", logPanics(handler))
	http.HandleFunc("/register", logPanics(handler))
	http.HandleFunc("/logout", logPanics(handler))
	http.HandleFunc("/joinRoom", logPanics(handler))
	http.HandleFunc("/exitRoom", logPanics(handler))
	http.HandleFunc("/getTable", logPanics(handler))
	http.HandleFunc("/ws/", logPanics(handler))

	//http.Handle("/socket.io/", corsMiddleware(server))
	//http.Handle("/", http.FileServer(http.Dir("./asset")))
	// 设置静态目录
	static := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", static))

	addr := flag.String("addr", fmt.Sprintf(":%d", conf.GameConf.HttpPort), "http service address")

	log.Printf("Serving at localhost:%s...\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

		w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

		// it can fix 403
		r.Header.Del("Origin")
		next.ServeHTTP(w, r)
	})
}
func logPanics(f HandleFunc) HandleFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				log.Printf("[%v] caught panic: %v", req.RemoteAddr, x)
				// 给页面一个错误信息, 如下示例返回一个500
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		// 跨域请求时，允许头部携带cookie，设置后Access-Control-Allow-Origin值不能是“*”
		//w.Header().Set("Access-Control-Allow-Credentials", "true")
		//w.Header().Set("Access-Control-Allow-Origin", "http://172.21.165.80:7456")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		if req.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
		} else {
			f(w, req)
		}
	}
}
