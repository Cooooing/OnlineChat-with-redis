package entity

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strings"
)

var Client *redis.Client
var SelfMessageReceiveNumChan chan int64
var ReceiveMessageChan chan Message

func init() {
	SelfMessageReceiveNumChan = make(chan int64, 20)
	ReceiveMessageChan = make(chan Message, 20)
}

// InitRedisClient 初始化 redis 连接
func InitRedisClient() {
	for {
		fmt.Print("please enter redis address:")
		reader := bufio.NewReader(os.Stdin)
		address, _ := reader.ReadString('\n')
		address = strings.TrimSpace(address)
		fmt.Print("please enter redis password:")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)
		if address == "" {
			address = "localhost:6379"
		}

		Client = redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password, // 没有密码，默认值
			DB:       0,        // 默认DB 0
			PoolSize: 10,       // 最大套接字连接数
		}) // Background返回一个非空的Context。它永远不会被取消，没有值，也没有截止日期。
		// 它通常由main函数、初始化和测试使用，并作为传入请求的顶级上下文
		ctx := context.Background()

		res, err := Client.Ping(ctx).Result()
		if err == nil {
			log.Println("redis connect succeeded : ", res)
			fmt.Println("redis connect succeeded.")
			break
		}
		log.Println("redis connect failed", err)
		fmt.Println("redis connect failed,please try again.")
	}
}

// CloseRedisClient 关闭 redis 连接
func CloseRedisClient() {
	err := Client.Close()
	if err != nil {
		log.Println("redis close failed", err)
		panic(err)
	}
	log.Println("redis closed")
}

// Subscriber 订阅者 处理订阅接收到的消息
func Subscriber(channel string) {
	ctx := context.Background()
	pubsub := Client.Subscribe(ctx, channel)
	defer func(pubsub *redis.PubSub) {
		err := pubsub.Close()
		if err != nil {
			panic(err)
		}
	}(pubsub)

	// 处理订阅接收到的消息
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			return
		}
		message := NewReceiveMessage(msg.Payload)
		message.Channel = channel
		ReceiveMessageChan <- message
		log.Printf("receive message:{%s} from channel:{%s},username:{%s}\n", message.Msg, message.Channel, message.Username)
	}
}

// Publisher 生产者 发布消息到频道
func Publisher(message chan Message) {
	ctx := context.Background()
	for {
		// 发布消息到频道
		msg := <-message
		marshal, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}
		res := Client.Publish(ctx, msg.Channel, string(marshal))
		SelfMessageReceiveNumChan <- res.Val()
		log.Printf("send message:{%s} to channel:{%s},receive num:{%d}\n", msg.Msg, msg.Channel, res.Val())
	}
}
