package api

import (
	"errors"
	"fmt"
	"landlords/internal/db"
	"landlords/internal/game/room"
	"strconv"
)

type LoginApi struct{}

// Msgs 出参
type Msgs map[string]interface{}

// RegisterHandler 注册接口
func (l *LoginApi) RegisterHandler(p interface{}) (d *db.User, err error) {
	m := p.(map[string]interface{})
	if m["account"] == nil || m["password"] == nil {
		return nil, errors.New("用户名或密码不能为空")
	}
	User := &db.User{
		NAME:     m["account"].(string),
		PASSWORD: m["password"].(string),
	}
	user, err := db.QueryUser(User)
	fmt.Println("QueryUser", *user, err)
	if User.NAME == user.NAME || err == nil {
		err = fmt.Errorf("账户已注册")
		return user, err
	}
	// 生成用户ID并注册
	db.InsertUser(User.NAME, User.PASSWORD)

	user, err = db.QueryUser(User)
	fmt.Println("QueryTwoUser", *user, err)
	return user, err
}

// LoginHandler 登录接口
func (l *LoginApi) LoginHandler(p interface{}) (d *db.User, err error) {
	m := p.(map[string]interface{})
	if m["account"] == nil || m["password"] == nil {
		return nil, errors.New("用户名或密码不能为空")
	}
	User := &db.User{
		NAME:     m["account"].(string),
		PASSWORD: m["password"].(string),
	}
	user, err := db.QueryUser(User)
	fmt.Println("LoginApi:", *user, *User)
	if err == nil && User.PASSWORD != user.PASSWORD {
		return nil, errors.New("密码错误")
	}
	if user.TABLEID > 0 && user.ROOMID > 0 {
		_, err = room.JoinRoom(user, strconv.Itoa(user.ROOMID))
	}
	return user, err
}

// LogoutHandler 退出
func (l *LoginApi) LogoutHandler(p interface{}) (d *db.User, err error) {
	m := p.(map[string]interface{})
	if m["id"] == nil || m["name"] == nil {
		return d, errors.New("ID或用户名错误")
	}
	User := &db.User{
		ID:   int(m["id"].(float64)),
		NAME: m["name"].(string),
	}
	var user *db.User
	user, err = db.QueryUser(User)
	fmt.Println(*user)
	if err == nil && User.ID != user.ID {
		err = fmt.Errorf("用户ID错误")
	}
	return user, err
}
