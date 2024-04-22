package api

import (
	"errors"
	"fmt"
	"landlords/internal/common"
	"landlords/internal/mysql"
	"math/rand"
	"time"
)

type LoginApi struct{}

// Msgs 出参
type Msgs map[string]interface{}

// RegisterHandler 注册接口
func (l *LoginApi) RegisterHandler(p interface{}) (d *mysql.User, err error) {
	m := p.(map[string]interface{})
	if m["account"] == nil || m["password"] == nil {
		return nil, errors.New("用户名或密码不能为空")
	}
	User := &mysql.User{
		NAME:     m["account"].(string),
		PASSWORD: m["password"].(string),
	}
	user, err := mysql.QueryUser(User)
	fmt.Println("QueryUser", *user, err)
	if User.NAME == user.NAME || err == nil {
		err = fmt.Errorf("账户已注册")
		return user, err
	}
	// 生成用户ID并注册
	mysql.InsertUser(User.NAME, User.PASSWORD)

	user, err = mysql.QueryUser(User)
	fmt.Println("QueryTwoUser", *user, err)
	return user, err
}

// LoginHandler 登录接口
func (l *LoginApi) LoginHandler(p interface{}) (d *mysql.User, err error) {
	m := p.(map[string]interface{})
	if m["account"] == nil || m["password"] == nil {
		return nil, errors.New("用户名或密码不能为空")
	}
	User := &mysql.User{
		NAME:     m["account"].(string),
		PASSWORD: m["password"].(string),
	}
	user, err := mysql.QueryUser(User)
	fmt.Println("LoginApi:", *user, *User)
	if err == nil && User.PASSWORD != user.PASSWORD {
		return nil, errors.New("密码错误")
	}
	return user, err
}

// LogoutHandler 退出
func (l *LoginApi) LogoutHandler(p interface{}) (d *mysql.User, err error) {
	m := p.(map[string]interface{})
	if m["id"] == nil || m["name"] == nil {
		return d, errors.New("ID或用户名错误")
	}
	User := &mysql.User{
		ID:   int(m["id"].(float64)),
		NAME: m["name"].(string),
	}
	fmt.Println("LogoutHandler:", *User)
	var user *mysql.User
	user, err = mysql.QueryUser(User)
	fmt.Println(*user)
	if err == nil && User.ID != user.ID {
		err = fmt.Errorf("用户ID错误")
	}
	return user, err
}

// NewUser 生成一个新玩家
func NewUser(userName string) *common.User {
	time.Sleep(10) // 添加睡眠，避免同步调用时生成相同的id
	rand.Seed(time.Now().UnixNano())
	var id = rand.Intn(100000)
	return &common.User{id, userName, 1}
}
