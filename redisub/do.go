package redisub

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/infernalfire72/acache/config"
	"github.com/infernalfire72/acache/leaderboards"
	"github.com/infernalfire72/acache/log"
	"github.com/infernalfire72/acache/tools"
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

		leaderboards.RemoveUser(i)
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

		leaderboards.AddUser(i)
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

		ss := strings.SplitN(msg.Payload, ",", 3)
		if len(ss) != 3 {
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

		gm, err := strconv.Atoi(ss[2])
		if err != nil {
			log.Error(err)
			continue
		}

		leaderboards.RemoveUserWithIdentifier(i, rx, byte(gm))
		log.Info("[WIPE] Wiped Cached Scores for", i, rx, byte(gm))
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

		sw := tools.Stopwatch{}
		sw.Start()
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
			md5  string
			mode byte
		)
		err = config.DB.QueryRow("SELECT "+table+".id, userid, score, pp, COALESCE(CONCAT('[', tag, '] ', username), username) AS username, max_combo, full_combo, mods, 300_count, 100_count, 50_count, katus_count, gekis_count, misses_count, time, play_mode, beatmap_md5 FROM "+table+" LEFT JOIN users ON users.id = userid LEFT JOIN clans ON clans.id = users.clan_id WHERE "+table+".id = ? AND (completed & 7) >= 3", i).Scan(
			&s.ID, &s.UserID, &s.Score, &s.Performance, &s.Username, &s.Combo, &s.FullCombo, &s.Mods, &s.N300, &s.N100, &s.N50, &s.NKatu, &s.NGeki, &s.NMiss, &s.Timestamp, &mode, &md5,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}

			log.Error(err)
			continue
		}

		lb := leaderboards.Get(leaderboards.Identifier{md5, mode, rx})
		if lb != nil {
			lb.AddScore(s)
			log.Info("Added Score", i, "to", md5)
		} else {
			log.Info("lb nil")
		}
		sw.Stop()
		log.Infof("Score took %s", sw.ElapsedReadable())
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

		id, err := strconv.Atoi(msg.Payload)
		if err != nil {
			log.Error(err)
			continue
		}

		var newUsername string
		err = config.DB.QueryRow("SELECT username FROM users WHERE id = ?", id).Scan(&newUsername)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}

			log.Error(err)
			continue
		}

		leaderboards.ChangeUsername(id, newUsername)
		log.Info("Changed Username", id, "->", newUsername)
	}
}
