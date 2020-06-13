package leaderboards

import (
	"sync"
	
	"github.com/infernalfire72/acache/config"
	"github.com/infernalfire72/acache/log"
)

type Identifier struct {
	Md5		string
	Mode	byte
	Relax	bool
}

var lmutex 	sync.Mutex
var Cache 	*LeaderboardCache

type LeaderboardCache struct {
	Leaderboards	map[Identifier]*Leaderboard
}

func (c *LeaderboardCache) Get(identifier Identifier) *Leaderboard {
	lbp := c.Leaderboards[identifier]

	if lbp != nil {
		return lbp
	} else {
		return c.UpdateCache(identifier)
	}
}

func (c *LeaderboardCache) UpdateCache(identifier Identifier) *Leaderboard {
	lb := &Leaderboard{
		BeatmapMd5:	identifier.Md5,
		Mode:		identifier.Mode,
		Relax:		identifier.Relax,
	}
	lb.UpdateCache()
	lmutex.Lock()
	c.Leaderboards[identifier] = lb
	lmutex.Unlock()
	return lb
}

func (c *LeaderboardCache) Clear() {
	c.Leaderboards = make(map[Identifier]*Leaderboard)
}

func (c *LeaderboardCache) RemoveUser(id int) {
	for _, a := range c.Leaderboards {
		a.RemoveUser(id)
	}
}

// For Wipe
func (c *LeaderboardCache) RemoveUserWithIdentifier(id int, rx bool) {
	for _, a := range c.Leaderboards {
		if a.Relax == rx {
			a.RemoveUser(id)
		}
	}
}

func (c *LeaderboardCache) AddUser(id int) {
	for i, a := range [...]string{"scores", "scores_relax"} {
		var relax bool
		if i == 1 {
			relax = true
		}

		rows, err := config.DB.Query("SELECT " + a + ".id, userid, score, pp, username, max_combo, full_combo, mods, 300_count, 100_count, 50_count, katus_count, gekis_count, misses_count, time, play_mode, beatmap_md5 FROM "+ a +" LEFT JOIN users ON users.id = userid WHERE userid = ? AND completed = 3", id)
		if err != nil {
			log.Error(err)
		}
		defer rows.Close()

		for rows.Next() {
			s := &Score{}
			var (
				md5		string
				mode	byte
			)
			err = rows.Scan(&s.ID, &s.UserID, &s.Score, &s.Performance, &s.Username, &s.Combo, &s.FullCombo, &s.Mods, &s.N300, &s.N100, &s.N50, &s.NKatu, &s.NGeki, &s.NMiss, &s.Timestamp, &mode, &md5)
			if err != nil {
				log.Error(err)
			}

			lbp := c.Leaderboards[Identifier{md5, mode, relax}]
			if lbp != nil {
				lbp.AddScore(s)
			}
		}
	}
}