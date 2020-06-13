package redisub

import (
	"strconv"
	"strings"

	"github.com/infernalfire72/acache/log"
	"github.com/infernalfire72/acache/leaderboards"
	"gopkg.in/redis.v5"
)

func ban(client *redis.Client) {
	ps, err := client.Subscribe("peppy:ban")
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Subscribed to peppy:unban")
	for {
		msg, err := ps.ReceiveMessage()
		if err != nil {
			log.Error(err)
			return
		}

		i, err := strconv.Atoi(msg.Payload)
		if err != nil {
			continue
		}
		
		leaderboards.Cache.RemoveUser(i)
	}
}

func unban(client *redis.Client) {
	ps, err := client.Subscribe("peppy:unban")
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Subscribed to peppy:unban")
	for {
		msg, err := ps.ReceiveMessage()
		if err != nil {
			log.Error(err)
			return
		}

		i, err := strconv.Atoi(msg.Payload)
		if err != nil {
			continue
		}
		
		leaderboards.Cache.AddUser(i)
	}
}

func wipe(client *redis.Client) {
	ps, err := client.Subscribe("peppy:wipe")
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Subscribed to peppy:wipe")
	for {
		msg, err := ps.ReceiveMessage()
		if err != nil {
			log.Error(err)
			return
		}

		ss := strings.SplitN(msg.Payload, ",", 2)
		if len(ss) != 2 {
			continue
		}

		i, err := strconv.Atoi(ss[0])
		if err != nil {
			log.Error(err)
			continue
		}

		rx, err := strconv.ParseBool(ss[1])
		if err != nil {
			log.Error(err)
			continue
		}
		
		leaderboards.Cache.RemoveUserWithIdentifier(i, rx)
	}
}

var Client *redis.Client

func Subscribe() {
	Client = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    _, err := Client.Ping().Result()
	if err != nil {
		log.Error(err)
		return
	}

	go ban(Client)
	go unban(Client)
	go wipe(Client)
}