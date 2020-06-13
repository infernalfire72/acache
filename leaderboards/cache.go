package leaderboards

type Identifier struct {
	Md5		string
	Mode	byte
	Relax	bool
}

var Cache *LeaderboardCache

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
	c.Leaderboards[identifier] = lb
	return lb
}

func (c LeaderboardCache) Clear() {
	c.Leaderboards = make(map[Identifier]*Leaderboard)
}