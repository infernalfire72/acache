package leaderboards

import "fmt"

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