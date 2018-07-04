package redis_pubsub

import (
	"github.com/abonec/web_pusher"
	"github.com/go-redis/redis"
	"github.com/abonec/web_pusher/logger"
	"strings"
)

type Backend struct {
	server      *web_pusher.Server
	redisClient *redis.Client
	logger      logger.Logger
}

func NewBackend(server *web_pusher.Server, logger logger.Logger) *Backend {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &Backend{server, client, logger}
}

func (back *Backend) log(format string, v ...interface{}) {
	back.logger.Printf("[REDI] "+format, v...)
}
func (back *Backend) Start(channels ...string) error {
	_, err := back.redisClient.Ping().Result()
	if err != nil {
		return err
	}
	pubSub := back.redisClient.PSubscribe(channels...)
	go func() {
		defer pubSub.Close()
		for {
			pubSubMessage, err := pubSub.ReceiveMessage()
			if err != nil {
				back.log("Error while reading pubsub message")
				continue
			}
			back.log("Got a new message: %s", pubSubMessage.Payload)

			channel := strings.Split(pubSubMessage.Channel, ":")
			if len(channel) != 2 {
				back.log("channel pattern has wrong format. It should be \"channelName:userId\". Got \"%s\"", pubSubMessage.Channel)
				continue
			}
			back.server.SendToUser(channel[1], []byte(pubSubMessage.Payload))
		}
	}()
	return nil
}
