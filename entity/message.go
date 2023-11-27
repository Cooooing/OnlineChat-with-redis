package entity

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Message 消息
type Message struct {
	SendTime string `json:"send_time"`
	Username string `json:"username"`
	Msg      string `json:"msg"`

	Channel string `json:"channel"`
}

// NewSendMessage 创建发送消息
func NewSendMessage(ctx context.Context, channel string, msg string) Message {
	sendTime := time.Now().Format("2006-01-02 15:04:05")
	username := fmt.Sprintf("%s", ctx.Value("username"))
	msg = strings.TrimSpace(msg)
	return Message{SendTime: sendTime, Username: username, Msg: msg, Channel: channel}
}

// NewReceiveMessage 创建接收消息
func NewReceiveMessage(msg string) Message {
	message := Message{}
	err := json.Unmarshal([]byte(msg), &message)
	if err != nil {
		return Message{SendTime: TimeHandle(time.Now()), Msg: msg, Username: "unknown"}
	}
	sendTime, err := time.ParseInLocation("2006-01-02 15:04:05", message.SendTime, time.Local)
	if err != nil {
		panic(err)
	}
	message.SendTime = TimeHandle(sendTime)
	return message
}
