package socketio

import (
	"fmt"
	"landlords/internal/api"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

// type Msg struct {
// 	UserId    string   `json:"userId"`
// 	Text      string   `json:"text"`
// 	State     string   `json:"state"`
// 	Namespace string   `json:"namespace"`
// 	Rooms     []string `json:"rooms"`
// }

// Msgs 出入参
type Msgs map[string]interface{}

type callBackData struct {
	status   int
	msgType  string
	msgIndex float64
	data     interface{}
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

	server.OnEvent("/", "notify", func(s socketio.Conn, msg Msgs) {
		var (
			status   = 0
			msgIndex = msg["msgIndex"]
			msgType  = msg["msgType"]
			data     = make(Msgs)
		)

		fmt.Println("notify:", msg)

		if len(msgType.(string)) == 0 {
			status = 4
			data["data"] = "请求错误"
		} else {
			if msgType == "login" {
				data["data"] = api.Login(msg["data"])
			}
		}
		data["status"] = status
		data["msgType"] = msgType
		data["msgIndex"] = msgIndex

		fmt.Println(data)
		s.Emit("notify", data)
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

	http.Handle("/socket.io/", corsMiddleware(server))
	// http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8528...")
	log.Fatal(http.ListenAndServe(":8528", nil))
}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

		w.Header().Set("Content-Type", "application/json")
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

		// it can fix 403
		r.Header.Del("Origin")
		next.ServeHTTP(w, r)
	})
}
