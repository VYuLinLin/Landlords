package socketio

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-socket.io"
	"io"
	"landlords/internal/api"
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

func parseResponse(r *http.Request) (interface{}, error) {
	var res RequestMsg
	b, err := io.ReadAll(r.Body)
	if err == nil {
		json.Unmarshal(b, &res)
	}
	return res, err
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := &responseMsg{0, "success", NullStruct{}}
	fmt.Println("URL", r.URL)
	fmt.Println("URL.Path", r.URL.Path)
	fmt.Println("Host", r.Host)

	if r.Method != "POST" || r.ParseForm() != nil {
		data.Code = 1
		data.Msg = "请求类型错误，请使用POST"
	} else {
		response, err1 := parseResponse(r)
		if err1 != nil {
			data.Code = 1
			data.Msg = err1.Error()
			log.Println(err1)
		}
		fmt.Println("request parameter:", response)
		var res any
		var err2 error
		var loginApi = &api.LoginApi{}
		var roomApi = &api.RoomApi{}
		switch r.URL.String() {
		case "/register":
			res, err2 = loginApi.RegisterHandler(response)
		case "/login":
			res, err2 = loginApi.LoginHandler(response)
		case "/logout":
			_, err2 = loginApi.LogoutHandler(response)
		case "/joinRoom":
			res, err2 = roomApi.JoinRoom(response)
		case "/getTable":
			res, err2 = roomApi.GetTable(response)
		}
		if err2 != nil {
			data.Code = 1
			data.Msg = err2.Error()
			log.Println(err2)
		} else {
			data.Data = res
		}
	}

	log.Println(data)
	strParams, _ := json.Marshal(data)
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

	// server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
	// 	s.SetContext(msg)
	// 	return "recv " + msg
	// })

	// server.OnEvent("/", "bye", func(s socketio.Conn) string {
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

	http.HandleFunc("/login", handler)
	http.HandleFunc("/register", handler)
	http.HandleFunc("/logout", handler)
	http.HandleFunc("/joinRoom", handler)
	http.HandleFunc("/getTable", handler)

	http.Handle("/socket.io/", corsMiddleware(server))
	//http.Handle("/", http.FileServer(http.Dir("./asset")))

	log.Println("Serving at localhost:8528...")
	log.Fatal(http.ListenAndServe(":8528", nil))
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
