package command

import (
	"log-record/server"
	"log-record/store"
)

type Type int

const (
	// AUTH 授权
	AUTH Type = 0
	// DATA 数据接收
	DATA Type = 1
	// PUSH 推送数据到客户端
	PUSH Type = 2
)

type Command struct {
	// Type 类型
	Type Type
	// Data 数据
	Data []byte
	// Header 标头
	Header map[string]string
}

func DoCommand(cmd *Command, conn *server.Conn) error {
	switch cmd.Type {
	case AUTH:
		if _, ok := conn.Attach["auth"]; ok {
			//授权过了
			break
		}
		username, ok := cmd.Header["username"]
		if !ok {
			_, _ = conn.Conn.Write(pushData("用户名不能为空"))
			_ = conn.Conn.Close()
		}
		password, ok := cmd.Header["password"]
		if !ok {
			_, _ = conn.Conn.Write(pushData("密码不能为空"))
			_ = conn.Conn.Close()
		}
		user := store.User{}
		if !user.Auth(username, password) {
			_, _ = conn.Conn.Write(pushData("用户名密码输入有误"))
			_ = conn.Conn.Close()
		}
		conn.Attach["auth"] = true
		break
	case DATA:
		if _, ok := conn.Attach["auth"]; ok {
			store.DataStore(
				cmd.Data,
				cmd.Header,
				conn,
			)
		} else {
			_, _ = conn.Conn.Write(pushData("请先登录"))
			_ = conn.Conn.Close()
		}
		break
	case PUSH:
		break
	}
	return nil
}

func pushData(msg string) []byte {
	return []byte(msg)
}
