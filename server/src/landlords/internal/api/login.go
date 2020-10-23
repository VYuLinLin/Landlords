package api

import (
	"fmt"
	"landlords/internal/common"
	"landlords/internal/mysql"
	"math/rand"
	"time"
)

// Msgs 出参
type Msgs map[string]interface{}

// Register 注册接口
func Register(p interface{}) (d *common.User, err error) {
	m := p.(map[string]interface{})
	account := m["account"].(string)
	password := m["password"].(string)
	user, err := mysql.QueryOneUser(account)

	fmt.Println("QueryOneUser", *user)
	if err == nil {
		err = fmt.Errorf("账户已注册")
		return d, err
	}
	// 注册
	d = NewUser(account)
	fmt.Println("d", d)
	mysql.InsertUser(d.UserID, account, password)
	return d, nil
}

// Login 登录接口
func Login(p interface{}) (d *common.User, err error) {
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
		mysql.InsertUser(d.UserID, "guest", "")
	}
	user, err := mysql.QueryOneUser(account)
	fmt.Println(*user)
	if err != nil {
		d = NewUser(account)
		mysql.InsertUser(d.UserID, account, password)
	} else {
		if password == user.PASSWORD {
			d = &common.User{user.ID, user.NAME, user.STATUS}
		} else {
			err = fmt.Errorf("密码错误")
		}
	}
	return d, err
}

// NewUser 生成一个新玩家
func NewUser(userName string) *common.User {
	time.Sleep(200) // 添加睡眠，避免同步调用时生成相同的id
	rand.Seed(time.Now().UnixNano())
	var id = rand.Intn(100000)
	return &common.User{id, userName, 1}
}
