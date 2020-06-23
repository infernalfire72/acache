package leaderboards

import (
	"sort"
	"sync"

	"github.com/infernalfire72/acache/beatmaps"
	"github.com/infernalfire72/acache/config"
	"github.com/infernalfire72/acache/log"
)

type Leaderboard struct {
	BeatmapMd5 string
	Scores     []Score
	Mode       byte
	Relax      bool
	Mutex      sync.RWMutex
}

func (l *Leaderboard) Map() *beatmaps.Beatmap {
	return beatmaps.Cache.Get(l.BeatmapMd5)
}

func (l *Leaderboard) Sort() {
	m := l.Map()
	if m == nil {
		return
	}
	l.Mutex.Lock()
	sort.Slice(l.Scores, func(i, j int) bool {
		if !l.Relax || m.Status == beatmaps.Loved {
			return l.Scores[i].Score > l.Scores[j].Score
		} else {
			return l.Scores[i].Performance > l.Scores[j].Performance
		}
	})
	l.Mutex.Unlock()
}

func (l *Leaderboard) AddScore(s *Score) {
	l.Mutex.Lock()
	for i, a := range l.Scores {
		if a.ID == s.ID {
			return
		}

		if a.UserID == s.UserID {
			l.RemoveScoreIndex(i)
		}
	}

	l.Scores = append(l.Scores, *s)
	l.Mutex.Unlock()
	l.Sort()
}

func (l *Leaderboard) RemoveScoreIndex(i int) {
	l.Scores = append(l.Scores[:i], l.Scores[i+1:]...)
}

func (l *Leaderboard) RemoveScore(id int) {
	l.Mutex.Lock()
	for i, a := range l.Scores {
		if a.ID == id {
			l.RemoveScoreIndex(i)
			break
		}
	}
	l.Mutex.Unlock()
}

func (l *Leaderboard) RemoveUser(id int) {
	l.Mutex.Lock()
	for i, a := range l.Scores {
		if a.UserID == id {
			l.RemoveScoreIndex(i)
			break
		}
	}
	l.Mutex.Unlock()
}

func (l *Leaderboard) UpdateCache() {
	if l.Map().Status < beatmaps.Ranked {
		return
	}

	l.Scores = make([]Score, 0)

	table := "scores"
	if l.Relax {
		table = "scores_relax"
	}

	tableSort := "score"
	if l.Relax && l.Map().Status != beatmaps.Loved {
		tableSort = "pp"
	}

	rows, err := config.DB.Query("SELECT "+table+".id, userid, score, pp, COALESCE(CONCAT('[', tag, '] ', username), username) AS username, max_combo, full_combo, mods, 300_count, 100_count, 50_count, katus_count, gekis_count, misses_count, time FROM "+table+" LEFT JOIN users ON users.id = userid LEFT JOIN clans ON clans.id = users.clan_id WHERE beatmap_md5 = ? AND completed = 3 AND play_mode = ? AND users.privileges & 1 ORDER BY "+tableSort+" DESC", l.BeatmapMd5, l.Mode)
	if err != nil {
		log.Error(err)
	}
	defer rows.Close()

	l.Mutex.Lock()
	for rows.Next() {
		var s Score
		err = rows.Scan(&s.ID, &s.UserID, &s.Score, &s.Performance, &s.Username, &s.Combo, &s.FullCombo, &s.Mods, &s.N300, &s.N100, &s.N50, &s.NKatu, &s.NGeki, &s.NMiss, &s.Timestamp)
		if err != nil {
			log.Error(err)
		}
		l.Scores = append(l.Scores, s)
	}
	l.Mutex.Unlock()
}

func (l *Leaderboard) FindUserScore(user int) (*Score, int) {
	l.Mutex.RLock()
	defer l.Mutex.RUnlock()
	for i, score := range l.Scores {
		if score.UserID == user {
			return &score, i
		}
	}
	return nil, -1
}