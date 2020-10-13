package api

import (
	"fmt"
	"landlords/internal/mysql"
	"math/rand"
	"time"
)

// Msgs 出参
type Msgs map[string]interface{}

// Register 注册接口
func Register(p interface{}) (d Msgs, err error) {
	m := p.(map[string]interface{})
	account := m["account"].(string)
	password := m["password"].(string)
	user, err := mysql.QueryOne(account)

	fmt.Println(*user)
	if err == nil {
		err = fmt.Errorf("账户已注册")
		return d, err
	}
	// 注册
	d = NewUser(account)
	mysql.InsertUser(d["userId"].(int), account, password)
	return d, nil
}

// Login 登录接口
func Login(p interface{}) (d Msgs, err error) {
	m := p.(map[string]interface{})
	isGuest := m["isGuest"]
	_, is := isGuest.(float64)
	if !is {
		fmt.Println("isGuest Assert Type error: ", is)
	}
	account := m["account"].(string)
	password := m["password"].(string)

	if isGuest == 1 {
		d = NewUser("guest")
		mysql.InsertUser(d["userId"].(int), "guest", "")
	}
	user, err := mysql.QueryOne(account)
	fmt.Println(*user)
	if err != nil {
		d = NewUser(account)
		mysql.InsertUser(d["userId"].(int), account, password)
	} else {
		d = make(Msgs, 3)
		if password == user.PASSWORD {
			d["userId"] = user.ID
			d["userName"] = user.NAME
			d["status"] = user.STATUS
		} else {
			err = fmt.Errorf("密码错误")
		}
	}
	return d, err
}

// NewUser 生成一个新玩家
func NewUser(userName string) (d Msgs) {
	rand.Seed(time.Now().Unix())
	var id = rand.Intn(100000)
	d = make(Msgs, 3)
	d["userId"] = id
	d["userName"] = userName
	d["status"] = 1
	return d
}
