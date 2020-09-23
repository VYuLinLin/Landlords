package api

import (
	"math/rand"
	"time"
)

// Msgs 出参
type Msgs map[string]interface{}

// Login 登录接口
func Login(p interface{}) (d Msgs) {
	d = make(Msgs, 2)
	m := p.(map[string]interface{})
	isGuest := m["isGuest"].(float64)

	if isGuest == 1 {
		rand.Seed(time.Now().Unix())
		d["userId"] = rand.Intn(100000)
		d["userName"] = "guest"
	}

	return d
}
