package beatmaps

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/infernalfire72/acache/config"
	"github.com/infernalfire72/acache/log"
)

type BeatmapStatus int

const (
	StatusUnknown BeatmapStatus = iota - 2
	StatusNotSubmitted
	StatusPending
	StatusNeedsUpdate
	StatusRanked
	StatusApproved
	StatusQualified
	StatusLoved
)

type Beatmap struct {
	Md5        string
	ID         int
	SetID      int
	Name       string
	Status     BeatmapStatus
	Playcount  int
	Passcount  int
	LastUpdate time.Time
}

func (b *Beatmap) FetchFromDb() {
	err := config.DB.QueryRow("SELECT beatmap_id, beatmapset_id, song_name, ranked, playcount, passcount FROM beatmaps WHERE beatmap_md5 = ?", b.Md5).Scan(
		&b.ID, &b.SetID, &b.Name, &b.Status, &b.Playcount, &b.Passcount,
	)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
	}

	b.LastUpdate = time.Now()
}

func (b *Beatmap) String(scoresCount int) string {
	return fmt.Sprintf("%d|false|%d|%d|%d\n0\n%s\n10.0\n", b.Status, b.ID, b.SetID, scoresCount, b.Name)
}
