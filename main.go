package main

import (
	bytes2 "bytes"
	"fmt"
	command2 "log-record/command"
	_ "log-record/http"
	"log-record/server"
	"log-record/utils"
	"strings"
	"time"
)

const (

	//HeaderSep 多条key value关联的连接符
	HeaderSep = "[<&>]"

	//HeaderLineSep key value之间的连接符
	HeaderLineSep = "="
)

func main() {
	server.Start(func(buffer *bytes2.Buffer, conn *server.Conn) error {
		attach := conn.Attach
		command, ok := attach["command"].(*command2.Command)
		if !ok {
			command = &command2.Command{
				Header: make(map[string]string, 0),
			}
		}
		resolveState, ok := attach["state"].(string)
		if !ok {
			resolveState = "Type"
		}
		for {
			//解析类型
			if strings.EqualFold(resolveState, "Type") {
				//读取type
				readType, err := buffer.ReadByte()
				if err != nil {
					break
				}
				command.Type = command2.Type(readType)
				//设置状态为读取标头长度
				attach["state"] = "HeaderLength"
				resolveState = "HeaderLength"
			}
			//解析标头长度
			if strings.EqualFold(resolveState, "HeaderLength") {
				//标识buffer中的长度满足int的长度
				if buffer.Len() >= 2 {
					//读取header的长度
					readShort, _ := utils.ReadShort(buffer)
					attach["headerLength"] = readShort
					attach["state"] = "HeaderBody"
					resolveState = "HeaderBody"
				} else {
					break
				}
			}
			//解析标头内容
			if strings.EqualFold(resolveState, "HeaderBody") {
				headerLength := attach["headerLength"].(int16)
				if buffer.Len() >= int(headerLength) {
					headerBody := make([]byte, headerLength)
					_, _ = buffer.Read(headerBody)
					headerArrays := strings.Split(string(headerBody), HeaderSep)
					for _, line := range headerArrays {
						headerKv := strings.Split(line, HeaderLineSep)
						if len(headerKv) == 2 {
							command.Header[strings.TrimSpace(headerKv[0])] = strings.TrimSpace(headerKv[1])
						} else {
							command.Header[strings.TrimSpace(headerKv[0])] = ""
						}
					}
					attach["state"] = "BodyLength"
					resolveState = "BodyLength"
				} else {
					break
				}
			}
			//读取body长度
			if strings.EqualFold(resolveState, "BodyLength") {
				//标识buffer中的长度满足int的长度
				if buffer.Len() >= 4 {
					//读取header的长度
					readInt, _ := utils.ReadInt(buffer)
					attach["bodyLength"] = readInt
					attach["state"] = "Body"
					resolveState = "Body"
				} else {
					break
				}
			}
			//读取body内容
			if strings.EqualFold(resolveState, "Body") {
				bodyLength := attach["bodyLength"].(int32)
				if buffer.Len() >= int(bodyLength) {
					body := make([]byte, bodyLength)
					_, _ = buffer.Read(body)
					command.Data = body
					//处理
					err := command2.DoCommand(command, conn)
					if err != nil {
						fmt.Println(err)
						return err
					}
					attach["command"] = &command2.Command{Header: make(map[string]string, 0)}
					attach["state"] = "Type"
					resolveState = "Type"
				}
				break
			}
		}
		return nil
	})

	time.Sleep(time.Hour * 24 * 30 * 1000)
}
