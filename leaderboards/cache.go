package leaderboards

import (
	"sync"

	"github.com/infernalfire72/acache/config"
	"github.com/infernalfire72/acache/log"
)

type Identifier struct {
	Md5   string
	Mode  byte
	Relax bool
}

var lMutex sync.RWMutex
var Cache *LeaderboardCache

type LeaderboardCache struct {
	Leaderboards map[Identifier]*Leaderboard
}

func (c *LeaderboardCache) Get(identifier Identifier) *Leaderboard {
	lMutex.RLock()
	lbp := c.Leaderboards[identifier]
	lMutex.RUnlock()
	if lbp != nil {
		return lbp
	} else {
		return c.UpdateCache(identifier)
	}
}

func (c *LeaderboardCache) UpdateCache(identifier Identifier) *Leaderboard {
	lb := &Leaderboard{
		BeatmapMd5: identifier.Md5,
		Mode:       identifier.Mode,
		Relax:      identifier.Relax,
	}
	lb.UpdateCache()
	lMutex.Lock()
	c.Leaderboards[identifier] = lb
	lMutex.Unlock()
	return lb
}

func (c *LeaderboardCache) Clear() {
	lMutex.Lock()
	c.Leaderboards = make(map[Identifier]*Leaderboard)
	lMutex.Unlock()
}

func (c *LeaderboardCache) RemoveUser(id int) {
	lMutex.RLock()
	for _, a := range c.Leaderboards {
		a.RemoveUser(id)
	}
	lMutex.RUnlock()
}

// For Wipe
func (c *LeaderboardCache) RemoveUserWithIdentifier(id int, rx bool) {
	lMutex.RLock()
	for _, a := range c.Leaderboards {
		if a.Relax == rx {
			a.RemoveUser(id)
		}
	}
	lMutex.RUnlock()
}

func (c *LeaderboardCache) AddUser(id int) {
	for i, a := range [...]string{"scores", "scores_relax"} {
		var relax bool
		if i == 1 {
			relax = true
		}

		rows, err := config.DB.Query("SELECT "+a+".id, userid, score, pp, COALESCE(CONCAT('[', tag, '] ', username), username) AS username, max_combo, full_combo, mods, 300_count, 100_count, 50_count, katus_count, gekis_count, misses_count, time, play_mode, beatmap_md5 FROM "+a+" LEFT JOIN users ON users.id = userid LEFT JOIN clans ON clans.id = users.clan_id WHERE userid = ? AND completed = 3", id)
		if err != nil {
			log.Error(err)
		}
		defer rows.Close()

		for rows.Next() {
			s := &Score{}
			var (
				md5  string
				mode byte
			)
			err = rows.Scan(&s.ID, &s.UserID, &s.Score, &s.Performance, &s.Username, &s.Combo, &s.FullCombo, &s.Mods, &s.N300, &s.N100, &s.N50, &s.NKatu, &s.NGeki, &s.NMiss, &s.Timestamp, &mode, &md5)
			if err != nil {
				log.Error(err)
			}

			lMutex.RLock()
			lbp := c.Leaderboards[Identifier{md5, mode, relax}]
			if lbp != nil {
				lbp.AddScore(s)
			}
			lMutex.RUnlock()
		}
	}
}
