package clients

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/nothollyhigh/kiss/log"
	"go_gate/config"
	"sync"
)

var once sync.Once
var redisClient *redis.Client

func RedisHandle() *redis.Client {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     config.GlobalXmlConfig.Options.RedisInfo.Address,
			Password: config.GlobalXmlConfig.Options.RedisInfo.Password, // no password set
			DB:       config.GlobalXmlConfig.Options.RedisInfo.DBNum,    // use default DB
		})
		pong, err := redisClient.Ping().Result()
		log.Debug(":%v :%v", pong, err)
		if nil != err {
			//self.DB.essence = nil //允许不使用redis
			log.Fatal("cann't to connect redis:%v  error:%v", config.GlobalXmlConfig.Options.RedisInfo.Address, err)
		}
	})
	return redisClient
}

// GetAccountKey 通过账号获取服务器地址
func GetAccountKey(acc string) string {
	return fmt.Sprintf("acc_%v", acc)
}

// GetAddressKey 通过客户端地址获取账号
func GetAddressKey(addr string) string {
	return fmt.Sprintf("addr_%v", addr)
}
