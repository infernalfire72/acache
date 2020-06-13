package main

import (
	"fmt"
	"sort"
)

type Score struct {
	ID			int
	UserID		int
	Username	string
	Score		int
	Performance	float32
	Combo		int
	FullCombo	bool
	N50			int16
	N100		int16
	N300		int16
	NMiss		int16
	NKatu		int16
	NGeki		int16
	Mods		int
	Timestamp	int
}

func (s *Score) String(displayScore bool, pos int) string {
	lbScore := int(s.Performance)
	if displayScore {
		lbScore = s.Score
	}

	fc := "False"
	if s.FullCombo {
		fc = "True"
	}
	return fmt.Sprintf("%d|%s|%d|%d|%d|%d|%d|%d|%d|%d|%s|%d|%d|%d|%d|1\n", s.ID, s.Username, lbScore, s.Combo, s.N50, s.N100, s.N300, s.NMiss, s.NKatu, s.NGeki, fc, s.Mods, s.UserID, pos, s.Timestamp)
}

type Leaderboard struct {
	BeatmapMd5	string
	Scores		[]Score
	Mode		byte
	Relax		bool
}

func (l *Leaderboard) Map() *Beatmap {
	return bcache.Get(l.BeatmapMd5)
}

func (l *Leaderboard) Sort() {
	m := *(l.Map())
	sort.Slice(l.Scores, func(i, j int) bool {
		if !l.Relax || m.Status == Loved {
			return l.Scores[i].Score > l.Scores[j].Score
		} else {
			return l.Scores[i].Performance > l.Scores[j].Performance
		}
	})
}

func (l *Leaderboard) AddScore(s Score) {
	l.Scores = append(l.Scores, s)
	l.Sort()
}

func (l *Leaderboard) RemoveScore(id int) {
	for i := 0; i < len(l.Scores); i++ {
		if l.Scores[i].ID == id {
			l.Scores = append(l.Scores[:i], l.Scores[i+1:]...)
			break
		}
	}
}

func (l *Leaderboard) UpdateCache() {
	if l.Map().Status < Ranked {
		return
	}

	l.Scores = make([]Score, 0)

	var table string
	if l.Relax {
		table = "_relax"
	}

	tableSort := "score"
	if l.Relax && l.Map().Status != Loved {
		tableSort = "pp"
	}

	rows, err := db.Query("SELECT scores" + table + ".id, userid, score, pp, username, max_combo, full_combo, mods, 300_count, 100_count, 50_count, katus_count, gekis_count, misses_count, time FROM scores" + table + " LEFT JOIN users ON users.id = userid WHERE beatmap_md5 = ? AND completed = 3 AND play_mode = ? AND users.privileges & 1 ORDER BY "+ tableSort +" DESC", l.BeatmapMd5, l.Mode)
	if err != nil {
		log.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		var s Score
		err = rows.Scan(&s.ID, &s.UserID, &s.Score, &s.Performance, &s.Username, &s.Combo, &s.FullCombo, &s.Mods, &s.N300, &s.N100, &s.N50, &s.NKatu, &s.NGeki, &s.NMiss, &s.Timestamp)
		if err != nil {
			log.Error(err)
		}
		l.Scores = append(l.Scores, s)
	}
}

func (l *Leaderboard) FindUserScore(user int) (*Score, int) {
	for i, score := range l.Scores {
		if score.UserID == user {
			return &score, i
		}
	}
	return nil, -1
}

type LBIdentifier struct {
	string
	byte
	bool
}

type LeaderboardCache struct {
	Leaderboards	map[LBIdentifier]*Leaderboard
}

func (c *LeaderboardCache) Get(identifier LBIdentifier) *Leaderboard {
	lbp := c.Leaderboards[identifier]

	if lbp != nil {
		return lbp
	} else {
		return c.UpdateCache(identifier)
	}
}

func (c *LeaderboardCache) UpdateCache(identifier LBIdentifier) *Leaderboard {
	lb := &Leaderboard{
		BeatmapMd5:	identifier.string,
		Mode:		identifier.byte,
		Relax:		identifier.bool,
	}
	lb.UpdateCache()
	c.Leaderboards[identifier] = lb
	return lb
}

func (c LeaderboardCache) Clear() {
	c.Leaderboards = make(map[LBIdentifier]*Leaderboard)
}