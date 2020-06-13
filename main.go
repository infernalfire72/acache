package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

var lcache *LeaderboardCache
var bcache *BeatmapCache

func init() {
	var err error
	db, err = sql.Open("mysql", "root:lol123@/ripple")
	if err != nil {
		log.Error(err)
	}

	bcache = &BeatmapCache{make(map[string]*Beatmap)}
	lcache = &LeaderboardCache{ make(map[LBIdentifier]*Leaderboard) }
}

func main() {
	err := db.Ping()
	if err != nil {
		log.Error(err)
	}
	log.Info("Connection to Database established")

	// Start da Api
	Start()
}