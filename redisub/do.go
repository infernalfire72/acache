package redisub

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/infernalfire72/acache/config"
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
	log.Info("Subscribed to peppy:ban")
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
		log.Info("[RESTRICT] Wiped Cached Scores for", i)
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
		log.Info("[UNRESTRICT] Added Scores to Cache for", i)
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
		log.Info("[WIPE] Wiped Cached Scores for", i, rx)
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
	go newScore(Client)
	go changeUsername(Client)
}

func newScore(client *redis.Client) {
	ps, err := client.Subscribe("api:score_submission")
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Subscribed to api:score_submission")

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

		table := "scores"
		if rx {
			table = "scores_relax"
		}

		s := &leaderboards.Score{}
		var (
			md5		string
			mode	byte
		)
		err = config.DB.QueryRow("SELECT " + table + ".id, userid, score, pp, username, max_combo, full_combo, mods, 300_count, 100_count, 50_count, katus_count, gekis_count, misses_count, time, play_mode, beatmap_md5 FROM " + table + " LEFT JOIN users ON users.id = userid WHERE " + table + ".id = ? AND completed = 3", i).Scan(
			&s.ID, &s.UserID, &s.Score, &s.Performance, &s.Username, &s.Combo, &s.FullCombo, &s.Mods, &s.N300, &s.N100, &s.N50, &s.NKatu, &s.NGeki, &s.NMiss, &s.Timestamp, &mode, &md5,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}

			log.Error(err)
			continue
		}

		lbp := leaderboards.Cache.Leaderboards[leaderboards.Identifier{md5, mode, rx}]
		if lbp != nil {
			lbp.AddScore(s)
			log.Info("Added Score", i, "to", md5)
		}
	}
}

func changeUsername(client *redis.Client) {
	ps, err := client.Subscribe("api:change_username")
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Subscribed to api:change_username")

	for {
		msg, err := ps.ReceiveMessage()
		if err != nil {
			log.Error(err)
			return
		}

		i, err := strconv.Atoi(msg.Payload)
		if err != nil {
			log.Error(err)
			continue
		}

		var newUsername string
		err = config.DB.QueryRow("SELECT username FROM users WHERE id = ?", i).Scan(&newUsername)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}

			log.Error(err)
			continue
		}

		for _, lb := range leaderboards.Cache.Leaderboards {
			if s, _ := lb.FindUserScore(i); s != nil {
				s.Username = newUsername
			}
		}
	}
}