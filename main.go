package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/infernalfire72/acache/api"
	"github.com/infernalfire72/acache/config"
	"github.com/infernalfire72/acache/leaderboards"
	"github.com/infernalfire72/acache/beatmaps"
	"github.com/infernalfire72/acache/log"
)

func init() {
	var err error
	config.DB, err = sql.Open("mysql", "root:lol123@/ripple")
	if err != nil {
		log.Error(err)
	}

	beatmaps.Cache = &beatmaps.BeatmapCache{make(map[string]*beatmaps.Beatmap)}
	leaderboards.Cache = &leaderboards.LeaderboardCache{ make(map[leaderboards.Identifier]*leaderboards.Leaderboard) }
}

func main() {
	err := config.DB.Ping()
	if err != nil {
		log.Error(err)
	}
	log.Info("Connection to Database established")

	api.Start()
}