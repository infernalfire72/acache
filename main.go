package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/infernalfire72/acache/api"
	"github.com/infernalfire72/acache/beatmaps"
	"github.com/infernalfire72/acache/config"
	"github.com/infernalfire72/acache/leaderboards"
	"github.com/infernalfire72/acache/log"
	"github.com/infernalfire72/acache/redisub"
)

func init() {
	conf, err := config.Load()
	if err != nil {
		log.Error(err)
		return
	}

	config.DB, err = sql.Open("mysql", conf.Database.String())
	if err != nil {
		log.Error(err)
		return
	}

	beatmaps.Init()
	leaderboards.Init()
}

func main() {
	err := config.DB.Ping()
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Connection to Database established")
	redisub.Subscribe()
	api.Start()
}
