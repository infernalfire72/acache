package main

import "fmt"

type BeatmapStatus int

const (
	Unknown BeatmapStatus = iota - 2
	NotSubmitted
	Pending
	NeedsUpdate
	Ranked
	Approved
	Qualified
	Loved
)

type Beatmap struct {
	Md5			string
	ID			int
	SetID		int
	Name		string
	Status		BeatmapStatus
	Playcount	int
	Passcount	int
}

func (b *Beatmap) String(scoresCount int) string {
	return fmt.Sprintf("%d|false|%d|%d|%d\n0\n%s\n10.0\n", b.Status, b.ID, b.SetID, scoresCount, b.Name)
}

type BeatmapCache struct {
	Beatmaps	map[string]*Beatmap
}

func (c *BeatmapCache) Get(md5 string) *Beatmap {
	bmp := c.Beatmaps[md5]

	if bmp != nil {
		return bmp
	} else {
		return c.UpdateCache(md5)
	}
}

func (c *BeatmapCache) UpdateCache(md5 string) *Beatmap {
	b := &Beatmap{
		Md5:	md5,
	}

	err := db.QueryRow("SELECT beatmap_id, beatmapset_id, song_name, ranked, playcount, passcount FROM beatmaps WHERE beatmap_md5 = ?", md5).Scan(
		&b.ID, &b.SetID, &b.Name, &b.Status, &b.Playcount, &b.Passcount,
	)

	if err != nil {
		log.Error(err)
	}

	c.Beatmaps[md5] = b
	return b
}

func (c BeatmapCache) Clear() {
	c.Beatmaps = make(map[string]*Beatmap)
}